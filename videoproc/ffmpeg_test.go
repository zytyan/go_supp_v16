package videoproc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"testing"
)

const videoPath = `testdata/test.mp4`

func TestGetScreenshotAtTime(t *testing.T) {
	as := assert.New(t)
	const time = "00:00:11"
	buf, err := GetScreenshotAtSec(videoPath, timeStrToSec(time))
	as.Nil(err)
	as.NotEmpty(buf)
	as.Equal(buf[:2], []byte{0xff, 0xd8}) // JPEG magic number
	as.Equal(buf[len(buf)-2:], []byte{0xff, 0xd9})
}

func TestMakeScreenShotTile(t *testing.T) {
	as := assert.New(t)
	const tileWidth = 3
	const tileHeight = 3
	buf, err := MakeScreenShotTile(videoPath, tileWidth, tileHeight)
	as.Nil(err)
	as.NotEmpty(buf)
	as.Equal(buf[:2], []byte{0xff, 0xd8}) // JPEG magic number
	as.Equal(buf[len(buf)-2:], []byte{0xff, 0xd9})
}

func TestMain(m *testing.M) {
	FontFilePath = `../static/NotoSans-Regular.ttf`
	os.Exit(m.Run())
}

func TestProbe(t *testing.T) {
	as := assert.New(t)
	probe, err := ffprobe.ProbeURL(context.Background(), videoPath)
	as.Nil(err)
	as.NotNil(probe)
	as.NotEmpty(probe.Format)
	as.NotEmpty(probe.Format.DurationSeconds)
	as.NotNil(probe.FirstVideoStream())
	v := probe.FirstVideoStream()
	as.Equal("25.358333", v.Duration)
	as.Equal("h264", v.CodecName)
	as.Equal(1920, v.Width)
	as.Equal(1080, v.Height)
	as.Equal("yuv420p", v.PixFmt)

}
