package settings

import "github.com/google/uuid"

type ListRolesApiResponse struct {
	ID              uuid.UUID `json:"id"`
	RoleName        string    `json:"role_name"`
	Description     *string   `json:"description"`
	PermissionCount int64     `json:"permission_count"`
	EmployeeCount   int64     `json:"employee_count"`
}

type ListAllPermissionsApiResponse struct {
	PermissionID       uuid.UUID `json:"permission_id"`
	PermissionName     string    `json:"permission_name"`
	PermissionResource string    `json:"permission_resource"`
	DisplayName        string    `json:"display_name"`
	Description        *string   `json:"description"`
	SortOrder          int32     `json:"sort_order"`
}

type PermissionSectionResponse struct {
	SectionKey   string                          `json:"section_key"`
	SectionLabel string                          `json:"section_label"`
	Permissions  []ListAllPermissionsApiResponse `json:"permissions"`
}

type PermissionGroupResponse struct {
	GroupKey   string                      `json:"group_key"`
	GroupLabel string                      `json:"group_label"`
	Sections   []PermissionSectionResponse `json:"sections"`
}

type ListAllRolePermissionsApiResponse struct {
	RoleID             uuid.UUID `json:"role_id"`
	PermissionID       uuid.UUID `json:"permission_id"`
	PermissionName     string    `json:"permission_name"`
	PermissionResource string    `json:"permission_resource"`
}

type AssignRoleToEmployeeParams struct {
	RoleID uuid.UUID `json:"role_id"`
}

type AssignRoleToEmployeeApiResponse struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	RoleID     uuid.UUID `json:"role_id"`
}

type RoleInfo struct {
	RoleID   uuid.UUID `json:"id"`
	RoleName string    `json:"name"`
}

type PermissionInfo struct {
	PermissionID       uuid.UUID `json:"id"`
	PermissionName     string    `json:"name"`
	PermissionResource string    `json:"resource"`
}

type PermissionOverrideInfo struct {
	PermissionID       uuid.UUID `json:"id"`
	PermissionName     string    `json:"name"`
	PermissionResource string    `json:"resource"`
}

type ListUserRolesAndPermissionsApiResponse struct {
	Role                 *RoleInfo                `json:"role"`
	InheritedPermissions []PermissionInfo         `json:"inherited_permissions"`
	OverrideAllows       []PermissionOverrideInfo `json:"override_allows"`
	OverrideDenies       []PermissionOverrideInfo `json:"override_denies"`
	EffectivePermissions []PermissionInfo         `json:"effective_permissions"`
}

type ReplaceUserPermissionOverridesRequest struct {
	AllowPermissionIDs []uuid.UUID `json:"allow_permission_ids"`
	DenyPermissionIDs  []uuid.UUID `json:"deny_permission_ids"`
}

type ReplaceUserPermissionOverridesResponse struct {
	EmployeeID         uuid.UUID   `json:"employee_id"`
	AllowPermissionIDs []uuid.UUID `json:"allow_permission_ids"`
	DenyPermissionIDs  []uuid.UUID `json:"deny_permission_ids"`
}

type AddPermissionsToRoleRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required"`
}

type AddPermissionsToRoleResponse struct {
	RoleID        uuid.UUID   `json:"role_id"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

type CreateRoleRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type CreateRoleResponse struct {
	RoleID      uuid.UUID `json:"role_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}
