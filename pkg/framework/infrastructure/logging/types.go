package logging

import (
	"io"
	"os"
)

// LogLevel represents logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// LogFormat represents log output format
type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

// LoggerType represents different logger implementations
type LoggerType string

const (
	LoggerTypeZerolog LoggerType = "zerolog"
	LoggerTypeZap     LoggerType = "zap"
	LoggerTypeLogrus  LoggerType = "logrus"
)

// LoggerConfig holds configuration for logger
type LoggerConfig struct {
	Type      LoggerType
	Level     LogLevel
	Output    io.Writer
	Format    LogFormat
	AddCaller bool
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Type:      LoggerTypeZerolog,
		Level:     LogLevelInfo,
		Output:    os.Stdout,
		Format:    LogFormatJSON,
		AddCaller: false,
	}
}
