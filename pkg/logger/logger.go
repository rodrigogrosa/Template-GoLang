package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger wraps zerolog logger
type Logger struct {
	logger zerolog.Logger
}

// New creates a new logger instance
func New(level, format string) *Logger {
	var l zerolog.Logger

	// Set log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Set format
	if format == "console" {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	} else {
		l = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	return &Logger{logger: l}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logger.Info().Fields(fields).Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	l.logger.Error().Err(err).Fields(fields).Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...interface{}) {
	l.logger.Fatal().Err(err).Fields(fields).Msg(msg)
}

// With creates a child logger with additional fields
func (l *Logger) With(fields map[string]interface{}) *Logger {
	return &Logger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}

// GetZerolog returns the underlying zerolog logger
func (l *Logger) GetZerolog() zerolog.Logger {
	return l.logger
}

// Global logger
var Global *Logger

func init() {
	Global = New("info", "json")
	log.Logger = Global.logger
}
