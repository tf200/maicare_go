// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRegistrationFormRoutes(baseRouter *gin.RouterGroup) {

	rfRoutes := baseRouter.Group("/registration_form")
	{
		rfRoutes.POST("", server.CreateRegistrationFormApi)
		rfRoutes.GET("", AuthMiddleware(server.tokenMaker), RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListRegistrationFormsApi)
		rfRoutes.GET("/:id", AuthMiddleware(server.tokenMaker), RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetRegistrationFormApi)
		rfRoutes.PUT("/:id", AuthMiddleware(server.tokenMaker), RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateRegistrationFormApi)
		rfRoutes.DELETE("/:id", AuthMiddleware(server.tokenMaker), RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteRegistrationFormApi)
		rfRoutes.POST("/:id/status", AuthMiddleware(server.tokenMaker), RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateRegistrationFormStatusApi)

	}

}
