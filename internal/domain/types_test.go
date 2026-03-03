package domain_test

import (
	"encoding/json"
	"strings"
	"testing"

	"athema/internal/domain"
)

func TestCompanionID_NewAndString(t *testing.T) {
	id := domain.NewCompanionID()
	s := id.String()
	if s == "" {
		t.Fatal("expected non-empty string from CompanionID")
	}
	if len(s) != 36 { // UUID v4 format: 8-4-4-4-12
		t.Fatalf("expected UUID string length 36, got %d: %q", len(s), s)
	}
}

func TestParseCompanionID_Valid(t *testing.T) {
	original := domain.NewCompanionID()
	parsed, err := domain.ParseCompanionID(original.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed != original {
		t.Fatalf("expected %v, got %v", original, parsed)
	}
}

func TestParseCompanionID_Invalid(t *testing.T) {
	_, err := domain.ParseCompanionID("not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestMessageID_NewAndString(t *testing.T) {
	id := domain.NewMessageID()
	s := id.String()
	if s == "" {
		t.Fatal("expected non-empty string from MessageID")
	}
	if len(s) != 36 {
		t.Fatalf("expected UUID string length 36, got %d: %q", len(s), s)
	}
}

func TestParseMessageID_Valid(t *testing.T) {
	original := domain.NewMessageID()
	parsed, err := domain.ParseMessageID(original.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed != original {
		t.Fatalf("expected %v, got %v", original, parsed)
	}
}

func TestParseMessageID_Invalid(t *testing.T) {
	_, err := domain.ParseMessageID("bad")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestConversationID_NewAndString(t *testing.T) {
	id := domain.NewConversationID()
	s := id.String()
	if s == "" {
		t.Fatal("expected non-empty string from ConversationID")
	}
	if len(s) != 36 {
		t.Fatalf("expected UUID string length 36, got %d: %q", len(s), s)
	}
}

func TestParseConversationID_Valid(t *testing.T) {
	original := domain.NewConversationID()
	parsed, err := domain.ParseConversationID(original.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed != original {
		t.Fatalf("expected %v, got %v", original, parsed)
	}
}

func TestParseConversationID_Invalid(t *testing.T) {
	_, err := domain.ParseConversationID("nope")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestRole_Constants(t *testing.T) {
	if domain.RoleUser != "user" {
		t.Fatalf("expected RoleUser to be %q, got %q", "user", domain.RoleUser)
	}
	if domain.RoleCompanion != "companion" {
		t.Fatalf("expected RoleCompanion to be %q, got %q", "companion", domain.RoleCompanion)
	}
}

// TestCompanionID_JSONMarshalAsUUIDString verifies IDs serialize as UUID strings, not byte arrays.
func TestCompanionID_JSONMarshalAsUUIDString(t *testing.T) {
	id := domain.NewCompanionID()

	type wrapper struct {
		ID domain.CompanionID `json:"id"`
	}

	data, err := json.Marshal(wrapper{ID: id})
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	s := string(data)
	// Must contain a UUID string, not an array of integers.
	if strings.Contains(s, "[") {
		t.Fatalf("CompanionID serialized as array, not UUID string: %s", s)
	}

	expected := id.String()
	if !strings.Contains(s, expected) {
		t.Fatalf("expected JSON to contain UUID string %q, got %s", expected, s)
	}

	// Round-trip
	var decoded wrapper
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if decoded.ID != id {
		t.Fatalf("round-trip mismatch: got %v, want %v", decoded.ID, id)
	}
}

func TestMessageID_JSONMarshalAsUUIDString(t *testing.T) {
	id := domain.NewMessageID()

	type wrapper struct {
		ID domain.MessageID `json:"id"`
	}

	data, err := json.Marshal(wrapper{ID: id})
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	s := string(data)
	if strings.Contains(s, "[") {
		t.Fatalf("MessageID serialized as array, not UUID string: %s", s)
	}

	var decoded wrapper
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if decoded.ID != id {
		t.Fatalf("round-trip mismatch: got %v, want %v", decoded.ID, id)
	}
}

func TestKnowledgeNodeID_NewAndParse(t *testing.T) {
	id := domain.NewKnowledgeNodeID()
	s := id.String()
	if len(s) != 36 {
		t.Fatalf("expected UUID string length 36, got %d: %q", len(s), s)
	}

	parsed, err := domain.ParseKnowledgeNodeID(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed != id {
		t.Fatalf("expected %v, got %v", id, parsed)
	}
}

func TestSnapshotID_NewAndParse(t *testing.T) {
	id := domain.NewSnapshotID()
	s := id.String()
	if len(s) != 36 {
		t.Fatalf("expected UUID string length 36, got %d: %q", len(s), s)
	}

	parsed, err := domain.ParseSnapshotID(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed != id {
		t.Fatalf("expected %v, got %v", id, parsed)
	}
}

func TestEmotionalStateID_NewAndParse(t *testing.T) {
	id := domain.NewEmotionalStateID()
	s := id.String()
	if len(s) != 36 {
		t.Fatalf("expected UUID string length 36, got %d: %q", len(s), s)
	}

	parsed, err := domain.ParseEmotionalStateID(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed != id {
		t.Fatalf("expected %v, got %v", id, parsed)
	}
}
