package memory

import (
	"encoding/json"
	"time"

	"athema/internal/domain"
)

// KnowledgeNode represents a node in the knowledge graph.
type KnowledgeNode struct {
	ID          domain.MessageID   `json:"id"`
	CompanionID domain.CompanionID `json:"companionId"`
	NodeType    string             `json:"nodeType"`
	Payload     json.RawMessage    `json:"payload"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

// KnowledgeEdge represents a directed edge between two knowledge nodes.
type KnowledgeEdge struct {
	ID        domain.MessageID `json:"id"`
	FromID    domain.MessageID `json:"fromId"`
	ToID      domain.MessageID `json:"toId"`
	Type      string           `json:"type"`
	Weight    float64          `json:"weight"`
	CreatedAt time.Time        `json:"createdAt"`
}
