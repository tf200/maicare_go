package repository

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

type RoleRepository struct {
	queries *db.Queries
}

func NewRoleRepository(queries *db.Queries) *RoleRepository {
	return &RoleRepository{queries: queries}
}

func (r *RoleRepository) ListRoles(ctx context.Context) ([]domain.Role, error) {
	rows, err := r.queries.ListRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	roles := make([]domain.Role, len(rows))
	for i, row := range rows {
		roles[i] = domain.Role{
			ID:              row.ID,
			Name:            row.Name,
			Description:     row.Description,
			PermissionCount: row.PermissionCount,
			EmployeeCount:   row.EmployeeCount,
		}
	}
	return roles, nil
}

func (r *RoleRepository) ListAllPermissions(ctx context.Context) ([]domain.SystemPermission, error) {
	rows, err := r.queries.ListAllPermissions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	perms := make([]domain.SystemPermission, len(rows))
	for i, row := range rows {
		perms[i] = domain.SystemPermission{
			ID:          row.ID,
			Name:        row.Name,
			Resource:    row.Resource,
			DisplayName: row.DisplayName,
			Description: row.Description,
			SortOrder:   row.SortOrder,
			GroupKey:    row.GroupKey,
			SectionKey:  row.SectionKey,
		}
	}
	return perms, nil
}

func (r *RoleRepository) ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain.RolePermission, error) {
	rows, err := r.queries.ListAllRolePermissions(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list role permissions: %w", err)
	}

	perms := make([]domain.RolePermission, len(rows))
	for i, row := range rows {
		perms[i] = domain.RolePermission{
			RoleID:         roleID,
			PermissionID:   row.PermissionID,
			PermissionName: row.PermissionName,
			Resource:       row.Resource,
		}
	}
	return perms, nil
}

func (r *RoleRepository) GetUserIDByEmployeeID(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	userID, err := r.queries.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user id by employee id: %w", err)
	}
	return userID, nil
}

func (r *RoleRepository) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	return r.queries.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	})
}

func (r *RoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]domain.UserRole, error) {
	rows, err := r.queries.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roles := make([]domain.UserRole, len(rows))
	for i, row := range rows {
		roles[i] = domain.UserRole{
			ID:   row.ID,
			Name: row.Name,
		}
	}
	return roles, nil
}

func (r *RoleRepository) ListInheritedUserPermissions(ctx context.Context, userID uuid.UUID) ([]domain.UserPermission, error) {
	rows, err := r.queries.ListInheritedUserPermissions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inherited user permissions: %w", err)
	}

	perms := make([]domain.UserPermission, len(rows))
	for i, row := range rows {
		perms[i] = domain.UserPermission{
			PermissionID:   row.PermissionID,
			PermissionName: row.PermissionName,
			Resource:       row.Resource,
		}
	}
	return perms, nil
}

func (r *RoleRepository) ListUserPermissionOverrides(ctx context.Context, userID uuid.UUID) ([]domain.UserPermissionOverride, error) {
	rows, err := r.queries.ListUserPermissionOverrides(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user permission overrides: %w", err)
	}

	overrides := make([]domain.UserPermissionOverride, len(rows))
	for i, row := range rows {
		overrides[i] = domain.UserPermissionOverride{
			PermissionID:   row.PermissionID,
			PermissionName: row.PermissionName,
			Resource:       row.Resource,
			Effect:         string(row.Effect),
		}
	}
	return overrides, nil
}

func (r *RoleRepository) ListEffectiveUserPermissions(ctx context.Context, userID uuid.UUID) ([]domain.UserPermission, error) {
	rows, err := r.queries.ListEffectiveUserPermissions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list effective user permissions: %w", err)
	}

	perms := make([]domain.UserPermission, len(rows))
	for i, row := range rows {
		perms[i] = domain.UserPermission{
			PermissionID:   row.PermissionID,
			PermissionName: row.PermissionName,
			Resource:       row.Resource,
		}
	}
	return perms, nil
}

func (r *RoleRepository) DeleteUserPermissionOverrides(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteUserPermissionOverrides(ctx, userID)
}

func (r *RoleRepository) AddUserPermissionOverrides(ctx context.Context, userID uuid.UUID, permissionIDs []uuid.UUID, effect string) error {
	var eff db.PermissionOverrideEffect
	switch effect {
	case "allow":
		eff = db.PermissionOverrideEffectAllow
	case "deny":
		eff = db.PermissionOverrideEffectDeny
	default:
		return fmt.Errorf("invalid permission override effect: %s", effect)
	}

	return r.queries.AddUserPermissionOverrides(ctx, db.AddUserPermissionOverridesParams{
		UserID:        userID,
		PermissionIds: permissionIDs,
		Effect:        eff,
	})
}

func (r *RoleRepository) RemovePermissionsFromRole(ctx context.Context, roleID uuid.UUID) error {
	return r.queries.RemovePermissionsFromRole(ctx, roleID)
}

func (r *RoleRepository) AddPermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return r.queries.AddPermissionsToRole(ctx, db.AddPermissionsToRoleParams{
		RoleID:        roleID,
		PermissionIds: permissionIDs,
	})
}

func (r *RoleRepository) CreateRole(ctx context.Context, params domain.CreateRoleParams) (*domain.Role, error) {
	role, err := r.queries.CreateRole(ctx, db.CreateRoleParams{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &domain.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}, nil
}
