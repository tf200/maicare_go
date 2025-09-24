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
//
//go:generate mockgen -source=service.go -destination=../mocks/mock_auth_service.go -package=mocks
type AuthService interface {
	Login(req LoginUserRequest, clientIP string, userAgent string, ctx context.Context) (*LoginUserResponse, error)
	RefreshToken(req RefreshTokenRequest, ctx context.Context) (*RefreshTokenResponse, error)
	SetupTwoFA(userID int64, ctx context.Context) (*Setup2FAResponse, error)
	VerifyTwoFAToken(req Verify2FARequest, ctx context.Context) (*LoginUserResponse, error)
	Logout(req LogoutRequest, ctx context.Context) error
	ChangePassword(req ChangePasswordRequest, userID int64, ctx context.Context) error
	EnableTwoFA(req EnableTwoFARequest, ctx context.Context) (*EnableTwoFAResult, error)
}

type authService struct {
	*deps.ServiceDependencies
}

func NewAuthService(deps *deps.ServiceDependencies) AuthService {
	return &authService{
		ServiceDependencies: deps,
	}
}
