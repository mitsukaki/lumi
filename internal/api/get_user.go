package api

import (
	"net/http"

	"github.com/mitsukaki/lumi/models"
)

func (apiServer *APIServer) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value("user").(*models.DBUser)
	if !ok {
		JsonWriteError(w, r, StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	JsonWriteOk(w, r, user.User)
}
