package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

func (server *Server) requestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Process request
		ctx.Next()

		duration := time.Since(startTime)

		// Only log errors or in development
		shouldLog := server.config.Environment != "production" || ctx.Writer.Status() >= 400

		if shouldLog {
			server.logger.Info("HTTP Request",
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.Request.URL.Path),
				zap.Int("status", ctx.Writer.Status()),
				zap.Duration("latency", duration),
				zap.String("client_ip", ctx.ClientIP()),
				zap.String("user_agent", ctx.Request.UserAgent()),
			)
		}
	}
}

// ALTERNATIVE: Optimized logging helper for handlers
func (server *Server) logBusinessEvent(level LogLevel, operation, message string, fields ...zap.Field) {
	// Add operation context to all business logs
	commonFields := []zap.Field{
		zap.String("operation", operation),
		zap.String("service", "maicare-api"),
		zap.Int64("timestamp", time.Now().Unix()),
	}

	allFields := append(commonFields, fields...)

	switch level {
	case LogLevelError:
		server.logger.Error(message, allFields...)
	case LogLevelWarn:
		server.logger.Warn(message, allFields...)
	case LogLevelInfo:
		server.logger.Info(message, allFields...)
	}
}
