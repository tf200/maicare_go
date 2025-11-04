package api

import (
	"fmt"
	"net/http"
	"strconv"

	"maicare_go/service/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// @Summary List roles
// @Description List all roles
// @Tags roles
// @Produce json
// @Success 200 {object} Response[[]auth.ListRolesApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles [get]
func (server *Server) ListRolesApi(ctx *gin.Context) {
	response, err := server.businessService.AuthService.ListRoles(ctx)
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
// @Success 200 {object} Response[[]auth.ListAllPermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /permissions [get]
func (server *Server) ListAllPermissionsApi(ctx *gin.Context) {
	response, err := server.businessService.AuthService.ListAllPermissions(ctx)
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
// @Param role_id path int true "Role ID"
// @Success 200 {object} Response[[]auth.ListAllRolePermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/{role_id}/permissions [get]
func (server *Server) ListAllRolePermissionsApi(ctx *gin.Context) {
	roleID, err := strconv.ParseInt(ctx.Param("role_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid role_id parameter")))
		return
	}

	response, err := server.businessService.AuthService.ListAllRolePermissions(ctx, int32(roleID))
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
// @Param input body auth.AssignRoleToEmployeeParams true "Assign role to user"
// @Success 200 {object} Response[auth.AssignRoleToEmployeeApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/roles [post]
func (server *Server) AssignRoleToEmployeeApi(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}
	var req auth.AssignRoleToEmployeeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	response, err := server.businessService.AuthService.AssignRoleToEmployee(ctx, employeeID, &req)
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
// @Success 200 {object} Response[auth.ListUserRolesAndPermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/roles_permissions [get]
func (server *Server) ListUserRolesAndPermissionsApi(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}
	response, err := server.businessService.AuthService.ListUserRolesAndPermissionsApi(ctx, employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list user roles and permissions")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "User roles and permissions retrieved successfully"))
}

// @Summary Grant user permissions
// @Description Grant specific permissions to a user by employee ID
// @Tags roles
// @Accept json
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Param input body auth.GrantUserPermissionsRequest true "Grant user permissions"
// @Success 200 {object} Response[auth.GrantUserPermissionsResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/permissions [post]
func (server *Server) GrantUserPermissionsApi(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("employee_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}

	var req auth.GrantUserPermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	response, err := server.businessService.AuthService.GrantUserPermission(ctx, employeeID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to grant user permissions")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "User permissions granted successfully"))
}

// @Summary Create a new role
// @Description Create a new role with the specified name
// @Tags roles
// @Accept json
// @Produce json
// @Param input body auth.CreateRoleRequest true "Create role"
// @Success 200 {object} Response[auth.CreateRoleResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles [post]
func (server *Server) CreateRoleApi(ctx *gin.Context) {
	var req auth.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	response, err := server.businessService.AuthService.CreateRole(ctx, &req)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateRoleApi", "Failed to create role", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create role")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role created successfully"))
}

// @Summary Add permissions to a role
// @Description Add specific permissions to a role by role ID
// @Tags roles
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param input body auth.AddPermissionsToRoleRequest true "Add permissions to role"
// @Success 200 {object} Response[auth.AddPermissionsToRoleResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/{role_id}/permissions [post]
func (server *Server) AddPermissionsToRoleApi(ctx *gin.Context) {
	roleID, err := strconv.ParseInt(ctx.Param("role_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid role_id parameter")))
		return
	}
	var req auth.AddPermissionsToRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	response, err := server.businessService.AuthService.AddPermissionsToRole(ctx, int32(roleID), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to add permissions to role")))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Permissions added to role successfully"))
}
