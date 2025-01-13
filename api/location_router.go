// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupLocationRoutes(baseRouter *gin.RouterGroup) {

	locationGroup := baseRouter.Group("/locations")
	locationGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		locationGroup.GET("", server.ListLocationsApi)
		locationGroup.POST("", RBACMiddleware(server.store, "LOCATION.CREATE"), server.CreateLocationApi)
		locationGroup.PUT("/:id", RBACMiddleware(server.store, "LOCATION.UPDATE"), server.UpdateLocationApi)
	}

}
