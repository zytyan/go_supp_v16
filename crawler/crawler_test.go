package crawler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetArticles(t *testing.T) {
	as := assert.New(t)
	articles, err := GetArticles()
	as.NoError(err)
	as.NotEmpty(articles)
}

func TestGetMagnetsFromLink(t *testing.T) {
	as := assert.New(t)
	magnets, err := GetMagnetsFromLink("https://www.hacg.mov/wp/99360.html")
	as.NoError(err)
	as.NotEmpty(magnets)
	as.Contains(magnets, "eb074e1e5840c3499b475514a9fd19246ee0ce2b")
}

func TestGetMagnetsFromHtml(t *testing.T) {
	as := assert.New(t)
	html := []byte(`<html><body>
<script>abcd1e5840c3499b475514a9fd19246ee0ce2c</script>
01234e1e5840c3499b475514a9fd19246ee0ce2b this is other text
01234e1e5840c3499b475514a9fd19246ee0ce2c11 this len &gt; 40</body></html>`)
	magnets, err := getMagnetsFromHtml(html)
	as.NoError(err)
	as.NotEmpty(magnets)

	as.Contains(magnets, "01234e1e5840c3499b475514a9fd19246ee0ce2b")
	as.NotContains(magnets, "abcd1e5840c3499b475514a9fd19246ee0ce2c")
	as.NotContains(magnets, "01234e1e5840c3499b475514a9fd19246ee0ce2c")
	as.NotContains(magnets, "234e1e5840c3499b475514a9fd19246ee0ce2c11")
	as.NotContains(magnets, "01234e1e5840c3499b475514a9fd19246ee0ce2c11")
	as.Len(magnets, 1)
	fmt.Println(magnets)

}

func TestParseHomePageArticles(t *testing.T) {
	as := assert.New(t)
	html, err := os.ReadFile("test.html")
	as.NoError(err)
	articles, err := parseHomePageArticles(html)
	as.NoError(err)
	as.NotEmpty(articles)
	as.Equal(articles[0].IdTag(), "#wp99401")
	as.Equal(articles[0].Url, "https://www.hacg.mov/wp/99401.html")
	as.Equal(articles[0].Title, "[えるぴーすたじお] C→I→Mカップ! ～あの娘のヒミツはミルク味!?～")
	as.Equal(articles[0].Author, "多啦H萌")
	as.Equal(articles[0].PostTime, "2024年8月7日")
	as.Contains(articles[0].Tags, "全彩")
	as.Contains(articles[0].Tags, "学园")
	as.NotContains(articles[0].Tags, "标签为")
}
