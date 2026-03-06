package api

import "github.com/gin-gonic/gin"

func (server *Server) setupSettingsRoutes(baseRouter *gin.RouterGroup) {
	settings := baseRouter.Group("/settings")
	settings.Use(server.AuthMiddleware())
	{
		settings.GET("/departments", server.RBACMiddleware("SETTINGS.DEPARTMENT.VIEW"), server.ListDepartmentsApi)
		settings.POST("/departments", server.RBACMiddleware("SETTINGS.DEPARTMENT.CREATE"), server.CreateDepartmentApi)
		settings.PUT("/departments/:id", server.RBACMiddleware("SETTINGS.DEPARTMENT.UPDATE"), server.UpdateDepartmentApi)
		settings.GET("/organization-profile", server.RBACMiddleware("SETTINGS.ORGANIZATION_PROFILE.VIEW"), server.GetOrganizationProfileApi)
		settings.PUT("/organization-profile", server.RBACMiddleware("SETTINGS.ORGANIZATION_PROFILE.UPDATE"), server.UpdateOrganizationProfileApi)
	}
}
