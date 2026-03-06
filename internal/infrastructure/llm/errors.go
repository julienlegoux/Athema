package llm

import "errors"

// Sentinel errors for the LLM provider abstraction layer.
var (
	ErrProviderUnavailable  = errors.New("llm: provider unavailable")
	ErrRateLimited          = errors.New("llm: rate limited")
	ErrAuthenticationFailed = errors.New("llm: authentication failed")
	ErrInvalidRequest       = errors.New("llm: invalid request")
	ErrContextTooLong       = errors.New("llm: context too long")
	ErrTimeout              = errors.New("llm: request timeout")
	ErrEmbeddingNotSupported = errors.New("llm: embedding not supported")
)
