package api

import (
	"database/sql"
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
	ID              int32  `json:"id"`
	RoleName        string `json:"role_name"`
	PermissionCount int64  `json:"permission_count"`
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
			ID:              role.ID,
			RoleName:        role.Name,
			PermissionCount: role.PermissionCount,
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

// ListAllRolePermissionsApiResponse represents a response for ListAllRolePermissionsApi
type ListAllRolePermissionsApiResponse struct {
	RoleID             int32  `json:"role_id"`
	PermissionID       int32  `json:"permission_id"`
	PermissionName     string `json:"permission_name"`
	PermissionResource string `json:"permission_resource"`
}

// @Summary List all permissions for a role
// @Description List all permissions associated with a specific role
// @Tags roles
// @Produce json
// @Param role_id path int true "Role ID"
// @Success 200 {object} Response[[]ListAllRolePermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/{role_id}/permissions [get]
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
			PermissionResource: perm.Resource,
		})
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role permissions retrieved successfully"))
}

// AssignRoleToUserParams represents a request for AssignRoleToUserApi
type AssignRoleToUserParams struct {
	RoleID int32 `json:"role_id"`
}

// AssignRoleToUserApiResponse represents a response for AssignRoleToUserApi
type AssignRoleToUserApiResponse struct {
	EmployeeID int64 `json:"employee_id"`
	RoleID     int32 `json:"role_id"`
}

// @Summary Assign role to user
// @Description Assign a role to a user
// @Tags roles
// @Param employee_id query int true "Employee ID"
// @Accept json
// @Produce json
// @Param input body AssignRoleToUserParams true "Assign role to user"
// @Success 200 {object} Response[AssignRoleToUserApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/roles [post]
func (server *Server) AssignRoleToEmployeeApi(ctx *gin.Context) {
	employeeID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelWarn, "AssignRoleToUserApi", "Invalid employee_id parameter", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}
	var req AssignRoleToUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelWarn, "AssignRoleToUserApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	userID, err := server.store.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			server.logBusinessEvent(LogLevelWarn, "AssignRoleToUserApi", "Employee not found", zap.Int64("employee_id", employeeID))
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("employee with ID %d not found", employeeID)))
			return
		}
		server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to get user ID by employee ID", zap.Error(err), zap.Int64("employee_id", employeeID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get user ID by employee ID")))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to begin transaction", zap.Error(err), zap.Int64("employee_id", employeeID), zap.Int32("role_id", req.RoleID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to begin transaction")))
		return
	}

	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to rollback transaction", zap.Error(rollbackErr), zap.Int64("employee_id", employeeID), zap.Int32("role_id", req.RoleID))
		}
	}()

	qtx := server.store.WithTx(tx)

	err = qtx.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: req.RoleID,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to assign role to user", zap.Error(err), zap.Int64("employee_id", employeeID), zap.Int32("role_id", req.RoleID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to assign role to user")))
		return
	}

	err = qtx.DeleteUserPermissions(ctx, userID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to delete user permissions", zap.Error(err), zap.Int64("employee_id", employeeID), zap.Int32("role_id", req.RoleID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete user permissions")))
		return
	}

	err = qtx.GrantRolePermissionsToUser(ctx, db.GrantRolePermissionsToUserParams{
		UserID: userID,
		RoleID: req.RoleID,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to grant role permissions to user", zap.Error(err), zap.Int64("employee_id", employeeID), zap.Int32("role_id", req.RoleID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to grant role permissions to user")))
		return
	}
	err = tx.Commit(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AssignRoleToUserApi", "Failed to commit transaction", zap.Error(err), zap.Int64("employee_id", employeeID), zap.Int32("role_id", req.RoleID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to commit transaction")))
		return
	}

	response := AssignRoleToUserApiResponse{
		EmployeeID: employeeID,
		RoleID:     req.RoleID,
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role assigned to user successfully"))
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

// @Summary List user roles and permissions
// @Description List roles and permissions for a user by employee ID
// @Tags roles
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Success 200 {object} Response[ListUserRolesAndPermissionsApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/roles_permissions [get]
func (server *Server) ListUserRolesAndPermissionsApi(ctx *gin.Context) {
	employeeID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
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

	roles, err := server.store.GetUserRoles(ctx, int64(userID))
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
	response := ListUserRolesAndPermissionsApiResponse{
		Roles: struct {
			RoleID   int32  `json:"id"`
			RoleName string `json:"name"`
		}{
			RoleID:   roles.ID,
			RoleName: roles.Name,
		},
		Permissions: []struct {
			PermissionID       int32  `json:"id"`
			PermissionName     string `json:"name"`
			PermissionResource string `json:"resource"`
		}{},
	}

	for _, perm := range permissions {
		response.Permissions = append(response.Permissions, struct {
			PermissionID       int32  `json:"id"`
			PermissionName     string `json:"name"`
			PermissionResource string `json:"resource"`
		}{
			PermissionID:       perm.PermissionID,
			PermissionName:     perm.PermissionName,
			PermissionResource: perm.Resource,
		})
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "User roles and permissions retrieved successfully"))
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

// @Summary Grant user permissions
// @Description Grant specific permissions to a user by employee ID
// @Tags roles
// @Accept json
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Param input body GrantUserPermissionsRequest true "Grant user permissions"
// @Success 200 {object} Response[GrantUserPermissionsResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /employees/{employee_id}/permissions [post]
func (server *Server) GrantUserPermissionsApi(ctx *gin.Context) {
	employeeID, err := strconv.ParseInt(ctx.Param("employee_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelWarn, "GrantUserPermissionsApi", "Invalid employee_id parameter", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee_id parameter")))
		return
	}

	var req GrantUserPermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelWarn, "GrantUserPermissionsApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	userID, err := server.store.GetUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			server.logBusinessEvent(LogLevelWarn, "GrantUserPermissionsApi", "Employee not found", zap.Int64("employee_id", employeeID))
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("employee with ID %d not found", employeeID)))
			return
		}
		server.logBusinessEvent(LogLevelError, "GrantUserPermissionsApi", "Failed to get user ID by employee ID", zap.Error(err), zap.Int64("employee_id", employeeID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get user ID by employee ID")))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GrantUserPermissionsApi", "Failed to begin transaction", zap.Error(err), zap.Int64("employee_id", employeeID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to begin transaction")))
		return
	}

	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			server.logBusinessEvent(LogLevelError, "GrantUserPermissionsApi", "Failed to rollback transaction", zap.Error(rollbackErr), zap.Int64("employee_id", employeeID))
		}
	}()

	qtx := server.store.WithTx(tx)

	err = qtx.DeleteUserPermissions(ctx, userID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GrantUserPermissionsApi", "Failed to delete user permissions", zap.Error(err), zap.Int64("employee_id", employeeID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete user permissions")))
		return
	}

	err = qtx.GrantUserPermissions(ctx, db.GrantUserPermissionsParams{
		UserID:        userID,
		PermissionIds: req.PermissionIDs,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GrantUserPermissionsApi", "Failed to grant user permissions", zap.Error(err), zap.Int64("employee_id", employeeID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to grant user permissions")))
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GrantUserPermissionsApi", "Failed to commit transaction", zap.Error(err), zap.Int64("employee_id", employeeID))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to commit transaction")))
		return
	}
	response := GrantUserPermissionsResponse{
		EmployeeID:    employeeID,
		PermissionIDs: req.PermissionIDs,
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "User permissions granted successfully"))
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

// @Summary Create a new role
// @Description Create a new role with the specified name
// @Tags roles
// @Accept json
// @Produce json
// @Param input body CreateRoleRequest true "Create role"
// @Success 200 {object} Response[CreateRoleResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles [post]
func (server *Server) CreateRoleApi(ctx *gin.Context) {
	var req CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelWarn, "CreateRoleApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	role, err := server.store.CreateRole(ctx, req.Name)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateRoleApi", "Failed to create role", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create role")))
		return
	}

	response := CreateRoleResponse{
		RoleID: role.ID,
		Name:   role.Name,
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role created successfully"))
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

// @Summary Add permissions to a role
// @Description Add specific permissions to a role by role ID
// @Tags roles
// @Accept json
// @Produce json
// @Param role_id path int true "Role ID"
// @Param input body AddPermissionsToRoleRequest true "Add permissions to role"
// @Success 200 {object} Response[AddPermissionsToRoleResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/{role_id}/permissions [post]
func (server *Server) AddPermissionsToRoleApi(ctx *gin.Context) {
	roleID, err := strconv.ParseInt(ctx.Param("role_id"), 10, 32)
	if err != nil {
		server.logBusinessEvent(LogLevelWarn, "AddPermissionsToRoleApi", "Invalid role_id parameter", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid role_id parameter")))
		return
	}
	var req AddPermissionsToRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelWarn, "AddPermissionsToRoleApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AddPermissionsToRoleApi", "Failed to begin transaction", zap.Error(err), zap.Int32("role_id", int32(roleID)))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to begin transaction")))
		return
	}

	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			server.logBusinessEvent(LogLevelError, "AddPermissionsToRoleApi", "Failed to rollback transaction", zap.Error(rollbackErr), zap.Int32("role_id", int32(roleID)))
		}
	}()

	qtx := server.store.WithTx(tx)
	err = qtx.RemovePermissionsFromRole(ctx, int32(roleID))
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AddPermissionsToRoleApi", "Failed to remove existing permissions from role", zap.Error(err), zap.Int32("role_id", int32(roleID)))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to remove existing permissions from role")))
		return
	}

	err = qtx.AddPermissionsToRole(ctx, db.AddPermissionsToRoleParams{
		RoleID:        int32(roleID),
		PermissionIds: req.PermissionIDs,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AddPermissionsToRoleApi", "Failed to add permissions to role", zap.Error(err), zap.Int32("role_id", int32(roleID)))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to add permissions to role")))
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "AddPermissionsToRoleApi", "Failed to commit transaction", zap.Error(err), zap.Int32("role_id", int32(roleID)))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to commit transaction")))
		return
	}

	response := AddPermissionsToRoleResponse{
		RoleID:        int32(roleID),
		PermissionIDs: req.PermissionIDs,
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response, "Permissions added to role successfully"))
}
