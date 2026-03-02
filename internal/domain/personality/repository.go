package personality

import (
	"context"

	"athema/internal/domain"
)

// PersonalityRepository defines the persistence port for personality snapshots.
// Stub interface — full API designed in later epics.
type PersonalityRepository interface {
	WithTx(ctx context.Context, fn func(PersonalityRepository) error) error
	CreateSnapshot(ctx context.Context, snapshot PersonalitySnapshot) error
	GetLatestSnapshot(ctx context.Context, companionID domain.CompanionID) (*PersonalitySnapshot, error)
}
