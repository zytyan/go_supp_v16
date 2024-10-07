#!/usr/bin/env bash
set -e

FILE_DIR=$(dirname "$0")
FILE_DIR=$(realpath "$FILE_DIR")
OPT_YES=false
BIN_NAME="suppbot"
TMP_DIR="/tmp"

check_commands() {
    local missing_cmds=""

    for cmd in "$@"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing_cmds="$missing_cmds\n$cmd"
        fi
    done
    if [ -n "$missing_cmds" ]; then
        echo -e "The following commands are missing:\n$missing_cmds"
        return 1
    else
        echo "All commands are available."
        return 0
    fi
}

check_yes() {
    if [ "$OPT_YES" = true ]; then
        return 0
    fi

    read -rp "$1" yn
    case "$yn" in
    [$2]*) return 0 ;;
    *) return 1 ;;
    esac
}
build() {
    go get
    go build -o "$BIN_NAME"
    echo "Compile done at" "$(date '+%Y-%m-%d %H:%M:%S')"
}
compile_and_restart() {
    build
    # 检查是否传入了 -y 参数
    if check_yes "是否重启服务 [y/N] " "yY"; then
        sudo systemctl daemon-reload
        sudo systemctl restart "$BIN_NAME"
    fi
}

install() {
    if ! check_commands "go" "systemctl" "sed" "ffmpeg" "ffprobe" "rar"; then
        return 1
    fi
    build
    cp "$FILE_DIR/$BIN_NAME.service.temp" "$FILE_DIR/$BIN_NAME.service"
    sudo sed -i "s|VAR_CUR_PATH|$FILE_DIR|g" "$FILE_DIR/$BIN_NAME.service"
    sudo sed -i "s|VAR_TMP_DIR|$TMP_DIR|g" "$FILE_DIR/$BIN_NAME.service"
    sudo sed -i "s|VAR_BIN_NAME|$BIN_NAME|g" "$FILE_DIR/$BIN_NAME.service"
    sudo cp "$FILE_DIR/$BIN_NAME.service" "/etc/systemd/system/$BIN_NAME.service"
    sudo systemctl daemon-reload
    sudo systemctl enable "$BIN_NAME"
    # check group exists
    if ! getent group tgbots >/dev/null; then
        sudo groupadd -r tgbots
    fi
    # check user exists
    if ! getent passwd tgbotapi >/dev/null; then
        sudo useradd -rNM -s /bin/false -d /nonexistent -g tgbotapi -c "Telegram Bots" tgbotapi
    fi

    sudo chown -R tgbotapi:tgbots "$(dirname "$FILE_DIR/..")"
}

check_opts() {
    # 遍历所有命令行参数，直到找到选项
    while [[ "$1" != "" ]]; do
      if [[ "$1" == -* ]]; then
        # 遇到选项，跳出循环，交给 getopts 处理
        break
      fi
      # 如果不是选项，跳过这个参数
      shift
    done
    # -y -t tmp_dir
    while getopts "yt:" opt; do
        case $opt in
        y)
            OPT_YES=true
            ;;
        t)
            TMP_DIR=$OPTARG
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            ;;
        esac
    done
}

main() {
    # 检查各种参数
    check_opts "$@"
    echo "OPT_YES = $OPT_YES"
    echo "TMP_DIR = $TMP_DIR"
    case "$1" in
    restart)
        compile_and_restart "${@:2}"
        ;;
    install)
        install "${@:2}"
        ;;
    *)
        echo "Usage: $0 {restart|install} [-y] [-t tmp_dir]"
        exit 1
        ;;
    esac
}

# 调用函数并传递所有参数
main "$@"
