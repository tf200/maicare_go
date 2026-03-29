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
		clientsGroup.GET("/status-counts", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientStatusCountsApi)
	}
	{
		clientsGroup.GET("/:id", server.RBACMiddleware("CLIENT.VIEW"), server.AuditMiddleware(), server.GetClientApi)
		clientsGroup.PUT("/:id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateClientApi)

		clientsGroup.GET("/:id/addresses", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientAddressesApi)

		clientsGroup.PUT("/:id/status", server.RBACMiddleware("CLIENT.STATUS.UPDATE"), server.UpdateClientStatusApi)
		clientsGroup.PUT("/:id/put-in-care", server.RBACMiddleware("CLIENT.STATUS.UPDATE"), server.PutClientInCareApi)
		clientsGroup.PUT("/:id/put-out-of-care", server.RBACMiddleware("CLIENT.STATUS.UPDATE"), server.PutClientOutOfCareApi)
		clientsGroup.GET("/:id/status_history", server.RBACMiddleware("CLIENT.VIEW"), server.ListStatusHistoryApi)

		clientsGroup.POST("/:id/documents", server.RBACMiddleware("CLIENT.CREATE"), server.AddClientDocumentApi)
		clientsGroup.GET("/:id/documents", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientDocumentsApi)
		clientsGroup.DELETE("/:id/documents/:doc_id", server.RBACMiddleware("CLIENT.VIEW"), server.DeleteClientDocumentApi)

		clientsGroup.GET("/:id/missing_documents", server.RBACMiddleware("CLIENT.CREATE"), server.GetMissingClientDocumentsApi)

		clientsGroup.GET("/:id/evaluations/bootstrap", server.RBACMiddleware("CLIENT.VIEW"), server.GetGoalEvaluationBootstrapApi)
		clientsGroup.POST("/:id/goals", server.RBACMiddleware("CLIENT.UPDATE"), server.CreateClientGoalApi)
		clientsGroup.PATCH("/:id/goals/:goal_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateClientGoalApi)
		clientsGroup.GET("/:id/goals", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientGoalsForEvaluationPageApi)
		clientsGroup.GET("/:id/goals/:goal_id/history", server.RBACMiddleware("CLIENT.VIEW"), server.ListGoalEvaluationHistoryApi)
		clientsGroup.GET("/:id/evaluations/submitted", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientSubmittedEvaluationsApi)
		clientsGroup.POST("/:id/evaluations", server.RBACMiddleware("CLIENT.UPDATE"), server.CreateGoalEvaluationApi)

		clientsGroup.POST("/:id/location_transfer", server.RBACMiddleware("CLIENT.UPDATE"), server.RequestLocationTransferApi)
		clientsGroup.POST("/location_transfer/approve_reject", server.RBACMiddleware("CLIENT.UPDATE"), server.ApproveOrRejectClientLocationTransferApi)
		clientsGroup.GET("/location_transfer", server.RBACMiddleware("CLIENT.VIEW"), server.ListLocationTransferRequestsApi)
	}

	{
		evaluationsGroup.GET("/upcoming", server.RBACMiddleware("CLIENT.VIEW"), server.ListUpcomingEvaluationsApi)
		evaluationsGroup.GET("/recent-submitted", server.RBACMiddleware("CLIENT.VIEW"), server.ListRecentSubmittedEvaluationsApi)
		evaluationsGroup.GET("/recent-drafts", server.RBACMiddleware("CLIENT.VIEW"), server.ListRecentDraftEvaluationsApi)
		evaluationsGroup.GET("/:evaluation_id", server.RBACMiddleware("CLIENT.VIEW"), server.GetGoalEvaluationApi)
	}
}
