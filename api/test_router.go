package api

import "github.com/gin-gonic/gin"

// setupTestRoutes configures all test-related routes
func (server *Server) setupTestRoutes(baseRouter *gin.RouterGroup) {
	testGroup := baseRouter.Group("/test")
	{
		// Basic health check endpoint
		testGroup.GET("/health", RBACMiddleware(server.store, "TEST_VIEW"), server.handleHealthCheck)

		// Echo endpoint to test request handling
		testGroup.POST("/echo", server.handleEcho)

		// Simulated latency endpoint for testing timeouts
		testGroup.GET("/latency/:ms", server.handleLatency)

		testGroup.GET("/send-email", RBACMiddleware(server.store, "TEST_VIEW"), server.EmailAndAsynq)

	}
}
