package api

import "github.com/gin-gonic/gin"

func (server *Server) setupIntakeFormRoutes(baseRouter *gin.RouterGroup) {

	intakeFormGroup := baseRouter.Group("/intake_form")

	{

		intakeFormGroup.POST("/upload", server.IntakeFormUploadHandlerApi)
		intakeFormGroup.POST("/token", server.GenerateIntakeFormToken)
		intakeFormGroup.GET("/token/:token", server.VerifyIntakeFormToken)
		intakeFormGroup.POST("", server.CreateIntakeFormApi)
	}

}
