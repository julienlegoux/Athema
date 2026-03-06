package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"athema/internal/domain"
	"athema/internal/domain/conversation"
)

// querier abstracts pgxpool.Pool and pgx.Tx so repository methods work
// inside and outside transactions.
type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// ConversationRepository implements conversation.ConversationRepository using PostgreSQL.
type ConversationRepository struct {
	q      querier
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewConversationRepository creates a new ConversationRepository backed by the given pool.
func NewConversationRepository(pool *pgxpool.Pool, logger *slog.Logger) *ConversationRepository {
	return &ConversationRepository{
		q:      pool,
		pool:   pool,
		logger: logger.With("subsystem", "conversation_repo"),
	}
}

func (r *ConversationRepository) WithTx(ctx context.Context, fn func(conversation.ConversationRepository) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("conversation_repo.WithTx: begin: %w", err)
	}

	txRepo := &ConversationRepository{
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
			r.logger.Error("conversation_repo.WithTx: rollback failed", "error", rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("conversation_repo.WithTx: commit: %w", err)
	}
	return nil
}

func (r *ConversationRepository) CreateConversation(ctx context.Context, conv conversation.Conversation) error {
	const query = `INSERT INTO conversations (id, companion_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.q.Exec(ctx, query,
		conv.ID.String(),
		conv.CompanionID.String(),
		conv.CreatedAt,
		conv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("conversation_repo.CreateConversation: %w", err)
	}
	return nil
}

func (r *ConversationRepository) CreateMessage(ctx context.Context, msg conversation.Message) error {
	const query = `INSERT INTO messages (id, conversation_id, companion_id, role, content, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.q.Exec(ctx, query,
		msg.ID.String(),
		msg.ConversationID.String(),
		msg.CompanionID.String(),
		string(msg.Role),
		msg.Content,
		msg.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("conversation_repo.CreateMessage: %w", err)
	}
	return nil
}

func (r *ConversationRepository) GetConversation(ctx context.Context, companionID domain.CompanionID, conversationID domain.ConversationID) (*conversation.Conversation, error) {
	const query = `SELECT id, companion_id, created_at, updated_at
		FROM conversations WHERE companion_id = $1 AND id = $2`

	var conv conversation.Conversation
	err := r.q.QueryRow(ctx, query, companionID.String(), conversationID.String()).
		Scan(&conv.ID, &conv.CompanionID, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, conversation.ErrConversationNotFound
		}
		return nil, fmt.Errorf("conversation_repo.GetConversation: %w", err)
	}
	return &conv, nil
}

func (r *ConversationRepository) ListMessages(ctx context.Context, companionID domain.CompanionID, conversationID domain.ConversationID) ([]conversation.Message, error) {
	const query = `SELECT id, conversation_id, companion_id, role, content, created_at
		FROM messages WHERE companion_id = $1 AND conversation_id = $2
		ORDER BY created_at ASC`

	rows, err := r.q.Query(ctx, query, companionID.String(), conversationID.String())
	if err != nil {
		return nil, fmt.Errorf("conversation_repo.ListMessages: %w", err)
	}
	defer rows.Close()

	var messages []conversation.Message
	for rows.Next() {
		var msg conversation.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.CompanionID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("conversation_repo.ListMessages: scan: %w", err)
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("conversation_repo.ListMessages: rows: %w", err)
	}
	return messages, nil
}

func (r *ConversationRepository) GetActiveConversation(ctx context.Context, companionID domain.CompanionID) (*conversation.Conversation, error) {
	const query = `SELECT id, companion_id, created_at, updated_at
		FROM conversations WHERE companion_id = $1
		ORDER BY updated_at DESC LIMIT 1`

	var conv conversation.Conversation
	err := r.q.QueryRow(ctx, query, companionID.String()).
		Scan(&conv.ID, &conv.CompanionID, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, conversation.ErrConversationNotFound
		}
		return nil, fmt.Errorf("conversation_repo.GetActiveConversation: %w", err)
	}
	return &conv, nil
}
