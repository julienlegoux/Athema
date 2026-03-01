package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"athema/internal/infrastructure/config"
)

// DB wraps a pgx connection pool.
type DB struct {
	Pool   *pgxpool.Pool
	logger *slog.Logger
}

// New creates a new database connection pool.
func New(ctx context.Context, cfg config.DBConfig, logger *slog.Logger) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parsing database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	logger.Info("database connected", "host", cfg.Host, "database", cfg.Database)
	return &DB{Pool: pool, logger: logger}, nil
}

// Close closes the database connection pool.
func (db *DB) Close() {
	db.logger.Info("database connection closing")
	db.Pool.Close()
}

// Health checks the database connection.
func (db *DB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
