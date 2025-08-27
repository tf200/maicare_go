// // auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupMaturityMatrixRoutes(baseRouter *gin.RouterGroup) {
	mmGroup := baseRouter.Group("")
	mmGroup.Use(server.AuthMiddleware())
	{
		mmGroup.POST("/clients/:id/assessments", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateClientMaturityMatrixAssessmentApi)
		mmGroup.GET("/clients/:id/assessments", server.RBACMiddleware("CARE_PLAN.VIEW"), server.ListClientMaturityMatrixAssessmentsApi)

		// Careplan routes
		mmGroup.GET("/care_plans/:care_plan_id", server.RBACMiddleware("CARE_PLAN.VIEW"), server.GetCarePlanOverviewApi)
		mmGroup.PUT("/care_plans/:care_plan_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanOverviewApi)
		mmGroup.DELETE("/care_plans/:care_plan_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanApi)

		// Careplan Objectives routes
		mmGroup.POST("/care_plans/:care_plan_id/objectives", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateCarePlanObjectiveApi)
		mmGroup.GET("/care_plans/:care_plan_id/objectives", server.RBACMiddleware("CARE_PLAN.VIEW"), server.GetCarePlanObjectivesApi)
		mmGroup.PUT("objectives/:objective_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanObjectiveApi)
		mmGroup.DELETE("objectives/:objective_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanObjectiveApi)

		// Careplan Actions routes
		mmGroup.POST("/objectives/:objective_id/actions", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateCarePlanActionsApi)
		mmGroup.PUT("/actions/:action_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanActionsApi)
		mmGroup.DELETE("/actions/:action_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanActionApi)

		// Careplan Interventions routes
		mmGroup.POST("/care_plans/:care_plan_id/interventions", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateCarePlanInterventionApi)
		mmGroup.GET("/care_plans/:care_plan_id/interventions", server.RBACMiddleware("CARE_PLAN.VIEW"), server.GetCarePlanInterventionsApi)
		mmGroup.PUT("/interventions/:intervention_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanInterventionApi)
		mmGroup.DELETE("/interventions/:intervention_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanInterventionApi)

		// Careplan Success Metrics routes
		mmGroup.POST("/care_plans/:care_plan_id/success_metrics", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateCarePlanSuccessMetricsApi)
		mmGroup.GET("/care_plans/:care_plan_id/success_metrics", server.RBACMiddleware("CARE_PLAN.VIEW"), server.GetCarePlanSuccessMetricsApi)
		mmGroup.PUT("/success_metrics/:metric_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanSuccessMetricsApi)
		mmGroup.DELETE("/success_metrics/:metric_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanSuccessMetricApi)

		// Careplan Risks routes
		mmGroup.POST("/care_plans/:care_plan_id/risks", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateCarePlanRisksApi)
		mmGroup.GET("/care_plans/:care_plan_id/risks", server.RBACMiddleware("CARE_PLAN.VIEW"), server.GetCarePlanRisksApi)
		mmGroup.PUT("/risks/:risk_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanRisksApi)
		mmGroup.DELETE("/risks/:risk_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanRiskApi)

		// Careplan Supportnetwork routes
		mmGroup.POST("/care_plans/:care_plan_id/support_network", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateCareplanSupportNetworkApi)
		mmGroup.GET("/care_plans/:care_plan_id/support_network", server.RBACMiddleware("CARE_PLAN.VIEW"), server.GetCarePlanSupportNetworkApi)
		mmGroup.PUT("/support_network/:support_network_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanSupportNetworkApi)
		mmGroup.DELETE("/support_network/:support_network_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanSupportNetworkApi)

		// Careplan Resources routes
		mmGroup.GET("/care_plans/:care_plan_id/resources", server.RBACMiddleware("CARE_PLAN.VIEW"), server.GetCarePlanResourcesApi)
		mmGroup.POST("/care_plans/:care_plan_id/resources", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateCarePlanResourcesApi)
		mmGroup.PUT("/resources/:resource_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanResourcesApi)
		mmGroup.DELETE("/resources/:resource_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanResourcesApi)

		// Careplan Reports routes
		mmGroup.POST("/care_plans/:care_plan_id/reports", server.RBACMiddleware("CARE_PLAN.CREATE"), server.CreateCarePlanReportApi)
		mmGroup.GET("/care_plans/:care_plan_id/reports", server.RBACMiddleware("CARE_PLAN.VIEW"), server.ListCarePlanReportsApi)
		mmGroup.PUT("/care_plans/reports/:report_id", server.RBACMiddleware("CARE_PLAN.UPDATE"), server.UpdateCarePlanReportApi)
		mmGroup.DELETE("/care_plans/reports/:report_id", server.RBACMiddleware("CARE_PLAN.DELETE"), server.DeleteCarePlanReportApi)

		// mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id/objectives/generate", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GenerateObjectivesApi)
	}

	baseRouter.GET("/maturity_matrix", server.ListMaturityMatrixApi)
}
