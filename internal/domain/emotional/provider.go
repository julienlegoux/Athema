package emotional

import (
	"context"

	"athema/internal/domain"
)

// EmotionalStateProvider is the narrow interface that other subsystems use
// to read the current emotional state of a companion.
type EmotionalStateProvider interface {
	CurrentState(ctx context.Context, companionID domain.CompanionID) (*EmotionalState, error)
}
