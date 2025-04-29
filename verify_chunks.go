package main

import (
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

func main() {
	chunkMapPath := "chunk_maps/MyResume-1.0.0.zip_chunkmap.json" // <-- Change filename here if needed

	// Read chunk map
	data, err := os.ReadFile(chunkMapPath)
	if err != nil {
		fmt.Println("❌ Failed to read chunk map:", err)
		return
	}

	var cmap ChunkMap
	err = json.Unmarshal(data, &cmap)
	if err != nil {
		fmt.Println("❌ Failed to parse chunk map:", err)
		return
	}

	fmt.Println("🔎 Verifying chunks for:", cmap.FileName)

	totalChunks := 0
	okChunks := 0

	// Loop through all chunks
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

		n, _ := io.Copy(io.Discard, resp.Body)

		if n == 0 {
			fmt.Printf("❌ Chunk %s: Empty chunk received from %s\n", indexStr, peer)
			continue
		}

		fmt.Printf("✅ Chunk %s: Found (%d bytes) from %s\n", indexStr, n, peer)
		okChunks++
	}

	fmt.Println()
	fmt.Printf("📊 Chunk Verification Result: %d/%d Chunks OK\n", okChunks, totalChunks)
}
