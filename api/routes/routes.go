package routes

import (
	"distributed-file-sharing/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/chunk", handlers.ChunkFileHandler)
	r.GET("/download/:fileName", handlers.DownloadFileHandler)
	r.GET("/merkle/:fileName", handlers.GenerateMerkleTreeHandler)
	r.POST("/upload-chunk", handlers.UploadChunkHandler)
}
