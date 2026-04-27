package handler

import "github.com/gin-gonic/gin"

func RegisterWebSocketRoutes(
	rg *gin.RouterGroup,
	handler *WebSocketHandler,
	auth gin.HandlerFunc,
) {
	// GET /ws — WebSocket upgrade endpoint (authenticated via ticket, no middleware)
	rg.GET("/ws", handler.HandleWebSocket)

	// POST /auth/ws-ticket — Create a one-time WebSocket ticket
	rg.POST("/auth/ws-ticket", auth, handler.CreateWebSocketTicket)
}
