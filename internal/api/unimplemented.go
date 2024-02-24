package api

import (
	"net/http"
)

func (apiServer *APIServer) Unimplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
