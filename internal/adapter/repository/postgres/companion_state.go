package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"athema/internal/domain"
)

// CompanionStateRepository implements domain.CompanionStateRepository using PostgreSQL.
type CompanionStateRepository struct {
	q      querier
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewCompanionStateRepository creates a new CompanionStateRepository backed by the given pool.
func NewCompanionStateRepository(pool *pgxpool.Pool, logger *slog.Logger) *CompanionStateRepository {
	return &CompanionStateRepository{
		q:      pool,
		pool:   pool,
		logger: logger.With("subsystem", "companion_state_repo"),
	}
}

func (r *CompanionStateRepository) WithTx(ctx context.Context, fn func(domain.CompanionStateRepository) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("companion_state_repo.WithTx: begin: %w", err)
	}

	txRepo := &CompanionStateRepository{
		q:      tx,
		pool:   r.pool,
		logger: r.logger,
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(txRepo); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			r.logger.Error("companion_state_repo.WithTx: rollback failed", "error", rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("companion_state_repo.WithTx: commit: %w", err)
	}
	return nil
}

func (r *CompanionStateRepository) GetState(ctx context.Context, companionID domain.CompanionID) (*domain.CompanionState, error) {
	const query = `SELECT companion_id, state, created_at, updated_at
		FROM companion_state WHERE companion_id = $1`

	var state domain.CompanionState
	err := r.q.QueryRow(ctx, query, companionID.String()).
		Scan(&state.ID, &state.State, &state.CreatedAt, &state.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("companion_state_repo.GetState: %w", err)
	}
	return &state, nil
}

func (r *CompanionStateRepository) SaveState(ctx context.Context, state domain.CompanionState) error {
	const query = `INSERT INTO companion_state (companion_id, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (companion_id) DO UPDATE SET state = EXCLUDED.state, updated_at = EXCLUDED.updated_at`

	_, err := r.q.Exec(ctx, query,
		state.ID.String(),
		state.State,
		state.CreatedAt,
		state.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("companion_state_repo.SaveState: %w", err)
	}
	return nil
}
