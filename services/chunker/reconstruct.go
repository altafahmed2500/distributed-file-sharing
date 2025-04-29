package chunker

// ... other imports
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReconstructFile(meta FileMeta, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Use MergeSort instead of sort.Slice
	sortedChunks := MergeSortChunks(meta.Chunks)

	for _, chunk := range sortedChunks {
		chunkPath := "output/" + meta.FileName + "/chunk_" + strconv.Itoa(chunk.Index)
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			return err
		}
		_, err = outFile.Write(data)
		if err != nil {
			return err
		}
	}

	return verifyRootHash(meta, outputPath)
}

func verifyRootHash(meta FileMeta, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024*256)
	hashes := []byte{}

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		hash := sha256.Sum256(buffer[:n])
		hashes = append(hashes, hash[:]...)
	}

	rootHash := sha256.Sum256(hashes)
	rootHashStr := hex.EncodeToString(rootHash[:])

	if rootHashStr != meta.RootHash {
		return fmt.Errorf("file verification failed: expected %s but got %s", meta.RootHash, rootHashStr)
	}

	return nil
}
