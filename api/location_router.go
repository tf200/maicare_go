// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupLocationRoutes(baseRouter *gin.RouterGroup) {
	organizationGroup := baseRouter.Group("")
	organizationGroup.Use(server.AuthMiddleware())
	{
		organizationGroup.POST("/organisations", server.RBACMiddleware("LOCATION.CREATE"), server.CreateOrganisationApi)
		organizationGroup.GET("/organisations", server.RBACMiddleware("LOCATION.VIEW"), server.ListOrganisationsApi)
		organizationGroup.GET("/organizations/count", server.RBACMiddleware("LOCATION.VIEW"), server.GetGlobalOrganisationCountApi)
		organizationGroup.GET("/organisations/:id", server.RBACMiddleware("LOCATION.VIEW"), server.GetOrganisationApi)
		organizationGroup.GET("/organisations/:id/counts", server.RBACMiddleware("LOCATION.VIEW"), server.GetOrganisationCountApi)
		organizationGroup.PUT("/organisations/:id", server.RBACMiddleware("LOCATION.UPDATE"), server.UpdateOrganisationApi)
		organizationGroup.DELETE("/organisations/:id", server.RBACMiddleware("LOCATION.DELETE"), server.DeleteOrganisationApi)

		organizationGroup.POST("/organisations/:id/locations", server.RBACMiddleware("LOCATION.CREATE"), server.CreateLocationApi)
		organizationGroup.GET("/organisations/:id/locations", server.RBACMiddleware("LOCATION.VIEW"), server.ListLocationsApi)
		organizationGroup.GET("/locations", server.RBACMiddleware("LOCATION.VIEW"), server.ListAllLocationsApi)
		organizationGroup.GET("/locations/:id", server.RBACMiddleware("LOCATION.VIEW"), server.GetLocationApi)
		organizationGroup.PUT("/locations/:id", server.RBACMiddleware("LOCATION.UPDATE"), server.UpdateLocationApi)
		organizationGroup.DELETE("/locations/:id", server.RBACMiddleware("LOCATION.DELETE"), server.DeleteLocationApi)
	}
}
