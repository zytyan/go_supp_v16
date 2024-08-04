package videoproc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"os/exec"
	"strings"
)

var FontFilePath = "static/NotoSans-Regular.ttf"

func SecToTimeStr(sec int) string {
	hour := sec / 3600
	sec = sec % 3600
	minute := sec / 60
	sec = sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, sec)
}

func TimeStrToSec(time string) int {
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
	time := SecToTimeStr(sec)
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

func GetDuration(videoPath string) (int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
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
		return 0, err
	}
	var duration int
	_, err = fmt.Fscanf(stdout, "%d", &duration)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func MakeScreenShotTile(videoPath string, tileWidth, tileHeight int) ([]byte, error) {
	duration, err := GetDuration(videoPath)
	if err != nil {
		return nil, err
	}
	count := tileWidth * tileHeight
	// 不截图前后 5% 的时间，避免首帧和尾帧可能的黑屏
	durF := float64(duration)
	duration = int(durF * 0.9)
	offset := int(durF * 0.05)
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
		sec := offset + i*duration/count
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
