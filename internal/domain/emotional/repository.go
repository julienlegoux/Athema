package emotional

import (
	"context"

	"athema/internal/domain"
)

// EmotionalRepository defines the persistence port for emotional state.
// Stub interface — full API designed in later epics.
type EmotionalRepository interface {
	WithTx(ctx context.Context, fn func(EmotionalRepository) error) error
	SaveState(ctx context.Context, state EmotionalState) error
	GetState(ctx context.Context, companionID domain.CompanionID) (*EmotionalState, error)
}
