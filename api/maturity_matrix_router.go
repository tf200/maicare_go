// // auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupMaturityMatrixRoutes(baseRouter *gin.RouterGroup) {
	mmGroup := baseRouter.Group("")
	mmGroup.Use(server.AuthMiddleware())
	{
		mmGroup.POST("/clients/:id/assessments", server.RBACMiddleware("CLIENT.CREATE"), server.CreateClientMaturityMatrixAssessmentApi)
		mmGroup.GET("/clients/:id/assessments", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientMaturityMatrixAssessmentsApi)

		// Careplan routes
		mmGroup.GET("/care_plans/:care_plan_id", server.RBACMiddleware("CLIENT.VIEW"), server.GetCarePlanOverviewApi)
		mmGroup.PUT("/care_plans/:care_plan_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanOverviewApi)
		mmGroup.DELETE("/care_plans/:care_plan_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanApi)

		// Careplan Objectives routes
		mmGroup.POST("/care_plans/:care_plan_id/objectives", server.RBACMiddleware("CLIENT.CREATE"), server.CreateCarePlanObjectiveApi)
		mmGroup.GET("/care_plans/:care_plan_id/objectives", server.RBACMiddleware("CLIENT.VIEW"), server.GetCarePlanObjectivesApi)
		mmGroup.PUT("objectives/:objective_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanObjectiveApi)
		mmGroup.DELETE("objectives/:objective_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanObjectiveApi)

		// Careplan Actions routes
		mmGroup.POST("/objectives/:objective_id/actions", server.RBACMiddleware("CLIENT.CREATE"), server.CreateCarePlanActionsApi)
		mmGroup.PUT("/actions/:action_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanActionsApi)
		mmGroup.DELETE("/actions/:action_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanActionApi)

		// Careplan Interventions routes
		mmGroup.POST("/care_plans/:care_plan_id/interventions", server.RBACMiddleware("CLIENT.CREATE"), server.CreateCarePlanInterventionApi)
		mmGroup.GET("/care_plans/:care_plan_id/interventions", server.RBACMiddleware("CLIENT.VIEW"), server.GetCarePlanInterventionsApi)
		mmGroup.PUT("/interventions/:intervention_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanInterventionApi)
		mmGroup.DELETE("/interventions/:intervention_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanInterventionApi)

		// Careplan Success Metrics routes
		mmGroup.POST("/care_plans/:care_plan_id/success_metrics", server.RBACMiddleware("CLIENT.CREATE"), server.CreateCarePlanSuccessMetricsApi)
		mmGroup.GET("/care_plans/:care_plan_id/success_metrics", server.RBACMiddleware("CLIENT.VIEW"), server.GetCarePlanSuccessMetricsApi)
		mmGroup.PUT("/success_metrics/:metric_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanSuccessMetricsApi)
		mmGroup.DELETE("/success_metrics/:metric_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanSuccessMetricApi)

		// Careplan Risks routes
		mmGroup.POST("/care_plans/:care_plan_id/risks", server.RBACMiddleware("CLIENT.CREATE"), server.CreateCarePlanRisksApi)
		mmGroup.GET("/care_plans/:care_plan_id/risks", server.RBACMiddleware("CLIENT.VIEW"), server.GetCarePlanRisksApi)
		mmGroup.PUT("/risks/:risk_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanRisksApi)
		mmGroup.DELETE("/risks/:risk_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanRiskApi)

		// Careplan Supportnetwork routes
		mmGroup.POST("/care_plans/:care_plan_id/support_network", server.RBACMiddleware("CLIENT.CREATE"), server.CreateCareplanSupportNetworkApi)
		mmGroup.GET("/care_plans/:care_plan_id/support_network", server.RBACMiddleware("CLIENT.VIEW"), server.GetCarePlanSupportNetworkApi)
		mmGroup.PUT("/support_network/:support_network_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanSupportNetworkApi)
		mmGroup.DELETE("/support_network/:support_network_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanSupportNetworkApi)

		// Careplan Resources routes
		mmGroup.GET("/care_plans/:care_plan_id/resources", server.RBACMiddleware("CLIENT.VIEW"), server.GetCarePlanResourcesApi)
		mmGroup.POST("/care_plans/:care_plan_id/resources", server.RBACMiddleware("CLIENT.CREATE"), server.CreateCarePlanResourcesApi)
		mmGroup.PUT("/resources/:resource_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanResourcesApi)
		mmGroup.DELETE("/resources/:resource_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanResourcesApi)

		// Careplan Reports routes
		mmGroup.POST("/care_plans/:care_plan_id/reports", server.RBACMiddleware("CLIENT.CREATE"), server.CreateCarePlanReportApi)
		mmGroup.GET("/care_plans/:care_plan_id/reports", server.RBACMiddleware("CLIENT.VIEW"), server.ListCarePlanReportsApi)
		mmGroup.PUT("/care_plans/reports/:report_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateCarePlanReportApi)
		mmGroup.DELETE("/care_plans/reports/:report_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteCarePlanReportApi)
		// old routes to be removed
		mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals", server.RBACMiddleware("CLIENT.CREATE"), server.CreateClientGoalsApi)
		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/goals", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientGoalsApi)
		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientGoalApi)

		mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id/objectives", server.RBACMiddleware("CLIENT.VIEW"), server.CreateGoalObjectiveApi)
		// mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id/objectives/generate", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GenerateObjectivesApi)
	}

	baseRouter.GET("/maturity_matrix", server.ListMaturityMatrixApi)
}
