package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientIncidentRoutes(baseRouter *gin.RouterGroup) {
	clientIncident := baseRouter.Group("/clients")
	clientIncident.Use(server.AuthMiddleware())
	{
		clientIncident.GET("/:id/incidents", server.RBACMiddleware("CLIENT.INCIDENT.VIEW"), server.ListIncidentsApi)
	}

	incidents := baseRouter.Group("/incidents")
	incidents.Use(server.AuthMiddleware())
	{
		incidents.POST("", server.RBACMiddleware("CLIENT.INCIDENT.CREATE"), server.CreateIncidentApi)
		incidents.GET("/:incident_id", server.RBACMiddleware("CLIENT.INCIDENT.VIEW"), server.GetIncidentApi)
		incidents.PUT("/:incident_id", server.RBACMiddleware("CLIENT.INCIDENT.UPDATE"), server.UpdateIncidentApi)
		incidents.DELETE("/:incident_id", server.RBACMiddleware("CLIENT.INCIDENT.DELETE"), server.DeleteIncidentApi)
		incidents.GET("/:incident_id/file", server.RBACMiddleware("CLIENT.INCIDENT.VIEW"), server.GenerateIncidentFileApi)
		incidents.PUT("/:incident_id/confirm", server.RBACMiddleware("CLIENT.INCIDENT.CONFIRM"), server.ConfirmIncidentApi)
	}
}
