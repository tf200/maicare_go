package api

import "github.com/gin-gonic/gin"

func (server *Server) setupScheduleRoutes(baseRouter *gin.RouterGroup) {
	schedulesRouter := baseRouter.Group("/schedules")
	schedulesRouter.Use(server.AuthMiddleware())
	{
		schedulesRouter.POST("", server.RBACMiddleware("SCHEDULE.CREATE"), server.CreateScheduleApi)
		schedulesRouter.DELETE("/:id", server.RBACMiddleware("SCHEDULE.DELETE"), server.DeleteScheduleApi)
		schedulesRouter.GET("/:id", server.RBACMiddleware("SCHEDULE.VIEW"), server.GetScheduleByIDApi)
		schedulesRouter.PUT("/:id", server.RBACMiddleware("SCHEDULE.UPDATE"), server.UpdateScheduleApi)
		schedulesRouter.POST("/auto_generate", server.RBACMiddleware("SCHEDULE.CREATE"), server.AutoGenerateSchedulesApi)
		schedulesRouter.POST("/save_generated", server.RBACMiddleware("SCHEDULE.CREATE"), server.SaveGeneratedSchedulesApi)
	}

	locationsRouter := baseRouter.Group("/locations")
	locationsRouter.Use(server.AuthMiddleware())
	{
		locationsRouter.GET("/:id/schedules", server.RBACMiddleware("SCHEDULE.VIEW"), server.GetSchedulesByLocationInRangeApi)
	}
}
