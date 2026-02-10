package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientRoutes(baseRouter *gin.RouterGroup) {
	clientsGroup := baseRouter.Group("/clients")
	clientsGroup.Use(server.AuthMiddleware())
	evaluationsGroup := baseRouter.Group("/evaluations")
	evaluationsGroup.Use(server.AuthMiddleware())
	{
		clientsGroup.POST("", server.RBACMiddleware("CLIENT.CREATE"), server.CreateClientApi)
		clientsGroup.GET("", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientsApi)
		clientsGroup.GET("/waiting-list", server.RBACMiddleware("CLIENT.VIEW"), server.ListWaitingListClientsApi)
		clientsGroup.GET("/in-care", server.RBACMiddleware("CLIENT.VIEW"), server.ListInCareClientsApi)
		clientsGroup.GET("/counts", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientsCountApi)
	}
	{
		clientsGroup.GET("/:id", server.RBACMiddleware("CLIENT.VIEW"), server.AuditMiddleware(), server.GetClientApi)
		clientsGroup.PUT("/:id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateClientApi)

		clientsGroup.GET("/:id/addresses", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientAddressesApi)

		clientsGroup.PUT("/:id/status", server.RBACMiddleware("CLIENT.STATUS.UPDATE"), server.UpdateClientStatusApi)
		clientsGroup.PUT("/:id/put-in-care", server.RBACMiddleware("CLIENT.STATUS.UPDATE"), server.PutClientInCareApi)
		clientsGroup.GET("/:id/status_history", server.RBACMiddleware("CLIENT.VIEW"), server.ListStatusHistoryApi)

		clientsGroup.POST("/:id/documents", server.RBACMiddleware("CLIENT.CREATE"), server.AddClientDocumentApi)
		clientsGroup.GET("/:id/documents", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientDocumentsApi)
		clientsGroup.DELETE("/:id/documents/:doc_id", server.RBACMiddleware("CLIENT.VIEW"), server.DeleteClientDocumentApi)

		clientsGroup.GET("/:id/missing_documents", server.RBACMiddleware("CLIENT.CREATE"), server.GetMissingClientDocumentsApi)

		clientsGroup.POST("/:id/appointments", server.RBACMiddleware("CLIENT.CREATE"), server.ListAppointmentsForClientApi)
		clientsGroup.GET("/:id/evaluations/bootstrap", server.RBACMiddleware("CLIENT.VIEW"), server.GetGoalEvaluationBootstrapApi)
		clientsGroup.POST("/:id/evaluations", server.RBACMiddleware("CLIENT.UPDATE"), server.CreateGoalEvaluationApi)

		clientsGroup.POST("/:id/location_transfer", server.RBACMiddleware("CLIENT.UPDATE"), server.RequestLocationTransferApi)
		clientsGroup.POST("/location_transfer/approve_reject", server.RBACMiddleware("CLIENT.UPDATE"), server.ApproveOrRejectClientLocationTransferApi)
		clientsGroup.GET("/location_transfer", server.RBACMiddleware("CLIENT.VIEW"), server.ListLocationTransferRequestsApi)
	}

	{
		evaluationsGroup.GET("/upcoming", server.RBACMiddleware("CLIENT.VIEW"), server.ListUpcomingEvaluationsApi)
		evaluationsGroup.GET("/recent-submitted", server.RBACMiddleware("CLIENT.VIEW"), server.ListRecentSubmittedEvaluationsApi)
		evaluationsGroup.GET("/recent-drafts", server.RBACMiddleware("CLIENT.VIEW"), server.ListRecentDraftEvaluationsApi)
		evaluationsGroup.PATCH("/:id/draft", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateEvaluationDraftApi)
		evaluationsGroup.POST("/:id/submit", server.RBACMiddleware("CLIENT.UPDATE"), server.SubmitEvaluationDraftApi)
	}
}
