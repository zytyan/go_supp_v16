package helper

import (
	"os"
	"path/filepath"
	"regexp"
)

var MaxSize int64 = 2000 * 1024 * 1024
var ArchiveSuffix = regexp.MustCompile(`\.(zip|rar|7z|zip\.\d{3}|7z\.\d{3})$`)

func fileNeedRar(file string) bool {
	if stat, err := os.Stat(file); err == nil {
		if stat.Size() > MaxSize {
			return true
		}
	}
	return false
}
func fileInDirNeedRar(file string) bool {
	if stat, err := os.Stat(file); err == nil {
		if stat.IsDir() {
			return true
		}
		if fileNeedRar(file) {
			return true
		}
		if !ArchiveSuffix.MatchString(file) {
			return true
		}
	}
	return false
}

func dirNeedRar(dir string) bool {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	if len(dirEntries) > 5 {
		return true
	}
	for _, dirEntry := range dirEntries {
		path := filepath.Join(dir, dirEntry.Name())
		if fileInDirNeedRar(path) {
			return true
		}
	}
	return false
}

func NeedRar(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if stat.IsDir() {
		return dirNeedRar(path)
	}
	return fileNeedRar(path)
}
