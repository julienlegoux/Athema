package llm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicProvider implements the Provider interface using the Anthropic API.
type AnthropicProvider struct {
	client  *anthropic.Client
	model   string
	logger  *slog.Logger
	limiter *RateLimiter
}

// NewAnthropicProvider creates a new Anthropic provider instance.
func NewAnthropicProvider(apiKey string, model string, logger *slog.Logger, limiter *RateLimiter) *AnthropicProvider {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicProvider{
		client:  &client,
		model:   model,
		logger:  logger,
		limiter: limiter,
	}
}

// Complete sends a completion request to the Anthropic API and returns the response.
func (p *AnthropicProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if req.MaxTokens <= 0 {
		return nil, fmt.Errorf("llm.Complete: %w: MaxTokens must be positive", ErrInvalidRequest)
	}
	if err := p.limiter.Acquire(ctx); err != nil {
		return nil, fmt.Errorf("llm.Complete: %w", err)
	}
	defer p.limiter.Release()

	params := p.buildParams(req)

	message, err := p.client.Messages.New(ctx, params)
	if err != nil {
		p.logger.Debug("anthropic API error", "error", err, "subsystem", req.SubsystemTag)
		return nil, fmt.Errorf("llm.Complete: %w", translateError(err))
	}

	content := ""
	if len(message.Content) > 0 {
		content = message.Content[0].AsText().Text
	}

	p.logger.Info("llm completion",
		"subsystem", req.SubsystemTag,
		"model", string(message.Model),
		"input_tokens", message.Usage.InputTokens,
		"output_tokens", message.Usage.OutputTokens,
	)

	return &CompletionResponse{
		Content:      content,
		InputTokens:  int(message.Usage.InputTokens),
		OutputTokens: int(message.Usage.OutputTokens),
		Model:        string(message.Model),
		StopReason:   string(message.StopReason),
	}, nil
}

// Stream sends a streaming completion request and returns a channel of events.
func (p *AnthropicProvider) Stream(ctx context.Context, req *CompletionRequest) (<-chan StreamEvent, error) {
	if req.MaxTokens <= 0 {
		return nil, fmt.Errorf("llm.Stream: %w: MaxTokens must be positive", ErrInvalidRequest)
	}
	if err := p.limiter.Acquire(ctx); err != nil {
		return nil, fmt.Errorf("llm.Stream: %w", err)
	}

	params := p.buildParams(req)
	stream := p.client.Messages.NewStreaming(ctx, params)

	events := make(chan StreamEvent, 64)

	go func() {
		defer close(events)
		defer p.limiter.Release()

		accum := anthropic.Message{}
		for stream.Next() {
			event := stream.Current()
			if err := accum.Accumulate(event); err != nil {
				p.logger.Debug("stream accumulate error", "error", err, "subsystem", req.SubsystemTag)
				events <- StreamEvent{IsFinal: true, Error: translateError(err).Error()}
				return
			}

			switch ev := event.AsAny().(type) {
			case anthropic.ContentBlockDeltaEvent:
				switch delta := ev.Delta.AsAny().(type) {
				case anthropic.TextDelta:
					events <- StreamEvent{ContentDelta: delta.Text}
				}
			}
		}

		if err := stream.Err(); err != nil {
			p.logger.Debug("anthropic stream error", "error", err, "subsystem", req.SubsystemTag)
			events <- StreamEvent{IsFinal: true, Error: translateError(err).Error()}
			return
		}

		events <- StreamEvent{
			IsFinal:      true,
			InputTokens:  int(accum.Usage.InputTokens),
			OutputTokens: int(accum.Usage.OutputTokens),
		}

		p.logger.Info("llm stream complete",
			"subsystem", req.SubsystemTag,
			"model", string(accum.Model),
			"input_tokens", accum.Usage.InputTokens,
			"output_tokens", accum.Usage.OutputTokens,
		)
	}()

	return events, nil
}

// Embed returns ErrEmbeddingNotSupported as Anthropic does not provide an embeddings API.
// Embedding provider (Voyage AI or OpenAI) will be added in Epic 3 when memory subsystem needs it.
func (p *AnthropicProvider) Embed(_ context.Context, _ *EmbeddingRequest) (*EmbeddingResponse, error) {
	return nil, ErrEmbeddingNotSupported
}

// buildParams constructs Anthropic API parameters from a normalized CompletionRequest.
func (p *AnthropicProvider) buildParams(req *CompletionRequest) anthropic.MessageNewParams {
	model := p.model
	if req.Model != "" {
		model = req.Model
	}

	messages := make([]anthropic.MessageParam, 0, len(req.Messages))
	for _, msg := range req.Messages {
		switch msg.Role {
		case "user":
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
		case "assistant":
			messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
		default:
			p.logger.Warn("unsupported message role ignored", "role", msg.Role, "subsystem", req.SubsystemTag)
		}
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(req.MaxTokens),
		Messages:  messages,
	}

	if req.SystemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: req.SystemPrompt},
		}
	}

	if req.Temperature > 0 {
		params.Temperature = anthropic.Float(req.Temperature)
	}

	return params
}

// translateError converts Anthropic SDK errors into domain sentinel errors.
func translateError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}

	var apierr *anthropic.Error
	if errors.As(err, &apierr) {
		switch apierr.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return ErrAuthenticationFailed
		case http.StatusTooManyRequests:
			return ErrRateLimited
		case http.StatusBadRequest:
			raw := strings.ToLower(apierr.RawJSON())
			if strings.Contains(raw, "too long") || strings.Contains(raw, "too many token") {
				return ErrContextTooLong
			}
			return ErrInvalidRequest
		default:
			if apierr.StatusCode >= 500 {
				return ErrProviderUnavailable
			}
		}
	}

	return ErrProviderUnavailable
}
