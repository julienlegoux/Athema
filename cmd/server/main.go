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

	"athema/internal/adapter/repository/postgres"
	"athema/internal/infrastructure/config"
	"athema/internal/infrastructure/database"
	"athema/internal/infrastructure/llm"
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

	// Create LLM provider.
	const defaultMaxConcurrent = 5
	llmLogger := baseLogger.With("subsystem", "llm")
	if cfg.LLM.Provider != "anthropic" {
		serverLogger.Error("unsupported LLM provider", "provider", cfg.LLM.Provider)
		os.Exit(1)
	}
	if cfg.LLM.APIKey == "" {
		serverLogger.Warn("LLM API key is empty — provider calls will fail")
	}
	rateLimiter := llm.NewRateLimiter(defaultMaxConcurrent)
	llmProvider := llm.NewAnthropicProvider(cfg.LLM.APIKey, cfg.LLM.Model, llmLogger, rateLimiter)
	_ = llmProvider // Provider will be injected into use cases in Story 1.6
	serverLogger.Info("llm provider initialized", "provider", cfg.LLM.Provider, "model", cfg.LLM.Model)

	// Initialize database connection pool.
	dbLogger := baseLogger.With("subsystem", "database")
	db, err := database.New(context.Background(), cfg.DB, dbLogger)
	if err != nil {
		serverLogger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run database migrations.
	if err := database.RunMigrations(cfg.DB.DSN(), "migrations", dbLogger); err != nil {
		serverLogger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Create repository instances.
	repoLogger := baseLogger.With("subsystem", "repository")
	conversationRepo := postgres.NewConversationRepository(db.Pool, repoLogger)
	companionStateRepo := postgres.NewCompanionStateRepository(db.Pool, repoLogger)
	_ = conversationRepo    // Will be injected into use cases in Story 1.6
	_ = companionStateRepo  // Will be injected into use cases in Story 1.6
	serverLogger.Info("database initialized and migrations applied")

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
