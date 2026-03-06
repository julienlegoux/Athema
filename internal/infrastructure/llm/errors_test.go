package llm_test

import (
	"errors"
	"fmt"
	"testing"

	"athema/internal/infrastructure/llm"
)

func TestSentinelErrorsAreDistinct(t *testing.T) {
	sentinels := []error{
		llm.ErrProviderUnavailable,
		llm.ErrRateLimited,
		llm.ErrAuthenticationFailed,
		llm.ErrInvalidRequest,
		llm.ErrContextTooLong,
		llm.ErrTimeout,
		llm.ErrEmbeddingNotSupported,
	}

	for i, a := range sentinels {
		for j, b := range sentinels {
			if i != j && errors.Is(a, b) {
				t.Errorf("sentinel errors %d and %d should be distinct, but errors.Is returned true", i, j)
			}
		}
	}
}

func TestErrorWrapping(t *testing.T) {
	wrapped := fmt.Errorf("llm.Complete: %w", llm.ErrRateLimited)
	if !errors.Is(wrapped, llm.ErrRateLimited) {
		t.Error("wrapped error should match ErrRateLimited via errors.Is")
	}
	if errors.Is(wrapped, llm.ErrTimeout) {
		t.Error("wrapped ErrRateLimited should not match ErrTimeout")
	}
}

func TestAllSentinelErrorsHavePrefix(t *testing.T) {
	sentinels := map[string]error{
		"ErrProviderUnavailable":  llm.ErrProviderUnavailable,
		"ErrRateLimited":          llm.ErrRateLimited,
		"ErrAuthenticationFailed": llm.ErrAuthenticationFailed,
		"ErrInvalidRequest":       llm.ErrInvalidRequest,
		"ErrContextTooLong":       llm.ErrContextTooLong,
		"ErrTimeout":              llm.ErrTimeout,
		"ErrEmbeddingNotSupported": llm.ErrEmbeddingNotSupported,
	}

	for name, err := range sentinels {
		msg := err.Error()
		if len(msg) < 4 || msg[:4] != "llm:" {
			t.Errorf("%s error message should start with 'llm:', got: %s", name, msg)
		}
	}
}
