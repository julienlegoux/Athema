package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// CompanionID uniquely identifies a companion instance.
type CompanionID uuid.UUID

func (id CompanionID) String() string                { return uuid.UUID(id).String() }
func (id CompanionID) MarshalText() ([]byte, error)  { return []byte(id.String()), nil }
func (id *CompanionID) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = CompanionID(u)
	return nil
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

func (id MessageID) String() string                { return uuid.UUID(id).String() }
func (id MessageID) MarshalText() ([]byte, error)  { return []byte(id.String()), nil }
func (id *MessageID) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = MessageID(u)
	return nil
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

func (id ConversationID) String() string                { return uuid.UUID(id).String() }
func (id ConversationID) MarshalText() ([]byte, error)  { return []byte(id.String()), nil }
func (id *ConversationID) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = ConversationID(u)
	return nil
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

// KnowledgeNodeID uniquely identifies a knowledge graph node.
type KnowledgeNodeID uuid.UUID

func (id KnowledgeNodeID) String() string                { return uuid.UUID(id).String() }
func (id KnowledgeNodeID) MarshalText() ([]byte, error)  { return []byte(id.String()), nil }
func (id *KnowledgeNodeID) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = KnowledgeNodeID(u)
	return nil
}

// ParseKnowledgeNodeID parses a string into a KnowledgeNodeID.
func ParseKnowledgeNodeID(s string) (KnowledgeNodeID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return KnowledgeNodeID{}, fmt.Errorf("domain.ParseKnowledgeNodeID: %w", err)
	}
	return KnowledgeNodeID(id), nil
}

// NewKnowledgeNodeID generates a new random KnowledgeNodeID.
func NewKnowledgeNodeID() KnowledgeNodeID {
	return KnowledgeNodeID(uuid.New())
}

// KnowledgeEdgeID uniquely identifies a knowledge graph edge.
type KnowledgeEdgeID uuid.UUID

func (id KnowledgeEdgeID) String() string                { return uuid.UUID(id).String() }
func (id KnowledgeEdgeID) MarshalText() ([]byte, error)  { return []byte(id.String()), nil }
func (id *KnowledgeEdgeID) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = KnowledgeEdgeID(u)
	return nil
}

// ParseKnowledgeEdgeID parses a string into a KnowledgeEdgeID.
func ParseKnowledgeEdgeID(s string) (KnowledgeEdgeID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return KnowledgeEdgeID{}, fmt.Errorf("domain.ParseKnowledgeEdgeID: %w", err)
	}
	return KnowledgeEdgeID(id), nil
}

// NewKnowledgeEdgeID generates a new random KnowledgeEdgeID.
func NewKnowledgeEdgeID() KnowledgeEdgeID {
	return KnowledgeEdgeID(uuid.New())
}

// SnapshotID uniquely identifies a personality snapshot.
type SnapshotID uuid.UUID

func (id SnapshotID) String() string                { return uuid.UUID(id).String() }
func (id SnapshotID) MarshalText() ([]byte, error)  { return []byte(id.String()), nil }
func (id *SnapshotID) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = SnapshotID(u)
	return nil
}

// ParseSnapshotID parses a string into a SnapshotID.
func ParseSnapshotID(s string) (SnapshotID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return SnapshotID{}, fmt.Errorf("domain.ParseSnapshotID: %w", err)
	}
	return SnapshotID(id), nil
}

// NewSnapshotID generates a new random SnapshotID.
func NewSnapshotID() SnapshotID {
	return SnapshotID(uuid.New())
}

// EmotionalStateID uniquely identifies an emotional state record.
type EmotionalStateID uuid.UUID

func (id EmotionalStateID) String() string                { return uuid.UUID(id).String() }
func (id EmotionalStateID) MarshalText() ([]byte, error)  { return []byte(id.String()), nil }
func (id *EmotionalStateID) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = EmotionalStateID(u)
	return nil
}

// ParseEmotionalStateID parses a string into an EmotionalStateID.
func ParseEmotionalStateID(s string) (EmotionalStateID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return EmotionalStateID{}, fmt.Errorf("domain.ParseEmotionalStateID: %w", err)
	}
	return EmotionalStateID(id), nil
}

// NewEmotionalStateID generates a new random EmotionalStateID.
func NewEmotionalStateID() EmotionalStateID {
	return EmotionalStateID(uuid.New())
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
