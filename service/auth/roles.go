package auth

import (
	"context"
	"database/sql"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *authService) ListRoles(ctx context.Context) ([]ListRolesApiResponse, error) {
	roles, err := s.Store.ListRoles(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListRoles", "Failed to list roles", zap.Error(err))
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	response := []ListRolesApiResponse{}
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
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListAllPermissions", "Failed to list all permissions", zap.Error(err))
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	response := []ListAllPermissionsApiResponse{}
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
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListAllRolePermissions", "Failed to list all role permissions", zap.Error(err))
		return nil, fmt.Errorf("failed to list role permissions: %w", err)
	}

	response := []ListAllRolePermissionsApiResponse{}
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
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignRoleToEmployee", "Failed to get user ID by employee ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignRoleToEmployee", "Failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignRoleToEmployee", "Failed to rollback transaction", zap.Error(rollbackErr), zap.String("employee_id", employeeID.String()), zap.String("role_id", req.RoleID.String()))
		}
	}()

	qtx := s.Store.WithTx(tx)

	err = qtx.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: req.RoleID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignRoleToEmployee", "Failed to assign role to user", zap.Error(err), zap.String("employee_id", employeeID.String()), zap.String("role_id", req.RoleID.String()))
		return nil, fmt.Errorf("failed to assign role to user: %w", err)
	}

	err = qtx.DeleteUserPermissions(ctx, userID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignRoleToEmployee", "Failed to delete user permissions", zap.Error(err), zap.String("employee_id", employeeID.String()), zap.String("role_id", req.RoleID.String()))
		return nil, fmt.Errorf("failed to delete user permissions: %w", err)
	}

	err = qtx.GrantRolePermissionsToUser(ctx, db.GrantRolePermissionsToUserParams{
		UserID: userID,
		RoleID: req.RoleID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignRoleToEmployee", "Failed to grant role permissions to user", zap.Error(err), zap.String("employee_id", employeeID.String()), zap.String("role_id", req.RoleID.String()))
		return nil, fmt.Errorf("failed to grant role permissions to user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignRoleToEmployee", "Failed to commit transaction", zap.Error(err), zap.String("employee_id", employeeID.String()), zap.String("role_id", req.RoleID.String()))
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
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListUserRolesAndPermissionsApi", "Failed to get user ID by employee ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	roles, err := s.Store.GetUserRoles(ctx, userID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListUserRolesAndPermissionsApi", "Failed to get user roles", zap.Error(err))
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	permissions, err := s.Store.ListUserPermissions(ctx, userID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListUserRolesAndPermissionsApi", "Failed to list user permissions", zap.Error(err))
		return nil, fmt.Errorf("failed to list user permissions: %w", err)
	}

	response := ListUserRolesAndPermissionsApiResponse{
		Roles: struct {
			RoleID   uuid.UUID `json:"id"`
			RoleName string    `json:"name"`
		}{
			RoleID:   roles[0].ID,
			RoleName: roles[0].Name,
		},
		Permissions: []struct {
			PermissionID       uuid.UUID `json:"id"`
			PermissionName     string    `json:"name"`
			PermissionResource string    `json:"resource"`
		}{},
	}
	for _, perm := range permissions {
		response.Permissions = append(response.Permissions, struct {
			PermissionID       uuid.UUID `json:"id"`
			PermissionName     string    `json:"name"`
			PermissionResource string    `json:"resource"`
		}{
			PermissionID:       perm.PermissionID,
			PermissionName:     perm.PermissionName,
			PermissionResource: perm.Resource,
		})
	}

	return &response, nil
}

func (s *authService) GrantUserPermission(ctx context.Context, employeeID uuid.UUID, req *GrantUserPermissionsRequest) (*GrantUserPermissionsResponse, error) {
	userID, err := s.Store.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GrantUserPermission", "Failed to get user ID by employee ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GrantUserPermission", "Failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GrantUserPermission", "Failed to rollback transaction", zap.Error(rollbackErr), zap.String("employee_id", employeeID.String()))
		}
	}()

	qtx := s.Store.WithTx(tx)

	err = qtx.DeleteUserPermissions(ctx, userID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GrantUserPermission", "Failed to delete user permissions", zap.Error(err), zap.String("employee_id", employeeID.String()))
		return nil, fmt.Errorf("failed to delete user permissions: %w", err)
	}

	err = qtx.GrantUserPermissions(ctx, db.GrantUserPermissionsParams{
		UserID:        userID,
		PermissionIds: req.PermissionIDs,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GrantUserPermission", "Failed to grant user permissions", zap.Error(err), zap.String("employee_id", employeeID.String()))
		return nil, fmt.Errorf("failed to grant user permissions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GrantUserPermission", "Failed to commit transaction", zap.Error(err), zap.String("employee_id", employeeID.String()))
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
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddPermissionsToRole", "Failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddPermissionsToRole", "Failed to rollback transaction", zap.Error(rollbackErr), zap.String("role_id", roleID.String()))
		}
	}()

	qtx := s.Store.WithTx(tx)

	err = qtx.RemovePermissionsFromRole(ctx, roleID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddPermissionsToRole", "Failed to remove existing permissions from role", zap.Error(err), zap.String("role_id", roleID.String()))
		return nil, fmt.Errorf("failed to remove existing permissions from role: %w", err)
	}

	err = qtx.AddPermissionsToRole(ctx, db.AddPermissionsToRoleParams{
		RoleID:        roleID,
		PermissionIds: req.PermissionIDs,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddPermissionsToRole", "Failed to add permissions to role", zap.Error(err), zap.String("role_id", roleID.String()))
		return nil, fmt.Errorf("failed to add permissions to role: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddPermissionsToRole", "Failed to commit transaction", zap.Error(err), zap.String("role_id", roleID.String()))
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
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateRole", "Failed to create role", zap.Error(err), zap.String("role_name", req.Name))
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
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetUserRoles", "Failed to get user roles", zap.Error(err), zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames, nil
}
