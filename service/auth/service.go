package auth

import (
	"context"
	"fmt"
	"maicare_go/service/deps"
)

var (
	ErrInvalidCredentials  = fmt.Errorf("invalid credentials")
	ErrUserNotFound        = fmt.Errorf("user not found")
	ErrSessionNotFound     = fmt.Errorf("session not found")
	ErrUnauthorized        = fmt.Errorf("unauthorized")
	ErrTwoFaAlreadyEnabled = fmt.Errorf("two-factor authentication already enabled")
	ErrTwoFARequired       = fmt.Errorf("two-factor authentication required")
	ErrInvalidTwoFACode    = fmt.Errorf("invalid two-factor authentication code")
)

// AuthService Interface and implementation
type AuthService interface {
	Login(req LoginRequest, ctx context.Context) (*LoginResult, error)
	RefreshToken(req RefreshTokenRequest, ctx context.Context) (*RefreshTokenResult, error)
}

type authService struct {
	*deps.ServiceDependencies
}

func NewAuthService(deps *deps.ServiceDependencies) AuthService {
	return &authService{
		ServiceDependencies: deps,
	}
}
