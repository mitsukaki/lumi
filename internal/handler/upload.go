package handler

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadHandle(c *gin.Context) {
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	files := form.File["files"]
	for _, file := range files {
		filename := filepath.Base(file.Filename)

		// save to "CDN" directory
		outPath := path.Join("CDN/", filename)
		if err := c.SaveUploadedFile(file, outPath); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
	}

	c.String(http.StatusOK, "Uploaded successfully %d files", len(files))
}
