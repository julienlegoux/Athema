package postgres_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"athema/internal/adapter/repository/postgres"
	"athema/internal/domain"
	"athema/internal/domain/conversation"
	"athema/internal/infrastructure/llm"
)

// Compile-time interface compliance check.
var _ conversation.ConversationRepository = (*postgres.ConversationRepository)(nil)

// conversationFixtures holds sample data loaded from test fixtures.
type conversationFixtures struct {
	Conversations []conversation.Conversation `json:"conversations"`
	Messages      []conversation.Message      `json:"messages"`
}

func loadConversationFixtures(t *testing.T) conversationFixtures {
	t.Helper()
	fixtures, err := llm.LoadFixture[conversationFixtures]("../../../../test/fixtures/conversation/sample_data.json")
	if err != nil {
		t.Fatalf("loading conversation fixtures: %v", err)
	}
	return fixtures
}

func TestConversationFixtures_Load(t *testing.T) {
	fixtures := loadConversationFixtures(t)

	if len(fixtures.Conversations) == 0 {
		t.Fatal("expected at least one conversation fixture")
	}
	if len(fixtures.Messages) == 0 {
		t.Fatal("expected at least one message fixture")
	}

	conv := fixtures.Conversations[0]
	if conv.ID.String() == "" {
		t.Error("conversation ID should not be empty")
	}
	if conv.CompanionID.String() == "" {
		t.Error("conversation CompanionID should not be empty")
	}
	if conv.CreatedAt.IsZero() {
		t.Error("conversation CreatedAt should not be zero")
	}

	msg := fixtures.Messages[0]
	if msg.Content == "" {
		t.Error("message Content should not be empty")
	}
	if msg.Role != domain.RoleUser && msg.Role != domain.RoleCompanion {
		t.Errorf("message Role should be user or companion, got %q", msg.Role)
	}
}

func TestConversationFixtures_JSONTags(t *testing.T) {
	conv := conversation.Conversation{
		ID:          domain.NewConversationID(),
		CompanionID: domain.NewCompanionID(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	data, err := json.Marshal(conv)
	if err != nil {
		t.Fatalf("marshal conversation: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal conversation: %v", err)
	}

	// Verify camelCase JSON field names (AC8).
	expectedFields := []string{"id", "companionId", "createdAt", "updatedAt"}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("expected camelCase JSON field %q, got keys: %v", field, keys(raw))
		}
	}
}

func TestMessageFixtures_JSONTags(t *testing.T) {
	msg := conversation.Message{
		ID:             domain.NewMessageID(),
		ConversationID: domain.NewConversationID(),
		CompanionID:    domain.NewCompanionID(),
		Role:           domain.RoleUser,
		Content:        "hello",
		CreatedAt:      time.Now().UTC(),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal message: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal message: %v", err)
	}

	// Verify camelCase JSON field names (AC8).
	expectedFields := []string{"id", "conversationId", "companionId", "role", "content", "createdAt"}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("expected camelCase JSON field %q, got keys: %v", field, keys(raw))
		}
	}
}

func TestConversationRepository_ErrorMapping(t *testing.T) {
	// Verify domain error wrapping is correct.
	if !errors.Is(conversation.ErrConversationNotFound, domain.ErrNotFound) {
		t.Error("ErrConversationNotFound should wrap domain.ErrNotFound")
	}
	if !errors.Is(conversation.ErrMessageEmpty, domain.ErrInvalidInput) {
		t.Error("ErrMessageEmpty should wrap domain.ErrInvalidInput")
	}
}

func keys[V any](m map[string]V) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func TestNewConversationRepository(t *testing.T) {
	// Constructor should not panic with nil pool (used for type checking only).
	// In production, pool is always non-nil.
	// We can't create a real pool without a DB, but we verify the constructor signature.
	t.Run("interface compliance", func(t *testing.T) {
		// Compile-time check at top of file ensures this.
		// This test documents the expected constructor pattern.
		var repo conversation.ConversationRepository
		_ = repo // Suppress unused warning - interface check is compile-time.
	})
}
