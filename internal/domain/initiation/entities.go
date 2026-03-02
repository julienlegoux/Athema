package initiation

import (
	"time"

	"athema/internal/domain"
)

// InitiationEvent is a stub entity for companion-initiated contact.
// Expanded in Epic 7.
type InitiationEvent struct {
	CompanionID domain.CompanionID `json:"companionId"`
	CreatedAt   time.Time          `json:"createdAt"`
}

// UrgeState is a stub entity for urge accumulation tracking.
// Expanded in Epic 7.
type UrgeState struct {
	CompanionID domain.CompanionID `json:"companionId"`
	CreatedAt   time.Time          `json:"createdAt"`
}
