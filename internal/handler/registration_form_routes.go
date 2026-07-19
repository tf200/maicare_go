package handler

import "github.com/gin-gonic/gin"

func RegisterRegistrationFormRoutes(
	rg *gin.RouterGroup,
	handler *RegistrationFormHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	// Public routes
	rg.POST("/public/registration-upload-sessions", handler.StartUploadSession)
	rg.POST("/public/registration-uploads/init", handler.InitRegistrationUpload)
	rg.POST("/registration_forms", handler.CreateRegistrationForm)
	rg.GET("/public/intake-options/:token", handler.GetPublicIntakeOptions)
	rg.POST("/public/intake-options/:token/confirm", handler.SelectIntakeDate)

	// Protected routes
	rg.GET("/registration_forms", auth, requirePermission("REGISTRATION_FORM.VIEW"), handler.ListRegistrationForms)
	rg.GET("/registration_forms/:id", auth, requirePermission("REGISTRATION_FORM.VIEW"), handler.GetRegistrationForm)
	rg.PUT("/registration_forms/:id", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.UpdateRegistrationForm)
	rg.DELETE("/registration_forms/:id", auth, requirePermission("REGISTRATION_FORM.DELETE"), handler.DeleteRegistrationForm)
	rg.POST("/registration_forms/:id/status", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.UpdateRegistrationFormStatus)
	rg.POST("/registration_forms/:id/process", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.ProcessRegistrationForm)
}
