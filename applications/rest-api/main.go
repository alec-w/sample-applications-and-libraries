//go:generate go tool oapi-codegen -config oapi-codegen.config.yaml openapi.yaml

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alec-w/sample-applications-and-libraries/applications/rest-api/internal/api"
	"github.com/alec-w/sample-applications-and-libraries/applications/rest-api/internal/database"
	"github.com/alec-w/sample-applications-and-libraries/libraries/logging"
)

const (
	EXIT_CODE_SUCCESS int = iota
	EXIT_CODE_FAILURE
)

const (
	// TODO - make these passed through as args / flags / env vars / config file / fetched/mounted from secret stores
	defaultPort            = 8080
	defaultShutdownTimeout = 5 * time.Second
	logLevel               = slog.LevelDebug
	databaseHost           = "postgres"
	databaseUser           = "postgres"
	databaseName           = "postgres"
	databasePassword       = "postgres"
	databasePort           = 5432
)

func main() {
	// Delegating control to app means any deferred calls within app always happen
	// and we can still always set the exit code correctly
	os.Exit(app(context.Background()))
}

func app(ctx context.Context) int {
	// Components for app
	logger := logging.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})))
	database, err := database.NewDatabase(ctx, databaseHost, databaseUser, databaseName, databasePassword, databasePort, logger)
	if err != nil {
		logger.WithError(err).Error("failed to instantiate database")
		return EXIT_CODE_FAILURE
	}
	server := api.NewServer(8080, logger, database)

	// Listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server, if it fails to start catch this in the erred channel
	erred := make(chan struct{})
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("http server failed to listen on port")
			erred <- struct{}{}
		}
	}()
	logger.Info("server starting")

	// Wait for shutdown signal or server to error out and stop listening
	select {
	case <-shutdown:
		logger.Info("server stopping")

		// Attempt graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(ctx, defaultShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.WithError(err).Error("server shutdown failed")
			return EXIT_CODE_FAILURE
		}

		logger.Info("server shutdown cleanly")
		return EXIT_CODE_SUCCESS
	case <-erred:
		logger.Error("server failed to start, exiting")
		return EXIT_CODE_FAILURE
	}
}
