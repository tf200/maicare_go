package auth

// ListRolesApiResponse represents a response for ListRolesApi
type ListRolesApiResponse struct {
	ID              int32  `json:"id"`
	RoleName        string `json:"role_name"`
	PermissionCount int64  `json:"permission_count"`
}

// ListAllPermissionsApiResponse represents a response for ListAllPermissionsApi
type ListAllPermissionsApiResponse struct {
	PermissionID       int32  `json:"permission_id"`
	PermissionName     string `json:"permission_name"`
	PermissionResource string `json:"permission_resource"`
}

// ListAllRolePermissionsApiResponse represents a response for ListAllRolePermissionsApi
type ListAllRolePermissionsApiResponse struct {
	RoleID             int32  `json:"role_id"`
	PermissionID       int32  `json:"permission_id"`
	PermissionName     string `json:"permission_name"`
	PermissionResource string `json:"permission_resource"`
}

// AssignRoleToUserParams represents a request for AssignRoleToUserApi
type AssignRoleToEmployeeParams struct {
	RoleID int32 `json:"role_id"`
}

// AssignRoleToUserApiResponse represents a response for AssignRoleToUserApi
type AssignRoleToEmployeeApiResponse struct {
	EmployeeID int64 `json:"employee_id"`
	RoleID     int32 `json:"role_id"`
}

// ListUserRolesAndPermissionsApiResponse represents a response for ListUserRolesAndPermissionsApi
type ListUserRolesAndPermissionsApiResponse struct {
	Roles struct {
		RoleID   int32  `json:"id"`
		RoleName string `json:"name"`
	} `json:"roles"`
	Permissions []struct {
		PermissionID       int32  `json:"id"`
		PermissionName     string `json:"name"`
		PermissionResource string `json:"resource"`
	} `json:"permissions"`
}

// GrantUserPermissionsRequest represents a request for GrantUserPermissionsApi
type GrantUserPermissionsRequest struct {
	PermissionIDs []int32 `json:"permission_ids"`
}

// GrantUserPermissionsResponse represents a response for GrantUserPermissionsApi
type GrantUserPermissionsResponse struct {
	EmployeeID    int64   `json:"employee_id"`
	PermissionIDs []int32 `json:"permission_ids"`
}

// AddPermissionsToRoleRequest represents a request for AddPermissionsToRoleApi
type AddPermissionsToRoleRequest struct {
	PermissionIDs []int32 `json:"permission_ids" binding:"required"`
}

// AddPermissionsToRoleResponse represents a response for AddPermissionsToRoleApi
type AddPermissionsToRoleResponse struct {
	RoleID        int32   `json:"role_id"`
	PermissionIDs []int32 `json:"permission_ids"`
}

// CreateRoleRequest represents a request for CreateRoleApi
type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateRoleResponse represents a response for CreateRoleApi
type CreateRoleResponse struct {
	RoleID int32  `json:"role_id"`
	Name   string `json:"name"`
}
