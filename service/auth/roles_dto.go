package auth

import "github.com/google/uuid"

// ListRolesApiResponse represents a response for ListRolesApi
type ListRolesApiResponse struct {
	ID              uuid.UUID `json:"id"`
	RoleName        string    `json:"role_name"`
	PermissionCount int64     `json:"permission_count"`
}

// ListAllPermissionsApiResponse represents a response for ListAllPermissionsApi
type ListAllPermissionsApiResponse struct {
	PermissionID       uuid.UUID `json:"permission_id"`
	PermissionName     string    `json:"permission_name"`
	PermissionResource string    `json:"permission_resource"`
}

// ListAllRolePermissionsApiResponse represents a response for ListAllRolePermissionsApi
type ListAllRolePermissionsApiResponse struct {
	RoleID             uuid.UUID `json:"role_id"`
	PermissionID       uuid.UUID `json:"permission_id"`
	PermissionName     string    `json:"permission_name"`
	PermissionResource string    `json:"permission_resource"`
}

// AssignRoleToUserParams represents a request for AssignRoleToUserApi
type AssignRoleToEmployeeParams struct {
	RoleID uuid.UUID `json:"role_id"`
}

// AssignRoleToUserApiResponse represents a response for AssignRoleToUserApi
type AssignRoleToEmployeeApiResponse struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	RoleID     uuid.UUID `json:"role_id"`
}

// ListUserRolesAndPermissionsApiResponse represents a response for ListUserRolesAndPermissionsApi
type ListUserRolesAndPermissionsApiResponse struct {
	Roles struct {
		RoleID   uuid.UUID `json:"id"`
		RoleName string    `json:"name"`
	} `json:"roles"`
	Permissions []struct {
		PermissionID       uuid.UUID `json:"id"`
		PermissionName     string    `json:"name"`
		PermissionResource string    `json:"resource"`
	} `json:"permissions"`
}

// GrantUserPermissionsRequest represents a request for GrantUserPermissionsApi
type GrantUserPermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

// GrantUserPermissionsResponse represents a response for GrantUserPermissionsApi
type GrantUserPermissionsResponse struct {
	EmployeeID    uuid.UUID   `json:"employee_id"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

// AddPermissionsToRoleRequest represents a request for AddPermissionsToRoleApi
type AddPermissionsToRoleRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required"`
}

// AddPermissionsToRoleResponse represents a response for AddPermissionsToRoleApi
type AddPermissionsToRoleResponse struct {
	RoleID        uuid.UUID   `json:"role_id"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

// CreateRoleRequest represents a request for CreateRoleApi
type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateRoleResponse represents a response for CreateRoleApi
type CreateRoleResponse struct {
	RoleID uuid.UUID `json:"role_id"`
	Name   string    `json:"name"`
}
