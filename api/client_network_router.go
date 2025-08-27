package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientNetworkRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientNetwork := baseRouter.Group("/clients")
	ClientNetwork.Use(server.AuthMiddleware())
	{
		ClientNetwork.GET("/:id/sender", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientSenderApi)

		ClientNetwork.POST("/:id/emergency_contacts", server.RBACMiddleware("CLIENT.EMERGENCY_CONTACT.CREATE"), server.CreateClientEmergencyContactApi)
		ClientNetwork.GET("/:id/emergency_contacts", server.RBACMiddleware("CLIENT.EMERGENCY_CONTACT.VIEW"), server.ListClientEmergencyContactsApi)
		ClientNetwork.GET("/:id/emergency_contacts/:contact_id", server.RBACMiddleware("CLIENT.EMERGENCY_CONTACT.VIEW"), server.GetClientEmergencyContactApi)
		ClientNetwork.PUT("/:id/emergency_contacts/:contact_id", server.RBACMiddleware("CLIENT.EMERGENCY_CONTACT.UPDATE"), server.UpdateClientEmergencyContactApi)
		ClientNetwork.DELETE("/:id/emergency_contacts/:contact_id", server.RBACMiddleware("CLIENT.EMERGENCY_CONTACT.DELETE"), server.DeleteClientEmergencyContactApi)

		ClientNetwork.POST("/:id/involved_employees", server.RBACMiddleware("CLIENT.INVOLVED_EMPLOYEE.CREATE"), server.AssignEmployeeApi)
		ClientNetwork.GET("/:id/involved_employees", server.RBACMiddleware("CLIENT.INVOLVED_EMPLOYEE.VIEW"), server.ListAssignedEmployeesApi)
		ClientNetwork.GET("/:id/involved_employees/:assign_id", server.RBACMiddleware("CLIENT.INVOLVED_EMPLOYEE.VIEW"), server.GetAssignedEmployeeApi)
		ClientNetwork.PUT("/:id/involved_employees/:assign_id", server.RBACMiddleware("CLIENT.INVOLVED_EMPLOYEE.UPDATE"), server.UpdateAssignedEmployeeApi)
		ClientNetwork.DELETE("/:id/involved_employees/:assign_id", server.RBACMiddleware("CLIENT.INVOLVED_EMPLOYEE.DELETE"), server.DeleteAssignedEmployeeApi)

		ClientNetwork.GET("/:id/related_emails", server.RBACMiddleware("CLIENT.VIEW"), server.GetClientRelatedEmailsApi)
	}

}
