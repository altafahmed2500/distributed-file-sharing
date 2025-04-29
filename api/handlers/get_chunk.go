package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetChunkHandler(c *gin.Context) {
	fileName := c.Query("fileName")
	index := c.Query("index")

	chunkPath := "chunks/" + fileName + "_chunk_" + index

	data, err := os.ReadFile(chunkPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chunk not found"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", data)
}
