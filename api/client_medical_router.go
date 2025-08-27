package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientMedicalRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientMedical := baseRouter.Group("/clients")
	ClientMedical.Use(server.AuthMiddleware())
	{

		ClientMedical.POST("/:id/diagnosis", server.RBACMiddleware("CLIENT.DIAGNOSIS.CREATE"), server.CreateClientDiagnosisApi)
		ClientMedical.GET("/:id/diagnosis", server.RBACMiddleware("CLIENT.DIAGNOSIS.VIEW"), server.ListClientDiagnosesApi)
		ClientMedical.GET("/:id/diagnosis/:diagnosis_id", server.RBACMiddleware("CLIENT.DIAGNOSIS.VIEW"), server.GetClientDiagnosisApi)
		ClientMedical.PUT("/:id/diagnosis/:diagnosis_id", server.RBACMiddleware("CLIENT.DIAGNOSIS.UPDATE"), server.UpdateClientDiagnosisApi)
		ClientMedical.DELETE("/:id/diagnosis/:diagnosis_id", server.RBACMiddleware("CLIENT.DIAGNOSIS.DELETE"), server.DeleteClientDiagnosisApi)

		ClientMedical.POST("/:id/diagnosis/:diagnosis_id/medications", server.RBACMiddleware("CLIENT.MEDICATION.CREATE"), server.CreateClientMedicationApi)
		ClientMedical.GET("/:id/diagnosis/:diagnosis_id/medications", server.RBACMiddleware("CLIENT.MEDICATION.VIEW"), server.ListClientMedicationsApi)
		ClientMedical.GET("/:id/diagnosis/:diagnosis_id/medications/:medication_id", server.RBACMiddleware("CLIENT.MEDICATION.VIEW"), server.GetClientMedicationApi)
		ClientMedical.PUT("/:id/diagnosis/:diagnosis_id/medications/:medication_id", server.RBACMiddleware("CLIENT.MEDICATION.UPDATE"), server.UpdateClientMedicationApi)
		ClientMedical.DELETE("/:id/medications/:medication_id", server.RBACMiddleware("CLIENT.MEDICATION.DELETE"), server.DeleteClientMedicationApi)

	}
}
