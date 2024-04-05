package api

import (
	"net/http"

	"github.com/mitsukaki/lumi/models"
)

func (apiServer *APIServer) GetAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	album, ok := ctx.Value("album").(*models.DBAlbum)
	if !ok {
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	JsonWriteOk(w, r, &models.Album{
		AlbumID:      album.AlbumID,
		AuthorUserID: album.AuthorUserID,
		CoverPhoto:   album.CoverPhoto,
		Date:         album.Date,
		Description:  album.Description,
		Title:        album.Title,
		Photos:       album.Photos,
	})
}
