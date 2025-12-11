// Package logger provides structured logging
package logger

import (
{%- if values.logging == "slog" %}
	"log/slog"
	"os"
{%- elif values.logging == "zap" %}
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
{%- elif values.logging == "zerolog" %}
	"os"
	"github.com/rs/zerolog"
{%- endif %}
)

{%- if values.logging == "slog" %}

// Log levels
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

var L *slog.Logger

// Init initializes the logger
func Init(level slog.Level, json bool) {
	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if json {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	L = slog.New(handler)
	slog.SetDefault(L)
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	L.Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	L.Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	L.Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	L.Error(msg, args...)
}

{%- elif values.logging == "zap" %}

var L *zap.Logger

// Init initializes the zap logger
func Init(level zapcore.Level, json bool) error {
	var config zap.Config
	if json {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}
	config.Level = zap.NewAtomicLevelAt(level)

	logger, err := config.Build()
	if err != nil {
		return err
	}
	L = logger
	zap.ReplaceGlobals(L)
	return nil
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	L.Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	L.Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	L.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	L.Error(msg, fields...)
}

{%- elif values.logging == "zerolog" %}

// Init initializes the zerolog logger
func Init(level zerolog.Level, json bool) {
	zerolog.SetGlobalLevel(level)
	if !json {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

// Debug logs a debug message
func Debug() *zerolog.Event {
	return log.Debug()
}

// Info logs an info message
func Info() *zerolog.Event {
	return log.Info()
}

// Warn logs a warning message
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error logs an error message
func Error() *zerolog.Event {
	return log.Error()
}
{%- endif %}
