package personality

import (
	"context"

	"athema/internal/domain"
)

// PersonalityProvider is the narrow interface that other subsystems use
// to read the current personality state of a companion.
type PersonalityProvider interface {
	CurrentPersonality(ctx context.Context, companionID domain.CompanionID) (*PersonalitySnapshot, error)
}
