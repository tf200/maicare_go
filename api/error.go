package api

import "github.com/gin-gonic/gin"

// ErrorResponse represents the standard error response
type ErrorResponse struct {
    Error string `json:"error" example:"error message"`
}

func errorResponse(err error) gin.H {
    return gin.H{
        "error": err.Error(),
    }
}