package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mitsukaki/lumi/internal/db"
	"github.com/mitsukaki/lumi/models"
)

type AlbumContextHandler struct {
	logger *zap.Logger

	albumTable *db.CouchTable
}

func CreateAlbumContextHandler(logger *zap.Logger, albumTable *db.CouchTable) *AlbumContextHandler {
	return &AlbumContextHandler{
		logger:     logger,
		albumTable: albumTable,
	}
}

func (albumCtx *AlbumContextHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		albumId := chi.URLParam(r, "albumId")

		doc, err := albumCtx.albumTable.Get(albumId)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			albumCtx.logger.Info("failed to get album", zap.Error(err))
			return
		}

		var album models.DBAlbum
		err = json.Unmarshal(doc, &album)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			albumCtx.logger.Info("failed to unmarshal album", zap.Error(err))
			return
		}

		ctx := context.WithValue(r.Context(), "album", &album)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
