package personality

import (
	"encoding/json"
	"time"

	"athema/internal/domain"
)

// PersonalitySnapshot represents a point-in-time capture of personality state.
type PersonalitySnapshot struct {
	ID          domain.MessageID   `json:"id"`
	CompanionID domain.CompanionID `json:"companionId"`
	Payload     json.RawMessage    `json:"payload"`
	CreatedAt   time.Time          `json:"createdAt"`
}
