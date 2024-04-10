package endpoints

import (
	"net/http"

	"github.com/mitsukaki/lumi/internal/net"
	"github.com/mitsukaki/lumi/models"
)

func (ep *EndpointHandler) Ping(w http.ResponseWriter, r *http.Request) {
	net.JsonWriteOk(w, r, models.StatusResponse{
		Ok:     true,
		Reason: "pong!",
	})
}
