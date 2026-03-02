package conversation

import "errors"

var (
	ErrConversationNotFound = errors.New("conversation.not_found")
	ErrMessageEmpty         = errors.New("conversation.message_empty")
)
