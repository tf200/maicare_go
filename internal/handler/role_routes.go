package handler

import "github.com/gin-gonic/gin"

func RegisterRoleRoutes(
	rg *gin.RouterGroup,
	handler *RoleHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	rg.GET("/roles", auth, requirePermission("ROLES.VIEW"), handler.ListRoles)
	rg.POST("/roles", auth, requirePermission("ROLES.CREATE"), handler.CreateRole)
	rg.GET("/roles/:role_id/permissions", auth, requirePermission("PERMISSIONS.VIEW"), handler.ListAllRolePermissions)
	rg.POST("/roles/:role_id/permissions", auth, requirePermission("PERMISSIONS.CREATE"), handler.AddPermissionsToRole)
	rg.GET("/permissions", auth, requirePermission("PERMISSIONS.VIEW"), handler.ListAllPermissions)

	rg.POST("/employees/:id/roles", auth, requirePermission("ROLES.ASSIGN"), handler.AssignRoleToEmployee)
	rg.GET("/employees/:id/roles_permissions", auth, requirePermission("PERMISSIONS.VIEW"), handler.ListUserRolesAndPermissions)
	rg.POST("/employees/:id/permissions", auth, requirePermission("PERMISSIONS.GRANT"), handler.ReplaceUserPermissionOverrides)
}
