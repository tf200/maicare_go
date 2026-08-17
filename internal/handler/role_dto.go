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

type permissionGrantRequest struct {
	PermissionID uuid.UUID               `json:"permission_id" binding:"required"`
	Scope        *domain.PermissionScope `json:"scope"`
}

type replaceRolePermissionsRequest struct {
	Permissions []permissionGrantRequest `json:"permissions" binding:"required"`
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
	PermissionID   uuid.UUID `json:"permission_id"`
	PermissionName string    `json:"permission_name"`
	DisplayName    string    `json:"display_name"`
	Description    *string   `json:"description"`
	IsScoped       bool      `json:"is_scoped"`
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
	RoleID         uuid.UUID               `json:"role_id"`
	PermissionID   uuid.UUID               `json:"permission_id"`
	PermissionName string                  `json:"permission_name"`
	IsScoped       bool                    `json:"is_scoped"`
	Scope          *domain.PermissionScope `json:"scope"`
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
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type userRolesAndPermissionsResponse struct {
	Role                 *roleInfoResponse        `json:"role"`
	EffectivePermissions []permissionInfoResponse `json:"effective_permissions"`
}

type replaceRolePermissionsResponse struct {
	RoleID      uuid.UUID                `json:"role_id"`
	Permissions []permissionGrantRequest `json:"permissions"`
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
		PermissionID:   perm.ID,
		PermissionName: perm.Name,
		DisplayName:    perm.DisplayName,
		Description:    perm.Description,
		IsScoped:       perm.IsScoped,
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
		RoleID:         rp.RoleID,
		PermissionID:   rp.PermissionID,
		PermissionName: rp.PermissionName,
		IsScoped:       rp.IsScoped,
		Scope:          rp.Scope,
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

	effective := make([]permissionInfoResponse, len(data.EffectivePermissions))
	for i, perm := range data.EffectivePermissions {
		effective[i] = permissionInfoResponse{
			ID:   perm.ID,
			Name: perm.Name,
		}
	}

	return userRolesAndPermissionsResponse{
		Role:                 role,
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

func toReplaceRolePermissionsParams(roleID uuid.UUID, req replaceRolePermissionsRequest) domain.ReplaceRolePermissionsParams {
	permissions := make([]domain.PermissionGrant, len(req.Permissions))
	for i, permission := range req.Permissions {
		permissions[i] = domain.PermissionGrant{
			PermissionID: permission.PermissionID,
			Scope:        permission.Scope,
		}
	}
	return domain.ReplaceRolePermissionsParams{RoleID: roleID, Permissions: permissions}
}
