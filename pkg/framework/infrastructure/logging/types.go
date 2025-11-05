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

// FieldType represents the type of a log field
type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInt     FieldType = "int"
	FieldTypeInt64   FieldType = "int64"
	FieldTypeFloat64 FieldType = "float64"
	FieldTypeBool    FieldType = "bool"
	FieldTypeError   FieldType = "error"
	FieldTypeAny     FieldType = "any"
)

// LogField represents a log field with its type
type LogField struct {
	Key   string
	Value interface{}
	Type  FieldType
}

// LoggerAdapter is the interface that different logger implementations must satisfy
type LoggerAdapter interface {
	// Log at different levels with message and fields
	Log(level LogLevel, msg string, fields []LogField)

	// Formatted logging
	Logf(level LogLevel, template string, args ...interface{})
}

// Logger is the main logging interface with builder pattern
type Logger interface {
	// Create log entries with different levels (builder pattern)
	Debug(msg string) Logger
	Info(msg string) Logger
	Warn(msg string) Logger
	Error(msg string) Logger
	Fatal(msg string) Logger

	// Builder methods - return Logger for chaining
	WithAny(key string, value interface{}) Logger
	WithString(key, value string) Logger
	WithInt(key string, value int) Logger
	WithInt64(key string, value int64) Logger
	WithFloat64(key string, value float64) Logger
	WithBool(key string, value bool) Logger
	WithError(err error) Logger
	WithFields(fields map[string]interface{}) Logger

	// Log method - writes the accumulated log entry
	Log()

	// Formatted logging methods
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}
