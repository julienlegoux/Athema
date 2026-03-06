package llm

import (
	"context"
	"sync"
)

// MockProvider implements the Provider interface for testing.
type MockProvider struct {
	mu sync.Mutex

	completeResp  *CompletionResponse
	completeErr   error
	streamEvents  []StreamEvent
	streamErr     error
	embedResp     *EmbeddingResponse
	embedErr      error

	completeCalls int
	lastCompleteReq *CompletionRequest
	streamCalls   int
	embedCalls    int
}

// NewMockProvider creates a new MockProvider.
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// SetCompleteResponse configures the response returned by Complete.
func (m *MockProvider) SetCompleteResponse(resp *CompletionResponse, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.completeResp = resp
	m.completeErr = err
}

// SetStreamEvents configures the events returned by Stream.
func (m *MockProvider) SetStreamEvents(events []StreamEvent, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.streamEvents = events
	m.streamErr = err
}

// SetEmbedResponse configures the response returned by Embed.
func (m *MockProvider) SetEmbedResponse(resp *EmbeddingResponse, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.embedResp = resp
	m.embedErr = err
}

// Complete returns the configured response.
func (m *MockProvider) Complete(_ context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.completeCalls++
	m.lastCompleteReq = req
	return m.completeResp, m.completeErr
}

// Stream returns a channel that emits the configured events.
func (m *MockProvider) Stream(_ context.Context, _ *CompletionRequest) (<-chan StreamEvent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.streamCalls++

	if m.streamErr != nil {
		return nil, m.streamErr
	}

	ch := make(chan StreamEvent, len(m.streamEvents))
	for _, e := range m.streamEvents {
		ch <- e
	}
	close(ch)
	return ch, nil
}

// Embed returns the configured response.
func (m *MockProvider) Embed(_ context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.embedCalls++
	return m.embedResp, m.embedErr
}

// CompleteCallCount returns the number of times Complete was called.
func (m *MockProvider) CompleteCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.completeCalls
}

// LastCompleteRequest returns the most recent CompletionRequest passed to Complete.
func (m *MockProvider) LastCompleteRequest() *CompletionRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastCompleteReq
}

// StreamCallCount returns the number of times Stream was called.
func (m *MockProvider) StreamCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.streamCalls
}

// EmbedCallCount returns the number of times Embed was called.
func (m *MockProvider) EmbedCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.embedCalls
}
