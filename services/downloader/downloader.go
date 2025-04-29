package downloader

import (
	"distributed-file-sharing/services/chunker"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	err = reconstructFileFromChunks(cmap, outputDir, finalFilePath)
	if err != nil {
		return err
	}

	// Step 5: Verify root hash
	return verifyRootHash(cmap.RootHash, finalFilePath)
}

// downloadChunk fetches a chunk from the appropriate peer
func downloadChunk(peerAddress, fileName string, index int, outputDir string) error {
	savePath := fmt.Sprintf("%s/chunk_%d", outputDir, index)

	// Skip download if chunk already exists and is not empty
	if stat, err := os.Stat(savePath); err == nil {
		if stat.Size() > 0 {
			fmt.Printf("🟡 Chunk %d already exists locally (%d bytes). Skipping download.\n", index, stat.Size())
			return nil
		}
		fmt.Printf("⚠️ Chunk %d exists but is empty. Re-downloading.\n", index)
	}

	url := fmt.Sprintf("%s/get-chunk?fileName=%s&chunkIndex=%d", peerAddress, fileName, index)
	fmt.Println("🌍 Fetching:", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET failed for %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch chunk: status %d", resp.StatusCode)
	}

	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	if err == nil {
		fmt.Printf("✅ Downloaded chunk %d (%d bytes)\n", index, n)
	}
	return err
}

// reconstructFileFromChunks merges downloaded chunks into a full file
func reconstructFileFromChunks(cmap ChunkMap, tempDir, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Build ChunkMeta list
	chunks := make([]chunker.ChunkMeta, 0, len(cmap.Chunks))
	for idx := range cmap.Chunks {
		chunks = append(chunks, chunker.ChunkMeta{Index: idx})
	}

	// Sort by index
	sortedChunks := chunker.MergeSortChunks(chunks)

	for _, chunk := range sortedChunks {
		chunkPath := fmt.Sprintf("%s/chunk_%d", tempDir, chunk.Index)
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			return fmt.Errorf("error reading chunk %d: %v", chunk.Index, err)
		}
		_, err = outFile.Write(data)
		if err != nil {
			return fmt.Errorf("error writing chunk %d: %v", chunk.Index, err)
		}
	}

	fmt.Println("📦 Reconstructed file at:", outputPath)
	return nil
}

// verifyRootHash checks if reconstructed file's root hash matches
func verifyRootHash(expectedHash, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024*256)
	all := []byte{}

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		chunkHash := chunker.HashChunk(buffer[:n])
		all = append(all, chunkHash...)
	}

	root := chunker.HashChunk(all)
	if root != expectedHash {
		return fmt.Errorf("❌ Root hash mismatch: expected %s, got %s", expectedHash, root)
	}

	fmt.Println("✅ Root hash matched:", root)
	return nil
}
