package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRolesRoutes(baseRouter *gin.RouterGroup) {
	rolesgroup := baseRouter.Group("/roles")
	rolesgroup.Use(AuthMiddleware(server.tokenMaker))
	{
		rolesgroup.GET("/:role_id/permissions", server.GetPermissionsByRoleIDApi)
		rolesgroup.GET("/user", server.GetRoleByIDApi)

		rolesgroup.PUT("/assign", RBACMiddleware(server.store, "ROLES.UPDATE"), server.AssignRoleToEmployeeApi)

	}

}
