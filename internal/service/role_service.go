package service

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

type RoleService struct {
	repo   domain.RoleRepository
	logger domain.Logger
}

func NewRoleService(repo domain.RoleRepository, logger domain.Logger) *RoleService {
	return &RoleService{repo: repo, logger: logger}
}

func (s *RoleService) ListRoles(ctx context.Context) ([]domain.Role, error) {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		s.logger.LogError(ctx, "ListRoles", "Failed to list roles", err)
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return roles, nil
}

func (s *RoleService) ListAllPermissions(ctx context.Context) ([]domain.PermissionGroup, error) {
	permissions, err := s.repo.ListAllPermissions(ctx)
	if err != nil {
		s.logger.LogError(ctx, "ListAllPermissions", "Failed to list all permissions", err)
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	groupOrder := make([]string, 0, 16)
	grouped := make(map[string]*domain.PermissionGroup, 16)

	for _, perm := range permissions {
		group, exists := grouped[perm.GroupKey]
		if !exists {
			groupOrder = append(groupOrder, perm.GroupKey)
			group = &domain.PermissionGroup{
				GroupKey:   perm.GroupKey,
				GroupLabel: humanizePermissionKey(perm.GroupKey),
				Sections:   make([]domain.PermissionSection, 0, 4),
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
			group.Sections = append(group.Sections, domain.PermissionSection{
				SectionKey:   perm.SectionKey,
				SectionLabel: humanizePermissionKey(perm.SectionKey),
				Permissions:  make([]domain.SystemPermission, 0, 8),
			})
			sectionIndex = len(group.Sections) - 1
		}

		group.Sections[sectionIndex].Permissions = append(group.Sections[sectionIndex].Permissions, perm)
	}

	result := make([]domain.PermissionGroup, 0, len(groupOrder))
	for _, key := range groupOrder {
		result = append(result, *grouped[key])
	}

	return result, nil
}

func (s *RoleService) ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain.RolePermission, error) {
	perms, err := s.repo.ListAllRolePermissions(ctx, roleID)
	if err != nil {
		s.logger.LogError(ctx, "ListAllRolePermissions", "Failed to list role permissions", err)
		return nil, fmt.Errorf("failed to list role permissions: %w", err)
	}
	return perms, nil
}

func (s *RoleService) AssignRoleToEmployee(ctx context.Context, employeeID uuid.UUID, params domain.AssignRoleToEmployeeParams) error {
	userID, err := s.repo.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		s.logger.LogError(ctx, "AssignRoleToEmployee", "Failed to get user ID by employee ID", err)
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	if err := s.repo.AssignRoleToUser(ctx, userID, params.RoleID); err != nil {
		s.logger.LogError(ctx, "AssignRoleToEmployee", "Failed to assign role to user", err)
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (s *RoleService) ListUserRolesAndPermissions(ctx context.Context, employeeID uuid.UUID) (*domain.UserRolesAndPermissions, error) {
	userID, err := s.repo.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		s.logger.LogError(ctx, "ListUserRolesAndPermissions", "Failed to get user ID by employee ID", err)
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	roles, err := s.repo.GetUserRoles(ctx, userID)
	if err != nil {
		s.logger.LogError(ctx, "ListUserRolesAndPermissions", "Failed to get user roles", err)
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	inherited, err := s.repo.ListInheritedUserPermissions(ctx, userID)
	if err != nil {
		s.logger.LogError(ctx, "ListUserRolesAndPermissions", "Failed to list inherited user permissions", err)
		return nil, fmt.Errorf("failed to list inherited user permissions: %w", err)
	}

	overrides, err := s.repo.ListUserPermissionOverrides(ctx, userID)
	if err != nil {
		s.logger.LogError(ctx, "ListUserRolesAndPermissions", "Failed to list user permission overrides", err)
		return nil, fmt.Errorf("failed to list user permission overrides: %w", err)
	}

	effective, err := s.repo.ListEffectiveUserPermissions(ctx, userID)
	if err != nil {
		s.logger.LogError(ctx, "ListUserRolesAndPermissions", "Failed to list effective user permissions", err)
		return nil, fmt.Errorf("failed to list effective user permissions: %w", err)
	}

	inheritedList := make([]domain.PermissionInfo, 0, len(inherited))
	for _, perm := range inherited {
		inheritedList = append(inheritedList, domain.PermissionInfo{
			ID:       perm.PermissionID,
			Name:     perm.PermissionName,
			Resource: perm.Resource,
		})
	}

	allowOverrides := make([]domain.PermissionOverrideInfo, 0)
	denyOverrides := make([]domain.PermissionOverrideInfo, 0)
	for _, override := range overrides {
		item := domain.PermissionOverrideInfo{
			ID:       override.PermissionID,
			Name:     override.PermissionName,
			Resource: override.Resource,
		}
		if override.Effect == "allow" {
			allowOverrides = append(allowOverrides, item)
		} else {
			denyOverrides = append(denyOverrides, item)
		}
	}

	effectiveList := make([]domain.PermissionInfo, 0, len(effective))
	for _, perm := range effective {
		effectiveList = append(effectiveList, domain.PermissionInfo{
			ID:       perm.PermissionID,
			Name:     perm.PermissionName,
			Resource: perm.Resource,
		})
	}

	var roleInfo *domain.RoleInfo
	if len(roles) > 0 {
		roleInfo = &domain.RoleInfo{
			ID:   roles[0].ID,
			Name: roles[0].Name,
		}
	}

	return &domain.UserRolesAndPermissions{
		Role:                 roleInfo,
		InheritedPermissions: inheritedList,
		OverrideAllows:       allowOverrides,
		OverrideDenies:       denyOverrides,
		EffectivePermissions: effectiveList,
	}, nil
}

func (s *RoleService) ReplaceUserPermissionOverrides(ctx context.Context, params domain.ReplaceUserPermissionOverridesParams) error {
	if hasPermissionOverlap(params.AllowPermissionIDs, params.DenyPermissionIDs) {
		return fmt.Errorf("allow and deny permission ids cannot overlap")
	}

	userID, err := s.repo.GetUserIDByEmployeeID(ctx, params.EmployeeID)
	if err != nil {
		s.logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to get user ID by employee ID", err)
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	if err := s.repo.DeleteUserPermissionOverrides(ctx, userID); err != nil {
		s.logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to delete user permission overrides", err)
		return fmt.Errorf("failed to replace user permission overrides: %w", err)
	}

	if len(params.AllowPermissionIDs) > 0 {
		if err := s.repo.AddUserPermissionOverrides(ctx, userID, params.AllowPermissionIDs, "allow"); err != nil {
			s.logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to add allow permission overrides", err)
			return fmt.Errorf("failed to replace user permission overrides: %w", err)
		}
	}

	if len(params.DenyPermissionIDs) > 0 {
		if err := s.repo.AddUserPermissionOverrides(ctx, userID, params.DenyPermissionIDs, "deny"); err != nil {
			s.logger.LogError(ctx, "ReplaceUserPermissionOverrides", "Failed to add deny permission overrides", err)
			return fmt.Errorf("failed to replace user permission overrides: %w", err)
		}
	}

	return nil
}

func (s *RoleService) AddPermissionsToRole(ctx context.Context, params domain.AddPermissionsToRoleParams) error {
	if err := s.repo.RemovePermissionsFromRole(ctx, params.RoleID); err != nil {
		s.logger.LogError(ctx, "AddPermissionsToRole", "Failed to remove permissions from role", err)
		return fmt.Errorf("failed to add permissions to role: %w", err)
	}

	if err := s.repo.AddPermissionsToRole(ctx, params.RoleID, params.PermissionIDs); err != nil {
		s.logger.LogError(ctx, "AddPermissionsToRole", "Failed to add permissions to role", err)
		return fmt.Errorf("failed to add permissions to role: %w", err)
	}

	return nil
}

func (s *RoleService) CreateRole(ctx context.Context, params domain.CreateRoleParams) (*domain.Role, error) {
	role, err := s.repo.CreateRole(ctx, params)
	if err != nil {
		s.logger.LogError(ctx, "CreateRole", "Failed to create role", err)
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return role, nil
}

func humanizePermissionKey(value string) string {
	if value == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(value))

	capitalizeNext := true
	lastWasSpace := false

	for _, r := range value {
		if r == '_' || r == '.' || r == ' ' {
			if b.Len() > 0 && !lastWasSpace {
				b.WriteByte(' ')
				lastWasSpace = true
			}
			capitalizeNext = true
		} else {
			if capitalizeNext {
				b.WriteRune(unicode.ToUpper(r))
				capitalizeNext = false
			} else {
				b.WriteRune(unicode.ToLower(r))
			}
			lastWasSpace = false
		}
	}

	res := b.String()
	if len(res) > 0 && res[len(res)-1] == ' ' {
		return res[:len(res)-1]
	}
	return res
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
