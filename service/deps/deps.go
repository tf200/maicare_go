package deps

import (
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/token"
	"maicare_go/util"
)

type ServiceDependencies struct {
	Store      *db.Store
	TokenMaker token.Maker
	Logger     logger.Logger
	Config     *util.Config
}

func NewServiceDependencies(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config) *ServiceDependencies {
	return &ServiceDependencies{
		Store:      store,
		TokenMaker: tokenMaker,
		Logger:     logger,
		Config:     config,
	}
}
