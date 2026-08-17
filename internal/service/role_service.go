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

	effective, err := s.repo.ListEffectiveUserPermissions(ctx, userID)
	if err != nil {
		s.logger.LogError(ctx, "ListUserRolesAndPermissions", "Failed to list effective user permissions", err)
		return nil, fmt.Errorf("failed to list effective user permissions: %w", err)
	}

	effectiveList := make([]domain.PermissionInfo, 0, len(effective))
	for _, perm := range effective {
		effectiveList = append(effectiveList, domain.PermissionInfo{
			ID:   perm.PermissionID,
			Name: perm.PermissionName,
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
		EffectivePermissions: effectiveList,
	}, nil
}

func (s *RoleService) ReplaceRolePermissions(ctx context.Context, params domain.ReplaceRolePermissionsParams) error {
	permissions, err := s.repo.ListAllPermissions(ctx)
	if err != nil {
		s.logger.LogError(ctx, "ReplaceRolePermissions", "Failed to list permissions", err)
		return fmt.Errorf("failed to replace role permissions: %w", err)
	}

	byID := make(map[uuid.UUID]domain.SystemPermission, len(permissions))
	for _, permission := range permissions {
		byID[permission.ID] = permission
	}

	seen := make(map[uuid.UUID]struct{}, len(params.Permissions))
	for _, grant := range params.Permissions {
		permission, exists := byID[grant.PermissionID]
		if !exists {
			return fmt.Errorf("%w: unknown permission_id %s", domain.ErrInvalidRolePermissions, grant.PermissionID)
		}
		if _, exists := seen[grant.PermissionID]; exists {
			return fmt.Errorf("%w: duplicate permission_id %s", domain.ErrInvalidRolePermissions, grant.PermissionID)
		}
		seen[grant.PermissionID] = struct{}{}

		if permission.IsScoped {
			if grant.Scope == nil {
				return fmt.Errorf("%w: scope is required for permission %s", domain.ErrInvalidRolePermissions, permission.Name)
			}
			if !grant.Scope.IsValid() {
				return fmt.Errorf("%w: invalid scope %q for permission %s", domain.ErrInvalidRolePermissions, *grant.Scope, permission.Name)
			}
		} else if grant.Scope != nil {
			return fmt.Errorf("%w: scope must be null for permission %s", domain.ErrInvalidRolePermissions, permission.Name)
		}
	}

	if err := s.repo.ReplaceRolePermissions(ctx, params.RoleID, params.Permissions); err != nil {
		s.logger.LogError(ctx, "ReplaceRolePermissions", "Failed to replace permissions for role", err)
		return fmt.Errorf("failed to replace role permissions: %w", err)
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
