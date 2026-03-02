package conversation

import (
	"context"

	"athema/internal/domain"
)

// ConversationRepository defines the persistence port for conversations and messages.
type ConversationRepository interface {
	WithTx(ctx context.Context, fn func(ConversationRepository) error) error
	CreateConversation(ctx context.Context, conv Conversation) error
	CreateMessage(ctx context.Context, msg Message) error
	GetConversation(ctx context.Context, companionID domain.CompanionID, conversationID domain.ConversationID) (*Conversation, error)
	ListMessages(ctx context.Context, companionID domain.CompanionID, conversationID domain.ConversationID) ([]Message, error)
	GetActiveConversation(ctx context.Context, companionID domain.CompanionID) (*Conversation, error)
}
