package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientNetworkRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientNetwork := baseRouter.Group("/clients")
	ClientNetwork.Use(server.AuthMiddleware())
	{
		ClientNetwork.GET("/:id/sender", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientSenderApi)

		ClientNetwork.POST("/:id/emergency_contacts", server.RBACMiddleware("CLIENT.CREATE"), server.CreateClientEmergencyContactApi)
		ClientNetwork.GET("/:id/emergency_contacts", server.RBACMiddleware("CLIENT.VIEW"), server.ListClientEmergencyContactsApi)
		ClientNetwork.GET("/:id/emergency_contacts/:contact_id", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientEmergencyContactApi)
		ClientNetwork.PUT("/:id/emergency_contacts/:contact_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateClientEmergencyContactApi)
		ClientNetwork.DELETE("/:id/emergency_contacts/:contact_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteClientEmergencyContactApi)

		ClientNetwork.POST("/:id/involved_employees", server.RBACMiddleware("CLIENT.CREATE"), server.AssignEmployeeApi)
		ClientNetwork.GET("/:id/involved_employees", server.RBACMiddleware("CLIENT.VIEW"), server.ListAssignedEmployeesApi)
		ClientNetwork.GET("/:id/involved_employees/:assign_id", server.RBACMiddleware("CLIENT.VIEW"), server.GetAssignedEmployeeApi)
		ClientNetwork.PUT("/:id/involved_employees/:assign_id", server.RBACMiddleware("CLIENT.UPDATE"), server.UpdateAssignedEmployeeApi)
		ClientNetwork.DELETE("/:id/involved_employees/:assign_id", server.RBACMiddleware("CLIENT.DELETE"), server.DeleteAssignedEmployeeApi)

		ClientNetwork.GET("/:id/related_emails", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientRelatedEmailsApi)
	}

}
