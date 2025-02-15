package api

import "github.com/gin-gonic/gin"

func (server *Server) setupProgressReportsRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ProgressReports := baseRouter.Group("/clients")
	ProgressReports.Use(AuthMiddleware(server.tokenMaker))
	{
		ProgressReports.POST("/:id/progress_reports", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateProgressReportApi)
		ProgressReports.GET("/:id/progress_reports", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListProgressReportsApi)
		ProgressReports.GET("/:id/progress_reports/:report_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetProgressReportApi)
		ProgressReports.PUT("/:id/progress_reports/:report_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateProgressReportApi)

		ProgressReports.POST("/:id/ai_progress_reports", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GenerateAutoReportsApi)
		ProgressReports.POST("/:id/ai_progress_reports/confirm", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ConfirmProgressReportApi)
		ProgressReports.GET("/:id/ai_progress_reports", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListAiGeneratedReportsApi)

	}

}
