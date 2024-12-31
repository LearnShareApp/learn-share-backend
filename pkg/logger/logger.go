package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new logger instance
func New(level string, development bool) (*zap.Logger, error) {
	var config zap.Config

	if development {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		config.DisableCaller = false
		config.DisableStacktrace = true
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.DisableCaller = false
		config.DisableStacktrace = false
	}

	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	config.Level = zap.NewAtomicLevelAt(logLevel)

	return config.Build()
}

// NewDefault creates a default production logger
func NewDefault() *zap.Logger {
	logger, _ := New("info", false)
	return logger
}

// NewDevelopment creates a default development logger with pretty console output
func NewDevelopment() *zap.Logger {
	logger, _ := New("debug", true)
	return logger
}
