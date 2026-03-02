package memory

import (
	"context"

	"athema/internal/domain"
)

// MemoryContextProvider is the narrow interface that other subsystems use
// to retrieve relevant memory context for a companion.
type MemoryContextProvider interface {
	RelevantContext(ctx context.Context, companionID domain.CompanionID, query string) ([]KnowledgeNode, error)
}
