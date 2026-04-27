package handler

import "github.com/gin-gonic/gin"

func RegisterNotificationRoutes(
	rg *gin.RouterGroup,
	handler *NotificationHandler,
	auth gin.HandlerFunc,
) {
	rg.GET("/notifications", auth, handler.ListNotifications)
	rg.POST("/notifications/:id/read", auth, handler.MarkNotificationAsRead)
}
