package domain

import (
	"context"

	"github.com/google/uuid"
)

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
	Resource    string
	DisplayName string
	Description *string
	SortOrder   int32
	GroupKey    string
	SectionKey  string
}

// RolePermission represents a permission assigned to a role
type RolePermission struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	PermissionName string
	Resource       string
}

// UserRole represents a role assigned to a user
type UserRole struct {
	ID   uuid.UUID
	Name string
}

// UserPermission represents a permission associated with a user
type UserPermission struct {
	PermissionID uuid.UUID
	PermissionName string
	Resource       string
}

// UserPermissionOverride represents a permission override for a user
type UserPermissionOverride struct {
	PermissionID uuid.UUID
	PermissionName string
	Resource       string
	Effect         string // "allow" or "deny"
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
	Resource string
}

// PermissionOverrideInfo is a minimal permission override representation
type PermissionOverrideInfo struct {
	ID       uuid.UUID
	Name     string
	Resource string
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
	InheritedPermissions []PermissionInfo
	OverrideAllows       []PermissionOverrideInfo
	OverrideDenies       []PermissionOverrideInfo
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

// ReplaceUserPermissionOverridesParams parameters for replacing permission overrides
type ReplaceUserPermissionOverridesParams struct {
	EmployeeID         uuid.UUID
	AllowPermissionIDs []uuid.UUID
	DenyPermissionIDs  []uuid.UUID
}

// AddPermissionsToRoleParams parameters for adding permissions to a role
type AddPermissionsToRoleParams struct {
	RoleID        uuid.UUID
	PermissionIDs []uuid.UUID
}

// RoleRepository defines role and permission persistence operations
type RoleRepository interface {
	ListRoles(ctx context.Context) ([]Role, error)
	ListAllPermissions(ctx context.Context) ([]SystemPermission, error)
	ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]RolePermission, error)
	GetUserIDByEmployeeID(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error)
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error)
	ListInheritedUserPermissions(ctx context.Context, userID uuid.UUID) ([]UserPermission, error)
	ListUserPermissionOverrides(ctx context.Context, userID uuid.UUID) ([]UserPermissionOverride, error)
	ListEffectiveUserPermissions(ctx context.Context, userID uuid.UUID) ([]UserPermission, error)
	DeleteUserPermissionOverrides(ctx context.Context, userID uuid.UUID) error
	AddUserPermissionOverrides(ctx context.Context, userID uuid.UUID, permissionIDs []uuid.UUID, effect string) error
	RemovePermissionsFromRole(ctx context.Context, roleID uuid.UUID) error
	AddPermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	CreateRole(ctx context.Context, params CreateRoleParams) (*Role, error)
}

// RoleService defines role and permission business operations
type RoleService interface {
	ListRoles(ctx context.Context) ([]Role, error)
	ListAllPermissions(ctx context.Context) ([]PermissionGroup, error)
	ListAllRolePermissions(ctx context.Context, roleID uuid.UUID) ([]RolePermission, error)
	AssignRoleToEmployee(ctx context.Context, employeeID uuid.UUID, params AssignRoleToEmployeeParams) error
	ListUserRolesAndPermissions(ctx context.Context, employeeID uuid.UUID) (*UserRolesAndPermissions, error)
	ReplaceUserPermissionOverrides(ctx context.Context, params ReplaceUserPermissionOverridesParams) error
	AddPermissionsToRole(ctx context.Context, params AddPermissionsToRoleParams) error
	CreateRole(ctx context.Context, params CreateRoleParams) (*Role, error)
}
