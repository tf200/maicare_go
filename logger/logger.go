package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//go:generate mockgen -source=logger.go -destination=../mocks/mock_logger.go -package=mocks
type Logger interface {
	LogBusinessEvent(level LogLevel, operation, message string, fields ...zap.Field)
}

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LoggerImpl struct {
	logger *zap.Logger
}

func SetupLogger(environment string) (Logger, error) {
	var config zap.Config
	if environment == "production" {
		config = zap.NewProductionConfig()
		config.DisableCaller = true
		config.DisableStacktrace = true
		config.OutputPaths = []string{"stdout"}
	} else {
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout"}
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	if environment == "production" {
		logDir := "/var/log/maicare"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}

		// All logs file
		allLogsWriter := &lumberjack.Logger{
			Filename:   filepath.Join(logDir, "app.log"),
			MaxSize:    100, // MB
			MaxBackups: 3,
			MaxAge:     7, // days - shorter for all logs
			Compress:   true,
		}

		// Error logs file (separate)
		errorLogsWriter := &lumberjack.Logger{
			Filename:   filepath.Join(logDir, "error.log"),
			MaxSize:    50, // MB - smaller, errors are less frequent
			MaxBackups: 10, // keep more error log backups
			MaxAge:     30, // days - keep errors longer
			Compress:   true,
		}

		// Create cores with different level filters
		allLogsCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(allLogsWriter),
			zap.InfoLevel, // All levels >= Info
		)

		errorLogsCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(errorLogsWriter),
			zap.ErrorLevel, // Only Error and above
		)

		// Combine: stdout + all.log + error.log
		core := zapcore.NewTee(
			logger.Core(), // stdout
			allLogsCore,   // app.log (all logs)
			errorLogsCore, // error.log (errors only)
		)

		logger = zap.New(core)
	}

	return &LoggerImpl{logger: logger}, nil
}

func (l *LoggerImpl) LogBusinessEvent(level LogLevel, operation, message string, fields ...zap.Field) {
	// Add operation context to all business logs
	commonFields := []zap.Field{
		zap.String("operation", operation),
		zap.String("service", "maicare-api"),
		zap.Int64("timestamp", time.Now().Unix()),
	}

	allFields := append(commonFields, fields...)

	switch level {
	case LogLevelError:
		l.logger.Error(message, allFields...)
	case LogLevelWarn:
		l.logger.Warn(message, allFields...)
	case LogLevelInfo:
		l.logger.Info(message, allFields...)
	}
}
