package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
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
	// Builder methods - return Logger for chaining (mutate self)
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

// loggerImpl is the implementation of Logger interface with builder pattern
type loggerImpl struct {
	adapter LoggerAdapter
	config  LoggerConfig
	mu      sync.RWMutex
	fields  []LogField // accumulated fields from builder (mutable)
	level   LogLevel   // log level for builder
	msg     string     // message for builder
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

// Builder methods - mutate self and return self

func (l *loggerImpl) WithAny(key string, value interface{}) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeAny,
	})
	return l
}

func (l *loggerImpl) WithString(key, value string) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeString,
	})
	return l
}

func (l *loggerImpl) WithInt(key string, value int) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeInt,
	})
	return l
}

func (l *loggerImpl) WithInt64(key string, value int64) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeInt64,
	})
	return l
}

func (l *loggerImpl) WithFloat64(key string, value float64) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeFloat64,
	})
	return l
}

func (l *loggerImpl) WithBool(key string, value bool) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fields = append(l.fields, LogField{
		Key:   key,
		Value: value,
		Type:  FieldTypeBool,
	})
	return l
}

func (l *loggerImpl) WithError(err error) Logger {
	if err != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.fields = append(l.fields, LogField{
			Key:   "error",
			Value: err,
			Type:  FieldTypeError,
		})
	}
	return l
}

func (l *loggerImpl) WithFields(fields map[string]interface{}) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	for k, v := range fields {
		l.fields = append(l.fields, LogField{
			Key:   k,
			Value: v,
			Type:  FieldTypeAny,
		})
	}
	return l
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
	
	l.mu.RLock()
	fields := make([]LogField, len(l.fields))
	copy(fields, l.fields)
	level := l.level
	msg := l.msg
	l.mu.RUnlock()
	
	if !l.shouldLog(level) {
		return
	}
	
	l.adapter.Log(level, msg, fields)
	
	// Clear fields after logging
	l.mu.Lock()
	l.fields = make([]LogField, 0)
	l.msg = ""
	l.level = ""
	l.mu.Unlock()
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
// These create a new logger builder with message and level set

func Debug(msg string) Logger {
	logger := GetLogger().(*loggerImpl)
	logger.mu.Lock()
	logger.level = LogLevelDebug
	logger.msg = msg
	logger.fields = make([]LogField, 0)
	logger.mu.Unlock()
	return logger
}

func Info(msg string) Logger {
	logger := GetLogger().(*loggerImpl)
	logger.mu.Lock()
	logger.level = LogLevelInfo
	logger.msg = msg
	logger.fields = make([]LogField, 0)
	logger.mu.Unlock()
	return logger
}

func Warn(msg string) Logger {
	logger := GetLogger().(*loggerImpl)
	logger.mu.Lock()
	logger.level = LogLevelWarn
	logger.msg = msg
	logger.fields = make([]LogField, 0)
	logger.mu.Unlock()
	return logger
}

func Error(msg string) Logger {
	logger := GetLogger().(*loggerImpl)
	logger.mu.Lock()
	logger.level = LogLevelError
	logger.msg = msg
	logger.fields = make([]LogField, 0)
	logger.mu.Unlock()
	return logger
}

func Fatal(msg string) Logger {
	logger := GetLogger().(*loggerImpl)
	logger.mu.Lock()
	logger.level = LogLevelFatal
	logger.msg = msg
	logger.fields = make([]LogField, 0)
	logger.mu.Unlock()
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

// ZerologAdapter implementation
type zerologAdapter struct {
	logger zerolog.Logger
}

func newZerologAdapter(config LoggerConfig) (LoggerAdapter, error) {
	var output io.Writer = config.Output
	if output == nil {
		output = os.Stdout
	}

	// Set log level
	var level zerolog.Level
	switch config.Level {
	case LogLevelDebug:
		level = zerolog.DebugLevel
	case LogLevelInfo:
		level = zerolog.InfoLevel
	case LogLevelWarn:
		level = zerolog.WarnLevel
	case LogLevelError:
		level = zerolog.ErrorLevel
	case LogLevelFatal:
		level = zerolog.FatalLevel
	default:
		level = zerolog.InfoLevel
	}

	// Create logger with appropriate format
	var logger zerolog.Logger
	if config.Format == LogFormatJSON {
		logger = zerolog.New(output).Level(level).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(output).Level(level).Output(zerolog.ConsoleWriter{Out: output}).With().Timestamp().Logger()
	}

	return &zerologAdapter{
		logger: logger,
	}, nil
}

// Log logs a message at the specified level with fields
func (z *zerologAdapter) Log(level LogLevel, msg string, fields []LogField) {
	event := z.logger.WithLevel(z.toZerologLevel(level))
	
	// Add fields using appropriate zerolog methods based on field type
	for _, field := range fields {
		switch field.Type {
		case FieldTypeString:
			if strVal, ok := field.Value.(string); ok {
				event = event.Str(field.Key, strVal)
			} else {
				event = event.Interface(field.Key, field.Value)
			}
		case FieldTypeInt:
			if intVal, ok := field.Value.(int); ok {
				event = event.Int(field.Key, intVal)
			} else {
				event = event.Interface(field.Key, field.Value)
			}
		case FieldTypeInt64:
			if int64Val, ok := field.Value.(int64); ok {
				event = event.Int64(field.Key, int64Val)
			} else {
				event = event.Interface(field.Key, field.Value)
			}
		case FieldTypeFloat64:
			if float64Val, ok := field.Value.(float64); ok {
				event = event.Float64(field.Key, float64Val)
			} else {
				event = event.Interface(field.Key, field.Value)
			}
		case FieldTypeBool:
			if boolVal, ok := field.Value.(bool); ok {
				event = event.Bool(field.Key, boolVal)
			} else {
				event = event.Interface(field.Key, field.Value)
			}
		case FieldTypeError:
			if errVal, ok := field.Value.(error); ok {
				event = event.Err(errVal)
			} else {
				event = event.Interface(field.Key, field.Value)
			}
		case FieldTypeAny:
			fallthrough
		default:
			event = event.Interface(field.Key, field.Value)
		}
	}
	
	event.Msg(msg)
}

// Logf logs a formatted message at the specified level
func (z *zerologAdapter) Logf(level LogLevel, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	z.logger.WithLevel(z.toZerologLevel(level)).Msg(msg)
}

// toZerologLevel converts LogLevel to zerolog.Level
func (z *zerologAdapter) toZerologLevel(level LogLevel) zerolog.Level {
	switch level {
	case LogLevelDebug:
		return zerolog.DebugLevel
	case LogLevelInfo:
		return zerolog.InfoLevel
	case LogLevelWarn:
		return zerolog.WarnLevel
	case LogLevelError:
		return zerolog.ErrorLevel
	case LogLevelFatal:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
