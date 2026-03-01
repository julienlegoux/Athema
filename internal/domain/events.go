package domain

import "time"

// Event is the base interface for all domain events.
type Event interface {
	EventType() string
	OccurredAt() time.Time
}

// BaseEvent provides common event fields.
type BaseEvent struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func (e BaseEvent) EventType() string    { return e.Type }
func (e BaseEvent) OccurredAt() time.Time { return e.Timestamp }

// NewBaseEvent creates a new BaseEvent with the current timestamp.
func NewBaseEvent(eventType string) BaseEvent {
	return BaseEvent{
		Type:      eventType,
		Timestamp: time.Now().UTC(),
	}
}
