package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mitsukaki/lumi/models"
)

func (apiServer *APIServer) AlbumContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		albumId := chi.URLParam(r, "albumId")

		doc, err := apiServer.AlbumTable.Get(albumId)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			apiServer.logger.Info("failed to get album", zap.Error(err))
			return
		}

		var album models.Album
		err = json.Unmarshal(doc, &album)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			apiServer.logger.Info("failed to unmarshal album", zap.Error(err))
			return
		}

		ctx := context.WithValue(r.Context(), "album", &album)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
