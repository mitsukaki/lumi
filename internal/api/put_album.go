package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mitsukaki/lumi/internal/db"
	"github.com/mitsukaki/lumi/models"
	"go.uber.org/zap"
)

func (apiServer *APIServer) PutAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(*models.DBUser)

	// get album request from JSON body
	var albumReq models.AlbumPutRequest
	err := json.NewDecoder(r.Body).Decode(&albumReq)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		apiServer.logger.Info("failed to decode album", zap.Error(err))
		return
	}

	// create the album
	album := models.Album{
		AuthorUserID: user.ID,
		AlbumID:      uuid.New().String(),
		CoverPhoto:   albumReq.CoverPhoto,
		Date:         albumReq.Date,
		Description:  albumReq.Description,
		Title:        albumReq.Title,
	}

	// update the user data
	if albumReq.Private {
		user.PrivateAlbums = append(user.PrivateAlbums, album.AlbumID)
	} else {
		user.PublicData.Albums = append(user.PublicData.Albums, album.AlbumID)
	}

	// update the user in the database
	if updateTable(user, apiServer.UserTable, w, r, apiServer.logger) == false {
		return
	}

	// put the album in the database
	if updateTable(album, apiServer.AlbumTable, w, r, apiServer.logger) == false {
		return
	}

	// return success
	JsonWriteOk(w, r, StatusResponse{
		Ok: true,
	})
}

/**
 * Utility function to update a table in the database
 *
 * @param data interface{} The data to put in the table
 * @param table *db.CouchTable The table to put the data in
 * @param w http.ResponseWriter The http response writer
 * @param r *http.Request The http request
 * @param logger *zap.Logger The console logger
 * @return bool Whether the operation was successful
 */
func updateTable(
	data interface{},
	table *db.CouchTable,
	w http.ResponseWriter,
	r *http.Request,
	logger *zap.Logger,
) bool {
	resp, err := table.Put(data)
	if err != nil {
		JsonWrite(http.StatusInternalServerError, w, r, StatusResponse{
			Ok:     false,
			Reason: "Failed to make album",
		})
		logger.Info("failed to put album", zap.Error(err))
		return false
	}

	if resp.Error != "" {
		JsonWrite(http.StatusInternalServerError, w, r, StatusResponse{
			Ok:     false,
			Reason: "Failed to make album",
		})
		logger.Info("failed to put album", zap.String("error", resp.Error), zap.String("reason", resp.Reason))
		return false
	}

	return true
}
