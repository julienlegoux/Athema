package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"athema/internal/infrastructure/config"
	"athema/internal/infrastructure/server"
)

func main() {
	// Load .env file if present (ignore error — .env is optional).
	_ = godotenv.Load()

	// Set up structured JSON logging.
	logLevel := slog.LevelInfo
	baseLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(baseLogger)

	serverLogger := baseLogger.With("subsystem", "server")

	// Load configuration.
	cfg, err := config.Load("config/default.yaml")
	if err != nil {
		serverLogger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Apply log level from config.
	switch cfg.Log.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	baseLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(baseLogger)
	serverLogger = baseLogger.With("subsystem", "server")

	serverLogger.Info("configuration loaded",
		"server_addr", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		"db_host", cfg.DB.Host,
		"log_level", cfg.Log.Level,
	)

	// Create HTTP server.
	srv := server.New(cfg.Server, serverLogger)

	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		serverLogger.Info("server started", "port", cfg.Server.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		serverLogger.Info("shutdown signal received")
	case err := <-errCh:
		serverLogger.Error("server error", "error", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		serverLogger.Error("shutdown error", "error", err)
		os.Exit(1)
	}

	serverLogger.Info("server stopped gracefully")
}
