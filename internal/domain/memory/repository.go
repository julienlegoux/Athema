package memory

import (
	"context"

	"athema/internal/domain"
)

// MemoryRepository defines the persistence port for the knowledge graph.
// Stub interface — full API designed in Epic 3.
type MemoryRepository interface {
	WithTx(ctx context.Context, fn func(MemoryRepository) error) error
	CreateNode(ctx context.Context, node KnowledgeNode) error
	GetNodeByID(ctx context.Context, companionID domain.CompanionID, nodeID domain.MessageID) (*KnowledgeNode, error)
}
