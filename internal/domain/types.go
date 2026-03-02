package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// CompanionID uniquely identifies a companion instance.
type CompanionID uuid.UUID

// String returns the string representation of a CompanionID.
func (id CompanionID) String() string {
	return uuid.UUID(id).String()
}

// ParseCompanionID parses a string into a CompanionID.
func ParseCompanionID(s string) (CompanionID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return CompanionID{}, fmt.Errorf("domain.ParseCompanionID: %w", err)
	}
	return CompanionID(id), nil
}

// NewCompanionID generates a new random CompanionID.
func NewCompanionID() CompanionID {
	return CompanionID(uuid.New())
}

// MessageID uniquely identifies a message.
type MessageID uuid.UUID

// String returns the string representation of a MessageID.
func (id MessageID) String() string {
	return uuid.UUID(id).String()
}

// ParseMessageID parses a string into a MessageID.
func ParseMessageID(s string) (MessageID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return MessageID{}, fmt.Errorf("domain.ParseMessageID: %w", err)
	}
	return MessageID(id), nil
}

// NewMessageID generates a new random MessageID.
func NewMessageID() MessageID {
	return MessageID(uuid.New())
}

// ConversationID uniquely identifies a conversation.
type ConversationID uuid.UUID

// String returns the string representation of a ConversationID.
func (id ConversationID) String() string {
	return uuid.UUID(id).String()
}

// ParseConversationID parses a string into a ConversationID.
func ParseConversationID(s string) (ConversationID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ConversationID{}, fmt.Errorf("domain.ParseConversationID: %w", err)
	}
	return ConversationID(id), nil
}

// NewConversationID generates a new random ConversationID.
func NewConversationID() ConversationID {
	return ConversationID(uuid.New())
}

// Role represents the author of a message.
type Role string

const (
	RoleUser      Role = "user"
	RoleCompanion Role = "companion"
)

// UserID uniquely identifies a user.
type UserID string

// SessionID uniquely identifies a conversation session.
type SessionID string
