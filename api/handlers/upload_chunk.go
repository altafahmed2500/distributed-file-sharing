package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// UploadChunkHandler handles chunk upload at each peer
func UploadChunkHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	// Create folder to save chunks if not exists
	err = os.MkdirAll("chunks", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chunks folder"})
		return
	}

	// Save file to 'chunks' folder
	savePath := "chunks/" + file.Filename
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chunk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Chunk uploaded and stored successfully",
	})
}
