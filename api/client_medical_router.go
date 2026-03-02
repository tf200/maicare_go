package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientMedicalRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientMedical := baseRouter.Group("/clients")
	ClientMedical.Use(server.AuthMiddleware())
	{
		// Overview
		ClientMedical.GET("/:id/medical/overview", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientMedicalOverviewApi)

		// Diagnoses
		ClientMedical.POST("/:id/medical/diagnoses", server.RBACMiddleware("CLIENT.DIAGNOSIS.CREATE"), server.CreateClientDiagnosisApi)
		ClientMedical.GET("/:id/medical/diagnoses", server.RBACMiddleware("CLIENT.DIAGNOSIS.VIEW"), server.ListClientDiagnosesApi)
		ClientMedical.GET("/:id/medical/diagnoses/:diagnosis_id", server.RBACMiddleware("CLIENT.DIAGNOSIS.VIEW"), server.GetClientDiagnosisApi)
		ClientMedical.PUT("/:id/medical/diagnoses/:diagnosis_id", server.RBACMiddleware("CLIENT.DIAGNOSIS.UPDATE"), server.UpdateClientDiagnosisApi)
		ClientMedical.DELETE("/:id/medical/diagnoses/:diagnosis_id", server.RBACMiddleware("CLIENT.DIAGNOSIS.DELETE"), server.DeleteClientDiagnosisApi)

		// Medication orders (RBAC uses existing CLIENT.MEDICATION.* permissions)
		ClientMedical.POST("/:id/medical/medication-orders", server.RBACMiddleware("CLIENT.MEDICATION.CREATE"), server.CreateClientMedicationOrderApi)
		ClientMedical.GET("/:id/medical/medication-orders", server.RBACMiddleware("CLIENT.MEDICATION.VIEW"), server.ListClientMedicationOrdersApi)
		ClientMedical.GET("/:id/medical/medication-orders/:order_id", server.RBACMiddleware("CLIENT.MEDICATION.VIEW"), server.GetClientMedicationOrderApi)
		ClientMedical.PUT("/:id/medical/medication-orders/:order_id", server.RBACMiddleware("CLIENT.MEDICATION.UPDATE"), server.UpdateClientMedicationOrderApi)
		ClientMedical.DELETE("/:id/medical/medication-orders/:order_id", server.RBACMiddleware("CLIENT.MEDICATION.DELETE"), server.DeleteClientMedicationOrderApi)

	}
}
