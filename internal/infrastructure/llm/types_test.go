package llm_test

import (
	"testing"

	"athema/internal/infrastructure/llm"
)

func TestCompletionRequestConstruction(t *testing.T) {
	req := &llm.CompletionRequest{
		Model:        "claude-sonnet-4-20250514",
		SystemPrompt: "You are a helpful assistant.",
		Messages: []llm.Message{
			{Role: "user", Content: "Hello"},
			{Role: "assistant", Content: "Hi there!"},
		},
		MaxTokens:    1024,
		Temperature:  0.7,
		CompanionID:  "test-companion-id",
		SubsystemTag: "conversation",
	}

	if req.Model != "claude-sonnet-4-20250514" {
		t.Errorf("expected model claude-sonnet-4-20250514, got %s", req.Model)
	}
	if req.SystemPrompt != "You are a helpful assistant." {
		t.Errorf("expected system prompt, got %s", req.SystemPrompt)
	}
	if len(req.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(req.Messages))
	}
	if req.Messages[0].Role != "user" {
		t.Errorf("expected role user, got %s", req.Messages[0].Role)
	}
	if req.Messages[0].Content != "Hello" {
		t.Errorf("expected content Hello, got %s", req.Messages[0].Content)
	}
	if req.MaxTokens != 1024 {
		t.Errorf("expected max tokens 1024, got %d", req.MaxTokens)
	}
	if req.Temperature != 0.7 {
		t.Errorf("expected temperature 0.7, got %f", req.Temperature)
	}
	if req.CompanionID != "test-companion-id" {
		t.Errorf("expected companion ID test-companion-id, got %s", req.CompanionID)
	}
	if req.SubsystemTag != "conversation" {
		t.Errorf("expected subsystem tag conversation, got %s", req.SubsystemTag)
	}
}

func TestMessageRoles(t *testing.T) {
	tests := []struct {
		role    string
		content string
	}{
		{"user", "Hello"},
		{"assistant", "Hi there!"},
		{"system", "You are helpful."},
	}

	for _, tt := range tests {
		msg := llm.Message{Role: tt.role, Content: tt.content}
		if msg.Role != tt.role {
			t.Errorf("expected role %s, got %s", tt.role, msg.Role)
		}
		if msg.Content != tt.content {
			t.Errorf("expected content %s, got %s", tt.content, msg.Content)
		}
	}
}

func TestCompletionResponseConstruction(t *testing.T) {
	resp := &llm.CompletionResponse{
		Content:      "Hello, how can I help?",
		InputTokens:  10,
		OutputTokens: 8,
		Model:        "claude-sonnet-4-20250514",
		StopReason:   "end_turn",
	}

	if resp.Content != "Hello, how can I help?" {
		t.Errorf("unexpected content: %s", resp.Content)
	}
	if resp.InputTokens != 10 {
		t.Errorf("expected 10 input tokens, got %d", resp.InputTokens)
	}
	if resp.OutputTokens != 8 {
		t.Errorf("expected 8 output tokens, got %d", resp.OutputTokens)
	}
}

func TestStreamEventConstruction(t *testing.T) {
	delta := llm.StreamEvent{ContentDelta: "Hello", IsFinal: false}
	if delta.ContentDelta != "Hello" {
		t.Errorf("unexpected delta: %s", delta.ContentDelta)
	}
	if delta.IsFinal {
		t.Error("expected non-final event")
	}

	final := llm.StreamEvent{IsFinal: true, InputTokens: 10, OutputTokens: 5}
	if !final.IsFinal {
		t.Error("expected final event")
	}
	if final.InputTokens != 10 {
		t.Errorf("expected 10 input tokens, got %d", final.InputTokens)
	}

	errEvent := llm.StreamEvent{IsFinal: true, Error: "llm: provider unavailable"}
	if !errEvent.IsFinal {
		t.Error("error event should be final")
	}
	if errEvent.Error != "llm: provider unavailable" {
		t.Errorf("unexpected error: %s", errEvent.Error)
	}
}

func TestEmbeddingTypes(t *testing.T) {
	req := &llm.EmbeddingRequest{Text: "hello world", Model: "voyage-3"}
	if req.Text != "hello world" {
		t.Errorf("unexpected text: %s", req.Text)
	}

	resp := &llm.EmbeddingResponse{
		Vector:      []float32{0.1, 0.2, 0.3},
		Model:       "voyage-3",
		InputTokens: 2,
	}
	if len(resp.Vector) != 3 {
		t.Errorf("expected 3 dimensions, got %d", len(resp.Vector))
	}
}
