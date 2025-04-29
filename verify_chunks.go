package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ChunkMap struct {
	FileName string            `json:"FileName"`
	RootHash string            `json:"RootHash"`
	Chunks   map[string]string `json:"Chunks"`
}

type MetaChunk struct {
	Index int    `json:"index"`
	Hash  string `json:"hash"`
	Size  int    `json:"size"`
}

type FileMeta struct {
	FileName string      `json:"fileName"`
	RootHash string      `json:"rootHash"`
	Chunks   []MetaChunk `json:"chunks"`
}

func main() {
	chunkMapPath := "chunk_maps/MyResume-1.0.0.zip_chunkmap.json" // <-- your chunk map
	metaPath := "output/MyResume-1.0.0.zip/meta.json"             // <-- your meta.json

	// Read chunk map
	chunkMapData, err := os.ReadFile(chunkMapPath)
	if err != nil {
		fmt.Println("❌ Failed to read chunk map:", err)
		return
	}

	var cmap ChunkMap
	err = json.Unmarshal(chunkMapData, &cmap)
	if err != nil {
		fmt.Println("❌ Failed to parse chunk map:", err)
		return
	}

	// Read meta.json
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		fmt.Println("❌ Failed to read meta.json:", err)
		return
	}

	var meta FileMeta
	err = json.Unmarshal(metaData, &meta)
	if err != nil {
		fmt.Println("❌ Failed to parse meta.json:", err)
		return
	}

	// Build lookup map for chunk hash
	metaHashes := make(map[int]string)
	for _, chunk := range meta.Chunks {
		metaHashes[chunk.Index] = chunk.Hash
	}

	fmt.Println("🔎 Verifying chunks for:", cmap.FileName)

	totalChunks := 0
	okChunks := 0

	// Verify each chunk
	for indexStr, peer := range cmap.Chunks {
		totalChunks++

		url := fmt.Sprintf("%s/get-chunk?fileName=%s&chunkIndex=%s", peer, cmap.FileName, indexStr)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("❌ Chunk %s: Failed to reach %s (%v)\n", indexStr, peer, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("❌ Chunk %s: Not found at %s (Status %d)\n", indexStr, peer, resp.StatusCode)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("❌ Chunk %s: Error reading response: %v\n", indexStr, err)
			continue
		}

		if len(data) == 0 {
			fmt.Printf("❌ Chunk %s: Empty chunk received from %s\n", indexStr, peer)
			continue
		}

		// Hash the downloaded chunk
		hash := sha256.Sum256(data)
		hashStr := hex.EncodeToString(hash[:])

		// Check against meta.json
		idx := atoi(indexStr)
		expectedHash, ok := metaHashes[idx]
		if !ok {
			fmt.Printf("❌ Chunk %s: Not found in meta.json\n", indexStr)
			continue
		}

		if hashStr != expectedHash {
			fmt.Printf("❌ Chunk %s: Hash mismatch! Expected %s, Got %s\n", indexStr, expectedHash, hashStr)
			continue
		}

		fmt.Printf("✅ Chunk %s: Verified (%d bytes)\n", indexStr, len(data))
		okChunks++
	}

	fmt.Println()
	fmt.Printf("📊 Chunk Verification Result: %d/%d Chunks OK\n", okChunks, totalChunks)
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
