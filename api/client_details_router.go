package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientRoutes(baseRouter *gin.RouterGroup) {
	clientsGroup := baseRouter.Group("/clients")
	clientsGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		clientsGroup.POST("", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientApi)
		clientsGroup.GET("", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientsApi)

	}
}
