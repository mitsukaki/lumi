package endpoints

import (
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/mitsukaki/lumi/internal/net"
	"github.com/mitsukaki/lumi/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func (ep *EndpointHandler) UploadAlbumPhoto(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	album, ok := ctx.Value("album").(*models.DBAlbum)
	if !ok {
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	// limit the size of the request to be no more than 35MB
	sizeLimit := int64(35 * 1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, sizeLimit)

	// parse the multipart form
	err := r.ParseMultipartForm(sizeLimit)
	if err != nil {
		ep.logger.Error("failed to parse multipart form", zap.Error(err))
		net.JsonWriteError(w, r, models.StatusResponse{
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
		ep.logger.Error("failed to parse row", zap.Error(err))
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	aspectRatio, err := strconv.ParseFloat(r.FormValue("aspect_ratio"), 64)
	if err != nil {
		ep.logger.Error("failed to parse aspect ratio", zap.Error(err))
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	// get the image from the request
	file, _, err := r.FormFile("image")
	if err != nil {
		ep.logger.Error("failed to get image from form", zap.Error(err))
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	defer file.Close()

	// generate a UUID for the file name
	photoID := uuid.New().String()

	// upload the file to S3
	_, err = ep.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("s3.bucket")),
		Key:    aws.String("albums/" + album.AlbumID + "/" + photoID),
		Body:   file,
	})
	if err != nil {
		ep.logger.Error("failed to upload file to S3", zap.Error(err))
		net.JsonWriteError(w, r, models.StatusResponse{
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
	resp, err := ep.albumTable.PutExisting(album.AlbumID, album.Rev, album)
	if err != nil {
		ep.logger.Error("failed to upload album", zap.Error(err))
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	if resp.Error != "" {
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: resp.Reason,
		})

		return
	}

	// write the response
	net.JsonWriteOk(w, r, models.StatusResponse{
		Ok: true,
	})
}
