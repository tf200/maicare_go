package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientIncidentRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientIncident := baseRouter.Group("/clients")
	ClientIncident.Use(server.AuthMiddleware())
	{
		ClientIncident.POST("/:id/incidents", server.RBACMiddleware("CLIENT.INCIDENT.CREATE"), server.CreateIncidentApi)
		ClientIncident.GET("/:id/incidents", server.RBACMiddleware("CLIENT.INCIDENT.VIEW"), server.ListIncidentsApi)
		ClientIncident.GET("/:id/incidents/:incident_id", server.RBACMiddleware("CLIENT.INCIDENT.VIEW"), server.GetIncidentApi)
		ClientIncident.PUT("/:id/incidents/:incident_id", server.RBACMiddleware("CLIENT.INCIDENT.UPDATE"), server.UpdateIncidentApi)
		ClientIncident.DELETE("/:id/incidents/:incident_id", server.RBACMiddleware("CLIENT.INCIDENT.DELETE"), server.DeleteIncidentApi)

		ClientIncident.GET("/:id/incidents/:incident_id/file", server.RBACMiddleware("CLIENT.INCIDENT.VIEW"), server.GenerateIncidentFileApi)
		ClientIncident.PUT("/:id/incidents/:incident_id/confirm", server.RBACMiddleware("CLIENT.INCIDENT.CONFIRM"), server.ConfirmIncidentApi)
	}

}
