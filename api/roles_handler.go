package api

import (
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetPermissionsByRoleIDApiResponse represents a response for GetPermissionsByRoleIDApi
type GetPermissionsByRoleIDApiResponse struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Resource string `json:"resource"`
	Method   string `json:"method"`
}

// @Summary Get permissions by role ID
// @Description Get all permissions associated with a role ID
// @Tags roles
// @Param role_id path int true "Role ID"
// @Produce json
// @Success 200 {object} Response[[]GetPermissionsByRoleIDApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/{role_id}/permissions [get]
func (server *Server) GetPermissionsByRoleIDApi(ctx *gin.Context) {
	id := ctx.Param("role_id")
	roleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	permissions, err := server.store.GetPermissionsByRoleID(ctx, int32(roleID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var response []GetPermissionsByRoleIDApiResponse
	for _, p := range permissions {
		response = append(response, GetPermissionsByRoleIDApiResponse{
			ID:       p.ID,
			Name:     p.Name,
			Resource: p.Resource,
			Method:   p.Method,
		})
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Permissions retrieved successfully"))

}

// GetRoleByIDApiResponse represents a response for GetRoleByIDApi
type GetRoleByIDApiResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// @Summary Get role by ID
// @Description Get a role by ID
// @Tags roles
// @Produce json
// @Success 200 {object} Response[GetRoleByIDApiResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /roles/user [get]
func (server *Server) GetRoleByIDApi(ctx *gin.Context) {
	// Get the auth payload from context
	value, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("authorization payload not found")))
		return
	}

	// Type assert to get the payload
	payload, ok := value.(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("invalid authorization payload type")))
		return
	}

	// Get role ID from payload
	roleID := payload.RoleID

	role, err := server.store.GetRoleByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("role with ID %d not found", roleID)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := GetRoleByIDApiResponse{
		ID:   role.ID,
		Name: role.Name,
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role retrieved successfully"))
}

// AssignRoleToUserParams represents a request for AssignRoleToUserApi
type AssignRoleToUserParams struct {
	ID     int64 `json:"id"`
	RoleID int32 `json:"role_id"`
}

// AssignRoleToUserApiResponse represents a response for AssignRoleToUserApi
type AssignRoleToUserApiResponse struct {
	ID     int64 `json:"id"`
	RoleID int32 `json:"role_id"`
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.AssignRoleToUserParams{
		ID:     req.ID,
		RoleID: req.RoleID,
	}

	updatedUser, err := server.store.AssignRoleToUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := AssignRoleToUserApiResponse{
		ID:     updatedUser.ID,
		RoleID: updatedUser.RoleID,
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Role assigned to user successfully"))
}

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
	response := make([]ListRolesApiResponse, 0)
	for i, role := range roles {
		response[i] = ListRolesApiResponse{
			ID:   role.ID,
			Name: role.Name,
		}
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Roles retrieved successfully"))
}
