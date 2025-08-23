package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRolesRoutes(baseRouter *gin.RouterGroup) {
	rolesgroup := baseRouter.Group("")
	rolesgroup.Use(server.AuthMiddleware())
	{
		rolesgroup.GET("/roles", server.ListRolesApi)
		rolesgroup.POST("/roles", server.CreateRoleApi)
		rolesgroup.GET("/roles/:role_id/permissions", server.ListAllRolePermissionsApi)
		rolesgroup.POST("/roles/:role_id/permissions", server.AddPermissionsToRoleApi)
		rolesgroup.GET("/permissions", server.ListAllPermissionsApi)

		rolesgroup.POST("/employees/:id/roles", server.AssignRoleToEmployeeApi)
		rolesgroup.GET("/employees/:id/roles_permissions", server.ListUserRolesAndPermissionsApi)

		rolesgroup.POST("/employees/:id/permissions", server.GrantUserPermissionsApi)

	}

}
