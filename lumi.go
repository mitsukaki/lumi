package main

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/mitsukaki/lumi/internal/handler"
)

// TODO: set GIN_MODE=release
func main() {
	exePath, _ := os.Executable()
	binaryPath := filepath.Dir(exePath)

	router := gin.Default()

	// api routes
	api := router.Group("/api")
	{
		api.POST("/upload", handler.UploadHandle)
	}

	// serve static content
	router.Static("/", filepath.Join(binaryPath, "public"))

	// direct all 404s to the index.html (for client-side routing)
	router.NoRoute(func(c *gin.Context) {
		filePath := filepath.Join(binaryPath, "public", "index.html")
		c.File(filePath)
	})

	// listen and serve
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	router.Run(":" + port) // listen and serve on specified port
}
