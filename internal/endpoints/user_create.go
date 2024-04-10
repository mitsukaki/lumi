package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mitsukaki/lumi/internal/net"
	"github.com/mitsukaki/lumi/models"
	"go.uber.org/zap"
)

func (ep *EndpointHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the userCreateReq data
	var userCreateReq models.UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&userCreateReq)
	if err != nil {
		ep.logger.Error("Failed to decode user create request", zap.Error(err))
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "Invalid request",
		})

		return
	}

	// Validate the userCreateReq data
	if userCreateReq.Email == "" || userCreateReq.Password == "" || userCreateReq.Username == "" {
		ep.logger.Error("Invalid request made from: " + r.RemoteAddr)
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "Invalid request",
		})

		return
	}

	// Generate a new user ID
	userID := uuid.New().String()

	// Create the user DB object
	userDB := models.DBUser{
		ID:       userID,
		Email:    userCreateReq.Email,
		Password: userCreateReq.Password,
		PublicData: models.UserData{
			Username: userCreateReq.Username,
		},
	}

	// Insert the user into the database
	resp, err := ep.userTable.PutNew(userID, userDB)
	if err != nil {
		ep.logger.Error("Failed to insert user into database", zap.Error(err))
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "Internal error",
		})

		return
	}

	if resp.Ok == false {
		ep.logger.Error("Failed to insert userinto database", zap.String("reason", resp.Reason))
		net.JsonWriteError(w, r, models.StatusResponse{
			Ok:     false,
			Reason: "Internal error",
		})

		return
	}

	// log user creation with the user id
	ep.logger.Info("User created", zap.String("user_id", userID))

	// Return success
	net.JsonWriteOk(w, r, models.StatusResponse{
		Ok: true,
	})
}
