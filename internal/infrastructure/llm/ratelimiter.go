package llm

import "context"

// RateLimiter controls concurrent access to LLM providers using a semaphore pattern.
type RateLimiter struct {
	sem chan struct{}
}

// NewRateLimiter creates a new RateLimiter with the given maximum concurrent permits.
func NewRateLimiter(maxConcurrent int) *RateLimiter {
	return &RateLimiter{
		sem: make(chan struct{}, maxConcurrent),
	}
}

// Acquire blocks until a permit is available or the context is canceled.
func (r *RateLimiter) Acquire(ctx context.Context) error {
	select {
	case r.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release returns a permit to the pool.
func (r *RateLimiter) Release() {
	<-r.sem
}
