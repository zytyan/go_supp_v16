package tphup

import (
	"fmt"
	"github.com/celestix/telegraph-go/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestEnsureImageSize(t *testing.T) {
	as := assert.New(t)
	const bigImg = "testdata/big.jpg"
	buf, err := ensureImageCanBeUploaded(bigImg)
	as.Nil(err)
	as.NotNil(buf)
	f, err := os.ReadFile(bigImg)
	as.Nil(err)
	as.NotEqual(f, buf)
	const smallImg = "testdata/small.jpg"
	buf, err = ensureImageCanBeUploaded(smallImg)
	as.Nil(err)
	as.NotNil(buf)
	f, err = os.ReadFile(smallImg)
	as.Nil(err)
	as.Equal(f, buf)
	buf, err = ensureImageCanBeUploaded("testdata/yellow_rose.lossy.webp")
	as.Nil(err)
	as.NotNil(buf)
	as.Equal(buf[:2], []byte{0xff, 0xd8}) // JPEG magic number
	as.Equal(buf[len(buf)-2:], []byte{0xff, 0xd9})
}

func TestUploadFolder(t *testing.T) {
	as := assert.New(t)
	const folder = "testdata"
	client := &telegraph.TelegraphClient{
		ApiUrl:     "https://api.telegra.ph/",
		HttpClient: &http.Client{},
	}
	images, err := UploadFolder(client, "9a887e31fdf662259b8d7911ed2263e6c29b6461e6e9919e8c91c698b3a7", folder)
	as.Nil(err)
	for img := range images {
		as.NotEmpty(img.Path)
		fmt.Println(img.Path)
	}
}
