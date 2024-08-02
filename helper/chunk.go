package helper

func Chunk[T any](slice []T, chunkSize int) [][]T {
	chunks := make([][]T, 0, (len(slice)+chunkSize-1)/chunkSize)
	i := 0
	for ; i < len(slice)-chunkSize; i += chunkSize {
		chunks = append(chunks, slice[i:i+chunkSize])
	}
	chunks = append(chunks, slice[i:])
	return chunks
}
