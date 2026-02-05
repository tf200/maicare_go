// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRegistrationFormRoutes(baseRouter *gin.RouterGroup) {
	rfRoutes := baseRouter.Group("/registration_forms")
	{
		rfRoutes.POST("", server.CreateRegistrationFormApi)
		rfRoutes.GET("", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.VIEW"), server.ListRegistrationFormsApi)
		rfRoutes.GET("/:id", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.VIEW"), server.GetRegistrationFormApi)
		rfRoutes.PUT("/:id", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.UpdateRegistrationFormApi)
		rfRoutes.DELETE("/:id", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.DELETE"), server.DeleteRegistrationFormApi)
		rfRoutes.POST("/:id/status", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.UpdateRegistrationFormStatusApi)
		rfRoutes.POST("/:id/process", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.ProcessRegistrationFormApi)
	}

	intakeMaturityRoutes := baseRouter.Group("/intake_maturity")
	intakeMaturityRoutes.Use(server.AuthMiddleware())
	{
		intakeMaturityRoutes.POST("/:assessment_id/generate_goals", server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.GenerateIntakeGoalsApi)
	}

	// intakeFormRoutes := baseRouter.Group("/intake_forms")
	// intakeFormRoutes.Use(server.AuthMiddleware())
	// {
	// 	intakeFormRoutes.POST("/:id/goals", server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.CreateIntakeFormGoalsApi)
	// 	intakeFormRoutes.POST("/:id/promote", server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.PromoteIntakeToClientApi)
	// }

	publicRoutes := baseRouter.Group("/public/intake-options")
	{
		publicRoutes.GET("/:token", server.GetPublicIntakeOptionsApi)
		publicRoutes.POST("/:token/confirm", server.SelectIntakeDateApi)
	}
}
