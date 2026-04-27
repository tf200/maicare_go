package handler

import "github.com/gin-gonic/gin"

func RegisterSettingsRoutes(
	rg *gin.RouterGroup,
	handler *SettingsHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	settings := rg.Group("/settings")
	{
		settings.GET("/departments", auth, requirePermission("SETTINGS.DEPARTMENT.VIEW"), handler.ListDepartments)
		settings.POST("/departments", auth, requirePermission("SETTINGS.DEPARTMENT.CREATE"), handler.CreateDepartment)
		settings.PUT("/departments/:id", auth, requirePermission("SETTINGS.DEPARTMENT.UPDATE"), handler.UpdateDepartment)
		settings.GET("/organization-profile", auth, requirePermission("SETTINGS.ORGANIZATION_PROFILE.VIEW"), handler.GetOrganizationProfile)
		settings.PUT("/organization-profile", auth, requirePermission("SETTINGS.ORGANIZATION_PROFILE.UPDATE"), handler.UpdateOrganizationProfile)
	}
}
