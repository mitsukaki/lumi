package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mitsukaki/lumi/models"
	"go.uber.org/zap"
)

func (apiServer *APIServer) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		photoId := chi.URLParam(r, "photoId")

		doc, err := apiServer.UserTable.Get(photoId)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			apiServer.logger.Info("failed to get user", zap.Error(err))
			return
		}

		var dbUser models.DBUser
		err = json.Unmarshal(doc, &dbUser)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			apiServer.logger.Info("failed to unmarshal db user", zap.Error(err))
			return
		}

		ctx := context.WithValue(r.Context(), "user", dbUser.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
