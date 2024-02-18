package main

import (
	"github.com/gin-gonic/gin"

	"github.com/mitsukaki/lumi/internal/handler"
)

func main() {
	router := gin.Default()
	router.Static("/", "./web/dist")

	router.POST("/upload", handler.UploadHandle)

	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}
