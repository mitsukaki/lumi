package endpoints

import (
	"net/http"
)

func (ep *EndpointHandler) DeleteAlbumPhoto(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
