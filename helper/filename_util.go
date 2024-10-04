package helper

import (
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	for !utf8.RuneStart(s[length]) {
		length--
	}
	return s[:length]
}

func Stem(s string) string {
	base := filepath.Base(s)
	ext := filepath.Ext(base)
	return base[:len(base)-len(ext)]
}

func IsVideoFile(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".mp4", ".mkv", ".avi", ".mov":
		return true
	}
	return false
}
