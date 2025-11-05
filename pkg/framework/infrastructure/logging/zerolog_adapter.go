package logging

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

// ZerologAdapter implements LoggerAdapter using zerolog library
type ZerologAdapter struct {
	logger zerolog.Logger
}

// NewZerologAdapter creates a new zerolog logger adapter
func NewZerologAdapter(config LoggerConfig) (LoggerAdapter, error) {
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

	return &ZerologAdapter{
		logger: logger,
	}, nil
}

// Log logs a message at the specified level with fields
func (z *ZerologAdapter) Log(level LogLevel, msg string, fields map[string]interface{}) {
	event := z.logger.WithLevel(z.toZerologLevel(level))
	
	// Add fields
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	
	event.Msg(msg)
}

// Logf logs a formatted message at the specified level
func (z *ZerologAdapter) Logf(level LogLevel, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	z.logger.WithLevel(z.toZerologLevel(level)).Msg(msg)
}

// toZerologLevel converts LogLevel to zerolog.Level
func (z *ZerologAdapter) toZerologLevel(level LogLevel) zerolog.Level {
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

// WithContext creates a new adapter with context (for future use)
func (z *ZerologAdapter) WithContext(ctx context.Context) LoggerAdapter {
	// Extract request ID or other context values if needed
	logger := z.logger.With().Logger()
	return &ZerologAdapter{logger: logger}
}

// WithFields creates a new adapter with additional fields (for future use)
func (z *ZerologAdapter) WithFields(fields map[string]interface{}) LoggerAdapter {
	ctx := z.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	logger := ctx.Logger()
	return &ZerologAdapter{logger: logger}
}
