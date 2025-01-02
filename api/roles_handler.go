package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPermissionsByRoleIDApiResponse represents a response for GetPermissionsByRoleIDApi
type GetPermissionsByRoleIDApiResponse struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Resource string `json:"resource"`
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
		})
	}

	ctx.JSON(http.StatusOK, SuccessResponse(response, "Permissions retrieved successfully"))

}
