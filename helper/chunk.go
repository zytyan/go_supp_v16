package helper

import "strconv"

func Chunk[T any](slice []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		panic("chunk size must be greater than 0, now: " + strconv.Itoa(chunkSize))
	}
	if len(slice) == 0 {
		return nil
	}
	chunks := make([][]T, 0, (len(slice)+chunkSize-1)/chunkSize)
	i := 0
	for ; i < len(slice)-chunkSize; i += chunkSize {
		chunks = append(chunks, slice[i:i+chunkSize])
	}
	chunks = append(chunks, slice[i:])
	return chunks
}
