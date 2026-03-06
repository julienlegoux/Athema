package llm

import "context"

// Provider defines the abstracted interface for LLM interactions.
// Implementations must translate all provider-specific errors into domain sentinel errors.
type Provider interface {
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	Stream(ctx context.Context, req *CompletionRequest) (<-chan StreamEvent, error)
	Embed(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
}
