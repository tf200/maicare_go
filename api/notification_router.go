package api

import "github.com/gin-gonic/gin"

func (server *Server) setupNotificationRoutes(baseRouter *gin.RouterGroup) {
	baseRouter.GET("/notifications", server.AuthMiddleware(), server.ListNotificationsApi)
	baseRouter.POST("/notifications/:id/read", server.AuthMiddleware(), server.MarkNotificationAsReadApi)
}
