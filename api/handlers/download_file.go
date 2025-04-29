package handlers

import (
	"distributed-file-sharing/services/chunker"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func DownloadFileHandler(c *gin.Context) {
	fileName := c.Param("fileName")
	metaPath := "output/" + fileName + "/meta.json"

	metaFile, err := os.ReadFile(metaPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metadata not found"})
		return
	}

	var meta chunker.FileMeta
	err = json.Unmarshal(metaFile, &meta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse metadata"})
		return
	}

	reconstructedPath := "output/" + fileName + "/reconstructed_" + fileName
	err = chunker.ReconstructFile(meta, reconstructedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reconstruct file"})
		return
	}

	c.File(reconstructedPath)
}
