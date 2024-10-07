package videoproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"gopkg.in/vansante/go-ffprobe.v2"
	"image"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var FontFilePath = "static/NotoSans-Regular.ttf"

func secToTimeStr(sec int) string {
	hour := sec / 3600
	sec = sec % 3600
	minute := sec / 60
	sec = sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, sec)
}

func timeStrToSec(time string) int {
	var hour, minute, sec int
	_, err := fmt.Sscanf(time, "%d:%d:%d", &hour, &minute, &sec)
	if err != nil {
		return 0
	}
	return hour*3600 + minute*60 + sec
}

func escapeFfmpegCmd(s string) string {
	//FIXME: 这里转义并不完善，但是考虑到我也就只用到了冒号，其实也没啥问题
	return strings.Replace(s, `:`, `\:`, -1)
}

func GetScreenshotAtSec(videoPath string, sec int) ([]byte, error) {
	time := secToTimeStr(sec)
	videoFilter := fmt.Sprintf(
		`drawtext=text='%s':fontfile='%s':fontcolor=white:fontsize=h/8:x=10:y=10:box=1:boxcolor=black@0.6:boxborderw=5`,
		escapeFfmpegCmd(time), FontFilePath)
	cmd := exec.Command("ffmpeg",
		"-y",
		"-loglevel", "error",
		"-ss", time,
		"-i", videoPath,
		"-vf", videoFilter,
		"-vframes", "1",
		"-c:v", "mjpeg",
		"-f", "image2pipe",
		"-",
	)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitErr.Stderr = stderr.Bytes()
		}
		return nil, err
	}
	return stdout.Bytes(), nil
}

func MakeScreenShotTile(videoPath string, tileWidth, tileHeight int) ([]byte, error) {
	probe, err := ffprobe.ProbeURL(context.Background(), videoPath)
	if err != nil {
		return nil, err
	}
	totalDur := probe.Format.DurationSeconds
	count := tileWidth * tileHeight
	// 不截图前后 2% 的时间，避免首帧和尾帧可能的黑屏
	shotDur := int(totalDur * 0.96)
	offset := int(totalDur * 0.02)
	screenshot, err := GetScreenshotAtSec(videoPath, offset)
	if err != nil {
		return nil, err
	}
	img, err := imaging.Decode(bytes.NewBuffer(screenshot))
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	canvas := imaging.New(tileWidth*bounds.Dx(), tileHeight*bounds.Dy(), color.NRGBA{})
	//draw img to canvas
	canvas = imaging.Paste(canvas, img, image.Pt(0, 0))
	for i := 1; i < count; i++ {
		sec := offset + i*shotDur/count
		screenshot, err = GetScreenshotAtSec(videoPath, sec)
		if err != nil {
			return nil, err
		}
		img, err = imaging.Decode(bytes.NewBuffer(screenshot))
		if err != nil {
			return nil, err
		}
		canvas = imaging.Paste(canvas, img, image.Pt((i%tileWidth)*bounds.Dx(), (i/tileWidth)*bounds.Dy()))
	}
	const maxLen = 1280
	if canvas.Bounds().Dx() > maxLen || canvas.Bounds().Dy() > maxLen {
		if canvas.Bounds().Dx() > canvas.Bounds().Dy() {
			canvas = imaging.Resize(canvas, maxLen, 0, imaging.Lanczos)
		} else {
			canvas = imaging.Resize(canvas, 0, maxLen, imaging.Lanczos)
		}
	}
	buf := bytes.NewBuffer(nil)
	err = imaging.Encode(buf, canvas, imaging.JPEG)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MakeScreenShotTileFile(videoPath string, tileWidth, tileHeight int) (string, error) {
	buf, err := MakeScreenShotTile(videoPath, tileWidth, tileHeight)
	if err != nil {
		return "", err
	}
	f, err := os.CreateTemp("", "screenshot*.jpg")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(buf)
	if err != nil {
		return "", err
	}
	return filepath.Abs(f.Name())
}
