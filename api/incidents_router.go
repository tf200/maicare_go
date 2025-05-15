package api

import "github.com/gin-gonic/gin"

func (server *Server) setupIncidentsAllRoutes(baseRouter *gin.RouterGroup) {

	incidents := baseRouter.Group("/incidents")

	{
		incidents.GET("", server.ListAllIncidentsApi)
	}
}
