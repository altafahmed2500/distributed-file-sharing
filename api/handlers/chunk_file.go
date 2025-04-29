package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"distributed-file-sharing/services/chunker"

	"github.com/gin-gonic/gin"
)

var peers = []string{
	"http://192.168.88.1:8080",
	"http://34.223.225.170:8080",
	// Add more instances if needed
}

// ChunkFileHandler handles file upload, chunking, and distribution
func ChunkFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	chunkSizeStr := c.DefaultPostForm("chunkSize", "256") // default 256 KB
	chunkSize, _ := strconv.Atoi(chunkSizeStr)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
		return
	}
	defer src.Close()

	// ✅ Step 1: Chunk the file
	meta, err := chunker.ChunkFile(src, file.Filename, chunkSize*1024) // KB -> bytes
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ Step 2: Save meta.json
	err = saveMeta(file.Filename, meta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metadata"})
		return
	}

	// ✅ Step 3: Distribute chunks
	chunkMap := make(map[int]string)

	for idx, chunk := range meta.Chunks {
		assignedPeer := peers[idx%len(peers)] // Round-robin peer assignment

		chunkPath := "output/" + file.Filename + "/chunk_" + strconv.Itoa(chunk.Index)
		chunkName := fmt.Sprintf("%s_chunk_%d", file.Filename, chunk.Index)

		err := uploadChunk(assignedPeer, chunkPath, chunkName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload chunk %d: %v", chunk.Index, err)})
			return
		}

		chunkMap[chunk.Index] = assignedPeer
	}

	// ✅ Step 4: Save chunk map properly
	err = saveChunkMap(file.Filename, chunkMap, meta.RootHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk map"})
		return
	}

	// ✅ Final success response
	c.JSON(http.StatusOK, gin.H{
		"message":  "File chunked and distributed successfully",
		"rootHash": meta.RootHash,
	})
}

// Save meta.json for the file
func saveMeta(fileName string, meta chunker.FileMeta) error {
	os.MkdirAll("output/"+fileName, os.ModePerm)
	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	metaPath := "output/" + fileName + "/meta.json"
	return os.WriteFile(metaPath, metaBytes, 0644)
}

// Upload a single chunk to assigned peer
func uploadChunk(peerAddress, chunkPath, chunkName string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(chunkPath)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", chunkName)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest("POST", peerAddress+"/upload-chunk", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed: %v", resp.Status)
	}

	return nil
}

// Save full chunk map with filename and root hash
func saveChunkMap(fileName string, chunkMap map[int]string, rootHash string) error {
	os.MkdirAll("chunk_maps", os.ModePerm)

	// Full structured chunk map
	fullMap := map[string]interface{}{
		"FileName": fileName,
		"RootHash": rootHash,
		"Chunks":   chunkMap,
	}

	data, err := json.MarshalIndent(fullMap, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("chunk_maps/"+fileName+"_chunkmap.json", data, 0644)
}
