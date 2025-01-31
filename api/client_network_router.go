package api

import "github.com/gin-gonic/gin"

func (server *Server) setupClientNetworkRoutes(baseRouter *gin.RouterGroup) {
	// Routes under /clients prefix
	ClientNetwork := baseRouter.Group("/clients")
	ClientNetwork.Use(AuthMiddleware(server.tokenMaker))
	{
		ClientNetwork.POST("/:id/emergency_contacts", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientEmergencyContactApi)
		ClientNetwork.GET("/:id/emergency_contacts", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListClientEmergencyContactsApi)
		ClientNetwork.GET("/:id/emergency_contacts/:contact_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetClientEmergencyContactApi)
		ClientNetwork.PUT("/:id/emergency_contacts/:contact_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateClientEmergencyContactApi)
		ClientNetwork.DELETE("/:id/emergency_contacts/:contact_id", RBACMiddleware(server.store, "CLIENT.DELETE"), server.DeleteClientEmergencyContactApi)

		ClientNetwork.POST("/:id/assigned_employees", RBACMiddleware(server.store, "CLIENT.CREATE"), server.AssignEmployeeApi)
		ClientNetwork.GET("/:id/assigned_employees", RBACMiddleware(server.store, "CLIENT.VIEW"), server.ListAssignedEmployeesApi)
		ClientNetwork.GET("/:id/assigned_employees/:assign_id", RBACMiddleware(server.store, "CLIENT.VIEW"), server.GetAssignedEmployeeApi)
		ClientNetwork.PUT("/:id/assigned_employees/:assign_id", RBACMiddleware(server.store, "CLIENT.UPDATE"), server.UpdateAssignedEmployeeApi)
	}

}
