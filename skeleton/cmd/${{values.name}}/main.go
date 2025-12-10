// Package main is the entry point for ${{values.name}} CLI
package main

import (
	"os"

	"github.com/fast-ish/${{values.name}}/internal/cli"
{%- if values.logging == "slog" %}
	"github.com/fast-ish/${{values.name}}/internal/logger"
{%- elif values.logging == "zap" %}
	"go.uber.org/zap"
{%- elif values.logging == "zerolog" %}
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
{%- endif %}
{%- if values.tracing %}
	"github.com/fast-ish/${{values.name}}/internal/telemetry"
{%- endif %}
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
{%- if values.logging == "slog" %}
	// Initialize structured logging
	logger.Init(logger.LevelInfo, false)
{%- elif values.logging == "zap" %}
	// Initialize zap logger
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
{%- elif values.logging == "zerolog" %}
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
{%- endif %}

{%- if values.tracing %}
	// Initialize OpenTelemetry tracing
	shutdown, err := telemetry.InitTracing("${{values.name}}", version)
	if err != nil {
{%- if values.logging == "slog" %}
		logger.Error("Failed to initialize tracing", "error", err)
{%- elif values.logging == "zap" %}
		zapLogger.Error("Failed to initialize tracing", zap.Error(err))
{%- elif values.logging == "zerolog" %}
		log.Error().Err(err).Msg("Failed to initialize tracing")
{%- endif %}
	}
	defer func() {
		if shutdown != nil {
			shutdown(context.Background())
		}
	}()
{%- endif %}

	// Set version info
	cli.SetVersion(version, buildTime, gitCommit)

	// Execute CLI
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
