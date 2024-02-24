package api

import (
	"net/http"
)

func (apiServer *APIServer) Ping(w http.ResponseWriter, r *http.Request) {
	JsonWriteOk(w, r, StatusResponse{
		Ok:     true,
		Reason: "pong!",
	})
}
