package api

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

func (server *Server) requestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Generate request ID if not present
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx.Set("request_id", requestID)
		ctx.Header("X-Request-ID", requestID)

		startTime := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()

		duration := time.Since(startTime)
		statusCode := ctx.Writer.Status()

		// Determine log level based on status code
		logLevel := server.getLogLevel(statusCode)

		// Build base fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", duration),
			zap.Int64("latency_ms", duration.Milliseconds()),
			zap.String("client_ip", ctx.ClientIP()),
			zap.Int("response_size", ctx.Writer.Size()),
		}

		// Add query string if present (be careful with sensitive data)
		if raw != "" && server.config.Environment != "production" {
			fields = append(fields, zap.String("query", raw))
		}

		// Add error information if present
		if len(ctx.Errors) > 0 {
			fields = append(fields, zap.String("error", ctx.Errors.String()))
		}

		// Add user context if authenticated
		if userID, exists := ctx.Get("user_id"); exists {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		// Log based on environment and status
		shouldLog := server.shouldLogRequest(statusCode)
		if shouldLog {
			message := "HTTP Request"
			switch logLevel {
			case zapcore.ErrorLevel:
				server.logger.Error(message, fields...)
			case zapcore.WarnLevel:
				server.logger.Warn(message, fields...)
			default:
				server.logger.Info(message, fields...)
			}
		}

		// Always log slow requests (>1s) even in production
		if duration > time.Second {
			fields = append(fields, zap.Bool("slow_request", true))
			server.logger.Warn("Slow HTTP Request", fields...)
		}
	}
}

func (server *Server) shouldLogRequest(statusCode int) bool {
	// In production: only log errors (4xx, 5xx)
	if server.config.Environment == "production" {
		return statusCode >= 400
	}
	// In development: log everything
	return true
}

func (server *Server) getLogLevel(statusCode int) zapcore.Level {
	switch {
	case statusCode >= 500:
		return zapcore.ErrorLevel
	case statusCode >= 400:
		return zapcore.WarnLevel
	default:
		return zapcore.InfoLevel
	}
}

// Add recovery middleware with proper logging
func (server *Server) recoveryLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := ctx.Get("request_id")

				server.logger.Error("Panic recovered",
					zap.String("request_id", requestID.(string)),
					zap.String("method", ctx.Request.Method),
					zap.String("path", ctx.Request.URL.Path),
					zap.Any("panic", err),
					zap.String("stack", string(debug.Stack())),
				)

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":      "Internal server error",
					"request_id": requestID,
				})
			}
		}()
		ctx.Next()
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
