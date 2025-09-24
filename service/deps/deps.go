package deps

import (
	"context"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
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
	b2Client   bucket.ObjectStorageInterface
}

func NewServiceDependencies(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config, b2Client bucket.ObjectStorageInterface) *ServiceDependencies {
	return &ServiceDependencies{
		Store:      store,
		TokenMaker: tokenMaker,
		Logger:     logger,
		Config:     config,
		b2Client:   b2Client,
	}
}

func (d *ServiceDependencies) GenerateResponsePresignedURL(fileKey *string, ctx context.Context) *string {
	if fileKey == nil {
		return nil
	}

	url, err := d.b2Client.GeneratePresignedURL(ctx, *fileKey, time.Minute*15)
	if err != nil {
		return nil
	}
	return &url
}
