package conversation_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"athema/internal/domain"
	"athema/internal/domain/conversation"
)

func TestMessage_JSONMarshalCamelCase(t *testing.T) {
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
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	expected := []string{"id", "conversationId", "companionId", "role", "content", "createdAt"}
	for _, key := range expected {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected camelCase key %q in JSON output", key)
		}
	}

	// Verify IDs are UUID strings, not arrays.
	for _, key := range []string{"id", "conversationId", "companionId"} {
		val := string(raw[key])
		if val[0] == '[' {
			t.Errorf("field %q serialized as array instead of UUID string: %s", key, val)
		}
		// UUID strings are 36 chars + quotes = 38 chars.
		if len(val) != 38 {
			t.Errorf("field %q expected UUID string (38 chars with quotes), got %d chars: %s", key, len(val), val)
		}
	}
}

func TestMessage_JSONRoundTrip(t *testing.T) {
	original := conversation.Message{
		ID:             domain.NewMessageID(),
		ConversationID: domain.NewConversationID(),
		CompanionID:    domain.NewCompanionID(),
		Role:           domain.RoleCompanion,
		Content:        "test content",
		CreatedAt:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded conversation.Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %v, want %v", decoded.ID, original.ID)
	}
	if decoded.ConversationID != original.ConversationID {
		t.Errorf("ConversationID mismatch: got %v, want %v", decoded.ConversationID, original.ConversationID)
	}
	if decoded.CompanionID != original.CompanionID {
		t.Errorf("CompanionID mismatch: got %v, want %v", decoded.CompanionID, original.CompanionID)
	}
	if decoded.Role != original.Role {
		t.Errorf("Role mismatch: got %v, want %v", decoded.Role, original.Role)
	}
	if decoded.Content != original.Content {
		t.Errorf("Content mismatch: got %v, want %v", decoded.Content, original.Content)
	}
	if !decoded.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", decoded.CreatedAt, original.CreatedAt)
	}
}

func TestConversationErrors_WrapDomainErrors(t *testing.T) {
	if !errors.Is(conversation.ErrConversationNotFound, domain.ErrNotFound) {
		t.Error("ErrConversationNotFound should wrap domain.ErrNotFound")
	}
	if !errors.Is(conversation.ErrMessageEmpty, domain.ErrInvalidInput) {
		t.Error("ErrMessageEmpty should wrap domain.ErrInvalidInput")
	}
}

func TestConversation_JSONMarshalCamelCase(t *testing.T) {
	conv := conversation.Conversation{
		ID:          domain.NewConversationID(),
		CompanionID: domain.NewCompanionID(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	data, err := json.Marshal(conv)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	expected := []string{"id", "companionId", "createdAt", "updatedAt"}
	for _, key := range expected {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected camelCase key %q in JSON output", key)
		}
	}
}
