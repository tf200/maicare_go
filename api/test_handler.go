package api

import (
	"maicare_go/tasks"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TestResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// handleHealthCheck returns basic service health information
func (server *Server) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, TestResponse{
		Status:    "healthy",
		Message:   "service is running",
		Timestamp: time.Now(),
	})
}

// handleEcho returns whatever JSON payload was sent
func (server *Server) handleEcho(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, TestResponse{
			Status:    "error",
			Message:   "invalid JSON payload",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"echo":      requestBody,
		"timestamp": time.Now(),
	})
}

// handleLatency simulates a delayed response
func (server *Server) handleLatency(c *gin.Context) {
	ms := c.Param("ms")
	delay, err := time.ParseDuration(ms + "ms")
	if err != nil {
		c.JSON(http.StatusBadRequest, TestResponse{
			Status:    "error",
			Message:   "invalid latency value",
			Timestamp: time.Now(),
		})
		return
	}

	time.Sleep(delay)
	c.JSON(http.StatusOK, TestResponse{
		Status:    "success",
		Message:   "delayed response",
		Timestamp: time.Now(),
	})
}

func (server *Server) EmailAndAsynq(c *gin.Context) {
	server.asynqClient.EnqueueEmailDelivery(tasks.EmailDeliveryPayload{
		To:           "farjiataha@gmail.com",
		UserEmail:    "farjiataha@gmail.com",
		UserPassword: "password",
	}, c)

	c.JSON(http.StatusOK, gin.H{
		"echo":      "Email sent",
		"timestamp": time.Now(),
	})
}
