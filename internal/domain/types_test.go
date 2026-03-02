package domain_test

import (
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
