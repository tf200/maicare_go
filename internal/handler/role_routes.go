package handler

import (
	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(
	rg *gin.RouterGroup,
	handler *RoleHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	rg.GET("/roles", auth, requirePermission(domain.PermRolesView.String()), handler.ListRoles)
	rg.POST("/roles", auth, requirePermission(domain.PermRolesCreate.String()), handler.CreateRole)
	rg.GET("/roles/:role_id/permissions", auth, requirePermission(domain.PermPermissionsView.String()), handler.ListAllRolePermissions)
	rg.POST("/roles/:role_id/permissions", auth, requirePermission(domain.PermPermissionsCreate.String()), handler.ReplaceRolePermissions)
	rg.GET("/permissions", auth, requirePermission(domain.PermPermissionsView.String()), handler.ListAllPermissions)

	rg.POST("/employees/:id/roles", auth, requirePermission(domain.PermRolesAssign.String()), handler.AssignRoleToEmployee)
	rg.GET("/employees/:id/roles_permissions", auth, requirePermission(domain.PermPermissionsView.String()), handler.ListUserRolesAndPermissions)
}
