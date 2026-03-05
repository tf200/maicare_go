package api

import "github.com/gin-gonic/gin"

func (server *Server) setupIncidentsAllRoutes(baseRouter *gin.RouterGroup) {
	incidents := baseRouter.Group("/incidents")
	incidents.Use(server.AuthMiddleware())

	{
		incidents.GET("", server.RBACMiddleware("CLIENT.INCIDENT.VIEW"), server.ListAllIncidentsApi)
		incidents.GET("/counts", server.RBACMiddleware("CLIENT.INCIDENT.VIEW"), server.GetIncidentCountsApi)
	}
}
