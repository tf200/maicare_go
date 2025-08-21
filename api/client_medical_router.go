package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientMedicalRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientMedical := baseRouter.Group("/clients")
	ClientMedical.Use(server.AuthMiddleware())
	{

		ClientMedical.POST("/:id/diagnosis", server.RBACMiddleware("CLIENT.CREATE"), server.CreateClientDiagnosisApi)
		ClientMedical.GET("/:id/diagnosis", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientDiagnosesApi)
		ClientMedical.GET("/:id/diagnosis/:diagnosis_id", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientDiagnosisApi)
		ClientMedical.PUT("/:id/diagnosis/:diagnosis_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateClientDiagnosisApi)
		ClientMedical.DELETE("/:id/diagnosis/:diagnosis_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteClientDiagnosisApi)

		ClientMedical.POST("/:id/diagnosis/:diagnosis_id/medications", server.RBACMiddleware("CLIENT.CREATE"), server.CreateClientMedicationApi)
		ClientMedical.GET("/:id/diagnosis/:diagnosis_id/medications", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientMedicationsApi)
		ClientMedical.GET("/:id/diagnosis/:diagnosis_id/medications/:medication_id", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientMedicationApi)
		ClientMedical.PUT("/:id/diagnosis/:diagnosis_id/medications/:medication_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateClientMedicationApi)
		ClientMedical.DELETE("/:id/medications/:medication_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteClientMedicationApi)

	}
}
