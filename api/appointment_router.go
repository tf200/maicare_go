package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAppointmentRoutes(baseRouter *gin.RouterGroup) {
	appointmentsRouter := baseRouter.Group("/appointments").Use(server.AuthMiddleware())
	{
		appointmentsRouter.POST("", server.CreateAppointmentApi)
		appointmentsRouter.GET("/:id", server.GetAppointmentApi)
		appointmentsRouter.PUT("/:id", server.UpdateAppointmentApi)
		appointmentsRouter.DELETE("/:id", server.DeleteAppointmentApi)

		appointmentsRouter.POST("/:id/participants", server.AddParticipantToAppointmentApi)
		appointmentsRouter.POST("/:id/clients", server.AddClientToAppointmentApi)

		appointmentsRouter.POST("/:id/confirm", server.ConfirmAppointmentApi)
	}
}
