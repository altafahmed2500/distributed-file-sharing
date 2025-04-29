package chunker

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strconv"
)

type ChunkMeta struct {
	Index int    `json:"index"`
	Hash  string `json:"hash"`
	Size  int    `json:"size"`
}

type FileMeta struct {
	FileName string      `json:"fileName"`
	RootHash string      `json:"rootHash"`
	Chunks   []ChunkMeta `json:"chunks"`
}

func ChunkFile(src io.Reader, fileName string, chunkSize int) (FileMeta, error) {
	chunks := []ChunkMeta{}
	buffer := make([]byte, chunkSize)
	index := 0

	err := os.MkdirAll("output/"+fileName, os.ModePerm)
	if err != nil {
		return FileMeta{}, err
	}

	for {
		n, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			return FileMeta{}, err
		}
		if n == 0 {
			break
		}

		hash := sha256.Sum256(buffer[:n])
		hashStr := hex.EncodeToString(hash[:])

		chunkPath := "output/" + fileName + "/chunk_" + strconv.Itoa(index)
		err = os.WriteFile(chunkPath, buffer[:n], 0644)
		if err != nil {
			return FileMeta{}, err
		}

		chunks = append(chunks, ChunkMeta{
			Index: index,
			Hash:  hashStr,
			Size:  n,
		})
		index++
	}

	var fileHashInput bytes.Buffer
	for _, chunk := range chunks {
		decodedHash, _ := hex.DecodeString(chunk.Hash)
		fileHashInput.Write(decodedHash)
	}
	rootHash := sha256.Sum256(fileHashInput.Bytes())
	rootHashStr := hex.EncodeToString(rootHash[:])

	return FileMeta{
		FileName: fileName,
		RootHash: rootHashStr,
		Chunks:   chunks,
	}, nil
}
