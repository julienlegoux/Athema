package memory

import (
	"encoding/json"
	"time"

	"athema/internal/domain"
)

// KnowledgeNode represents a node in the knowledge graph.
type KnowledgeNode struct {
	ID          domain.KnowledgeNodeID `json:"id"`
	CompanionID domain.CompanionID     `json:"companionId"`
	NodeType    string                 `json:"nodeType"`
	Payload     json.RawMessage        `json:"payload"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// KnowledgeEdge represents a directed edge between two knowledge nodes.
type KnowledgeEdge struct {
	ID        domain.KnowledgeEdgeID `json:"id"`
	FromID    domain.KnowledgeNodeID `json:"fromId"`
	ToID      domain.KnowledgeNodeID `json:"toId"`
	Type      string                 `json:"type"`
	Weight    float64                `json:"weight"`
	CreatedAt time.Time              `json:"createdAt"`
}
