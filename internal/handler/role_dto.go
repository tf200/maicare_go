package handler

import (
	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

// Request DTOs

type createRoleRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type assignRoleToEmployeeRequest struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

type replaceUserPermissionOverridesRequest struct {
	AllowPermissionIDs []uuid.UUID `json:"allow_permission_ids"`
	DenyPermissionIDs  []uuid.UUID `json:"deny_permission_ids"`
}

type addPermissionsToRoleRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required"`
}

// Response DTOs

type roleResponse struct {
	ID              uuid.UUID `json:"id"`
	RoleName        string    `json:"role_name"`
	Description     *string   `json:"description"`
	PermissionCount int64     `json:"permission_count"`
	EmployeeCount   int64     `json:"employee_count"`
}

type systemPermissionResponse struct {
	PermissionID       uuid.UUID `json:"permission_id"`
	PermissionName     string    `json:"permission_name"`
	PermissionResource string    `json:"permission_resource"`
	DisplayName        string    `json:"display_name"`
	Description        *string   `json:"description"`
	SortOrder          int32     `json:"sort_order"`
}

type permissionSectionResponse struct {
	SectionKey   string                     `json:"section_key"`
	SectionLabel string                     `json:"section_label"`
	Permissions  []systemPermissionResponse `json:"permissions"`
}

type permissionGroupResponse struct {
	GroupKey   string                      `json:"group_key"`
	GroupLabel string                      `json:"group_label"`
	Sections   []permissionSectionResponse `json:"sections"`
}

type rolePermissionResponse struct {
	RoleID             uuid.UUID `json:"role_id"`
	PermissionID       uuid.UUID `json:"permission_id"`
	PermissionName     string    `json:"permission_name"`
	PermissionResource string    `json:"permission_resource"`
}

type assignRoleToEmployeeResponse struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	RoleID     uuid.UUID `json:"role_id"`
}

type roleInfoResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type permissionInfoResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Resource string    `json:"resource"`
}

type permissionOverrideInfoResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Resource string    `json:"resource"`
}

type userRolesAndPermissionsResponse struct {
	Role                 *roleInfoResponse              `json:"role"`
	InheritedPermissions []permissionInfoResponse       `json:"inherited_permissions"`
	OverrideAllows       []permissionOverrideInfoResponse `json:"override_allows"`
	OverrideDenies       []permissionOverrideInfoResponse `json:"override_denies"`
	EffectivePermissions []permissionInfoResponse       `json:"effective_permissions"`
}

type replaceUserPermissionOverridesResponse struct {
	EmployeeID         uuid.UUID   `json:"employee_id"`
	AllowPermissionIDs []uuid.UUID `json:"allow_permission_ids"`
	DenyPermissionIDs  []uuid.UUID `json:"deny_permission_ids"`
}

type addPermissionsToRoleResponse struct {
	RoleID        uuid.UUID   `json:"role_id"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

type createRoleResponse struct {
	RoleID      uuid.UUID `json:"role_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

// Mappers

func toRoleResponse(role domain.Role) roleResponse {
	return roleResponse{
		ID:              role.ID,
		RoleName:        role.Name,
		Description:     role.Description,
		PermissionCount: role.PermissionCount,
		EmployeeCount:   role.EmployeeCount,
	}
}

func toSystemPermissionResponse(perm domain.SystemPermission) systemPermissionResponse {
	return systemPermissionResponse{
		PermissionID:       perm.ID,
		PermissionName:     perm.Name,
		PermissionResource: perm.Resource,
		DisplayName:        perm.DisplayName,
		Description:        perm.Description,
		SortOrder:          perm.SortOrder,
	}
}

func toPermissionGroupResponse(group domain.PermissionGroup) permissionGroupResponse {
	sections := make([]permissionSectionResponse, len(group.Sections))
	for i, section := range group.Sections {
		perms := make([]systemPermissionResponse, len(section.Permissions))
		for j, perm := range section.Permissions {
			perms[j] = toSystemPermissionResponse(perm)
		}
		sections[i] = permissionSectionResponse{
			SectionKey:   section.SectionKey,
			SectionLabel: section.SectionLabel,
			Permissions:  perms,
		}
	}
	return permissionGroupResponse{
		GroupKey:   group.GroupKey,
		GroupLabel: group.GroupLabel,
		Sections:   sections,
	}
}

func toRolePermissionResponse(rp domain.RolePermission) rolePermissionResponse {
	return rolePermissionResponse{
		RoleID:             rp.RoleID,
		PermissionID:       rp.PermissionID,
		PermissionName:     rp.PermissionName,
		PermissionResource: rp.Resource,
	}
}

func toUserRolesAndPermissionsResponse(data *domain.UserRolesAndPermissions) userRolesAndPermissionsResponse {
	var role *roleInfoResponse
	if data.Role != nil {
		role = &roleInfoResponse{
			ID:   data.Role.ID,
			Name: data.Role.Name,
		}
	}

	inherited := make([]permissionInfoResponse, len(data.InheritedPermissions))
	for i, perm := range data.InheritedPermissions {
		inherited[i] = permissionInfoResponse{
			ID:       perm.ID,
			Name:     perm.Name,
			Resource: perm.Resource,
		}
	}

	allowOverrides := make([]permissionOverrideInfoResponse, len(data.OverrideAllows))
	for i, perm := range data.OverrideAllows {
		allowOverrides[i] = permissionOverrideInfoResponse{
			ID:       perm.ID,
			Name:     perm.Name,
			Resource: perm.Resource,
		}
	}

	denyOverrides := make([]permissionOverrideInfoResponse, len(data.OverrideDenies))
	for i, perm := range data.OverrideDenies {
		denyOverrides[i] = permissionOverrideInfoResponse{
			ID:       perm.ID,
			Name:     perm.Name,
			Resource: perm.Resource,
		}
	}

	effective := make([]permissionInfoResponse, len(data.EffectivePermissions))
	for i, perm := range data.EffectivePermissions {
		effective[i] = permissionInfoResponse{
			ID:       perm.ID,
			Name:     perm.Name,
			Resource: perm.Resource,
		}
	}

	return userRolesAndPermissionsResponse{
		Role:                 role,
		InheritedPermissions: inherited,
		OverrideAllows:       allowOverrides,
		OverrideDenies:       denyOverrides,
		EffectivePermissions: effective,
	}
}

func toCreateRoleResponse(role *domain.Role) createRoleResponse {
	return createRoleResponse{
		RoleID:      role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

func toCreateRoleParams(req createRoleRequest) domain.CreateRoleParams {
	return domain.CreateRoleParams{
		Name:        req.Name,
		Description: req.Description,
	}
}

func toAssignRoleToEmployeeParams(req assignRoleToEmployeeRequest) domain.AssignRoleToEmployeeParams {
	return domain.AssignRoleToEmployeeParams{
		RoleID: req.RoleID,
	}
}

func toReplaceUserPermissionOverridesParams(employeeID uuid.UUID, req replaceUserPermissionOverridesRequest) domain.ReplaceUserPermissionOverridesParams {
	return domain.ReplaceUserPermissionOverridesParams{
		EmployeeID:         employeeID,
		AllowPermissionIDs: req.AllowPermissionIDs,
		DenyPermissionIDs:  req.DenyPermissionIDs,
	}
}

func toAddPermissionsToRoleParams(roleID uuid.UUID, req addPermissionsToRoleRequest) domain.AddPermissionsToRoleParams {
	return domain.AddPermissionsToRoleParams{
		RoleID:        roleID,
		PermissionIDs: req.PermissionIDs,
	}
}
