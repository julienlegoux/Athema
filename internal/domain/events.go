package domain

import "time"

// Event is the base interface for all domain events.
type Event interface {
	EventType() string
	OccurredAt() time.Time
	GetCompanionID() CompanionID
}

// BaseEvent provides common event fields.
type BaseEvent struct {
	Type        string      `json:"type"`
	Timestamp   time.Time   `json:"occurredAt"`
	CompanionID CompanionID `json:"companionId"`
}

func (e BaseEvent) EventType() string         { return e.Type }
func (e BaseEvent) OccurredAt() time.Time     { return e.Timestamp }
func (e BaseEvent) GetCompanionID() CompanionID { return e.CompanionID }

// NewBaseEvent creates a new BaseEvent with the current timestamp.
func NewBaseEvent(eventType string, companionID CompanionID) BaseEvent {
	return BaseEvent{
		Type:        eventType,
		Timestamp:   time.Now().UTC(),
		CompanionID: companionID,
	}
}
