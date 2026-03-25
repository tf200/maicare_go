package api

import "github.com/gin-gonic/gin"

func (server *Server) setupLateArrivalRoutes(baseRouter *gin.RouterGroup) {
	lateArrivalsRouter := baseRouter.Group("/late-arrivals")
	lateArrivalsRouter.Use(server.AuthMiddleware())
	{
		lateArrivalsRouter.POST("", server.RBACMiddleware("LATE_ARRIVAL.CREATE"), server.CreateLateArrivalApi)
		lateArrivalsRouter.POST("/admin", server.RBACMiddleware("LATE_ARRIVAL.CREATE_ALL"), server.CreateLateArrivalByAdminApi)
		lateArrivalsRouter.GET("/my", server.RBACMiddleware("LATE_ARRIVAL.VIEW"), server.ListMyLateArrivalsApi)
		lateArrivalsRouter.GET("", server.RBACMiddleware("LATE_ARRIVAL.VIEW_ALL"), server.ListLateArrivalsApi)
	}
}
