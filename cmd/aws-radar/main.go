package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nimishgj/aws-radar/internal/collector"
	"github.com/nimishgj/aws-radar/internal/config"
	"github.com/nimishgj/aws-radar/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup logging
	setupLogging(cfg.Logging)

	log.Info().
		Str("version", "1.0.0").
		Msg("Starting AWS Radar")

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start collector orchestrator
	orchestrator := collector.NewOrchestrator(
		cfg.AWS.Regions,
		cfg.Collection.Interval,
		cfg.Collection.Timeout,
		cfg.Collectors,
	)
	go orchestrator.Start(ctx)

	// Start HTTP server
	srv := server.New(&cfg.Server)
	go func() {
		if err := srv.Start(); err != nil {
			log.Error().Err(err).Msg("HTTP server error")
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Info().Msg("Shutdown signal received")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
	}

	cancel()
	log.Info().Msg("AWS Radar stopped")
}

func setupLogging(cfg config.LoggingConfig) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set output format
	if cfg.Format == "text" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}
