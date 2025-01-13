// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupLocationRoutes(baseRouter *gin.RouterGroup) {

	locationGroup := baseRouter.Group("")
	locationGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		locationGroup.GET("/locations", server.ListLocationsApi)
		locationGroup.POST("/locations", RBACMiddleware(server.store, "LOCATION.CREATE"), server.CreateLocationApi)
	}

}
