package api

import "github.com/gin-gonic/gin"

func (server *Server) setupScheduleRoutes(baseRouter *gin.RouterGroup) {
	schedule := baseRouter.Group("")
	schedule.Use(server.AuthMiddleware())
	{
		schedule.POST("/schedules", server.RBACMiddleware("SCHEDULE.CREATE"), server.CreateScheduleApi)

		schedule.GET("/locations/:id/schedules", server.RBACMiddleware("SCHEDULE.VIEW"), server.GetSchedulesByLocationInRangeApi)

		schedule.DELETE("/schedules/:id", server.RBACMiddleware("SCHEDULE.DELETE"), server.DeleteScheduleApi)
		schedule.GET("/schedules/:id", server.RBACMiddleware("SCHEDULE.VIEW"), server.GetScheduleByIDApi)
		schedule.PUT("/schedules/:id", server.RBACMiddleware("SCHEDULE.UPDATE"), server.UpdateScheduleApi)

		schedule.POST("/schedules/auto_generate", server.RBACMiddleware("SCHEDULE.CREATE"), server.AutoGenerateSchedulesApi)
		schedule.POST("/schedules/save_generated", server.RBACMiddleware("SCHEDULE.CREATE"), server.SaveGeneratedSchedulesApi)

	}
}
