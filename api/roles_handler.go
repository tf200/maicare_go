package api

import (
	"errors"
	"fmt"
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
