package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAppointmentRoutes(baseRouter *gin.RouterGroup) {
	eventsRouter := baseRouter.Group("/events").Use(server.AuthMiddleware())
	{
		eventsRouter.POST("", server.RBACMiddleware("APPOINTMENT.CREATE"), server.CreateEventApi)
		eventsRouter.POST("/list", server.RBACMiddleware("APPOINTMENT.VIEW"), server.ListEventsApi)
		eventsRouter.GET("/:id", server.RBACMiddleware("APPOINTMENT.VIEW"), server.GetEventApi)
		eventsRouter.PATCH("/:id", server.RBACMiddleware("APPOINTMENT.UPDATE"), server.UpdateEventApi)
		eventsRouter.DELETE("/:id", server.RBACMiddleware("APPOINTMENT.DELETE"), server.DeleteEventApi)
	}
}
