// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupLocationRoutes(baseRouter *gin.RouterGroup) {

	organisationGroup := baseRouter.Group("")
	organisationGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		organisationGroup.POST("/organisations", RBACMiddleware(server.store, "LOCATION.CREATE"), server.CreateOrganisationApi)
		organisationGroup.GET("/organisations", RBACMiddleware(server.store, "LOCATION.VIEW"), server.ListOrganisationsApi)
		organisationGroup.GET("/organisations/:id", RBACMiddleware(server.store, "LOCATION.VIEW"), server.GetOrganisationApi)
		organisationGroup.PUT("/organisations/:id", RBACMiddleware(server.store, "LOCATION.UPDATE"), server.UpdateOrganisationApi)
		organisationGroup.DELETE("/organisations/:id", RBACMiddleware(server.store, "LOCATION.DELETE"), server.DeleteOrganisationApi)

		organisationGroup.POST("/organisations/:id/locations", RBACMiddleware(server.store, "LOCATION.CREATE"), server.CreateLocationApi)
		organisationGroup.GET("/organisations/:id/locations", RBACMiddleware(server.store, "LOCATION.VIEW"), server.ListLocationsApi)
		organisationGroup.GET("/locations", RBACMiddleware(server.store, "LOCATION.VIEW"), server.ListAllLocationsApi)
		organisationGroup.GET("/locations/:id", RBACMiddleware(server.store, "LOCATION.VIEW"), server.GetLocationApi)
		organisationGroup.PUT("/locations/:id", RBACMiddleware(server.store, "LOCATION.UPDATE"), server.UpdateLocationApi)
		organisationGroup.DELETE("/locations/:id", RBACMiddleware(server.store, "LOCATION.DELETE"), server.DeleteLocationApi)
	}

}
