package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientIncidentRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientIncident := baseRouter.Group("/clients")
	ClientIncident.Use(AuthMiddleware(server.tokenMaker))
	{
		ClientIncident.POST("/:id/incidents", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateIncidentApi)
		ClientIncident.GET("/:id/incidents", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListIncidentsApi)

	}

}
