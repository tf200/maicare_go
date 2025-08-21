package api

import "github.com/gin-gonic/gin"

func (server *Server) setupProgressReportsRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ProgressReports := baseRouter.Group("/clients")
	ProgressReports.Use(server.AuthMiddleware())
	{
		ProgressReports.POST("/:id/progress_reports", server.RBACMiddleware("CLIENT.CREATE"), server.CreateProgressReportApi)
		ProgressReports.GET("/:id/progress_reports", server.RBACMiddleware("CLIENT.VIEW"), server.ListProgressReportsApi)
		ProgressReports.GET("/:id/progress_reports/:report_id", server.RBACMiddleware("CLIENT.VIEW"), server.GetProgressReportApi)
		ProgressReports.PUT("/:id/progress_reports/:report_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateProgressReportApi)
		ProgressReports.DELETE("/:id/progress_reports/:report_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteProgressReportApi)

		ProgressReports.POST("/:id/ai_progress_reports", server.RBACMiddleware("CLIENT.VIEW"), server.GenerateAutoReportsApi)
		ProgressReports.POST("/:id/ai_progress_reports/confirm", server.RBACMiddleware("CLIENT.VIEW"), server.ConfirmProgressReportApi)
		ProgressReports.GET("/:id/ai_progress_reports", server.RBACMiddleware("CLIENT.VIEW"), server.ListAiGeneratedReportsApi)

	}

}
