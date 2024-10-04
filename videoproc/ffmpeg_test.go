package videoproc

import (
	"github.com/stretchr/testify/assert"
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

func TestGetDuration(t *testing.T) {
	as := assert.New(t)
	duration, err := GetDuration(videoPath)
	as.Nil(err)
	as.Equal(duration, 25)
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

func TestGetSize(t *testing.T) {
	as := assert.New(t)
	w, h, err := GetSize(videoPath)
	as.Nil(err)
	as.Equal(w, 1920)
	as.Equal(h, 1080)
}
