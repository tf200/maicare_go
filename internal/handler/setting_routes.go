package handler

import (
	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
)

func RegisterSettingsRoutes(
	rg *gin.RouterGroup,
	handler *SettingsHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	settings := rg.Group("/settings")
	{
		settings.GET("/departments", auth, requirePermission(domain.PermSettingsDepartmentView.String()), handler.ListDepartments)
		settings.POST("/departments", auth, requirePermission(domain.PermSettingsDepartmentCreate.String()), handler.CreateDepartment)
		settings.PUT("/departments/:id", auth, requirePermission(domain.PermSettingsDepartmentUpdate.String()), handler.UpdateDepartment)
		settings.GET("/organization-profile", auth, requirePermission(domain.PermSettingsOrgProfileView.String()), handler.GetOrganizationProfile)
		settings.PUT("/organization-profile", auth, requirePermission(domain.PermSettingsOrgProfileUpdate.String()), handler.UpdateOrganizationProfile)
	}
}
