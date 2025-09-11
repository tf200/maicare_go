package service

import (
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/auth"
	"maicare_go/service/deps"
	"maicare_go/token"
	"maicare_go/util"
)

type BusinessService struct {
	*deps.ServiceDependencies
	AuthService auth.AuthService
}

func NewBusinessService(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config) *BusinessService {
	deps := deps.NewServiceDependencies(store, tokenMaker, logger, config)
	authService := auth.NewAuthService(deps)
	return &BusinessService{
		ServiceDependencies: deps,
		AuthService:         authService,
	}
}
