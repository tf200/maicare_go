package settings

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	db "maicare_go/db/sqlc"

	"github.com/google/uuid"
)

func (s *settingsService) ListRoles(ctx context.Context) ([]ListRolesApiResponse, error) {
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
			Description:     role.Description,
			PermissionCount: role.PermissionCount,
			EmployeeCount:   role.EmployeeCount,
		})
	}
	return response, nil
}

func (s *settingsService) ListAllPermissions(ctx context.Context) ([]PermissionGroupResponse, error) {
	permissions, err := s.Store.ListAllPermissions(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "ListAllPermissions", "Failed to list all permissions", err)
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	groupOrder := make([]string, 0)
	grouped := make(map[string]*PermissionGroupResponse)

	for _, perm := range permissions {
		permission := ListAllPermissionsApiResponse{
			PermissionID:       perm.ID,
			PermissionName:     perm.Name,
			PermissionResource: perm.Resource,
			DisplayName:        perm.DisplayName,
			Description:        perm.Description,
			SortOrder:          perm.SortOrder,
		}

		group, exists := grouped[perm.GroupKey]
		if !exists {
			groupOrder = append(groupOrder, perm.GroupKey)
			group = &PermissionGroupResponse{
				GroupKey:   perm.GroupKey,
				GroupLabel: humanizePermissionKey(perm.GroupKey),
				Sections:   make([]PermissionSectionResponse, 0),
			}
			grouped[perm.GroupKey] = group
		}

		sectionIndex := -1
		for i := range group.Sections {
			if group.Sections[i].SectionKey == perm.SectionKey {
				sectionIndex = i
				break
			}
		}

		if sectionIndex == -1 {
			group.Sections = append(group.Sections, PermissionSectionResponse{
				SectionKey:   perm.SectionKey,
				SectionLabel: humanizePermissionKey(perm.SectionKey),
				Permissions:  make([]ListAllPermissionsApiResponse, 0),
			})
			sectionIndex = len(group.Sections) - 1
		}

		group.Sections[sectionIndex].Permissions = append(group.Sections[sectionIndex].Permissions, permission)
	}

	response := make([]PermissionGroupResponse, 0, len(groupOrder))
	for _, key := range groupOrder {
		response = append(response, *grouped[key])
	}
	return response, nil
}

func (s *settingsService) ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]ListAllRolePermissionsApiResponse, error) {
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

func (s *settingsService) AssignRoleToEmployee(ctx context.Context, employeeID uuid.UUID, req *AssignRoleToEmployeeParams) (*AssignRoleToEmployeeApiResponse, error) {
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

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "AssignRoleToEmployee", "Failed to commit transaction", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &AssignRoleToEmployeeApiResponse{
		EmployeeID: employeeID,
		RoleID:     req.RoleID,
	}, nil
}

func (s *settingsService) ListUserRolesAndPermissionsApi(ctx context.Context, employeeID uuid.UUID) (*ListUserRolesAndPermissionsApiResponse, error) {
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

	inheritedPermissions, err := s.Store.ListInheritedUserPermissions(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "ListUserRolesAndPermissionsApi", "Failed to list inherited user permissions", err)
		return nil, fmt.Errorf("failed to list inherited user permissions: %w", err)
	}

	overrides, err := s.Store.ListUserPermissionOverrides(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "ListUserRolesAndPermissionsApi", "Failed to list user permission overrides", err)
		return nil, fmt.Errorf("failed to list user permission overrides: %w", err)
	}

	effectivePermissions, err := s.Store.ListEffectiveUserPermissions(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "ListUserRolesAndPermissionsApi", "Failed to list effective user permissions", err)
		return nil, fmt.Errorf("failed to list effective user permissions: %w", err)
	}

	inheritedList := make([]PermissionInfo, 0, len(inheritedPermissions))
	for _, perm := range inheritedPermissions {
		inheritedList = append(inheritedList, PermissionInfo{
			PermissionID:       perm.PermissionID,
			PermissionName:     perm.PermissionName,
			PermissionResource: perm.Resource,
		})
	}

	allowOverrides := make([]PermissionOverrideInfo, 0)
	denyOverrides := make([]PermissionOverrideInfo, 0)
	for _, override := range overrides {
		item := PermissionOverrideInfo{
			PermissionID:       override.PermissionID,
			PermissionName:     override.PermissionName,
			PermissionResource: override.Resource,
		}
		if override.Effect == db.PermissionOverrideEffectAllow {
			allowOverrides = append(allowOverrides, item)
			continue
		}
		denyOverrides = append(denyOverrides, item)
	}

	effectiveList := make([]PermissionInfo, 0, len(effectivePermissions))
	for _, perm := range effectivePermissions {
		effectiveList = append(effectiveList, PermissionInfo{
			PermissionID:       perm.PermissionID,
			PermissionName:     perm.PermissionName,
			PermissionResource: perm.Resource,
		})
	}

	var role *RoleInfo
	if len(roles) > 0 {
		role = &RoleInfo{
			RoleID:   roles[0].ID,
			RoleName: roles[0].Name,
		}
	}

	return &ListUserRolesAndPermissionsApiResponse{
		Role:                 role,
		InheritedPermissions: inheritedList,
		OverrideAllows:       allowOverrides,
		OverrideDenies:       denyOverrides,
		EffectivePermissions: effectiveList,
	}, nil
}

func (s *settingsService) ReplaceUserPermissionOverrides(ctx context.Context, employeeID uuid.UUID, req *ReplaceUserPermissionOverridesRequest) (*ReplaceUserPermissionOverridesResponse, error) {
	if hasPermissionOverlap(req.AllowPermissionIDs, req.DenyPermissionIDs) {
		return nil, fmt.Errorf("allow and deny permission ids cannot overlap")
	}

	userID, err := s.Store.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		s.Logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to get user ID by employee ID", err)
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to begin transaction", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to rollback transaction", rollbackErr)
		}
	}()

	qtx := s.Store.WithTx(tx)

	err = qtx.DeleteUserPermissionOverrides(ctx, userID)
	if err != nil {
		s.Logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to delete user permission overrides", err)
		return nil, fmt.Errorf("failed to delete user permission overrides: %w", err)
	}

	if len(req.AllowPermissionIDs) > 0 {
		err = qtx.AddUserPermissionOverrides(ctx, db.AddUserPermissionOverridesParams{
			UserID:        userID,
			PermissionIds: req.AllowPermissionIDs,
			Effect:        db.PermissionOverrideEffectAllow,
		})
		if err != nil {
			s.Logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to add allow overrides", err)
			return nil, fmt.Errorf("failed to add allow overrides: %w", err)
		}
	}

	if len(req.DenyPermissionIDs) > 0 {
		err = qtx.AddUserPermissionOverrides(ctx, db.AddUserPermissionOverridesParams{
			UserID:        userID,
			PermissionIds: req.DenyPermissionIDs,
			Effect:        db.PermissionOverrideEffectDeny,
		})
		if err != nil {
			s.Logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to add deny overrides", err)
			return nil, fmt.Errorf("failed to add deny overrides: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to commit transaction", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &ReplaceUserPermissionOverridesResponse{
		EmployeeID:         employeeID,
		AllowPermissionIDs: req.AllowPermissionIDs,
		DenyPermissionIDs:  req.DenyPermissionIDs,
	}, nil
}

func (s *settingsService) AddPermissionsToRole(ctx context.Context, roleID uuid.UUID, req *AddPermissionsToRoleRequest) (*AddPermissionsToRoleResponse, error) {
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

func (s *settingsService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error) {
	role, err := s.Store.CreateRole(ctx, db.CreateRoleParams{
		Name:        req.Name,
		Description: normalizeOptionalString(req.Description),
	})
	if err != nil {
		s.Logger.LogError(ctx, "CreateRole", "Failed to create role", err)
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &CreateRoleResponse{
		RoleID:      role.ID,
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func humanizePermissionKey(value string) string {
	parts := strings.Fields(strings.ReplaceAll(strings.ReplaceAll(value, "_", " "), ".", " "))
	for i, part := range parts {
		lower := strings.ToLower(part)
		parts[i] = strings.ToUpper(lower[:1]) + lower[1:]
	}
	return strings.Join(parts, " ")
}

func hasPermissionOverlap(allowIDs []uuid.UUID, denyIDs []uuid.UUID) bool {
	if len(allowIDs) == 0 || len(denyIDs) == 0 {
		return false
	}

	allowSet := make(map[uuid.UUID]struct{}, len(allowIDs))
	for _, id := range allowIDs {
		allowSet[id] = struct{}{}
	}

	for _, id := range denyIDs {
		if _, exists := allowSet[id]; exists {
			return true
		}
	}

	return false
}
