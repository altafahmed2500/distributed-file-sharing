package chunker

func MergeSortChunks(chunks []ChunkMeta) []ChunkMeta {
	if len(chunks) <= 1 {
		return chunks
	}

	mid := len(chunks) / 2
	left := MergeSortChunks(chunks[:mid])
	right := MergeSortChunks(chunks[mid:])

	return merge(left, right)
}

func merge(left, right []ChunkMeta) []ChunkMeta {
	result := []ChunkMeta{}
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i].Index < right[j].Index {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// Append leftovers
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}
