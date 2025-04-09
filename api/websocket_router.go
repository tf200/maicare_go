package api

import (
	"github.com/gin-gonic/gin"
)

// setupWebsocketRoutes defines the routes related to WebSocket connections.
func (server *Server) setupWebsocketRoutes(router *gin.RouterGroup) {
	// Apply authentication middleware to the WebSocket endpoint
	wsGroup := router.Group("/ws")                 // You can choose a different base path if needed
	wsGroup.Use(AuthMiddleware(server.tokenMaker)) // Use your existing middleware
	{
		// Handler for upgrading the connection
		wsGroup.GET("", server.handleWebSocket)
	}
}
