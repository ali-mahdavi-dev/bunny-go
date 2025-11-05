package logging

import (
	"fmt"
	"sync"
)

// loggerImpl is the implementation of Logger interface
type loggerImpl struct {
	adapter LoggerAdapter
	config  LoggerConfig
	fields  []LogField // accumulated fields from builder
	level   LogLevel   // log level for builder
	msg     string     // message for builder
}

var (
	globalLogger Logger
	globalOnce   sync.Once
)

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

// NewLogger creates a new logger instance
func NewLogger(config LoggerConfig) (Logger, error) {
	var adapter LoggerAdapter
	var err error

	switch config.Type {
	case LoggerTypeZerolog:
		adapter, err = newZerologAdapter(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create zerolog adapter: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported logger type: %s", config.Type)
	}

	return &loggerImpl{
		adapter: adapter,
		config:  config,
		fields:  make([]LogField, 0),
	}, nil
}

func (l *loggerImpl) Debug(msg string) Logger {
	l.level = LogLevelDebug
	l.msg = msg
	l.fields = make([]LogField, 0)
	return l
}

func (l *loggerImpl) Info(msg string) Logger {
	l.level = LogLevelInfo
	l.msg = msg
	l.fields = make([]LogField, 0)
	return l
}

func (l *loggerImpl) Warn(msg string) Logger {
	l.level = LogLevelWarn
	l.msg = msg
	l.fields = make([]LogField, 0)
	return l
}

func (l *loggerImpl) Error(msg string) Logger {
	l.level = LogLevelError
	l.msg = msg
	l.fields = make([]LogField, 0)
	return l
}

func (l *loggerImpl) Fatal(msg string) Logger {
	l.level = LogLevelFatal
	l.msg = msg
	l.fields = make([]LogField, 0)
	return l
}

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

// Logger builder methods - mutate self and return self

func (l *loggerImpl) WithAny(key string, value interface{}) Logger {
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeAny,
	})
	return l
}

func (l *loggerImpl) WithString(key, value string) Logger {
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeString,
	})
	return l
}

func (l *loggerImpl) WithInt(key string, value int) Logger {
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeInt,
	})
	return l
}

func (l *loggerImpl) WithInt64(key string, value int64) Logger {
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeInt64,
	})
	return l
}

func (l *loggerImpl) WithFloat64(key string, value float64) Logger {
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeFloat64,
	})
	return l
}

func (l *loggerImpl) WithBool(key string, value bool) Logger {
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeBool,
	})
	return l
}

func (l *loggerImpl) WithError(err error) Logger {
	if err != nil {
		l.fields = append(l.fields, LogField{
			Key:   "error",
			Value: err,
			Type:  FieldTypeError,
		})
	}
	return l
}

func (l *loggerImpl) WithFields(fields map[string]interface{}) Logger {
	for k, v := range fields {
		l.fields = append(l.fields, LogField{
			Key:   k,
			Value: v,
			Type:  FieldTypeAny,
		})
	}
	return l
}

func (l *loggerImpl) Log() {
	if l.msg == "" {
		return // No message set, can't log
	}

	if !l.shouldLog(l.level) {
		return
	}

	l.adapter.Log(l.level, l.msg, l.fields)

	// Clear fields after logging
	l.fields = make([]LogField, 0)
	l.msg = ""
	l.level = ""
}

// shouldLog checks if the message should be logged based on log level
func (l *loggerImpl) shouldLog(level LogLevel) bool {
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
func Debug(msg string) Logger {
	return GetLogger().Debug(msg)
}

func Info(msg string) Logger {
	return GetLogger().Info(msg)
}

func Warn(msg string) Logger {
	return GetLogger().Warn(msg)
}

func Error(msg string) Logger {
	return GetLogger().Error(msg)
}

func Fatal(msg string) Logger {
	return GetLogger().Fatal(msg)
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
