package api

import (
	"net/http"

	"github.com/mitsukaki/lumi/models"
)

func (apiServer *APIServer) GetAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	album, ok := ctx.Value("album").(*models.Album)
	if !ok {
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	JsonWriteOk(w, r, album)
}
