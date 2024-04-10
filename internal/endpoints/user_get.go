package endpoints

import (
	"net/http"

	"github.com/mitsukaki/lumi/internal/net"
	"github.com/mitsukaki/lumi/models"
)

func (ep *EndpointHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value("user").(*models.DBUser)
	if !ok {
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "internal error",
		})

		return
	}

	net.JsonWriteOk(w, r, user.PublicData)
}
