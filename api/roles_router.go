package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRolesRoutes(baseRouter *gin.RouterGroup) {
	rolesgroup := baseRouter.Group("/roles")
	rolesgroup.Use(server.AuthMiddleware())
	{
		rolesgroup.GET("", server.ListRolesApi)
		rolesgroup.GET("/:role_id/permissions", server.GetPermissionsByRoleIDApi)
		rolesgroup.GET("/user", server.GetRoleByIDApi)

		rolesgroup.PUT("/assign", server.RBACMiddleware("ROLES.UPDATE"), server.AssignRoleToEmployeeApi)

	}

}
