package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChunk(t *testing.T) {
	as := assert.New(t)
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	chunkSize := 3
	chunks := Chunk(slice, chunkSize)
	as.Equal([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, chunks)
	chunkSize = 4
	chunks = Chunk(slice, chunkSize)
	as.Equal([][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9}}, chunks)
}

func TestShouldPanicWhenChunkSizeIsZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got nil")
		}
	}()
	Chunk([]int{1, 2, 3}, 0)
}

func TestShouldEmptyWhenSliceIsEmpty(t *testing.T) {
	as := assert.New(t)
	chunks := Chunk([]int{}, 1)
	as.Empty(chunks)
}
