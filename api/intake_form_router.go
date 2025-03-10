package api

import "github.com/gin-gonic/gin"

func (server *Server) setupIntakeFormRoutes(baseRouter *gin.RouterGroup) {

	intakeFormGroup := baseRouter.Group("/intake_form")

	{

		intakeFormGroup.POST("/upload", server.IntakeFormUploadHandlerApi)
		intakeFormGroup.POST("", server.CreateIntakeFormApi)
		intakeFormGroup.GET("", server.ListIntakeFormsApi)
		intakeFormGroup.GET("/:id", server.GetIntakeFormApi)

		intakeFormGroup.POST("/:id/urgency_score", server.AddUrgencyScoreApi)

	}
}
