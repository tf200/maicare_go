package api

import "github.com/gin-gonic/gin"

func (server *Server) setupIntakeFormRoutes(baseRouter *gin.RouterGroup) {

	intakeFormGroup := baseRouter.Group("/intake_forms")

	{

		intakeFormGroup.POST("", server.CreateIntakeFormApi)
		intakeFormGroup.GET("", server.ListIntakeFormsApi)
		intakeFormGroup.GET("/:id", server.GetIntakeFormApi)
		intakeFormGroup.PUT("/:id/goals", server.CreateIntakeFormGoalsApi)

	}
}
