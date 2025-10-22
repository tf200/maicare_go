package deps

import (
	"context"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/hub"
	"maicare_go/logger"
	"maicare_go/token"
	"maicare_go/util"
	"time"
)

type ServiceDependencies struct {
	Store      *db.Store
	TokenMaker token.Maker
	Logger     logger.Logger
	Config     *util.Config
	B2Client   bucket.ObjectStorageInterface
	GrpcClient grpclient.GrpcClientInterface
	WsHub      *hub.Hub
}

func NewServiceDependencies(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config, b2Client bucket.ObjectStorageInterface, grpcClient grpclient.GrpcClientInterface, wsHub *hub.Hub) *ServiceDependencies {
	return &ServiceDependencies{
		Store:      store,
		TokenMaker: tokenMaker,
		Logger:     logger,
		Config:     config,
		B2Client:   b2Client,
		GrpcClient: grpcClient,
		WsHub:      wsHub,
	}
}

func (d *ServiceDependencies) GenerateResponsePresignedURL(fileKey *string, ctx context.Context) *string {
	if fileKey == nil {
		return nil
	}

	url, err := d.B2Client.GeneratePresignedURL(ctx, *fileKey, time.Minute*15)
	if err != nil {
		return nil
	}
	return &url
}
