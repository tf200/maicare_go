package handler

import (
	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
)

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
	rg.GET("/registration_forms", auth, requirePermission(domain.PermRegistrationFormView.String()), handler.ListRegistrationForms)
	rg.GET("/registration_forms/counts", auth, requirePermission(domain.PermRegistrationFormView.String()), handler.GetRegistrationFormCounts)
	rg.GET("/registration_forms/:id", auth, requirePermission(domain.PermRegistrationFormView.String()), handler.GetRegistrationForm)
	rg.PUT("/registration_forms/:id", auth, requirePermission(domain.PermRegistrationFormUpdate.String()), handler.UpdateRegistrationForm)
	rg.PUT("/registration_forms/:id/documents", auth, requirePermission(domain.PermRegistrationFormUpdate.String()), handler.ReplaceRegistrationFormDocument)
	rg.DELETE("/registration_forms/:id", auth, requirePermission(domain.PermRegistrationFormDelete.String()), handler.DeleteRegistrationForm)
	rg.POST("/registration_forms/:id/status", auth, requirePermission(domain.PermRegistrationFormUpdate.String()), handler.UpdateRegistrationFormStatus)
	rg.POST("/registration_forms/:id/process", auth, requirePermission(domain.PermRegistrationFormUpdate.String()), handler.ProcessRegistrationForm)
}
