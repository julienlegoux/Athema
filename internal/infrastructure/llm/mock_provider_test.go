package llm_test

import (
	"context"
	"errors"
	"testing"

	"athema/internal/infrastructure/llm"
)

func TestMockProviderReturnsConfiguredResponse(t *testing.T) {
	mock := llm.NewMockProvider()
	expected := &llm.CompletionResponse{
		Content:      "Hello!",
		InputTokens:  10,
		OutputTokens: 5,
		Model:        "test-model",
		StopReason:   "end_turn",
	}
	mock.SetCompleteResponse(expected, nil)

	resp, err := mock.Complete(context.Background(), &llm.CompletionRequest{
		Model:    "test-model",
		Messages: []llm.Message{{Role: "user", Content: "Hi"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "Hello!" {
		t.Errorf("expected Hello!, got %s", resp.Content)
	}
	if resp.InputTokens != 10 {
		t.Errorf("expected 10 input tokens, got %d", resp.InputTokens)
	}
}

func TestMockProviderReturnsConfiguredError(t *testing.T) {
	mock := llm.NewMockProvider()
	mock.SetCompleteResponse(nil, llm.ErrRateLimited)

	_, err := mock.Complete(context.Background(), &llm.CompletionRequest{})
	if !errors.Is(err, llm.ErrRateLimited) {
		t.Errorf("expected ErrRateLimited, got %v", err)
	}
}

func TestMockProviderRecordsCallCount(t *testing.T) {
	mock := llm.NewMockProvider()
	mock.SetCompleteResponse(&llm.CompletionResponse{}, nil)

	if mock.CompleteCallCount() != 0 {
		t.Error("expected 0 calls initially")
	}

	mock.Complete(context.Background(), &llm.CompletionRequest{Model: "m1"})
	mock.Complete(context.Background(), &llm.CompletionRequest{Model: "m2"})

	if mock.CompleteCallCount() != 2 {
		t.Errorf("expected 2 calls, got %d", mock.CompleteCallCount())
	}
}

func TestMockProviderRecordsLastRequest(t *testing.T) {
	mock := llm.NewMockProvider()
	mock.SetCompleteResponse(&llm.CompletionResponse{}, nil)

	req := &llm.CompletionRequest{
		Model:        "test-model",
		SystemPrompt: "Be helpful",
		Messages:     []llm.Message{{Role: "user", Content: "Hello"}},
		MaxTokens:    512,
		SubsystemTag: "test",
	}
	mock.Complete(context.Background(), req)

	last := mock.LastCompleteRequest()
	if last.Model != "test-model" {
		t.Errorf("expected test-model, got %s", last.Model)
	}
	if last.SubsystemTag != "test" {
		t.Errorf("expected subsystem tag test, got %s", last.SubsystemTag)
	}
}

func TestMockProviderStreamEvents(t *testing.T) {
	mock := llm.NewMockProvider()
	events := []llm.StreamEvent{
		{ContentDelta: "Hello"},
		{ContentDelta: " world"},
		{IsFinal: true, InputTokens: 5, OutputTokens: 2},
	}
	mock.SetStreamEvents(events, nil)

	ch, err := mock.Stream(context.Background(), &llm.CompletionRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var received []llm.StreamEvent
	for ev := range ch {
		received = append(received, ev)
	}

	if len(received) != 3 {
		t.Fatalf("expected 3 events, got %d", len(received))
	}
	if received[0].ContentDelta != "Hello" {
		t.Errorf("expected Hello delta, got %s", received[0].ContentDelta)
	}
	if !received[2].IsFinal {
		t.Error("expected final event")
	}
}

func TestMockProviderStreamError(t *testing.T) {
	mock := llm.NewMockProvider()
	mock.SetStreamEvents(nil, llm.ErrProviderUnavailable)

	_, err := mock.Stream(context.Background(), &llm.CompletionRequest{})
	if !errors.Is(err, llm.ErrProviderUnavailable) {
		t.Errorf("expected ErrProviderUnavailable, got %v", err)
	}
}

func TestMockProviderImplementsInterface(t *testing.T) {
	var _ llm.Provider = llm.NewMockProvider()
}
