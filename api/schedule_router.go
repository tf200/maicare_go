package api

import "github.com/gin-gonic/gin"

func (server *Server) setupScheduleRoutes(baseRouter *gin.RouterGroup) {
	schedule := baseRouter.Group("")
	schedule.Use(AuthMiddleware(server.tokenMaker))
	{
		schedule.POST("/schedules", RBACMiddleware(server.store, "SCHEDULE.CREATE"), server.CreateScheduleApi)

		schedule.GET("/locations/:id/monthly_schedules", RBACMiddleware(server.store, "SCHEDULE.VIEW"), server.GetMonthlySchedulesByLocationApi)
		schedule.GET("/locations/:id/daily_schedules", RBACMiddleware(server.store, "SCHEDULE.VIEW"), server.GetDailySchedulesByLocationApi)

		schedule.DELETE("/schedules/:id", RBACMiddleware(server.store, "SCHEDULE.DELETE"), server.DeleteScheduleApi)
		schedule.GET("/schedules/:id", RBACMiddleware(server.store, "SCHEDULE.VIEW"), server.GetScheduleByIDApi)
		schedule.PUT("/schedules/:id", RBACMiddleware(server.store, "SCHEDULE.UPDATE"), server.UpdateScheduleApi)
	}
}
