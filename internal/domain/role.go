package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidRolePermissions = errors.New("invalid role permissions")

// Role represents a system role with metadata
type Role struct {
	ID              uuid.UUID
	Name            string
	Description     *string
	PermissionCount int64
	EmployeeCount   int64
}

// SystemPermission represents a system permission
type SystemPermission struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Description *string
	GroupKey    string
	SectionKey  string
	IsScoped    bool
}

// RolePermission represents a permission assigned to a role
type RolePermission struct {
	RoleID         uuid.UUID
	PermissionID   uuid.UUID
	PermissionName string
	IsScoped       bool
	Scope          *PermissionScope
}

// UserRole represents a role assigned to a user
type UserRole struct {
	ID   uuid.UUID
	Name string
}

// UserPermission represents a permission associated with a user
type UserPermission struct {
	PermissionID   uuid.UUID
	PermissionName string
	IsScoped       bool
	Scope          *PermissionScope
}

func (p UserPermission) IsUsable() bool {
	if !p.IsScoped {
		return p.Scope == nil
	}
	return p.Scope != nil && p.Scope.IsValid()
}

// RoleInfo is a minimal role representation
type RoleInfo struct {
	ID   uuid.UUID
	Name string
}

// PermissionInfo is a minimal permission representation
type PermissionInfo struct {
	ID       uuid.UUID
	Name     string
	IsScoped bool
	Scope    *PermissionScope
}

// PermissionSection groups permissions by section
type PermissionSection struct {
	SectionKey   string
	SectionLabel string
	Permissions  []SystemPermission
}

// PermissionGroup groups permissions by group key
type PermissionGroup struct {
	GroupKey   string
	GroupLabel string
	Sections   []PermissionSection
}

// UserRolesAndPermissions aggregates a user's roles and permissions
type UserRolesAndPermissions struct {
	Role                 *RoleInfo
	EffectivePermissions []PermissionInfo
}

// CreateRoleParams parameters for creating a role
type CreateRoleParams struct {
	Name        string
	Description *string
}

// AssignRoleToEmployeeParams parameters for assigning a role
type AssignRoleToEmployeeParams struct {
	RoleID uuid.UUID
}

// PermissionGrant represents one permission and its optional scope on a role.
type PermissionGrant struct {
	PermissionID uuid.UUID
	Scope        *PermissionScope
}

// ReplaceRolePermissionsParams parameters for replacing a role's permission grants
type ReplaceRolePermissionsParams struct {
	RoleID      uuid.UUID
	Permissions []PermissionGrant
}

// RoleRepository defines role and permission persistence operations
type RoleRepository interface {
	ListRoles(ctx context.Context) ([]Role, error)
	ListAllPermissions(ctx context.Context) ([]SystemPermission, error)
	ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]RolePermission, error)
	GetUserIDByEmployeeID(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error)
	ListEffectiveUserPermissions(ctx context.Context, userID uuid.UUID) ([]UserPermission, error)
	GetEffectiveUserPermission(ctx context.Context, userID uuid.UUID, permissionName string) (*UserPermission, error)
	ReplaceRolePermissions(ctx context.Context, roleID uuid.UUID, permissions []PermissionGrant) error
	CreateRole(ctx context.Context, params CreateRoleParams) (*Role, error)
}

// RoleService defines role and permission business operations
type RoleService interface {
	ListRoles(ctx context.Context) ([]Role, error)
	ListAllPermissions(ctx context.Context) ([]PermissionGroup, error)
	ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]RolePermission, error)
	AssignRoleToEmployee(ctx context.Context, employeeID uuid.UUID, params AssignRoleToEmployeeParams) error
	ListUserRolesAndPermissions(ctx context.Context, employeeID uuid.UUID) (*UserRolesAndPermissions, error)
	ReplaceRolePermissions(ctx context.Context, params ReplaceRolePermissionsParams) error
	CreateRole(ctx context.Context, params CreateRoleParams) (*Role, error)
}
