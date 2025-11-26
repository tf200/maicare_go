package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientRoutes(baseRouter *gin.RouterGroup) {
	clientsGroup := baseRouter.Group("/clients")
	clientsGroup.Use(server.AuthMiddleware())
	{
		clientsGroup.POST("", server.RBACMiddleware("CLIENT.CREATE"), server.CreateClientApi)
		clientsGroup.GET("", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientsApi)
		clientsGroup.GET("/counts", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientsCountApi)
	}
	{
		clientsGroup.GET("/:id", server.RBACMiddleware("CLIENT.VIEW"), server.AuditMiddleware(), server.GetClientApi)
		clientsGroup.PUT("/:id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateClientApi)

		clientsGroup.PUT("/:id/profile_picture", server.RBACMiddleware("CLIENT.UPDATE"), server.SetClientProfilePictureApi)

		clientsGroup.GET("/:id/addresses", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientAddressesApi)

		clientsGroup.PUT("/:id/status", server.RBACMiddleware("CLIENT.STATUS.UPDATE"), server.UpdateClientStatusApi)
		clientsGroup.GET("/:id/status_history", server.RBACMiddleware("CLIENT.VIEW"), server.ListStatusHistoryApi)

		clientsGroup.POST("/:id/documents", server.RBACMiddleware("CLIENT.CREATE"), server.AddClientDocumentApi)
		clientsGroup.GET("/:id/documents", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientDocumentsApi)
		clientsGroup.DELETE("/:id/documents/:doc_id", server.RBACMiddleware("CLIENT.VIEW"), server.DeleteClientDocumentApi)

		clientsGroup.GET("/:id/missing_documents", server.RBACMiddleware("CLIENT.CREATE"), server.GetMissingClientDocumentsApi)

		clientsGroup.POST("/:id/appointments", server.RBACMiddleware("CLIENT.CREATE"), server.ListAppointmentsForClientApi)

		clientsGroup.POST("/:id/location_transfer", server.RBACMiddleware("CLIENT.UPDATE"), server.RequestLocationTransferApi)
		clientsGroup.POST("/location_transfer/approve_reject", server.RBACMiddleware("CLIENT.UPDATE"), server.ApproveOrRejectClientLocationTransferApi)
		clientsGroup.GET("/location_transfer", server.RBACMiddleware("CLIENT.VIEW"), server.ListLocationTransferRequestsApi)
	}
}
