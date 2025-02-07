package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientRoutes(baseRouter *gin.RouterGroup) {
	clientsGroup := baseRouter.Group("/clients")
	clientsGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		clientsGroup.POST("", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientApi)
		clientsGroup.GET("", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientsApi)
		clientsGroup.GET("/:id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientApi)

		clientsGroup.PUT("/:id/profile_picture", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.SetClientProfilePictureApi)

		clientsGroup.POST("/:id/documents", RBACMiddleware(server.store, "CLIENT.CREATE"), server.AddClientDocumentApi)
		clientsGroup.GET("/:id/documents", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientDocumentsApi)
	}
}
