package endpoints

import (
	"net/http"
)

func (ep *EndpointHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
