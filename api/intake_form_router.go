package api

import "github.com/gin-gonic/gin"

func (server *Server) setupIntakeFormRoutes(baseRouter *gin.RouterGroup) {

	intakeFormGroup := baseRouter.Group("/intake_forms")
	intakeFormGroup.Use(server.AuthMiddleware())

	{

		intakeFormGroup.POST("", server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.CreateIntakeFormApi)
		intakeFormGroup.GET("", server.RBACMiddleware("REGISTRATION_FORM.VIEW"), server.ListIntakeFormsApi)
		intakeFormGroup.GET("/:id", server.RBACMiddleware("REGISTRATION_FORM.VIEW"), server.GetIntakeFormApi)
		intakeFormGroup.PUT("/:id/goals", server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.CreateIntakeFormGoalsApi)
		intakeFormGroup.PATCH("/:id/conclusion", server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.UpdateIntakeConclusionApi)
		intakeFormGroup.POST("/:id/promote", server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.PromoteIntakeToClientApi)

	}
}
