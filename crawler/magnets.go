package crawler

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

var magnetRegexp = regexp.MustCompile("[a-zA-Z0-9]{40}")

func isHexDigit(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func isMagnetBoundaryValid(start, end int, html string) bool {
	if start > 0 && isHexDigit(html[start-1]) {
		return false
	}
	if end < len(html) && isHexDigit(html[end]) {
		return false
	}
	return true
}

func findMagnets(html string) []string {
	uniq := make(map[string]struct{})
	magnets := magnetRegexp.FindAllStringIndex(html, -1)
	out := make([]string, 0, len(magnets))
	for _, idx := range magnets {
		if !isMagnetBoundaryValid(idx[0], idx[1], html) {
			continue
		}
		hash := strings.ToLower(html[idx[0]:idx[1]])
		if _, ok := uniq[hash]; ok {
			continue
		}
		uniq[hash] = struct{}{}
		out = append(out, hash)
	}
	return out
}

func isVisibleElementName(name string) bool {
	switch name {
	case "script", "style", "head", "iframe", "noscript":
		return false
	}
	return true
}

func getVisibleText(s *goquery.Selection) string {
	var buf bytes.Buffer
	// 从 goquery 复制，稍微优化：不选择通常情况下不可见的标签
	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			// Keep newlines and spaces, like jQuery
			buf.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && !isVisibleElementName(c.Data) {
					continue // 跳过
				}
				f(c)
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
	}
	return buf.String()
}

func getMagnetsFromHtml(html []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}
	text := getVisibleText(doc.Find("body"))
	return findMagnets(text), nil
}

func GetMagnetsFromLink(link string) ([]string, error) {
	html2, err := getHtml(link)
	if err != nil {
		return nil, err
	}
	return getMagnetsFromHtml(html2)
}
