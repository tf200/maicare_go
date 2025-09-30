// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupLocationRoutes(baseRouter *gin.RouterGroup) {

	organisationGroup := baseRouter.Group("")
	organisationGroup.Use(server.AuthMiddleware())
	{
		organisationGroup.POST("/organisations", server.RBACMiddleware("LOCATION.CREATE"), server.CreateOrganisationApi)
		organisationGroup.GET("/organisations", server.RBACMiddleware("LOCATION.VIEW"), server.ListOrganisationsApi)
		organisationGroup.GET("/organisations/:id", server.RBACMiddleware("LOCATION.VIEW"), server.GetOrganisationApi)
		organisationGroup.GET("/organisations/:id/counts", server.RBACMiddleware("LOCATION.VIEW"), server.GetOrganisationCountApi)
		organisationGroup.PUT("/organisations/:id", server.RBACMiddleware("LOCATION.UPDATE"), server.UpdateOrganisationApi)
		organisationGroup.DELETE("/organisations/:id", server.RBACMiddleware("LOCATION.DELETE"), server.DeleteOrganisationApi)

		organisationGroup.POST("/organisations/:id/locations", server.RBACMiddleware("LOCATION.CREATE"), server.CreateLocationApi)
		organisationGroup.GET("/organisations/:id/locations", server.RBACMiddleware("LOCATION.VIEW"), server.ListLocationsApi)
		organisationGroup.GET("/locations", server.RBACMiddleware("LOCATION.VIEW"), server.ListAllLocationsApi)
		organisationGroup.GET("/locations/:id", server.RBACMiddleware("LOCATION.VIEW"), server.GetLocationApi)
		organisationGroup.PUT("/locations/:id", server.RBACMiddleware("LOCATION.UPDATE"), server.UpdateLocationApi)
		organisationGroup.DELETE("/locations/:id", server.RBACMiddleware("LOCATION.DELETE"), server.DeleteLocationApi)
	}

}
