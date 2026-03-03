package emotional

import (
	"encoding/json"
	"time"

	"athema/internal/domain"
)

// EmotionalState represents the current emotional state of a companion.
type EmotionalState struct {
	ID          domain.EmotionalStateID `json:"id"`
	CompanionID domain.CompanionID      `json:"companionId"`
	State       json.RawMessage         `json:"state"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}
