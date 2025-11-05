package logging

import "context"

// Logger is the main interface for logging in the application
type Logger interface {
	Init()

	// Fluent API methods - returns Entry for chaining
	Debug(msg string) *Entry
	Info(msg string) *Entry
	Warn(msg string) *Entry
	Error(msg string) *Entry
	Fatal(msg string) *Entry

	// Formatted methods
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	// Context and fields
	WithContext(ctx context.Context) Logger
	WithFields(fields map[string]interface{}) Logger
}
