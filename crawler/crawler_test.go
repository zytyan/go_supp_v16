package crawler

import (
	"github.com/stretchr/testify/assert"
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
