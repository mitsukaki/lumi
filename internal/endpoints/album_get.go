package endpoints

import (
	"net/http"

	"github.com/mitsukaki/lumi/internal/net"
	"github.com/mitsukaki/lumi/models"
)

func (ep *EndpointHandler) GetAlbum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	album, ok := ctx.Value("album").(*models.DBAlbum)
	if !ok {
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	net.JsonWriteOk(w, r, &models.Album{
		AlbumID:      album.AlbumID,
		AuthorUserID: album.AuthorUserID,
		CoverPhoto:   album.CoverPhoto,
		Date:         album.Date,
		Description:  album.Description,
		Title:        album.Title,
		Photos:       album.Photos,
	})
}
