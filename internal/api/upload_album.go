package api

import (
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/mitsukaki/lumi/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func (apiServer *APIServer) UploadAlbum(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	album, ok := ctx.Value("album").(*models.DBAlbum)
	if !ok {
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	// limit the size of the request to be no more than 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	// parse the multipart form
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		apiServer.logger.Error("failed to parse multipart form", zap.Error(err))
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	// get the metadata from the form
	title := r.FormValue("title")
	description := r.FormValue("description")
	row, err := strconv.Atoi(r.FormValue("row"))
	if err != nil {
		apiServer.logger.Error("failed to parse row", zap.Error(err))
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	aspectRatio, err := strconv.ParseFloat(r.FormValue("aspect_ratio"), 64)
	if err != nil {
		apiServer.logger.Error("failed to parse aspect ratio", zap.Error(err))
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	// get the image from the request
	file, _, err := r.FormFile("image")
	if err != nil {
		apiServer.logger.Error("failed to get image from form", zap.Error(err))
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	defer file.Close()

	// generate a UUID for the file name
	photoID := uuid.New().String()

	// upload the file to S3
	_, err = apiServer.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("s3.bucket")),
		Key:    aws.String("albums/" + album.AlbumID + "/" + photoID),
		Body:   file,
	})
	if err != nil {
		apiServer.logger.Error("failed to upload file to S3", zap.Error(err))
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	// add the photo to the album
	album.Photos = append(album.Photos, models.Photo{
		Title:       title,
		Description: description,
		PhotoID:     photoID,
		Ratio:       aspectRatio,
		Row:         row,
	})

	// upload the updated album
	resp, err := apiServer.AlbumTable.PutExisting(album.AlbumID, album.Rev, album)
	if err != nil {
		apiServer.logger.Error("failed to upload album", zap.Error(err))
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	if resp.Error != "" {
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: resp.Reason,
		})

		return
	}

	// write the response
	JsonWriteOk(w, r, StatusResponse{
		Ok: true,
	})
}
