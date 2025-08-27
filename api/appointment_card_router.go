package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAppointmentCardRoutes(baseRouter *gin.RouterGroup) {
	appointmentCardRouter := baseRouter.Group("/clients").Use(server.AuthMiddleware())
	{
		appointmentCardRouter.POST("/:id/appointment_cards", server.RBACMiddleware("APPOINTMENT_CARD.CREATE"), server.CreateAppointmentCardApi)
		appointmentCardRouter.GET("/:id/appointment_cards", server.RBACMiddleware("APPOINTMENT_CARD.VIEW"), server.GetAppointmentCardApi)
		appointmentCardRouter.PUT("/:id/appointment_cards", server.RBACMiddleware("APPOINTMENT_CARD.UPDATE"), server.UpdateAppointmentCardApi)
		appointmentCardRouter.POST("/:id/appointment_cards/generate_document", server.RBACMiddleware("APPOINTMENT_CARD.GENERATE_DOCUMENT"), server.GenerateAppointmentCardDocumentApi)
	}
}
