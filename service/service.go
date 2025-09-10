package service

import (
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/token"
	"maicare_go/util"
)

type BusinessService struct {
	AuthService AuthService
}

func NewBusinessService(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config) *BusinessService {
	authService := NewAuthService(tokenMaker, store, logger, config)
	return &BusinessService{
		AuthService: authService,
	}
}
