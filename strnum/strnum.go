package strnum

import (
	"slices"
	"strings"
)

type StrNum struct {
	Slice string
	IsNum bool
}
type Buf struct {
	Chunks []StrNum
	Raw    string
}

func compareNum(i, j StrNum) int {
	if len(i.Slice) > len(j.Slice) {
		return 1
	} else if len(i.Slice) < len(j.Slice) {
		return -1
	} else {
		return strings.Compare(i.Slice, j.Slice)
	}
}

func compare(i, j StrNum) int {
	if i.IsNum && j.IsNum {
		return compareNum(i, j)
	} else if !i.IsNum && !j.IsNum {
		return strings.Compare(i.Slice, j.Slice)
	} else if i.IsNum {
		return -1
	} else {
		return 1
	}
}
func compareVec(i, j Buf) int {
	left := i.Chunks
	right := j.Chunks
	for idx := 0; idx < len(left) && idx < len(right); idx++ {
		if ret := compare(left[idx], right[idx]); ret != 0 {
			return ret
		}
	}
	if len(left) > len(right) {
		return 1
	} else if len(left) < len(right) {
		return -1
	} else {
		return 0
	}
}
func isNum(c byte) bool {
	return c >= '0' && c <= '9'
}

func Split(s string) (ret Buf) {
	idx := 0
	chunks := make([]StrNum, 0, 8)
	for idx < len(s) {
		for idx < len(s) && s[idx] == '0' {
			idx++
		}
		start := idx
		for idx < len(s) && isNum(s[idx]) {
			idx++
		}
		if idx > start {
			chunks = append(chunks, StrNum{Slice: s[start:idx], IsNum: true})
		}
		start = idx
		for idx < len(s) && !isNum(s[idx]) {
			idx++
		}
		if idx > start {
			chunks = append(chunks, StrNum{Slice: s[start:idx], IsNum: false})
		}
	}
	ret.Chunks = chunks
	ret.Raw = s
	return
}

func Sort(v []Buf) {
	slices.SortFunc(v, compareVec)
}
