package service

import (
	"context"
	"database/sql"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/token"
	"maicare_go/util"
	"strings"

	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrUserNotFound       = fmt.Errorf("user not found")
)

// AuthService Interface and implementation
type AuthService interface {
	Login(req LoginRequest, ctx context.Context) (*LoginResult, error)
}

type authService struct {
	tokenMaker token.Maker
	store      *db.Store
	logger     logger.Logger
	config     *util.Config
}

func NewAuthService(tokenMaker token.Maker, store *db.Store, logger logger.Logger, config *util.Config) AuthService {
	return &authService{
		tokenMaker: tokenMaker,
		store:      store,
		logger:     logger,
		config:     config,
	}
}

// Methods and types for AuthService

type LoginRequest struct {
	Email     string
	Password  string
	ClientIP  string
	UserAgent string
}

type LoginResult struct {
	AccessToken   string
	RefreshToken  string
	RequiresTwoFA bool
	TempToken     string
}

func (s *authService) Login(req LoginRequest, ctx context.Context) (*LoginResult, error) {
	email := strings.ToLower(req.Email)

	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.LogBusinessEvent(logger.LogLevelWarn, "Login", "Failed login attempt: user not found", zap.String("email", email), zap.String("client_ip", req.ClientIP), zap.String("user_agent", req.UserAgent))
			return nil, ErrInvalidCredentials
		}
		s.logger.LogBusinessEvent(logger.LogLevelError, "Login", "Database error during login", zap.String("email", email), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		s.logger.LogBusinessEvent(logger.LogLevelWarn, "Login", "Failed login attempt: incorrect password", zap.String("email", email), zap.String("client_ip", req.ClientIP), zap.String("user_agent", req.UserAgent))
		return nil, ErrInvalidCredentials
	}

	if user.TwoFactorEnabled {
		tempToken, _, err := s.tokenMaker.CreateToken(user.ID, user.EmployeeID, s.config.TwoFATokenDuration, token.TwoFAToken)
		if err != nil {
			s.logger.LogBusinessEvent(logger.LogLevelError, "Login", "Failed to create 2FA token", zap.String("email", email), zap.String("error", err.Error()))
			return nil, fmt.Errorf("failed to create 2FA token: %v", err)
		}
		s.logger.LogBusinessEvent(logger.LogLevelInfo, "Login", "2FA required for user", zap.String("email", email), zap.String("client_ip", req.ClientIP), zap.String("user_agent", req.UserAgent))
		return &LoginResult{
			RequiresTwoFA: true,
			TempToken:     tempToken,
		}, nil
	}

	accessToken, _, err := s.tokenMaker.CreateToken(user.ID, user.EmployeeID, s.config.AccessTokenDuration, token.AccessToken)
	if err != nil {
		s.logger.LogBusinessEvent(logger.LogLevelError, "Login", "Failed to create access token", zap.String("email", email), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create access token: %v", err)
	}

	refreshToken, _, err := s.tokenMaker.CreateToken(user.ID, user.EmployeeID, s.config.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		s.logger.LogBusinessEvent(logger.LogLevelError, "Login", "Failed to create refresh token", zap.String("email", email), zap.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create refresh token: %v", err)
	}

	s.logger.LogBusinessEvent(logger.LogLevelInfo, "Login", "User logged in successfully", zap.String("email", email), zap.String("client_ip", req.ClientIP), zap.String("user_agent", req.UserAgent))

	return &LoginResult{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		RequiresTwoFA: false,
		TempToken:     "",
	}, nil
}
