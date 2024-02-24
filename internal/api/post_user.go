package api

import (
	"net/http"
)

func (apiServer *APIServer) PostUser(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// user, ok := ctx.Value("user").(*models.UserData)

	JsonWriteError(w, r, StatusResponse{
		Ok:     false,
		Reason: "unimplemented",
	})
}
