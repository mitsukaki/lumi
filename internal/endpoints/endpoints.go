package endpoints

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mitsukaki/lumi/internal/db"
	"go.uber.org/zap"
)

type EndpointHandler struct {
	logger *zap.Logger

	userTable  *db.CouchTable
	albumTable *db.CouchTable

	svc *s3.S3
}

func CreateEndpointHandler(
	logger *zap.Logger,
	userTable *db.CouchTable,
	albumTable *db.CouchTable,
	svc *s3.S3,
) *EndpointHandler {
	return &EndpointHandler{
		logger: logger,

		userTable:  userTable,
		albumTable: albumTable,

		svc: svc,
	}
}