package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRolesRoutes(baseRouter *gin.RouterGroup) {
	rolesgroup := baseRouter.Group("")
	rolesgroup.Use(server.AuthMiddleware())
	{
		rolesgroup.GET("/roles", server.RBACMiddleware("ROLES.VIEW"), server.ListRolesApi)
		rolesgroup.POST("/roles", server.RBACMiddleware("ROLES.CREATE"), server.CreateRoleApi)
		rolesgroup.GET("/roles/:role_id/permissions", server.RBACMiddleware("PERMISSIONS.VIEW"), server.ListAllRolePermissionsApi)
		rolesgroup.POST("/roles/:role_id/permissions", server.RBACMiddleware("PERMISSIONS.CREATE"), server.AddPermissionsToRoleApi)
		rolesgroup.GET("/permissions", server.RBACMiddleware("PERMISSIONS.VIEW"), server.ListAllPermissionsApi)

		rolesgroup.POST("/employees/:id/roles", server.RBACMiddleware("ROLES.ASSIGN"), server.AssignRoleToEmployeeApi)
		rolesgroup.GET("/employees/:id/roles_permissions", server.RBACMiddleware("PERMISSIONS.VIEW"), server.ListUserRolesAndPermissionsApi)

		rolesgroup.POST("/employees/:id/permissions", server.RBACMiddleware("PERMISSIONS.GRANT"), server.GrantUserPermissionsApi)

	}

}
