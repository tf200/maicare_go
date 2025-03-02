// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupMaturityMatrixRoutes(baseRouter *gin.RouterGroup) {

	mmGroup := baseRouter.Group("/clients")
	mmGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		mmGroup.POST("/:id/maturity_matrix_assessment", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientMaturityMatrixAssessmentApi)
		mmGroup.GET("/:id/maturity_matrix_assessment", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientMaturityMatrixAssessmentsApi)

		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientMaturityMatrixAssessmentApi)

		mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientGoalsApi)
		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/goals", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientGoalsApi)
		mmGroup.GET("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientGoalApi)

		mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id/objectives", RBACMiddleware(server.store, "CLIENT.VIEW"), server.CreateGoalObjectiveApi)
		mmGroup.POST("/:id/maturity_matrix_assessment/:assessment_id/goals/:goal_id/objectives/generate", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GenerateObjectivesApi)
	}

	baseRouter.GET("/maturity_matrix", server.ListMaturityMatrixApi)
}
