package conversation

import (
	"time"

	"athema/internal/domain"
)

// Message represents a single message in a conversation.
type Message struct {
	ID             domain.MessageID      `json:"id"`
	ConversationID domain.ConversationID `json:"conversationId"`
	CompanionID    domain.CompanionID    `json:"companionId"`
	Role           domain.Role           `json:"role"`
	Content        string                `json:"content"`
	CreatedAt      time.Time             `json:"createdAt"`
}

// Conversation represents an ongoing or completed dialogue session.
type Conversation struct {
	ID          domain.ConversationID `json:"id"`
	CompanionID domain.CompanionID    `json:"companionId"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}
