package llm

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

var testLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func testProvider() *AnthropicProvider {
	return &AnthropicProvider{model: "default-model", logger: testLogger}
}

func TestTranslateErrorNil(t *testing.T) {
	if err := translateError(nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestTranslateErrorDeadlineExceeded(t *testing.T) {
	err := translateError(context.DeadlineExceeded)
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("expected ErrTimeout, got %v", err)
	}
}

func TestTranslateErrorContextCanceled(t *testing.T) {
	err := translateError(context.Canceled)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestTranslateErrorHTTPStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   error
	}{
		{"unauthorized", http.StatusUnauthorized, ErrAuthenticationFailed},
		{"forbidden", http.StatusForbidden, ErrAuthenticationFailed},
		{"rate_limited", http.StatusTooManyRequests, ErrRateLimited},
		{"bad_request", http.StatusBadRequest, ErrInvalidRequest},
		{"server_error_500", http.StatusInternalServerError, ErrProviderUnavailable},
		{"server_error_502", http.StatusBadGateway, ErrProviderUnavailable},
		{"server_error_503", http.StatusServiceUnavailable, ErrProviderUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apierr := &anthropic.Error{StatusCode: tt.statusCode}
			result := translateError(apierr)
			if !errors.Is(result, tt.expected) {
				t.Errorf("status %d: expected %v, got %v", tt.statusCode, tt.expected, result)
			}
		})
	}
}

func TestTranslateErrorUnknownError(t *testing.T) {
	err := translateError(errors.New("something unknown"))
	if !errors.Is(err, ErrProviderUnavailable) {
		t.Errorf("expected ErrProviderUnavailable for unknown error, got %v", err)
	}
}

func TestBuildParamsSystemPrompt(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		SystemPrompt: "Be helpful",
		Messages:     []Message{{Role: "user", Content: "Hi"}},
		MaxTokens:    512,
	}

	params := p.buildParams(req)

	if len(params.System) != 1 {
		t.Fatalf("expected 1 system block, got %d", len(params.System))
	}
	if params.System[0].Text != "Be helpful" {
		t.Errorf("expected 'Be helpful', got %s", params.System[0].Text)
	}
}

func TestBuildParamsNoSystemPrompt(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 512,
	}

	params := p.buildParams(req)

	if len(params.System) != 0 {
		t.Errorf("expected no system blocks, got %d", len(params.System))
	}
}

func TestBuildParamsMessageRoles(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		Messages: []Message{
			{Role: "user", Content: "Hello"},
			{Role: "assistant", Content: "Hi there"},
			{Role: "user", Content: "How are you?"},
		},
		MaxTokens: 512,
	}

	params := p.buildParams(req)

	if len(params.Messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(params.Messages))
	}
	if params.Messages[0].Role != anthropic.MessageParamRoleUser {
		t.Errorf("expected user role, got %s", params.Messages[0].Role)
	}
	if params.Messages[1].Role != anthropic.MessageParamRoleAssistant {
		t.Errorf("expected assistant role, got %s", params.Messages[1].Role)
	}
}

func TestBuildParamsUnsupportedRoleIsDropped(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		Messages: []Message{
			{Role: "user", Content: "Hello"},
			{Role: "system", Content: "This should be dropped"},
			{Role: "user", Content: "World"},
		},
		MaxTokens: 512,
	}

	params := p.buildParams(req)

	if len(params.Messages) != 2 {
		t.Fatalf("expected 2 messages (system role dropped), got %d", len(params.Messages))
	}
}

func TestBuildParamsMaxTokens(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 2048,
	}

	params := p.buildParams(req)

	if params.MaxTokens != 2048 {
		t.Errorf("expected max tokens 2048, got %d", params.MaxTokens)
	}
}

func TestBuildParamsModelOverride(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		Model:     "override-model",
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 512,
	}

	params := p.buildParams(req)

	if string(params.Model) != "override-model" {
		t.Errorf("expected override-model, got %s", params.Model)
	}
}

func TestBuildParamsDefaultModel(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 512,
	}

	params := p.buildParams(req)

	if string(params.Model) != "default-model" {
		t.Errorf("expected default-model, got %s", params.Model)
	}
}

func TestBuildParamsTemperature(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		Messages:    []Message{{Role: "user", Content: "Hi"}},
		MaxTokens:   512,
		Temperature: 0.8,
	}

	params := p.buildParams(req)

	if params.Temperature.Value != 0.8 {
		t.Errorf("expected temperature 0.8, got %f", params.Temperature.Value)
	}
}

func TestBuildParamsTemperatureZeroNotSet(t *testing.T) {
	p := testProvider()
	req := &CompletionRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 512,
	}

	params := p.buildParams(req)

	if params.Temperature.Valid() {
		t.Error("expected temperature to not be set when zero")
	}
}

func TestCompleteRejectsZeroMaxTokens(t *testing.T) {
	p := testProvider()
	_, err := p.Complete(context.Background(), &CompletionRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 0,
	})
	if err == nil {
		t.Fatal("expected error for zero MaxTokens")
	}
	if !errors.Is(err, ErrInvalidRequest) {
		t.Errorf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestStreamRejectsZeroMaxTokens(t *testing.T) {
	p := testProvider()
	_, err := p.Stream(context.Background(), &CompletionRequest{
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 0,
	})
	if err == nil {
		t.Fatal("expected error for zero MaxTokens")
	}
	if !errors.Is(err, ErrInvalidRequest) {
		t.Errorf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestEmbedReturnsNotSupported(t *testing.T) {
	p := &AnthropicProvider{}
	_, err := p.Embed(context.Background(), &EmbeddingRequest{Text: "hello"})
	if !errors.Is(err, ErrEmbeddingNotSupported) {
		t.Errorf("expected ErrEmbeddingNotSupported, got %v", err)
	}
}

func TestAnthropicProviderImplementsInterface(t *testing.T) {
	var _ Provider = &AnthropicProvider{}
}
