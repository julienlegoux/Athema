package domain_test

import (
	"testing"
	"time"

	"athema/internal/domain"
)

func TestNewBaseEvent(t *testing.T) {
	companionID := domain.NewCompanionID()
	before := time.Now().UTC()
	event := domain.NewBaseEvent("test.event", companionID)
	after := time.Now().UTC()

	if event.EventType() != "test.event" {
		t.Fatalf("expected event type %q, got %q", "test.event", event.EventType())
	}

	if event.GetCompanionID() != companionID {
		t.Fatalf("expected companion ID %v, got %v", companionID, event.GetCompanionID())
	}

	occurred := event.OccurredAt()
	if occurred.Before(before) || occurred.After(after) {
		t.Fatalf("expected OccurredAt between %v and %v, got %v", before, after, occurred)
	}
}

func TestBaseEvent_ImplementsEventInterface(t *testing.T) {
	companionID := domain.NewCompanionID()
	event := domain.NewBaseEvent("test.type", companionID)

	// Verify interface satisfaction at compile time via assignment.
	var _ domain.Event = event
}
