package strnum

import (
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"testing"
)

func TestNumCompare(t *testing.T) {
	as := assert.New(t)
	as.Equal(compareVec(Split("01"), Split("1")), 0)
	as.Equal(compareVec(Split("01"), Split("0")), 1)
	as.Equal(compareVec(Split("0"), Split("01")), -1)
	as.Equal(compareVec(Split("1"), Split("2")), -1)
	as.Equal(compareVec(Split("1"), Split("002")), -1)
	as.Equal(compareVec(Split("a"), Split("002")), 1)
	as.Equal(compareVec(Split("1"), Split("a")), -1)
}

func TestFullOrder(t *testing.T) {
	as := assert.New(t)
	as.True(compareVec(Split("a"), Split("b")) < 0)
	as.True(compareVec(Split("9"), Split("20")) < 0)
	as.True(compareVec(Split("a"), Split("b")) < 0)
	as.True(compareVec(Split("999999"), Split("A")) < 0)
	as.True(compareVec(Split("A"), Split("B")) < 0)
	as.True(compareVec(Split("a9."), Split("a11")) < 0)
	as.True(compareVec(Split("a11"), Split("a12")) < 0)
	as.True(compareVec(Split("a11"), Split("a111")) < 0)
	as.True(compareVec(Split("a11"), Split("a111a")) < 0)
	as.True(compareVec(Split("a1"), Split("a1a")) < 0)
	as.True(compareVec(Split("a1"), Split("a2")) < 0)
	as.Equal(compareVec(Split("a1"), Split("a0001")), 0)
	as.Equal(compareVec(Split("01"), Split("0001")), 0)
	as.Equal(compareVec(Split("1"), Split("0001")), 0)
}
func TestSort(t *testing.T) {
	as := assert.New(t)
	strings := []string{
		"a01-1",
		"a01-2",
		"a01-3",
		"a01-10",
		"a01-14",
		"a01-19",
		"a01-20",
		"a01-21",
		"a01-42",
		"a01-100",
		"a02-1",
		"a02-2",
		"a02-3",
		"a02-101",
		"a02-102",
	}
	stringsCopy := make([]string, len(strings))
	copy(stringsCopy, strings)
	rand.Shuffle(len(stringsCopy), func(i, j int) {
		stringsCopy[i], stringsCopy[j] = stringsCopy[j], stringsCopy[i]
	})
	sortedStr := SortedStrings(stringsCopy)
	as.Equal(strings, sortedStr)
	sortedStr = SortedStrings(strings)
	as.Equal(strings, sortedStr)

	buf := make([]Buf, 0, len(strings))
	for _, s := range strings {
		buf = append(buf, Split(s))
	}
	sameBuf := make([]Buf, len(buf))
	copy(sameBuf, buf)
	Sort(buf)
	as.Equal(sameBuf, buf)
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	Sort(buf)
	as.Equal(sameBuf, buf)
}
