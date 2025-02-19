package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAppointmentCardRoutes(baseRouter *gin.RouterGroup) {
	appointmentCardRouter := baseRouter.Group("/clients").Use(AuthMiddleware(server.tokenMaker))
	{
		appointmentCardRouter.POST("/:id/appointment_cards", server.CreateAppointmentCardApi)
		appointmentCardRouter.GET("/:id/appointment_cards", server.GetAppointmentCardApi)
		appointmentCardRouter.PUT("/:id/appointment_cards", server.UpdateAppointmentCardApi)
		appointmentCardRouter.POST("/:id/appointment_cards/generate_document", server.GenerateAppointmentCardDocumentApi)
	}
}
