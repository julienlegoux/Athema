package conversation

import "athema/internal/domain"

const (
	EventMessageReceived       = "conversation.message_received"
	EventConversationCompleted = "conversation.completed"
)

// MessageReceivedEvent is published when a new message is added to a conversation.
type MessageReceivedEvent struct {
	domain.BaseEvent
	MessageID      domain.MessageID      `json:"messageId"`
	ConversationID domain.ConversationID `json:"conversationId"`
	Role           domain.Role           `json:"role"`
	Content        string                `json:"content"`
}

// ConversationCompletedEvent is published when a conversation session ends.
type ConversationCompletedEvent struct {
	domain.BaseEvent
	ConversationID domain.ConversationID `json:"conversationId"`
}
