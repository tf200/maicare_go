package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAppointmentRoutes(baseRouter *gin.RouterGroup) {
	appointmentsRouter := baseRouter.Group("/appointments").Use(server.AuthMiddleware())
	{
		appointmentsRouter.POST("", server.RBACMiddleware("APPOINTMENT.CREATE"), server.CreateAppointmentApi)
		appointmentsRouter.GET("/:id", server.RBACMiddleware("APPOINTMENT.VIEW"), server.GetAppointmentApi)
		appointmentsRouter.PUT("/:id", server.RBACMiddleware("APPOINTMENT.UPDATE"), server.UpdateAppointmentApi)
		appointmentsRouter.DELETE("/:id", server.RBACMiddleware("APPOINTMENT.DELETE"), server.DeleteAppointmentApi)

		appointmentsRouter.POST("/:id/participants", server.RBACMiddleware("APPOINTMENT.UPDATE"), server.AddParticipantToAppointmentApi)
		appointmentsRouter.POST("/:id/clients", server.RBACMiddleware("APPOINTMENT.UPDATE"), server.AddClientToAppointmentApi)

		appointmentsRouter.POST("/:id/confirm", server.RBACMiddleware("APPOINTMENT.CONFIRM"), server.ConfirmAppointmentApi)
	}
}
