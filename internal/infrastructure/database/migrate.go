package database

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending database migrations.
func RunMigrations(dsn string, migrationsPath string, logger *slog.Logger) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}

	version, dirty, vErr := m.Version()
	if vErr != nil {
		logger.Warn("migrations applied but version retrieval failed", "error", vErr)
	} else {
		logger.Info("migrations applied", "version", version, "dirty", dirty)
	}
	return nil
}
