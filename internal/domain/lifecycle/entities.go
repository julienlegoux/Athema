package lifecycle

import (
	"time"

	"athema/internal/domain"
)

// LifecycleTask is a stub entity for autonomous lifecycle processing.
// Expanded in Epic 4.
type LifecycleTask struct {
	CompanionID domain.CompanionID `json:"companionId"`
	CreatedAt   time.Time          `json:"createdAt"`
}

// ProcessingResult is a stub entity for lifecycle processing outcomes.
// Expanded in Epic 4.
type ProcessingResult struct {
	CompanionID domain.CompanionID `json:"companionId"`
	CreatedAt   time.Time          `json:"createdAt"`
}
