# Story 1.3: LLM Provider Abstraction

Status: done
Created: 2026-03-06
Epic: 1 - Foundation & First Conversation
Story Key: 1-3-llm-provider-abstraction
Dependencies: Story 1.1 (scaffold - done), Story 1.2 (domain foundation - done)
Blocks: Story 1.4 (conversation persistence), Story 1.6 (conversation engine)

## Story

As a developer,
I want an abstracted LLM provider interface with Anthropic as the initial implementation,
So that the companion can generate responses through any LLM provider without coupling to a specific vendor.

## Acceptance Criteria

1. **Given** the provider interface (`internal/infrastructure/llm/`)
   **When** I review it
   **Then** it defines `Complete`, `Stream`, and `Embed` methods with normalized request/response types that hide all provider-specific details

2. **Given** the Anthropic implementation
   **When** I send a completion request
   **Then** it returns a valid response using the Anthropic API, and all requests are made with ephemeral/non-training metadata where supported (NFR8)

3. **Given** the Anthropic implementation
   **When** I send a streaming request
   **Then** it delivers response tokens incrementally through a Go channel with zero application-layer buffering delay (NFR1)

4. **Given** the rate limiter
   **When** multiple subsystems make concurrent LLM requests
   **Then** requests are throttled via semaphore to prevent provider rate-limit exhaustion

5. **Given** token tracking
   **When** any LLM request completes
   **Then** token usage (input, output) is logged via `slog` with a subsystem tag for cost observability (NFR14)

6. **Given** a provider failure (network error, API error, timeout)
   **When** the LLM API returns an error
   **Then** the error is translated into a domain-level sentinel error — no raw provider SDK errors leak to callers (NFR13)

7. **Given** the abstraction
   **When** a new provider needs to be added
   **Then** only the Provider interface needs implementing and wiring in the composition root — zero changes to any subsystem code (NFR11, NFR12)

8. **Given** test fixtures (`test/fixtures/llm/`)
   **When** tests run
   **Then** no live LLM calls are made — a mock provider implementation and fixture-based response contracts are used

## Tasks / Subtasks

- [x] **Task 1: Define LLM domain types and errors** (AC: #1, #6)
  - [x] 1.1 Create `internal/infrastructure/llm/types.go` — define normalized request/response types:
    - `CompletionRequest` struct: Model, SystemPrompt (string), Messages ([]Message), MaxTokens, Temperature, CompanionID, SubsystemTag
    - `Message` struct: Role (string: "user"/"assistant"/"system"), Content (string)
    - `CompletionResponse` struct: Content (string), InputTokens (int), OutputTokens (int), Model (string), StopReason (string)
    - `StreamEvent` struct: ContentDelta (string), IsFinal (bool), InputTokens (int — only on final), OutputTokens (int — only on final)
    - `EmbeddingRequest` struct: Text (string), Model (string)
    - `EmbeddingResponse` struct: Vector ([]float32), Model (string), InputTokens (int)
  - [x] 1.2 Create `internal/infrastructure/llm/errors.go` — define LLM-specific sentinel errors:
    - `ErrProviderUnavailable` — provider API unreachable or returned 5xx
    - `ErrRateLimited` — provider returned 429
    - `ErrAuthenticationFailed` — invalid API key (401/403)
    - `ErrInvalidRequest` — malformed request (400)
    - `ErrContextTooLong` — input exceeds model context window
    - `ErrTimeout` — request exceeded deadline
    - `ErrEmbeddingNotSupported` — provider does not support embeddings
    - All errors use `fmt.Errorf("llm: %w", ...)` wrapping convention matching Story 1.2 patterns

- [x] **Task 2: Define Provider interface** (AC: #1, #7)
  - [x] 2.1 Create `internal/infrastructure/llm/provider.go` — define the `Provider` interface:
    - `Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)`
    - `Stream(ctx context.Context, req *CompletionRequest) (<-chan StreamEvent, error)` — returns a read-only channel; channel is closed when stream ends or errors
    - `Embed(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)`
  - [x] 2.2 Verify: the interface is in the `llm` package under infrastructure — it does NOT import any Anthropic SDK types

- [x] **Task 3: Implement Anthropic provider — Complete method** (AC: #2, #5, #6)
  - [x] 3.1 Add `github.com/anthropics/anthropic-sdk-go` dependency via `go get`
  - [x] 3.2 Create `internal/infrastructure/llm/anthropic.go` — implement `AnthropicProvider` struct:
    - Constructor: `NewAnthropicProvider(apiKey string, model string, logger *slog.Logger, limiter *RateLimiter) *AnthropicProvider`
    - Stores `*anthropic.Client` (created via `anthropic.NewClient(option.WithAPIKey(apiKey))`)
    - Stores model string, logger, limiter references
  - [x] 3.3 Implement `Complete(ctx, req)`:
    - Acquire rate limiter permit before calling API
    - Build `anthropic.MessageNewParams` from normalized `CompletionRequest`:
      - Map `req.SystemPrompt` → `params.System` (`[]anthropic.TextBlockParam`)
      - Map `req.Messages` → `params.Messages` (`[]anthropic.MessageParam`) using `anthropic.NewUserMessage()`/`anthropic.NewAssistantMessage()`
      - Set `params.MaxTokens`, `params.Model`
    - Call `client.Messages.New(ctx, params)`
    - Map `anthropic.Message` response → `CompletionResponse` (extract `Content[0].Text`, `Usage.InputTokens`, `Usage.OutputTokens`)
    - Log token usage via slog with subsystem tag from request
    - Release rate limiter permit
    - On error: translate `*anthropic.Error` → domain sentinel errors via `translateError()`

- [x] **Task 4: Implement Anthropic provider — Stream method** (AC: #3, #5, #6)
  - [x] 4.1 Implement `Stream(ctx, req)`:
    - Acquire rate limiter permit
    - Create output channel: `events := make(chan StreamEvent, 64)`
    - Call `client.Messages.NewStreaming(ctx, params)` — same param building as Complete
    - Launch goroutine to iterate stream:
      - `for stream.Next() { ... }` — process each event
      - On `ContentBlockDeltaEvent` with `TextDelta`: send `StreamEvent{ContentDelta: delta.Text}` to channel
      - After loop ends: send final `StreamEvent{IsFinal: true, InputTokens: ..., OutputTokens: ...}` from accumulated message
      - Close channel when done
      - On `stream.Err()`: send error event or close channel (caller checks channel close + returned error)
    - Log token usage after stream completes
    - Release rate limiter permit after stream completes (in goroutine defer)
    - Return channel and nil immediately (non-blocking for caller)

- [x] **Task 5: Implement Anthropic provider — Embed stub** (AC: #1)
  - [x] 5.1 Implement `Embed(ctx, req)` — return `nil, ErrEmbeddingNotSupported`
  - [x] 5.2 Add code comment: Anthropic does not provide an embeddings API; embedding provider (Voyage AI or OpenAI) will be added in Epic 3 when memory subsystem needs it

- [x] **Task 6: Implement error translation** (AC: #6)
  - [x] 6.1 Create private `translateError(err error) error` function in `anthropic.go`:
    - Use `errors.As(err, &apierr)` to detect `*anthropic.Error`
    - Map HTTP status codes: 401/403 → `ErrAuthenticationFailed`, 429 → `ErrRateLimited`, 400 → `ErrInvalidRequest`, 5xx → `ErrProviderUnavailable`
    - Map context deadline exceeded → `ErrTimeout`
    - Map context canceled → `context.Canceled` (pass through)
    - Default: wrap with `ErrProviderUnavailable` for unknown errors
  - [x] 6.2 Ensure raw `*anthropic.Error` details are logged at debug level before translation (for troubleshooting)

- [x] **Task 7: Implement rate limiter** (AC: #4)
  - [x] 7.1 Create `internal/infrastructure/llm/ratelimiter.go`:
    - `RateLimiter` struct with semaphore (buffered channel)
    - `NewRateLimiter(maxConcurrent int) *RateLimiter` constructor
    - `Acquire(ctx context.Context) error` — blocks until permit available or context canceled
    - `Release()` — returns permit to pool
  - [x] 7.2 Default max concurrent: read from config or use sensible default (e.g., 5)
  - [x] 7.3 Rate limiter is shared across all subsystems — single instance created in composition root

- [x] **Task 8: Wire provider in composition root** (AC: #7)
  - [x] 8.1 Update `cmd/server/main.go`:
    - Read `config.LLM` (Provider, APIKey, Model) — already in config struct
    - Create `RateLimiter` instance
    - Create `AnthropicProvider` instance (pass API key, model, tagged logger, limiter)
    - Log provider initialization at info level
  - [x] 8.2 Do NOT inject provider into any use case yet — Story 1.6 will wire provider to conversation engine
  - [x] 8.3 Validate: server starts successfully with valid config, logs provider creation

- [x] **Task 9: Create mock provider and test fixtures** (AC: #8)
  - [x] 9.1 Create `internal/infrastructure/llm/mock_provider.go`:
    - `MockProvider` struct implementing `Provider` interface
    - Configurable responses: `SetCompleteResponse(*CompletionResponse, error)`, `SetStreamEvents([]StreamEvent, error)`
    - Records calls: `CompleteCallCount()`, `LastCompleteRequest() *CompletionRequest`
  - [x] 9.2 Create `test/fixtures/llm/` directory with fixture JSON files:
    - `completion_response.json` — sample completion with realistic token counts
    - `stream_events.json` — sequence of stream events including final
    - `error_responses.json` — rate limit, auth failure, timeout scenarios
  - [x] 9.3 Create `internal/infrastructure/llm/fixtures.go` — helper to load fixture files: `LoadFixture[T any](path string) (T, error)`

- [x] **Task 10: Unit tests** (AC: all)
  - [x] 10.1 `internal/infrastructure/llm/types_test.go`:
    - Test CompletionRequest construction and field access
    - Test Message role validation
  - [x] 10.2 `internal/infrastructure/llm/errors_test.go`:
    - Test error wrapping: `errors.Is(err, ErrRateLimited)` works through wrapping chain
    - Test all sentinel errors are distinct
  - [x] 10.3 `internal/infrastructure/llm/anthropic_test.go`:
    - Test `translateError()`: each HTTP status maps to correct sentinel error
    - Test `translateError()`: context deadline → ErrTimeout
    - Test `translateError()`: context canceled → context.Canceled
    - Test `Complete()` param building: system prompt mapping, message role mapping, max tokens
    - Test `Stream()` channel behavior: events arrive in order, channel closes on completion
    - Test `Embed()` returns ErrEmbeddingNotSupported
  - [x] 10.4 `internal/infrastructure/llm/ratelimiter_test.go`:
    - Test basic acquire/release flow
    - Test concurrent access up to max (permits exhausted, next blocks)
    - Test context cancellation unblocks waiting acquire
    - Test release after acquire restores capacity
  - [x] 10.5 `internal/infrastructure/llm/mock_provider_test.go`:
    - Test mock returns configured responses
    - Test mock records call counts and last request
  - [x] 10.6 All tests use mock provider or direct struct construction — zero live API calls
  - [x] 10.7 Run `go vet ./...` and `go test ./...` — all pass, no issues

- [x] **Task 11: Verify constraints and clean up** (AC: #7)
  - [x] 11.1 Verify no subsystem imports `anthropic-sdk-go` — only `internal/infrastructure/llm/anthropic.go` imports it
  - [x] 11.2 Verify `internal/infrastructure/llm/provider.go` has zero Anthropic SDK imports
  - [x] 11.3 Remove `internal/infrastructure/llm/doc.go` (replaced by real files)
  - [x] 11.4 Run `go mod tidy` and verify go.sum is clean
  - [x] 11.5 Run `go build ./...` — succeeds
  - [x] 11.6 Run `go test ./...` — all tests pass (existing + new)

## Dev Notes

### Implementation Order Guidance

**Recommended execution order:** Task 1 (types) → Task 2 (interface) → Task 7 (rate limiter) → Task 6 (error translation) → Task 3 (Complete) → Task 4 (Stream) → Task 5 (Embed stub) → Task 9 (mock + fixtures) → Task 10 (tests) → Task 8 (wire composition root) → Task 11 (verify). Build the foundation types and interface first, then the supporting infrastructure (rate limiter, errors), then the implementation, then tests, then wiring.

### Critical Architecture Constraints

- **Clean Architecture layer:** The Provider interface and all implementations live in `internal/infrastructure/llm/`. Infrastructure is the outermost layer. Use cases will depend on the Provider interface, not concrete implementations.
- **Dependency direction:** `domain/ → usecase/ → adapter/ → infrastructure/`. The LLM package may import `domain/` for types like `CompanionID`, but domain NEVER imports infrastructure.
- **No subsystem coupling:** Only `internal/infrastructure/llm/anthropic.go` imports the Anthropic SDK. The Provider interface in `provider.go` must have ZERO SDK imports. Subsystems (conversation, memory, etc.) interact only with the Provider interface.
- **Constructor injection:** Same pattern as Story 1.1/1.2. Provider instance created in `cmd/server/main.go` and injected into consumers. No DI framework, no `init()` wiring.
- **Single shared rate limiter:** One `RateLimiter` instance coordinates all LLM calls across all subsystems. Created in composition root, passed to provider constructor.

### Anthropic Go SDK Reference

**Package:** `github.com/anthropics/anthropic-sdk-go` (requires Go 1.22+, project uses Go 1.26.0)
**Import:** `"github.com/anthropics/anthropic-sdk-go"` and `"github.com/anthropics/anthropic-sdk-go/option"`

**Client creation:**
```go
client := anthropic.NewClient(option.WithAPIKey(apiKey))
```

**Completion call:**
```go
message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
    Model:     anthropic.ModelClaude3_7SonnetLatest, // or use string from config
    MaxTokens: 1024,
    System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
    Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(text))},
})
// Response: message.Content[0] has text, message.Usage has InputTokens/OutputTokens
```

**Streaming call:**
```go
stream := client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{...})
accum := anthropic.Message{}
for stream.Next() {
    event := stream.Current()
    accum.Accumulate(event)
    switch ev := event.AsAny().(type) {
    case anthropic.ContentBlockDeltaEvent:
        switch delta := ev.Delta.AsAny().(type) {
        case anthropic.TextDelta:
            // delta.Text contains the incremental token
        }
    }
}
if stream.Err() != nil { /* handle */ }
// accum.Usage has final token counts
```

**Error handling:**
```go
var apierr *anthropic.Error
if errors.As(err, &apierr) {
    // apierr.StatusCode for HTTP status
    // apierr.DumpRequest(true), apierr.DumpResponse(true) for debugging
}
```

**Built-in retries:** SDK auto-retries connection errors, timeouts, 409, 429, 5xx. Default 2 retries. Configure via `option.WithMaxRetries(n)`.

**Timeouts:** Use `context.WithTimeout()` for overall deadline. Use `option.WithRequestTimeout()` for per-retry timeout.

### Anthropic Does NOT Provide Embeddings

Anthropic has no embeddings API. The `Embed()` method on the Anthropic provider returns `ErrEmbeddingNotSupported`. When memory subsystem needs embeddings (Epic 3, Story 3.2), a separate embedding provider (Voyage AI or OpenAI) will be added implementing the same `Provider` interface, or a dedicated `EmbeddingProvider` interface can be extracted at that point. For now, the interface includes `Embed()` for forward compatibility.

### Config Already in Place

`internal/infrastructure/config/loader.go` already defines:
```go
type LLMConfig struct {
    Provider string `yaml:"provider" env:"ATHEMA_LLM_PROVIDER"`
    APIKey   string `yaml:"api_key" env:"ATHEMA_LLM_API_KEY"`
    Model    string `yaml:"model"   env:"ATHEMA_LLM_MODEL"`
}
```

Default values in `config/default.yaml`:
```yaml
llm:
  provider: "anthropic"
  api_key: ""
  model: "claude-sonnet-4-20250514"
```

Environment override: set `ATHEMA_LLM_API_KEY` in `.env` file or environment.

### Existing Code Patterns (from Stories 1.1 & 1.2)

**Follow these established patterns:**
- Constructor injection: `NewXxx(logger *slog.Logger, ...) *Xxx`
- `*slog.Logger` passed as dependency, subsystem-tagged in composition root
- Explicit `json:"camelCase"` struct tags on JSON-facing structs
- `context.Context` as first parameter on all public methods
- Error wrapping: `fmt.Errorf("llm.Complete: %w", err)`
- Sentinel errors: `var ErrXxx = errors.New("llm: descriptive message")`
- Tests co-located with source (same package `_test` suffix)
- No build tags for unit tests — run by default with `go test ./...`

### What This Story Does NOT Include

- **No use case logic** — Story 1.6 implements conversation use case that calls Provider
- **No prompt assembly** — `adapter/presenter/prompt/` stays as stub (Story 1.6 or Epic 2 concern)
- **No database changes** — No migrations, no persistence
- **No WebSocket changes** — Story 1.5 handles WebSocket protocol
- **No TUI changes** — Story 1.7 handles TUI
- **No embedding provider implementation** — Only the Embed() stub; real implementation deferred to Epic 3
- **No use case injection** — Provider is created in main.go but NOT injected into any use case yet

### NFR Alignment

| NFR | Requirement | Story 1.3 Implementation |
|-----|-------------|------------------------|
| NFR1 | Zero added latency | Stream tokens via channel immediately as received from SDK |
| NFR7 | User data ownership | No data stored by provider layer; API key only in config/env |
| NFR8 | LLM data non-persistence | Document ephemeral flags; Anthropic SDK defaults to no training |
| NFR11 | Full LLM abstraction | Provider interface enables swapping without subsystem changes |
| NFR12 | Normalize provider features | CompletionRequest/Response/StreamEvent hide SDK specifics |
| NFR13 | Graceful provider failures | translateError() maps all SDK errors to domain sentinels |
| NFR14 | Token/cost observability | slog logging with subsystem tag after every request |

### Project Structure Notes

**New files to create:**
```
internal/infrastructure/llm/
    types.go              # Normalized request/response types
    errors.go             # LLM sentinel errors
    provider.go           # Provider interface definition
    anthropic.go          # Anthropic implementation
    ratelimiter.go        # Semaphore-based rate limiter
    mock_provider.go      # Mock for testing
    fixtures.go           # Test fixture loader helper
    types_test.go         # Type construction tests
    errors_test.go        # Error wrapping tests
    anthropic_test.go     # Anthropic provider tests (mocked HTTP)
    ratelimiter_test.go   # Rate limiter tests
    mock_provider_test.go # Mock provider tests

test/fixtures/llm/
    completion_response.json  # Sample completion fixture
    stream_events.json        # Sample stream event sequence
    error_responses.json      # Error scenario fixtures
```

**Files to modify:**
- `cmd/server/main.go` — wire AnthropicProvider creation
- `go.mod` / `go.sum` — add `github.com/anthropics/anthropic-sdk-go` dependency

**Files to delete:**
- `internal/infrastructure/llm/doc.go` — replaced by real files

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#LLM Abstraction] — Provider interface design, ~200 line wrapper target
- [Source: _bmad-output/planning-artifacts/architecture.md#Clean Architecture] — Layer rules, dependency direction
- [Source: _bmad-output/planning-artifacts/architecture.md#Testing] — Three-tier testing, fixture-based LLM tests
- [Source: _bmad-output/planning-artifacts/architecture.md#Configuration] — YAML + env var config pattern
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.3] — Acceptance criteria, user story, NFR alignment
- [Source: _bmad-output/planning-artifacts/prd.md#NFR11-14] — LLM integration non-functional requirements
- [Source: github.com/anthropics/anthropic-sdk-go] — Official Anthropic Go SDK documentation
- [Source: _bmad-output/implementation-artifacts/1-2-domain-foundation-and-event-bus.md] — Previous story patterns, learnings

### Previous Story Intelligence (from Story 1.2)

**Key learnings to build on:**
- `google/uuid` was the only external dependency added to domain — keep LLM SDK confined to infrastructure only
- Constructor injection pattern with `*slog.Logger` works well — follow same for provider
- Event bus has buffered channels (256) with panic recovery — similar channel pattern for streaming
- Tests co-located with source, no build tags for unit tests
- `go vet ./...` and `go test ./...` are the verification commands (no `make` on Windows)
- Code review found issues with JSON marshaling of UUID types — ensure CompletionResponse/StreamEvent JSON tags are correct
- Error wrapping convention established: `fmt.Errorf("subsystem.method: %w", err)` — follow for LLM errors

**Git intelligence (last 5 commits):**
- `af0007f` Add domain-specific ID types and behavioral fixes (Story 1.2 code review fixes)
- `adc1fa6` Add domain foundation, event bus, and scaffolds (Story 1.2 implementation)
- `cf2f300` Update scaffold docs, settings, and sprint status
- `b4257ad` Add project scaffold and dev environment (Story 1.1)
- `cd72fd5` Add epics and sprint-status artifacts

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6 (claude-opus-4-6)

### Debug Log References

No blocking issues encountered during implementation.

### Completion Notes List

- Implemented normalized LLM types (CompletionRequest, CompletionResponse, StreamEvent, EmbeddingRequest, EmbeddingResponse, Message) in types.go
- Defined 7 LLM-specific sentinel errors with "llm:" prefix convention in errors.go
- Created Provider interface with Complete, Stream, Embed methods — zero SDK imports
- Implemented AnthropicProvider with Complete (sync), Stream (async channel-based), and Embed (stub returning ErrEmbeddingNotSupported)
- Implemented translateError mapping all Anthropic SDK errors to domain sentinels (401/403→AuthFailed, 429→RateLimited, 400→InvalidRequest, 5xx→Unavailable, deadline→Timeout)
- Implemented semaphore-based RateLimiter with context-aware Acquire/Release
- Created MockProvider with configurable responses and call recording for consumer testing
- Created test fixture files (completion_response.json, stream_events.json, error_responses.json) and generic LoadFixture[T] helper
- Wired AnthropicProvider in cmd/server/main.go composition root (not injected into use cases yet per story scope)
- Removed doc.go placeholder, ran go mod tidy
- 42 unit tests covering all components — zero live API calls

### File List

**New files:**
- internal/infrastructure/llm/types.go
- internal/infrastructure/llm/errors.go
- internal/infrastructure/llm/provider.go
- internal/infrastructure/llm/anthropic.go
- internal/infrastructure/llm/ratelimiter.go
- internal/infrastructure/llm/mock_provider.go
- internal/infrastructure/llm/fixtures.go
- internal/infrastructure/llm/types_test.go
- internal/infrastructure/llm/errors_test.go
- internal/infrastructure/llm/anthropic_test.go
- internal/infrastructure/llm/ratelimiter_test.go
- internal/infrastructure/llm/mock_provider_test.go
- internal/infrastructure/llm/fixtures_test.go
- test/fixtures/llm/completion_response.json
- test/fixtures/llm/stream_events.json
- test/fixtures/llm/error_responses.json

**Modified files:**
- cmd/server/main.go
- go.mod
- go.sum

**Deleted files:**
- internal/infrastructure/llm/doc.go

### Change Log
| Change | Date | Reason |
|--------|------|--------|
| Story created | 2026-03-06 | Ultimate context engine analysis completed — comprehensive developer guide created |
| Story implemented | 2026-03-06 | All 11 tasks complete, 35 tests passing, full LLM provider abstraction with Anthropic implementation |
| Code review fixes | 2026-03-06 | Fixed 8 issues: stream error signaling (H2), unsupported role warning (H1), fixture tests (H3), Temperature mapping (M1), MaxTokens validation (M3), ErrContextTooLong mapping (M4), provider config validation (M5), magic number extraction (M2) — 42 tests passing |
