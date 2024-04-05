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
	albumId := uuid.New().String()
	album := models.Album{
		AuthorUserID: user.ID,
		AlbumID:      albumId,
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
	if updateTable(user.ID, user.Rev, user, apiServer.UserTable, w, r, apiServer.logger) == false {
		return
	}

	// put the album in the database
	if insertIntoTable(albumId, album, apiServer.AlbumTable, w, r, apiServer.logger) == false {
		return
	}

	// return success & the album ID
	JsonWrite(http.StatusOK, w, r, AlbumPutResponse{
		Ok: true,
		ID: albumId,
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
	id string,
	rev string,
	data interface{},
	table *db.CouchTable,
	w http.ResponseWriter,
	r *http.Request,
	logger *zap.Logger,
) bool {
	resp, err := table.PutExisting(id, rev, data)
	if err != nil {
		JsonWrite(http.StatusInternalServerError, w, r, StatusResponse{
			Ok:     false,
			Reason: "Failed to update to table",
		})
		logger.Info("failed to update to table", zap.Error(err))
		return false
	}

	if resp.Error != "" {
		JsonWrite(http.StatusInternalServerError, w, r, StatusResponse{
			Ok:     false,
			Reason: "Failed to update to table",
		})
		logger.Info("failed to update to table", zap.String("error", resp.Error), zap.String("reason", resp.Reason))
		return false
	}

	return true
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
func insertIntoTable(
	id string,
	data interface{},
	table *db.CouchTable,
	w http.ResponseWriter,
	r *http.Request,
	logger *zap.Logger,
) bool {
	resp, err := table.PutNew(id, data)
	if err != nil {
		JsonWrite(http.StatusInternalServerError, w, r, StatusResponse{
			Ok:     false,
			Reason: "Failed to insert to table",
		})
		logger.Info("failed to insert to table", zap.Error(err))
		return false
	}

	if resp.Error != "" {
		JsonWrite(http.StatusInternalServerError, w, r, StatusResponse{
			Ok:     false,
			Reason: "Failed to insert to table",
		})
		logger.Info("failed to insert to table", zap.String("error", resp.Error), zap.String("reason", resp.Reason))
		return false
	}

	return true
}
