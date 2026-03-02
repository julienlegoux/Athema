package domain

import (
	"context"
	"encoding/json"
	"time"
)

// CompanionState represents the persisted state of a companion instance.
type CompanionState struct {
	ID        CompanionID     `json:"id"`
	State     json.RawMessage `json:"state"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// CompanionStateRepository defines the persistence port for companion state.
type CompanionStateRepository interface {
	WithTx(ctx context.Context, fn func(CompanionStateRepository) error) error
	GetState(ctx context.Context, companionID CompanionID) (*CompanionState, error)
	SaveState(ctx context.Context, state CompanionState) error
}
