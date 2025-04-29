package downloader

import (
	"distributed-file-sharing/services/chunker"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	_ "strconv"
)

type ChunkMap struct {
	FileName string         `json:"FileName"`
	RootHash string         `json:"RootHash"`
	Chunks   map[int]string `json:"Chunks"` // map[index]peerAddress
}

// DownloadAndReconstruct downloads chunks from peers and reconstructs the file
func DownloadAndReconstruct(chunkMapPath, outputDir, finalFilePath string) error {
	// Step 1: Read chunk map
	data, err := os.ReadFile(chunkMapPath)
	if err != nil {
		return fmt.Errorf("failed to read chunk map: %v", err)
	}

	var cmap ChunkMap
	err = json.Unmarshal(data, &cmap)
	if err != nil {
		return fmt.Errorf("failed to parse chunk map: %v", err)
	}

	// Step 2: Create temp folder to store chunks
	os.MkdirAll(outputDir, os.ModePerm)

	// Step 3: Download all chunks
	for idx, peer := range cmap.Chunks {
		err := downloadChunk(peer, cmap.FileName, idx, outputDir)
		if err != nil {
			return fmt.Errorf("failed to download chunk %d: %v", idx, err)
		}
	}

	// Step 4: Reconstruct file from chunks
	return reconstructFileFromChunks(cmap, outputDir, finalFilePath)
}

func downloadChunk(peerAddress, fileName string, index int, outputDir string) error {
	url := fmt.Sprintf("%s/get-chunk?fileName=%s&index=%d", peerAddress, fileName, index)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch chunk: status %d", resp.StatusCode)
	}

	savePath := fmt.Sprintf("%s/chunk_%d", outputDir, index)
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func reconstructFileFromChunks(cmap ChunkMap, tempDir, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	chunks := make([]chunker.ChunkMeta, 0, len(cmap.Chunks))
	for idx := range cmap.Chunks {
		chunks = append(chunks, chunker.ChunkMeta{
			Index: idx,
		})
	}

	// Sort chunks using MergeSortChunks
	sortedChunks := chunker.MergeSortChunks(chunks)

	for _, chunk := range sortedChunks {
		chunkPath := fmt.Sprintf("%s/chunk_%d", tempDir, chunk.Index)
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			return err
		}
		_, err = outFile.Write(data)
		if err != nil {
			return err
		}
	}

	fmt.Println("✅ Reconstructed file at", outputPath)
	return nil
}
