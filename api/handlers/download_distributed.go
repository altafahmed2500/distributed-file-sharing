package handlers

import (
	"distributed-file-sharing/services/downloader"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func DownloadDistributedHandler(c *gin.Context) {
	fileName := c.Param("fileName")

	chunkMapPath := "chunk_maps/" + fileName + "_chunkmap.json"
	outputDir := "downloaded_chunks/" + fileName
	finalFilePath := "reconstructed/" + fileName

	// Ensure output folders exist
	os.MkdirAll(outputDir, os.ModePerm)
	os.MkdirAll("reconstructed", os.ModePerm)

	// Call downloader module
	err := downloader.DownloadAndReconstruct(chunkMapPath, outputDir, finalFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Download failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "✅ File downloaded and reconstructed successfully",
		"reconstructed": finalFilePath,
	})
}
