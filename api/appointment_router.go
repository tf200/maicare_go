package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAppointmentRoutes(baseRouter *gin.RouterGroup) {
	appointmentsRouter := baseRouter.Group("/appointments").Use(AuthMiddleware(server.tokenMaker))
	{
		appointmentsRouter.POST("", server.CreateAppointmentApi)
		appointmentsRouter.POST("/:appointment_id/participants", server.AddParticipantToAppointmentApi)
		appointmentsRouter.POST("/:appointment_id/clients", server.AddClientToAppointmentApi)

		appointmentsRouter.POST("/employee_list", server.ListAppointmentsForEmployee)

	}
}
