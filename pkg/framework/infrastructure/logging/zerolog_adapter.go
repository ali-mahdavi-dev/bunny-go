package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

// zerologAdapter is the Zerolog implementation of LoggerAdapter
type zerologAdapter struct {
	logger zerolog.Logger
}

// newZerologAdapter creates a new Zerolog adapter
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
