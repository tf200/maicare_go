package api

import "github.com/gin-gonic/gin"

// setupTestRoutes configures all test-related routes
func (server *Server) setupTestRoutes(baseRouter *gin.RouterGroup) {
	testGroup := baseRouter.Group("/test")
	{
		// Basic health check endpoint
		testGroup.GET("/health", AuthMiddleware(server.tokenMaker), RBACMiddleware(server.store, "CLIENT.VIEW"), server.handleHealthCheck)

		// Echo endpoint to test request handling
		testGroup.POST("/echo", server.handleEcho)

		// Simulated latency endpoint for testing timeouts
		testGroup.GET("/latency/:ms", server.handleLatency)

		testGroup.GET("/send-email", server.EmailAndAsynq)

	}
}
