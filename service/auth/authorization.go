package auth

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"

	"github.com/google/uuid"
)

func (s *authService) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	return s.Store.CheckUserPermission(ctx, db.CheckUserPermissionParams{
		UserID: userID,
		Name:   permission,
	})
}

func (s *authService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	roles, err := s.Store.GetUserRoles(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "GetUserRoles", "Failed to get user roles", err)
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames, nil
}
