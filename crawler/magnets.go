package crawler

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

var magnetRegexp = regexp.MustCompile("[a-zA-Z0-9]{40}")

func findMagnets(html string) []string {
	magnets := magnetRegexp.FindAllString(html, -1)
	for i := range magnets {
		magnets[i] = strings.ToLower(magnets[i])
	}
	return magnets
}

func GetMagnetsFromLink(link string) ([]string, error) {
	html, err := getHtml(link)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}
	text := doc.Find("body").Text()
	return findMagnets(text), nil
}
