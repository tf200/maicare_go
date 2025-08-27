package api

import "github.com/gin-gonic/gin"

func (server *Server) setupProgressReportsRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ProgressReports := baseRouter.Group("/clients")
	ProgressReports.Use(server.AuthMiddleware())
	{
		ProgressReports.POST("/:id/progress_reports", server.RBACMiddleware("CLIENT.PROGRESS_REPORT.CREATE"), server.CreateProgressReportApi)
		ProgressReports.GET("/:id/progress_reports", server.RBACMiddleware("CLIENT.PROGRESS_REPORT.VIEW"), server.ListProgressReportsApi)
		ProgressReports.GET("/:id/progress_reports/:report_id", server.RBACMiddleware("CLIENT.PROGRESS_REPORT.VIEW"), server.GetProgressReportApi)
		ProgressReports.PUT("/:id/progress_reports/:report_id", server.RBACMiddleware("CLIENT.PROGRESS_REPORT.UPDATE"), server.UpdateProgressReportApi)
		ProgressReports.DELETE("/:id/progress_reports/:report_id", server.RBACMiddleware("CLIENT.PROGRESS_REPORT.DELETE"), server.DeleteProgressReportApi)

		ProgressReports.POST("/:id/ai_progress_reports", server.RBACMiddleware("CLIENT.AI_PROGRESS_REPORT.GENERATE"), server.GenerateAutoReportsApi)
		ProgressReports.POST("/:id/ai_progress_reports/confirm", server.RBACMiddleware("CLIENT.AI_PROGRESS_REPORT.CONFIRM"), server.ConfirmProgressReportApi)
		ProgressReports.GET("/:id/ai_progress_reports", server.RBACMiddleware("CLIENT.AI_PROGRESS_REPORT.VIEW"), server.ListAiGeneratedReportsApi)

	}

}
