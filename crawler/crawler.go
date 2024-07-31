package crawler

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
)

func getHost() string {
	// read liuli.link file
	s, err := os.ReadFile("liuli.link")
	if err != nil {
		return "www.hacg.mov"
	}
	return string(s)
}
func setHost(host string) {
	// write liuli.link file
	_ = os.WriteFile("liuli.link", []byte(host), 0644)
}
func getHtml(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func getHomePage() ([]byte, error) {
	link := fmt.Sprintf("https://%s/wp/", getHost())
	return getHtml(link)
}

type Article struct {
	Url      string
	Title    string
	Author   string
	PostTime string
	ImgUrl   string
	Category string
	Tags     []string
}

func parseOneArticle(sel *goquery.Selection) (Article, error) {
	h1 := sel.Find("h1 > a").First()
	url, ok := h1.Attr("href")
	if !ok {
		return Article{}, fmt.Errorf("href not found")
	}
	title := h1.Text()
	author := sel.Find("span.author.vcard").Text()
	postTime := sel.Find("time").Text()
	imgUrl, _ := sel.Find("img").First().Attr("src")
	category := sel.Find("span.cat-links > a").Text()
	tags := sel.Find("span.tags-links > a").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	return Article{
		Url:      url,
		Title:    title,
		Author:   author,
		PostTime: postTime,
		ImgUrl:   imgUrl,
		Category: category,
		Tags:     tags,
	}, nil
}

func parseHomePageArticles(buf []byte) ([]Article, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	articles := make([]Article, 0)
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		article, err := parseOneArticle(s)
		if err != nil {
			return
		}
		articles = append(articles, article)
	})
	return articles, nil
}
func GetArticles() ([]Article, error) {
	buf, err := getHomePage()
	if err != nil {
		return nil, err
	}
	return parseHomePageArticles(buf)
}
