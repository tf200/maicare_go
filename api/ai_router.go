package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAiRoutes(baseRouter *gin.RouterGroup) {
	aiGroup := baseRouter.Group("/ai")
	aiGroup.Use(server.AuthMiddleware())
	{
		aiGroup.POST("/spelling_check", server.SpellingCheckApi)
	}

}
