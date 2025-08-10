// // auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupMaturityMatrixRoutes(baseRouter *gin.RouterGroup) {
	mmGroup := baseRouter.Group("")
	mmGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		mmGroup.POST("/clients/:id/maturity_matrix_assessment", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientMaturityMatrixAssessmentApi)
		mmGroup.GET("/clients/:id/maturity_matrix_assessment", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientMaturityMatrixAssessmentsApi)

		// Careplan routes
		mmGroup.GET("/care_plans/:care_plan_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetCarePlanOverviewApi)
		mmGroup.PUT("/care_plans/:care_plan_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateCarePlanOverviewApi)
		mmGroup.DELETE("/care_plans/:care_plan_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteCarePlanApi)

		// Careplan Objectives routes
		mmGroup.POST("/care_plans/:care_plan_id/objectives", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateCarePlanObjectiveApi)
		mmGroup.GET("/care_plans/:care_plan_id/objectives", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetCarePlanObjectivesApi)
		mmGroup.PUT("objectives/:objective_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateCarePlanObjectiveApi)
		mmGroup.DELETE("objectives/:objective_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteCarePlanObjectiveApi)

		// Careplan Actions routes
		mmGroup.POST("/objectives/:objective_id/actions", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateCarePlanActionsApi)
		mmGroup.PUT("/actions/:action_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateCarePlanActionsApi)
		mmGroup.DELETE("/actions/:action_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteCarePlanActionApi)

		// Careplan Interventions routes
		mmGroup.POST("/care_plans/:care_plan_id/interventions", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateCarePlanInterventionApi)
		mmGroup.GET("/care_plans/:care_plan_id/interventions", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetCarePlanInterventionsApi)
		mmGroup.PUT("/interventions/:intervention_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateCarePlanInterventionApi)
		mmGroup.DELETE("/interventions/:intervention_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteCarePlanInterventionApi)

		// Careplan Successes routes
		mmGroup.POST("/care_plans/:care_plan_id/success_metrics", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateCarePlanSuccessMetricsApi)
		mmGroup.GET("/care_plans/:care_plan_id/success_metrics", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetCarePlanSuccessMetricsApi)
		mmGroup.PUT("/success_metrics/:metric_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateCarePlanSuccessMetricsApi)
		mmGroup.DELETE("/success_metrics/:metric_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteCarePlanSuccessMetricApi)

		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/careplan/successes", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetCarePlanSuccessMetricsApi)

		// Careplan Risks routes
		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/careplan/risks", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetCarePlanRisksApi)

		// Careplan Supportnetwork routes
		// mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/careplan/supportnetwork", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetCarePlanSupportNetworkApi)

		// Careplan Resources routes
		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/careplan/resources", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetCarePlanResourcesApi)

		mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientGoalsApi)
		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/goals", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientGoalsApi)
		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientGoalApi)

		mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id/objectives", RBACMiddleware(server.store, "CLIENT.VIEW"), server.CreateGoalObjectiveApi)
		// mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id/objectives/generate", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GenerateObjectivesApi)
	}

	baseRouter.GET("/maturity_matrix", server.ListMaturityMatrixApi)
}
