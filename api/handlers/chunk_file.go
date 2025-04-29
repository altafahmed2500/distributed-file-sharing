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
	"http://INSTANCE1_IP:8080",
	"http://INSTANCE2_IP:8080",
	// Add more if needed
}

func ChunkFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	chunkSizeStr := c.DefaultPostForm("chunkSize", "256")
	chunkSize, _ := strconv.Atoi(chunkSizeStr)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
		return
	}
	defer src.Close()

	meta, err := chunker.ChunkFile(src, file.Filename, chunkSize*1024) // multiply KB to bytes
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save meta.json
	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal metadata"})
		return
	}

	metaPath := "output/" + file.Filename + "/meta.json"
	err = os.WriteFile(metaPath, metaBytes, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metadata"})
		return
	}

	// ✅ ✅ ✅ Now, distribute chunks to different instances
	chunkMap := make(map[int]string)

	for idx, chunk := range meta.Chunks {
		assignedPeer := peers[idx%len(peers)] // Round Robin

		chunkPath := "output/" + file.Filename + "/chunk_" + strconv.Itoa(chunk.Index)
		chunkName := fmt.Sprintf("%s_chunk_%d", file.Filename, chunk.Index)

		err := uploadChunk(assignedPeer, chunkPath, chunkName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload chunk %d: %v", chunk.Index, err)})
			return
		}

		chunkMap[chunk.Index] = assignedPeer
	}

	// Save ChunkMap
	err = saveChunkMap(file.Filename, chunkMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk map"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File chunked and distributed successfully",
		"rootHash": meta.RootHash,
	})
}

// Upload a chunk to assigned peer
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

// Save chunk map locally
func saveChunkMap(fileName string, chunkMap map[int]string) error {
	os.MkdirAll("chunk_maps", os.ModePerm)
	mapFile, err := os.Create("chunk_maps/" + fileName + "_chunkmap.json")
	if err != nil {
		return err
	}
	defer mapFile.Close()

	data, _ := json.MarshalIndent(chunkMap, "", "  ")
	mapFile.Write(data)
	return nil
}
