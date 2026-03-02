package auth

import (
	"context"
	"database/sql"
	"fmt"

	db "maicare_go/db/sqlc"

	"github.com/google/uuid"
)

func (s *authService) ListRoles(ctx context.Context) ([]ListRolesApiResponse, error) {
	roles, err := s.Store.ListRoles(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "ListRoles", "Failed to list roles", err)
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	response := make([]ListRolesApiResponse, 0, len(roles))
	for _, role := range roles {
		response = append(response, ListRolesApiResponse{
			ID:              role.ID,
			RoleName:        role.Name,
			PermissionCount: role.PermissionCount,
		})
	}
	return response, nil
}

func (s *authService) ListAllPermissions(ctx context.Context) ([]ListAllPermissionsApiResponse, error) {
	permissions, err := s.Store.ListAllPermissions(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "ListAllPermissions", "Failed to list all permissions", err)
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	response := make([]ListAllPermissionsApiResponse, 0, len(permissions))
	for _, perm := range permissions {
		response = append(response, ListAllPermissionsApiResponse{
			PermissionID:       perm.ID,
			PermissionName:     perm.Name,
			PermissionResource: perm.Resource,
		})
	}
	return response, nil
}

func (s *authService) ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]ListAllRolePermissionsApiResponse, error) {
	rolePermissions, err := s.Store.ListAllRolePermissions(ctx, roleID)
	if err != nil {
		s.Logger.LogError(ctx, "ListAllRolePermissions", "Failed to list all role permissions", err)
		return nil, fmt.Errorf("failed to list role permissions: %w", err)
	}

	response := make([]ListAllRolePermissionsApiResponse, 0, len(rolePermissions))
	for _, rp := range rolePermissions {
		response = append(response, ListAllRolePermissionsApiResponse{
			RoleID:             roleID,
			PermissionID:       rp.PermissionID,
			PermissionName:     rp.PermissionName,
			PermissionResource: rp.Resource,
		})
	}
	return response, nil
}

func (s *authService) AssignRoleToEmployee(ctx context.Context, employeeID uuid.UUID, req *AssignRoleToEmployeeParams) (*AssignRoleToEmployeeApiResponse, error) {
	userID, err := s.Store.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		s.Logger.LogError(ctx, "AssignRoleToEmployee", "Failed to get user ID by employee ID", err)
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "AssignRoleToEmployee", "Failed to begin transaction", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogError(ctx, "AssignRoleToEmployee", "Failed to rollback transaction", rollbackErr)
		}
	}()

	qtx := s.Store.WithTx(tx)

	err = qtx.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: req.RoleID,
	})
	if err != nil {
		s.Logger.LogError(ctx, "AssignRoleToEmployee", "Failed to assign role to user", err)
		return nil, fmt.Errorf("failed to assign role to user: %w", err)
	}

	err = qtx.DeleteUserPermissions(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "AssignRoleToEmployee", "Failed to delete user permissions", err)
		return nil, fmt.Errorf("failed to delete user permissions: %w", err)
	}

	err = qtx.GrantRolePermissionsToUser(ctx, db.GrantRolePermissionsToUserParams{
		UserID: userID,
		RoleID: req.RoleID,
	})
	if err != nil {
		s.Logger.LogError(ctx, "AssignRoleToEmployee", "Failed to grant role permissions to user", err)
		return nil, fmt.Errorf("failed to grant role permissions to user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "AssignRoleToEmployee", "Failed to commit transaction", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &AssignRoleToEmployeeApiResponse{
		EmployeeID: employeeID,
		RoleID:     req.RoleID,
	}, nil
}

func (s *authService) ListUserRolesAndPermissionsApi(ctx context.Context, employeeID uuid.UUID) (*ListUserRolesAndPermissionsApiResponse, error) {
	userID, err := s.Store.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		s.Logger.LogError(ctx, "ListUserRolesAndPermissionsApi", "Failed to get user ID by employee ID", err)
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	roles, err := s.Store.GetUserRoles(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "ListUserRolesAndPermissionsApi", "Failed to get user roles", err)
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("user has no role assigned")
	}

	permissions, err := s.Store.ListUserPermissions(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "ListUserRolesAndPermissionsApi", "Failed to list user permissions", err)
		return nil, fmt.Errorf("failed to list user permissions: %w", err)
	}

	permissionList := make([]PermissionInfo, 0, len(permissions))
	for _, perm := range permissions {
		permissionList = append(permissionList, PermissionInfo{
			PermissionID:       perm.PermissionID,
			PermissionName:     perm.PermissionName,
			PermissionResource: perm.Resource,
		})
	}

	return &ListUserRolesAndPermissionsApiResponse{
		Roles: RoleInfo{
			RoleID:   roles[0].ID,
			RoleName: roles[0].Name,
		},
		Permissions: permissionList,
	}, nil
}

func (s *authService) GrantUserPermission(ctx context.Context, employeeID uuid.UUID, req *GrantUserPermissionsRequest) (*GrantUserPermissionsResponse, error) {
	userID, err := s.Store.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		s.Logger.LogError(ctx, "GrantUserPermission", "Failed to get user ID by employee ID", err)
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "GrantUserPermission", "Failed to begin transaction", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogError(ctx, "GrantUserPermission", "Failed to rollback transaction", rollbackErr)
		}
	}()

	qtx := s.Store.WithTx(tx)

	err = qtx.DeleteUserPermissions(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "GrantUserPermission", "Failed to delete user permissions", err)
		return nil, fmt.Errorf("failed to delete user permissions: %w", err)
	}

	err = qtx.GrantUserPermissions(ctx, db.GrantUserPermissionsParams{
		UserID:        userID,
		PermissionIds: req.PermissionIDs,
	})
	if err != nil {
		s.Logger.LogError(ctx, "GrantUserPermission", "Failed to grant user permissions", err)
		return nil, fmt.Errorf("failed to grant user permissions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "GrantUserPermission", "Failed to commit transaction", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &GrantUserPermissionsResponse{
		EmployeeID:    employeeID,
		PermissionIDs: req.PermissionIDs,
	}, nil
}

func (s *authService) AddPermissionsToRole(ctx context.Context, roleID uuid.UUID, req *AddPermissionsToRoleRequest) (*AddPermissionsToRoleResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "AddPermissionsToRole", "Failed to begin transaction", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogError(ctx, "AddPermissionsToRole", "Failed to rollback transaction", rollbackErr)
		}
	}()

	qtx := s.Store.WithTx(tx)

	err = qtx.RemovePermissionsFromRole(ctx, roleID)
	if err != nil {
		s.Logger.LogError(ctx, "AddPermissionsToRole", "Failed to remove existing permissions from role", err)
		return nil, fmt.Errorf("failed to remove existing permissions from role: %w", err)
	}

	err = qtx.AddPermissionsToRole(ctx, db.AddPermissionsToRoleParams{
		RoleID:        roleID,
		PermissionIds: req.PermissionIDs,
	})
	if err != nil {
		s.Logger.LogError(ctx, "AddPermissionsToRole", "Failed to add permissions to role", err)
		return nil, fmt.Errorf("failed to add permissions to role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "AddPermissionsToRole", "Failed to commit transaction", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &AddPermissionsToRoleResponse{
		RoleID:        roleID,
		PermissionIDs: req.PermissionIDs,
	}, nil
}

func (s *authService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error) {
	role, err := s.Store.CreateRole(ctx, req.Name)
	if err != nil {
		s.Logger.LogError(ctx, "CreateRole", "Failed to create role", err)
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &CreateRoleResponse{
		RoleID: role.ID,
		Name:   role.Name,
	}, nil
}

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
