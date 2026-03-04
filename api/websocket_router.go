package api

import (
	"github.com/gin-gonic/gin"
)

// setupWebsocketRoutes defines the routes related to WebSocket connections.
func (server *Server) setupWebsocketRoutes(router *gin.RouterGroup) {
	wsGroup := router.Group("/ws")
	{
		// Handler for upgrading the connection
		wsGroup.GET("", server.handleWebSocket)
	}
}
