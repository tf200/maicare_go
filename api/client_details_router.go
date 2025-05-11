package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientRoutes(baseRouter *gin.RouterGroup) {
	clientsGroup := baseRouter.Group("/clients")
	clientsGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		clientsGroup.POST("", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientApi)
		clientsGroup.GET("", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientsApi)
		clientsGroup.GET("/counts", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientsCountApi)
		clientsGroup.GET("/:id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientApi)
		clientsGroup.PUT("/:id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateClientApi)

		clientsGroup.PUT("/:id/profile_picture", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.SetClientProfilePictureApi)

		clientsGroup.GET("/:id/addresses", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientAddressesApi)

		clientsGroup.PUT("/:id/status", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateClientStatusApi)
		clientsGroup.GET("/:id/status_history", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListStatusHistoryApi)

		clientsGroup.POST("/:id/documents", RBACMiddleware(server.store, "CLIENT.CREATE"), server.AddClientDocumentApi)
		clientsGroup.GET("/:id/documents", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientDocumentsApi)
		clientsGroup.DELETE("/:id/documents/:doc_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.DeleteClientDocumentApi)

		clientsGroup.GET("/:id/missing_documents", RBACMiddleware(server.store, "CLIENT.CREATE"), server.GetMissingClientDocumentsApi)

		clientsGroup.POST("/:id/appointments", RBACMiddleware(server.store, "CLIENT.CREATE"), server.ListAppointmentsForClientApi)
	}
}
