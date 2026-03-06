package llm

// CompletionRequest represents a normalized LLM completion request.
type CompletionRequest struct {
	Model        string    `json:"model"`
	SystemPrompt string    `json:"systemPrompt"`
	Messages     []Message `json:"messages"`
	MaxTokens    int       `json:"maxTokens"`
	Temperature  float64   `json:"temperature"`
	CompanionID  string    `json:"companionId"`
	SubsystemTag string    `json:"subsystemTag"`
}

// Message represents a single message in a conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse represents a normalized LLM completion response.
type CompletionResponse struct {
	Content     string `json:"content"`
	InputTokens int    `json:"inputTokens"`
	OutputTokens int   `json:"outputTokens"`
	Model       string `json:"model"`
	StopReason  string `json:"stopReason"`
}

// StreamEvent represents a single event in a streaming LLM response.
// When Error is non-empty, the stream encountered a failure; IsFinal will be true.
type StreamEvent struct {
	ContentDelta string `json:"contentDelta,omitempty"`
	IsFinal      bool   `json:"isFinal"`
	InputTokens  int    `json:"inputTokens,omitempty"`
	OutputTokens int    `json:"outputTokens,omitempty"`
	Error        string `json:"error,omitempty"`
}

// EmbeddingRequest represents a normalized embedding request.
type EmbeddingRequest struct {
	Text  string `json:"text"`
	Model string `json:"model"`
}

// EmbeddingResponse represents a normalized embedding response.
type EmbeddingResponse struct {
	Vector      []float32 `json:"vector"`
	Model       string    `json:"model"`
	InputTokens int       `json:"inputTokens"`
}
