package logging

import (
	"context"
	"fmt"
	"sync"
)

// loggerImpl is the default implementation of Logger interface with builder pattern
type loggerImpl struct {
	adapter LoggerAdapter
	config  LoggerConfig
	mu      sync.RWMutex
	fields  map[string]interface{} // accumulated fields from builder
	level   LogLevel               // log level for builder
	msg     string                 // message for builder
}

var (
	globalLogger Logger
	globalOnce   sync.Once
)

// NewLogger creates a new logger instance
func NewLogger(config LoggerConfig) (Logger, error) {
	var adapter LoggerAdapter
	var err error

	switch config.Type {
	case LoggerTypeZerolog:
		adapter, err = NewZerologAdapter(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create zerolog adapter: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported logger type: %s", config.Type)
	}

	return &loggerImpl{
		adapter: adapter,
		config:  config,
		fields:  make(map[string]interface{}),
	}, nil
}

// Init initializes the global logger instance
func Init(config LoggerConfig) error {
	var err error
	globalOnce.Do(func() {
		globalLogger, err = NewLogger(config)
	})
	return err
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		globalOnce.Do(func() {
			globalLogger, _ = NewLogger(DefaultLoggerConfig())
		})
	}
	return globalLogger
}

// SetLogger sets the global logger instance (useful for testing)
func SetLogger(l Logger) {
	globalLogger = l
}

// copy creates a copy of the logger with its own fields map
func (l *loggerImpl) copy() *loggerImpl {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	return &loggerImpl{
		adapter: l.adapter,
		config:  l.config,
		fields:  newFields,
		level:   l.level,
		msg:     l.msg,
	}
}

// Builder methods - return Logger for chaining

func (l *loggerImpl) WithAny(key string, value interface{}) Logger {
	newLogger := l.copy()
	newLogger.fields[key] = value
	return newLogger
}

func (l *loggerImpl) WithString(key, value string) Logger {
	return l.WithAny(key, value)
}

func (l *loggerImpl) WithInt(key string, value int) Logger {
	return l.WithAny(key, value)
}

func (l *loggerImpl) WithInt64(key string, value int64) Logger {
	return l.WithAny(key, value)
}

func (l *loggerImpl) WithFloat64(key string, value float64) Logger {
	return l.WithAny(key, value)
}

func (l *loggerImpl) WithBool(key string, value bool) Logger {
	return l.WithAny(key, value)
}

func (l *loggerImpl) WithError(err error) Logger {
	if err != nil {
		return l.WithAny("error", err.Error())
	}
	return l
}

func (l *loggerImpl) WithFields(fields map[string]interface{}) Logger {
	newLogger := l.copy()
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

func (l *loggerImpl) WithContext(ctx context.Context) Logger {
	// For now, return the same logger
	// Can be enhanced to extract context values
	return l
}

// Log method - writes the accumulated log entry
func (l *loggerImpl) Log() {
	if l.msg == "" {
		return // No message set, can't log
	}
	
	if !l.shouldLog(l.level) {
		return
	}
	
	l.adapter.Log(l.level, l.msg, l.fields)
}

// Direct logging methods (for simple cases)

func (l *loggerImpl) Debug(msg string) {
	if !l.shouldLog(LogLevelDebug) {
		return
	}
	l.adapter.Log(LogLevelDebug, msg, nil)
}

func (l *loggerImpl) Info(msg string) {
	if !l.shouldLog(LogLevelInfo) {
		return
	}
	l.adapter.Log(LogLevelInfo, msg, nil)
}

func (l *loggerImpl) Warn(msg string) {
	if !l.shouldLog(LogLevelWarn) {
		return
	}
	l.adapter.Log(LogLevelWarn, msg, nil)
}

func (l *loggerImpl) Error(msg string) {
	if !l.shouldLog(LogLevelError) {
		return
	}
	l.adapter.Log(LogLevelError, msg, nil)
}

func (l *loggerImpl) Fatal(msg string) {
	l.adapter.Log(LogLevelFatal, msg, nil)
}

// Formatted logging methods

func (l *loggerImpl) Debugf(template string, args ...interface{}) {
	if l.shouldLog(LogLevelDebug) {
		l.adapter.Logf(LogLevelDebug, template, args...)
	}
}

func (l *loggerImpl) Infof(template string, args ...interface{}) {
	if l.shouldLog(LogLevelInfo) {
		l.adapter.Logf(LogLevelInfo, template, args...)
	}
}

func (l *loggerImpl) Warnf(template string, args ...interface{}) {
	if l.shouldLog(LogLevelWarn) {
		l.adapter.Logf(LogLevelWarn, template, args...)
	}
}

func (l *loggerImpl) Errorf(template string, args ...interface{}) {
	if l.shouldLog(LogLevelError) {
		l.adapter.Logf(LogLevelError, template, args...)
	}
}

func (l *loggerImpl) Fatalf(template string, args ...interface{}) {
	l.adapter.Logf(LogLevelFatal, template, args...)
}

// shouldLog checks if the message should be logged based on log level
func (l *loggerImpl) shouldLog(level LogLevel) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	levels := map[LogLevel]int{
		LogLevelDebug: 0,
		LogLevelInfo:  1,
		LogLevelWarn:  2,
		LogLevelError: 3,
		LogLevelFatal: 4,
	}

	return levels[level] >= levels[l.config.Level]
}

// Convenience functions for global logger with builder pattern
// These return Logger builders that can be chained

func Debug(msg string) Logger {
	logger := GetLogger().(*loggerImpl).copy()
	logger.level = LogLevelDebug
	logger.msg = msg
	return logger
}

func Info(msg string) Logger {
	logger := GetLogger().(*loggerImpl).copy()
	logger.level = LogLevelInfo
	logger.msg = msg
	return logger
}

func Warn(msg string) Logger {
	logger := GetLogger().(*loggerImpl).copy()
	logger.level = LogLevelWarn
	logger.msg = msg
	return logger
}

func Error(msg string) Logger {
	logger := GetLogger().(*loggerImpl).copy()
	logger.level = LogLevelError
	logger.msg = msg
	return logger
}

func Fatal(msg string) Logger {
	logger := GetLogger().(*loggerImpl).copy()
	logger.level = LogLevelFatal
	logger.msg = msg
	return logger
}

func Debugf(template string, args ...interface{}) {
	GetLogger().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	GetLogger().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	GetLogger().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	GetLogger().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	GetLogger().Fatalf(template, args...)
}
