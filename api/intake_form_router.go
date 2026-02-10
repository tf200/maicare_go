package api

import "github.com/gin-gonic/gin"

func (server *Server) setupIntakeFormRoutes(baseRouter *gin.RouterGroup) {

	intakeFormGroup := baseRouter.Group("/intake_form")

	{

		intakeFormGroup.POST("", server.CreateIntakeFormApi)
		intakeFormGroup.GET("", server.ListIntakeFormsApi)
		intakeFormGroup.PUT("/:id/outcome", server.CompleteIntakeFormApi)

	}
}
