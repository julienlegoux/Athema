package conversation

import (
	"fmt"

	"athema/internal/domain"
)

var (
	ErrConversationNotFound = fmt.Errorf("conversation: %w", domain.ErrNotFound)
	ErrMessageEmpty         = fmt.Errorf("conversation: message empty: %w", domain.ErrInvalidInput)
)
