package handlers

import (
	"distributed-file-sharing/services/merkle"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GenerateMerkleTreeHandler(c *gin.Context) {
	fileName := c.Param("fileName")
	metaPath := "output/" + fileName + "/meta.json"

	hashes, err := merkle.LoadChunkHashes(metaPath)
	if err != nil {
		c.String(http.StatusNotFound, "Metadata file not found")
		return
	}

	root := merkle.BuildMerkleTree(hashes)

	// 🚀 Return JSON, not HTML
	c.JSON(http.StatusOK, root)
}
