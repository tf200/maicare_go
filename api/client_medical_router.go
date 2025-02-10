package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientMedicalRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientMedical := baseRouter.Group("/clients")
	ClientMedical.Use(AuthMiddleware(server.tokenMaker))
	{
		ClientMedical.POST("/:id/allergies", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientAllergyApi)
		ClientMedical.GET("/:id/allergies", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientAllergiesApi)
		ClientMedical.GET("/:id/allergies/:allergy_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientAllergyApi)
		ClientMedical.PUT("/:id/allergies/:allergy_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateClientAllergyApi)
		ClientMedical.DELETE("/:id/allergies/:allergy_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteClientAllergyApi)

		ClientMedical.POST("/:id/diagnosis", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientDiagnosisApi)
		ClientMedical.GET("/:id/diagnosis", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientDiagnosesApi)
		ClientMedical.GET("/:id/diagnosis/:diagnosis_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientDiagnosisApi)
		ClientMedical.PUT("/:id/diagnosis/:diagnosis_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateClientDiagnosisApi)
		ClientMedical.DELETE("/:id/diagnosis/:diagnosis_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteClientDiagnosisApi)

		ClientMedical.POST("/:id/medications", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientMedicationApi)
		ClientMedical.GET("/:id/medications", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientMedicationsApi)
		ClientMedical.GET("/:id/medications/:medication_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientMedicationApi)
		ClientMedical.PUT("/:id/medications/:medication_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateClientMedicationApi)
		ClientMedical.DELETE("/:id/medications/:medication_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteClientMedicationApi)

	}
}
