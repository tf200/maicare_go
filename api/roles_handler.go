package api

import (
	"fmt"
	"net/http"

	"maicare_go/service/settings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// @Summary List roles
// @Description List all roles
// @Tags roles
// @Produce json
// @Success 200 {object} Response[[]settings.ListRolesApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles [get]
func (server *Server) ListRolesApi(ctx *gin.Context) {
	response, err := server.businessService.SettingsService.ListRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list roles")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Roles retrieved successfully"))
}

// @Summary List all permissions
// @Description List all permissions
// @Tags roles
// @Produce json
// @Success 200 {object} Response[[]settings.PermissionGroupResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /permissions [get]
func (server *Server) ListAllPermissionsApi(ctx *gin.Context) {
	response, err := server.businessService.SettingsService.ListAllPermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list permissions")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Permissions retrieved successfully"))
}

// @Summary List all permissions for a role
// @Description List all permissions associated with a specific role
// @Tags roles
// @Produce json
// @Param role_id path uuid true "Role ID"
// @Success 200 {object} Response[[]settings.ListAllRolePermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/{role_id}/permissions [get]
func (server *Server) ListAllRolePermissionsApi(ctx *gin.Context) {
	roleID, err := uuid.Parse(ctx.Param("role_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid role_id parameter")))
		return
	}

	response, err := server.businessService.SettingsService.ListAllRolePermissions(ctx, roleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list role permissions")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role permissions retrieved successfully"))
}

// @Summary Assign role to user
// @Description Assign a role to a user
// @Tags roles
// @Param employee_id query int true "Employee ID"
// @Accept json
// @Produce json
// @Param input body settings.AssignRoleToEmployeeParams true "Assign role to user"
// @Success 200 {object} Response[settings.AssignRoleToEmployeeApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/roles [post]
func (server *Server) AssignRoleToEmployeeApi(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}
	var req settings.AssignRoleToEmployeeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	response, err := server.businessService.SettingsService.AssignRoleToEmployee(ctx, employeeID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to assign role to user")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role assigned to user successfully"))
}

// @Summary List user roles and permissions
// @Description List roles and permissions for a user by employee ID
// @Tags roles
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Success 200 {object} Response[settings.ListUserRolesAndPermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/roles_permissions [get]
func (server *Server) ListUserRolesAndPermissionsApi(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}
	response, err := server.businessService.SettingsService.ListUserRolesAndPermissionsApi(ctx, employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list user roles and permissions")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "User roles and permissions retrieved successfully"))
}

// @Summary Replace user permission overrides
// @Description Replace explicit allow and deny permission overrides for a user by employee ID
// @Tags roles
// @Accept json
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Param input body settings.ReplaceUserPermissionOverridesRequest true "Replace user permission overrides"
// @Success 200 {object} Response[settings.ReplaceUserPermissionOverridesResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/permissions [post]
func (server *Server) GrantUserPermissionsApi(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}

	var req settings.ReplaceUserPermissionOverridesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	response, err := server.businessService.SettingsService.ReplaceUserPermissionOverrides(ctx, employeeID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to replace user permission overrides")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "User permission overrides replaced successfully"))
}

// @Summary Create a new role
// @Description Create a new role with the specified name
// @Tags roles
// @Accept json
// @Produce json
// @Param input body settings.CreateRoleRequest true "Create role"
// @Success 200 {object} Response[settings.CreateRoleResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles [post]
func (server *Server) CreateRoleApi(ctx *gin.Context) {
	var req settings.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	response, err := server.businessService.SettingsService.CreateRole(ctx, &req)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateRoleApi", "Failed to create role", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create role")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role created successfully"))
}

// @Summary Replace role permissions
// @Description Replace the permission set assigned to a role by role ID
// @Tags roles
// @Accept json
// @Produce json
// @Param role_id path uuid true "Role ID"
// @Param input body settings.AddPermissionsToRoleRequest true "Add permissions to role"
// @Success 200 {object} Response[settings.AddPermissionsToRoleResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/{role_id}/permissions [post]
func (server *Server) AddPermissionsToRoleApi(ctx *gin.Context) {
	roleID, err := uuid.Parse(ctx.Param("role_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid role_id parameter")))
		return
	}
	var req settings.AddPermissionsToRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	response, err := server.businessService.SettingsService.AddPermissionsToRole(ctx, roleID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to replace permissions for role")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role permissions replaced successfully"))
}
