package logger

import (
	"fmt"
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

		config.OutputPaths = []string{
			"stdout",
			"/var/log/maicare/app.log",
		}
	} else {
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{
			"stdout",
		}
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	if environment == "production" {
		fileWritter := &lumberjack.Logger{
			Filename:   "/var/log/maicare/app.log",
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   // days
			Compress:   true, // compress log files
		}

		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(config.EncoderConfig),
				zapcore.AddSync(fileWritter),
				zap.InfoLevel,
			),
			logger.Core(),
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
