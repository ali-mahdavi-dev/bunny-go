package logging

import "context"

// Logger is the main logging interface with builder pattern
type Logger interface {
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
	
	// Direct logging methods (for simple cases)
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	
	// Formatted logging methods
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	
	// Context support
	WithContext(ctx context.Context) Logger
}

// LoggerAdapter is the interface that different logger implementations must satisfy
type LoggerAdapter interface {
	// Log at different levels with message and fields
	Log(level LogLevel, msg string, fields map[string]interface{})
	
	// Formatted logging
	Logf(level LogLevel, template string, args ...interface{})
}
