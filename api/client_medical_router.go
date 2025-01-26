package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientMedicalRoutes(baseRouter *gin.RouterGroup) {
	ClientMedical := baseRouter.Group("/clients")
	ClientMedical.Use(AuthMiddleware(server.tokenMaker))
	{
		ClientMedical.POST("/:id/client_allergies", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientAllergyApi)
		ClientMedical.GET("/:id/client_allergies", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientAllergiesApi)

	}
}
