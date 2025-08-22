package api

import (
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// ListRolesApiResponse represents a response for ListRolesApi
type ListRolesApiResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// @Summary List roles
// @Description List all roles
// @Tags roles
// @Produce json
// @Success 200 {object} Response[[]ListRolesApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles [get]
func (server *Server) ListRolesApi(ctx *gin.Context) {
	roles, err := server.store.ListRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := make([]ListRolesApiResponse, len(roles))
	for i, role := range roles {
		response[i] = ListRolesApiResponse{
			ID:   role.ID,
			Name: role.Name,
		}
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Roles retrieved successfully"))
}

// ListAllPermissionsApiResponse represents a response for ListAllPermissionsApi
type ListAllPermissionsApiResponse struct {
	PermissionID       int32  `json:"permission_id"`
	PermissionName     string `json:"permission_name"`
	PermissionResource string `json:"permission_resource"`
}

// @Summary List all permissions
// @Description List all permissions
// @Tags roles
// @Produce json
// @Success 200 {object} Response[[]ListAllPermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /permissions [get]
func (server *Server) ListAllPermissionsApi(ctx *gin.Context) {
	permissions, err := server.store.ListAllPermissions(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListAllPermissionsApi", "Failed to list all permissions", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list permissions")))
		return
	}

	var response []ListAllPermissionsApiResponse
	for _, perm := range permissions {
		response = append(response, ListAllPermissionsApiResponse{
			PermissionID:       perm.ID,
			PermissionName:     perm.Name,
			PermissionResource: perm.Resource,
		})
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Permissions retrieved successfully"))
}

type ListAllRolePermissionsApiResponse struct {
	RoleID             int32  `json:"role_id"`
	PermissionID       int32  `json:"permission_id"`
	PermissionName     string `json:"permission_name"`
	PermissionResource string `json:"permission_resource"`
}

func (server *Server) ListAllRolePermissionsApi(ctx *gin.Context) {
	roleID, err := strconv.ParseInt(ctx.Param("role_id"), 10, 32)
	if err != nil {
		server.logBusinessEvent(LogLevelWarn, "ListAllRolePermissionsApi", "Invalid role_id parameter", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid role_id parameter")))
		return
	}

	permissions, err := server.store.ListAllRolePermissions(ctx, int32(roleID))
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListAllRolePermissionsApi", "Failed to list role permissions", zap.Error(err), zap.Int32("role_id", int32(roleID)))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list role permissions")))
		return
	}
	var response []ListAllRolePermissionsApiResponse
	for _, perm := range permissions {
		response = append(response, ListAllRolePermissionsApiResponse{
			RoleID:             int32(roleID),
			PermissionID:       perm.PermissionID,
			PermissionName:     perm.PermissionName,
			PermissionResource: perm.PermissionResource,
		})
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role permissions retrieved successfully"))
}

// AssignRoleToUserParams represents a request for AssignRoleToUserApi
type AssignRoleToUserParams struct {
	EmployeeID int64 `json:"employee_id"`
	RoleID     int32 `json:"role_id"`
}

// AssignRoleToUserApiResponse represents a response for AssignRoleToUserApi
type AssignRoleToUserApiResponse struct {
	EmployeeID int64 `json:"employee_id"`
	RoleID     int32 `json:"role_id"`
}

// @Summary Assign role to user
// @Description Assign a role to a user
// @Tags roles
// @Accept json
// @Produce json
// @Param input body AssignRoleToUserParams true "Assign role to user"
// @Success 200 {object} Response[AssignRoleToUserApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/assign [put]
func (server *Server) AssignRoleToEmployeeApi(ctx *gin.Context) {
	var req AssignRoleToUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelWarn, "AssignRoleToUserApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	userID, err := server.store.GetUserIDByEmployeeID(ctx, req.EmployeeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			server.logBusinessEvent(LogLevelWarn, "AssignRoleToUserApi", "Employee not found", zap.Int64("employee_id", req.EmployeeID))
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("employee with ID %d not found", req.EmployeeID)))
			return
		}
		server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to get user ID by employee ID", zap.Error(err), zap.Int64("employee_id", req.EmployeeID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get user ID by employee ID")))
		return
	}

	err = server.store.GrantRoleToUser(ctx, db.GrantRoleToUserParams{
		UserID: userID,
		RoleID: req.RoleID,
	})

	if err != nil {
		server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to assign role to user", zap.Error(err), zap.Int64("employee_id", req.EmployeeID), zap.Int32("role_id", req.RoleID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to assign role to user")))
		return
	}

	response := AssignRoleToUserApiResponse{
		EmployeeID: req.EmployeeID,
		RoleID:     req.RoleID,
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role assigned to user successfully"))
}

// ListUserRolesAndPermissionsApiResponse represents a response for ListUserRolesAndPermissionsApi
type ListUserRolesAndPermissionsApiResponse struct {
	Roles []struct {
		RoleID   int32  `json:"id"`
		RoleName string `json:"name"`
	} `json:"roles"`
	Permissions []struct {
		PermissionID       int32  `json:"id"`
		PermissionName     string `json:"name"`
		PermissionResource string `json:"resource"`
	} `json:"permissions"`
}

// @Summary List user roles and permissions
// @Description List roles and permissions for a user by employee ID
// @Tags roles
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Success 200 {object} Response[ListUserRolesAndPermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/roles_and_permissions [get]
func (server *Server) ListUserRolesAndPermissionsApi(ctx *gin.Context) {
	employeeID, err := strconv.ParseInt(ctx.Param("employee_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelWarn, "ListUserRolesAndPermissionsApi", "Invalid employee_id parameter", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}

	userID, err := server.store.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			server.logBusinessEvent(LogLevelWarn, "ListUserRolesAndPermissionsApi", "Employee not found", zap.Int64("employee_id", employeeID))
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("employee with ID %d not found", employeeID)))
			return
		}
		server.logBusinessEvent(LogLevelError, "ListUserRolesAndPermissionsApi", "Failed to get user ID by employee ID", zap.Error(err), zap.Int64("employee_id", employeeID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get user ID by employee ID")))
		return
	}

	roles, err := server.store.ListUserRoles(ctx, int64(userID))
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListUserRolesAndPermissionsApi", "Failed to list user roles", zap.Error(err), zap.Int64("user_id", userID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list user roles")))
		return
	}
	permissions, err := server.store.ListUserPermissions(ctx, int64(userID))
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListUserRolesAndPermissionsApi", "Failed to list user permissions", zap.Error(err), zap.Int64("user_id", userID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list user permissions")))
		return
	}
	var response ListUserRolesAndPermissionsApiResponse
	for _, role := range roles {
		response.Roles = append(response.Roles, struct {
			RoleID   int32  `json:"id"`
			RoleName string `json:"name"`
		}{
			RoleID:   role.ID,
			RoleName: role.Name,
		})
	}
	for _, perm := range permissions {
		response.Permissions = append(response.Permissions, struct {
			PermissionID       int32  `json:"id"`
			PermissionName     string `json:"name"`
			PermissionResource string `json:"resource"`
		}{
			PermissionID:       perm.PermissionID,
			PermissionName:     perm.PermissionName,
			PermissionResource: perm.PermissionResource,
		})
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "User roles and permissions retrieved successfully"))
}
