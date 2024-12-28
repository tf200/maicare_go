package api

import "github.com/gin-gonic/gin"

func (server *Server) setupSenderRoutes(baseRouter *gin.RouterGroup) {
	senderGroup := baseRouter.Group("/")
	senderGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		baseRouter.POST("/sender", server.CreateSenderApi)
		baseRouter.GET("/sender", server.ListSendersAPI)
	}
}
