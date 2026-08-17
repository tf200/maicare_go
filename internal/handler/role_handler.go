package handler

import (
	"errors"
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	service domain.RoleService
}

func NewRoleHandler(service domain.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) ListRoles(ctx *gin.Context) {
	roles, err := h.service.ListRoles(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list roles", ""))
		return
	}

	items := make([]roleResponse, len(roles))
	for i, role := range roles {
		items[i] = toRoleResponse(role)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(items, "Roles retrieved successfully"))
}

func (h *RoleHandler) ListAllPermissions(ctx *gin.Context) {
	groups, err := h.service.ListAllPermissions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list permissions", ""))
		return
	}

	items := make([]permissionGroupResponse, len(groups))
	for i, group := range groups {
		items[i] = toPermissionGroupResponse(group)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(items, "Permissions retrieved successfully"))
}

func (h *RoleHandler) ListAllRolePermissions(ctx *gin.Context) {
	roleID, err := uuid.Parse(ctx.Param("role_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid role_id parameter", ""))
		return
	}

	perms, err := h.service.ListAllRolePermissions(ctx.Request.Context(), roleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list role permissions", ""))
		return
	}

	items := make([]rolePermissionResponse, len(perms))
	for i, perm := range perms {
		items[i] = toRolePermissionResponse(perm)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(items, "Role permissions retrieved successfully"))
}

func (h *RoleHandler) AssignRoleToEmployee(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid employee_id parameter", ""))
		return
	}

	var req assignRoleToEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	if err := h.service.AssignRoleToEmployee(ctx.Request.Context(), employeeID, toAssignRoleToEmployeeParams(req)); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to assign role to user", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(assignRoleToEmployeeResponse{
		EmployeeID: employeeID,
		RoleID:     req.RoleID,
	}, "Role assigned to user successfully"))
}

func (h *RoleHandler) ListUserRolesAndPermissions(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid employee_id parameter", ""))
		return
	}

	result, err := h.service.ListUserRolesAndPermissions(ctx.Request.Context(), employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list user roles and permissions", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toUserRolesAndPermissionsResponse(result), "User roles and permissions retrieved successfully"))
}

func (h *RoleHandler) CreateRole(ctx *gin.Context) {
	var req createRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	role, err := h.service.CreateRole(ctx.Request.Context(), toCreateRoleParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create role", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toCreateRoleResponse(role), "Role created successfully"))
}

func (h *RoleHandler) ReplaceRolePermissions(ctx *gin.Context) {
	roleID, err := uuid.Parse(ctx.Param("role_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid role_id parameter", ""))
		return
	}

	var req replaceRolePermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
		return
	}

	if err := h.service.ReplaceRolePermissions(ctx.Request.Context(), toReplaceRolePermissionsParams(roleID, req)); err != nil {
		if errors.Is(err, domain.ErrInvalidRolePermissions) {
			ctx.JSON(http.StatusBadRequest, httpapi.Fail(err.Error(), ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to replace permissions for role", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(replaceRolePermissionsResponse{
		RoleID:      roleID,
		Permissions: req.Permissions,
	}, "Role permissions replaced successfully"))
}
