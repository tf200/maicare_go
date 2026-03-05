package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"maicare_go/service/deps"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials  = fmt.Errorf("invalid credentials")
	ErrUserNotFound        = fmt.Errorf("user not found")
	ErrSessionNotFound     = fmt.Errorf("session not found")
	ErrUnauthorized        = fmt.Errorf("unauthorized")
	ErrTwoFaAlreadyEnabled = fmt.Errorf("two-factor authentication already enabled")
	ErrTwoFARequired       = fmt.Errorf("two-factor authentication required")
	ErrInvalidTwoFACode    = fmt.Errorf("invalid two-factor authentication code")
	ErrTooManyAttempts     = fmt.Errorf("too many attempts")
)

// AuthService Interface and implementation
//
//go:generate mockgen -source=service.go -destination=../mocks/mock_auth_service.go -package=mocks
type AuthService interface {
	Login(req LoginUserRequest, clientIP string, userAgent string, ctx context.Context) (*LoginUserResponse, error)
	RefreshToken(req RefreshTokenRequest, ctx context.Context) (*RefreshTokenResponse, error)
	SetupTwoFA(req Setup2FARequest, userID uuid.UUID, ctx context.Context) (*Setup2FAResponse, error)
	VerifyTwoFAToken(req Verify2FARequest, clientIP string, userAgent string, ctx context.Context) (*LoginUserResponse, error)
	Logout(req LogoutRequest, ctx context.Context) error
	ChangePassword(req ChangePasswordRequest, userID uuid.UUID, ctx context.Context) error
	EnableTwoFA(req Enable2FARequest, userID uuid.UUID, ctx context.Context) (*Enable2FAResponse, error)

	// role methods
	ListRoles(ctx context.Context) ([]ListRolesApiResponse, error)
	ListAllPermissions(ctx context.Context) ([]ListAllPermissionsApiResponse, error)
	ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]ListAllRolePermissionsApiResponse, error)
	AssignRoleToEmployee(ctx context.Context, employeeID uuid.UUID, req *AssignRoleToEmployeeParams) (*AssignRoleToEmployeeApiResponse, error)
	ListUserRolesAndPermissionsApi(ctx context.Context, employeeID uuid.UUID) (*ListUserRolesAndPermissionsApiResponse, error)
	GrantUserPermission(ctx context.Context, employeeID uuid.UUID, req *GrantUserPermissionsRequest) (*GrantUserPermissionsResponse, error)
	AddPermissionsToRole(ctx context.Context, roleID uuid.UUID, req *AddPermissionsToRoleRequest) (*AddPermissionsToRoleResponse, error)
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error)
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type authService struct {
	*deps.ServiceDependencies
	mu            sync.Mutex
	loginAttempts map[string]attemptState
	twoFAAttempts map[string]attemptState
}

func NewAuthService(deps *deps.ServiceDependencies) AuthService {
	return &authService{
		ServiceDependencies: deps,
		loginAttempts:       make(map[string]attemptState),
		twoFAAttempts:       make(map[string]attemptState),
	}
}

type attemptState struct {
	Count       int
	FirstFailed time.Time
	LockedUntil time.Time
}
