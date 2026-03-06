package llm_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"athema/internal/infrastructure/llm"
)

func fixturesDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "test", "fixtures", "llm")
}

func TestLoadFixtureCompletionResponse(t *testing.T) {
	path := filepath.Join(fixturesDir(), "completion_response.json")
	resp, err := llm.LoadFixture[llm.CompletionResponse](path)
	if err != nil {
		t.Fatalf("failed to load fixture: %v", err)
	}
	if resp.Content != "Hello! How can I help you today?" {
		t.Errorf("unexpected content: %s", resp.Content)
	}
	if resp.InputTokens != 15 {
		t.Errorf("expected 15 input tokens, got %d", resp.InputTokens)
	}
	if resp.OutputTokens != 9 {
		t.Errorf("expected 9 output tokens, got %d", resp.OutputTokens)
	}
	if resp.Model != "claude-sonnet-4-20250514" {
		t.Errorf("unexpected model: %s", resp.Model)
	}
	if resp.StopReason != "end_turn" {
		t.Errorf("unexpected stop reason: %s", resp.StopReason)
	}
}

func TestLoadFixtureStreamEvents(t *testing.T) {
	path := filepath.Join(fixturesDir(), "stream_events.json")
	events, err := llm.LoadFixture[[]llm.StreamEvent](path)
	if err != nil {
		t.Fatalf("failed to load fixture: %v", err)
	}
	if len(events) != 7 {
		t.Fatalf("expected 7 stream events, got %d", len(events))
	}
	if events[0].ContentDelta != "Hello" {
		t.Errorf("expected first delta 'Hello', got %s", events[0].ContentDelta)
	}
	if events[0].IsFinal {
		t.Error("first event should not be final")
	}
	last := events[len(events)-1]
	if !last.IsFinal {
		t.Error("last event should be final")
	}
	if last.InputTokens != 15 {
		t.Errorf("expected 15 input tokens on final, got %d", last.InputTokens)
	}
}

func TestLoadFixtureErrorResponses(t *testing.T) {
	type errorFixture struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}
	path := filepath.Join(fixturesDir(), "error_responses.json")
	fixtures, err := llm.LoadFixture[map[string]errorFixture](path)
	if err != nil {
		t.Fatalf("failed to load fixture: %v", err)
	}
	if _, ok := fixtures["rateLimited"]; !ok {
		t.Error("missing rateLimited fixture")
	}
	if fixtures["rateLimited"].StatusCode != 429 {
		t.Errorf("expected 429, got %d", fixtures["rateLimited"].StatusCode)
	}
	if _, ok := fixtures["authFailure"]; !ok {
		t.Error("missing authFailure fixture")
	}
	if fixtures["serverError"].StatusCode != 500 {
		t.Errorf("expected 500, got %d", fixtures["serverError"].StatusCode)
	}
}

func TestLoadFixtureInvalidPath(t *testing.T) {
	_, err := llm.LoadFixture[llm.CompletionResponse]("nonexistent.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
