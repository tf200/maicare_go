// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRegistrationFormRoutes(baseRouter *gin.RouterGroup) {

	rfRoutes := baseRouter.Group("/registration_form")
	{
		rfRoutes.POST("", server.CreateRegistrationFormApi)
		rfRoutes.GET("", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.VIEW"), server.ListRegistrationFormsApi)
		rfRoutes.GET("/:id", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.VIEW"), server.GetRegistrationFormApi)
		rfRoutes.PUT("/:id", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.UpdateRegistrationFormApi)
		rfRoutes.DELETE("/:id", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.DELETE"), server.DeleteRegistrationFormApi)
		rfRoutes.POST("/:id/status", server.AuthMiddleware(), server.RBACMiddleware("REGISTRATION_FORM.UPDATE"), server.UpdateRegistrationFormStatusApi)

	}

}
