# Story 1.2: Domain Foundation & Event Bus

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a developer,
I want core domain entities, repository interfaces, service interfaces, and an in-process event bus,
so that subsystems can communicate through well-defined contracts without importing each other.

## Acceptance Criteria

1. **Given** the domain package (`internal/domain/`)
   **When** I define entities
   **Then** Companion, Message, and Conversation entities exist with proper fields and companion ID scoping on all state

2. **Given** the domain package
   **When** I define repository interfaces
   **Then** ConversationRepository and CompanionStateRepository ports are defined with standard CRUD and query methods

3. **Given** the domain package
   **When** I define subsystem business interfaces
   **Then** each subsystem exposes narrow domain interfaces for its public operations (e.g., MemoryContextProvider, EmotionalStateProvider, PersonalityProvider) — these are distinct from the lifecycle `app.Service` interface which already exists in `internal/app/lifecycle.go`

4. **Given** the domain package
   **When** I check imports
   **Then** it has zero project imports — pure domain, no infrastructure, no lifecycle concerns

5. **Given** the event bus (`internal/infrastructure/eventbus/`)
   **When** a subsystem publishes an event
   **Then** all subscribed handlers receive the event asynchronously via in-process channel dispatcher

6. **Given** domain events
   **When** I define the Event interface
   **Then** events implement `EventType() string` and `OccurredAt() time.Time` and carry `CompanionID`

7. **Given** domain events
   **When** I define initial event types
   **Then** `conversation.message_received` and `conversation.completed` events are defined as past-tense semantic events

8. **Given** domain errors
   **When** I define sentinel errors
   **Then** subsystem-specific error types exist for consistent error handling across layers

## Tasks / Subtasks

- [x] **Task 1: Expand domain value types and ID types** (AC: #1, #6)
  - [x] Add `github.com/google/uuid` dependency (`go get github.com/google/uuid`)
  - [x] Refactor `CompanionID` to UUID-based: `type CompanionID uuid.UUID` with `String() string` method and `ParseCompanionID(s string) (CompanionID, error)` constructor
  - [x] Add ID types: `MessageID`, `ConversationID` (same UUID-based pattern with parse constructors)
  - [x] Add `Role` type (string enum: `"user"`, `"companion"`)
  - [x] Keep `UserID`, `SessionID` as string-based (not needed as UUID yet)
  - [x] `google/uuid` is the ONLY allowed external dependency in `internal/domain/` — verify after changes
  - [x] Run `go mod tidy` after adding the dependency

- [x] **Task 2: Create conversation domain entities** (AC: #1)
  - [x] Create `internal/domain/conversation/entities.go`
  - [x] Define `Message` entity: ID, ConversationID, CompanionID, Role, Content, CreatedAt
  - [x] Define `Conversation` entity: ID, CompanionID, CreatedAt, UpdatedAt
  - [x] All fields use explicit `json:"camelCase"` struct tags
  - [x] All timestamps are `time.Time` (TIMESTAMPTZ, UTC in DB)
  - [x] CompanionID scoping on every entity

- [x] **Task 3: Create remaining subsystem domain entities** (AC: #1)
  - [x] `internal/domain/memory/entities.go` — KnowledgeNode (ID, CompanionID, NodeType, Payload as `json.RawMessage`, CreatedAt, UpdatedAt), KnowledgeEdge (ID, FromID, ToID, Type, Weight, CreatedAt)
  - [x] `internal/domain/personality/entities.go` — PersonalitySnapshot (ID, CompanionID, Payload as `json.RawMessage`, CreatedAt)
  - [x] `internal/domain/emotional/entities.go` — EmotionalState (ID, CompanionID, State as `json.RawMessage`, CreatedAt, UpdatedAt)
  - [x] `internal/domain/lifecycle/entities.go` — stub entities only: type declaration + CompanionID + CreatedAt fields, no methods, no tests. Expanded in Epic 4.
  - [x] `internal/domain/initiation/entities.go` — stub entities only: type declaration + CompanionID + CreatedAt fields, no methods, no tests. Expanded in Epic 7.
  - [x] `internal/domain/companion.go` — CompanionState entity (ID as CompanionID, State as `json.RawMessage`, CreatedAt, UpdatedAt)

- [x] **Task 4: Define repository interfaces (ports)** (AC: #2)
  - [x] `internal/domain/conversation/repository.go` — ConversationRepository interface:
    - `WithTx(ctx context.Context, fn func(ConversationRepository) error) error`
    - `CreateConversation(ctx, conversation) error`
    - `CreateMessage(ctx, message) error`
    - `GetConversation(ctx, companionID, conversationID) (*Conversation, error)`
    - `ListMessages(ctx, companionID, conversationID) ([]Message, error)`
    - `GetActiveConversation(ctx, companionID) (*Conversation, error)`
  - [x] `internal/domain/companion.go` — CompanionStateRepository interface:
    - `WithTx(ctx context.Context, fn func(CompanionStateRepository) error) error`
    - `GetState(ctx, companionID) (*CompanionState, error)`
    - `SaveState(ctx, state) error`
  - [x] `internal/domain/memory/repository.go` — MemoryRepository interface (stub: `WithTx` + `CreateNode` + `GetNodeByID` only — full API designed in Epic 3)
  - [x] `internal/domain/personality/repository.go` — PersonalityRepository interface (stub: `WithTx` + `CreateSnapshot` + `GetLatestSnapshot` only)
  - [x] `internal/domain/emotional/repository.go` — EmotionalRepository interface (stub: `WithTx` + `SaveState` + `GetState` only)
  - [x] Every interface method takes `context.Context` as first parameter
  - [x] Every query method scopes by CompanionID

- [x] **Task 5: Define narrow cross-subsystem interfaces** (AC: #3)
  - [x] `internal/domain/memory/provider.go` — MemoryContextProvider interface: `RelevantContext(ctx, companionID, query string) ([]KnowledgeNode, error)`
  - [x] `internal/domain/emotional/provider.go` — EmotionalStateProvider interface: `CurrentState(ctx, companionID) (*EmotionalState, error)`
  - [x] `internal/domain/personality/provider.go` — PersonalityProvider interface: `CurrentPersonality(ctx, companionID) (*PersonalitySnapshot, error)`
  - [x] These are the interfaces conversation uses to read from other subsystems — NO direct subsystem imports

- [x] **Task 6: Expand domain events with CompanionID** (AC: #6, #7)
  - [x] Update `BaseEvent` in `internal/domain/events.go` to include `CompanionID CompanionID` field
  - [x] Update `NewBaseEvent` to accept companionID parameter
  - [x] Define conversation events in `internal/domain/conversation/events.go`:
    - `MessageReceivedEvent` (CompanionID, MessageID, ConversationID, Role, Content)
    - `ConversationCompletedEvent` (CompanionID, ConversationID)
    - String types: `conversation.message_received`, `conversation.completed`
  - [x] Define placeholder event type CONSTANTS in each publishing subsystem's domain package (not centralized — subscribers import the publisher's domain package to reference the constant):
    - `internal/domain/memory/events.go`: `EventKnowledgeExtracted = "memory.knowledge_extracted"`, `EventNodePruned`, `EventPatternPromoted`
    - `internal/domain/personality/events.go`: `EventDriftDetected = "personality.drift_detected"`, `EventSnapshotTaken`
    - `internal/domain/emotional/events.go`: `EventStateShifted = "emotional.state_shifted"`, `EventNeglectDetected`, `EventGravityChanged`
    - `internal/domain/lifecycle/events.go`: `EventCycleCompleted = "lifecycle.cycle_completed"`, `EventArtifactProduced`
    - `internal/domain/initiation/events.go`: `EventThresholdCrossed = "initiation.threshold_crossed"`, `EventUrgeAccumulated`
    - These are string constants ONLY — no full event structs yet (those come in each subsystem's epic)

- [x] **Task 7: Expand sentinel errors per subsystem** (AC: #8)
  - [x] Keep generic errors in `internal/domain/errors.go` (ErrNotFound, etc.)
  - [x] Add `internal/domain/conversation/errors.go` — conversation-specific errors (ErrConversationNotFound, ErrMessageEmpty)
  - [x] Add `internal/domain/memory/errors.go` — memory-specific errors (ErrNodeNotFound, ErrEdgeNotFound)
  - [x] Add `internal/domain/personality/errors.go` — personality-specific errors (ErrSnapshotNotFound)
  - [x] Add `internal/domain/emotional/errors.go` — emotional-specific errors (ErrStateNotFound)
  - [x] All use `fmt.Errorf("subsystem.method: %w", err)` wrapping convention
  - [x] Error codes follow UPPER_SNAKE_CASE with subsystem prefix when exposed via API

- [x] **Task 8: Implement event bus** (AC: #5)
  - [x] Create `internal/infrastructure/eventbus/bus.go`
  - [x] Implement `Bus` struct with channel-based dispatcher
  - [x] `NewBus(logger *slog.Logger) *Bus` constructor (takes slog dependency)
  - [x] `Publish(event domain.Event)` — fans out to all matching subscribers
  - [x] `Subscribe(eventType string, handler func(domain.Event))` — registers handler
  - [x] `Close()` — graceful shutdown of dispatcher goroutines
  - [x] Async delivery: publishing does not block the publisher
  - [x] Buffered channels per subscriber (configurable size, default 256) — log warning on overflow, drop event (do NOT block publisher)
  - [x] Handler panics must be recovered (one bad handler cannot crash the bus)
  - [x] Log handler errors and panic recoveries via slog
  - [x] Bus must be thread-safe — concurrent Publish calls from multiple goroutines

- [x] **Task 9: Unit tests** (AC: all)
  - [x] `internal/domain/types_test.go` — test ID type creation and string conversion
  - [x] `internal/domain/events_test.go` — test BaseEvent creation, EventType(), OccurredAt(), CompanionID presence
  - [x] `internal/domain/conversation/entities_test.go` — test entity JSON marshaling/unmarshaling (camelCase verification)
  - [x] `internal/infrastructure/eventbus/bus_test.go` — HIGHEST TEST PRIORITY (80% of test effort here):
    - publish/subscribe basic flow
    - multiple subscribers on same event type
    - event type filtering (subscriber only receives matching types)
    - handler panic recovery (bus continues after handler panic)
    - close behavior (graceful shutdown, no goroutine leaks)
    - **concurrent publish from multiple goroutines** (thread safety — use sync.WaitGroup + multiple publishers + verify all events received)
    - buffer overflow behavior (publish more events than buffer size, verify warning logged, no block)
  - [x] All tests co-located with source files (no build tags — unit tests run by default)
  - [x] No integration or e2e tests in this story — all tests are unit-level
  - [x] No live LLM calls, no database calls in any test

- [x] **Task 10: Verify constraints and clean up** (AC: #4)
  - [x] Run `go vet ./...` — no issues
  - [x] Run `go test ./...` — all tests pass
  - [x] Verify `internal/domain/` has zero project imports: run `go list -f '{{.Imports}}' ./internal/domain/...` and confirm no `athema/internal` imports appear (only stdlib + google/uuid)
  - [x] Verify no subsystem package imports another subsystem package (e.g., `domain/conversation` must not import `domain/memory`)
  - [x] Run `go mod tidy` and verify go.sum is clean
  - [x] Remove doc.go stubs that have been replaced by real files (keep doc.go only if subsystem has no other files yet)
  - [x] Verify `go build ./...` succeeds (make not available on Windows)
  - [x] Verify `go test ./...` passes

## Dev Notes

### Implementation Order Guidance

**Recommended execution order:** Complete the `conversation` subsystem first (Tasks 1 → 2 → 4 conversation repo → 6 conversation events → 7 conversation errors) since it's the most fleshed out. Then use it as the pattern for all other subsystems. This prevents inconsistency across subsystems. Then do event bus (Task 8) and tests (Task 9) last.

### Critical Architecture Constraints

- **Domain layer (`internal/domain/`) must have ZERO project imports.** Only Go stdlib and `github.com/google/uuid` are allowed. No infrastructure, no lifecycle, no event bus imports in domain.
- **Clean Architecture dependency direction is STRICT:** `domain/ → usecase/ → adapter/ → infrastructure/`. Outer layers import inward only.
- **Subsystem packages NEVER import each other.** `domain/conversation/` must not import `domain/memory/`. Cross-subsystem communication is via narrow provider interfaces defined in the consumed subsystem's package.
- **`internal/app/lifecycle.go` already exists** with the Service interface (Start, Stop, Health, Ready — all with ctx). AC #3 says subsystems have service interfaces — these are the USE CASE service interfaces (business logic), NOT the lifecycle interface. The lifecycle `app.Service` interface is for startup/shutdown. The domain service interfaces are for subsystem business operations.
- **Companion ID scoping on ALL entities from day one.** Every entity carries `CompanionID`. Every query method takes companion ID. This is non-negotiable.
- **Repository interfaces belong in the DOMAIN layer** (they are ports). Implementations go in `internal/adapter/repository/postgres/` (Story 1.4+).
- **Event bus implementation belongs in INFRASTRUCTURE** (`internal/infrastructure/eventbus/`). Domain only defines the Event interface and event types.

### Existing Code to Build On (from Story 1.1)

**Files to MODIFY (not recreate):**
- `internal/domain/types.go` — Currently has `CompanionID string`, `UserID string`, `SessionID string`. Expand with UUID-based types and new ID types.
- `internal/domain/events.go` — Currently has `Event` interface, `BaseEvent`, `NewBaseEvent()`. Add `CompanionID` field to `BaseEvent`, update constructor.
- `internal/domain/errors.go` — Currently has 5 generic sentinel errors. Keep these, add subsystem-specific errors in subsystem packages.
- `internal/infrastructure/eventbus/doc.go` — Replace stub with full implementation.

**Subsystem doc.go stubs to replace with real files:**
- `internal/domain/conversation/doc.go` → entities.go, repository.go, events.go, errors.go
- `internal/domain/memory/doc.go` → entities.go, repository.go, errors.go
- `internal/domain/personality/doc.go` → entities.go, repository.go, provider.go, errors.go
- `internal/domain/emotional/doc.go` → entities.go, repository.go, provider.go, errors.go
- `internal/domain/lifecycle/doc.go` → entities.go (stub)
- `internal/domain/initiation/doc.go` → entities.go (stub)

**Files to NOT touch:**
- `cmd/server/main.go` — Composition root, modified only when wiring new services
- `internal/app/lifecycle.go` — Service interface already correct
- `internal/infrastructure/config/` — Config system complete
- `internal/infrastructure/database/` — DB connection complete
- `internal/infrastructure/server/` — HTTP server complete
- All `internal/usecase/*/doc.go` — Not in scope for this story
- All `internal/adapter/*/doc.go` — Not in scope for this story
- `tui/` — Not in scope for this story

### Event System Design

**Event interface (already exists, needs expansion):**
```go
type Event interface {
    EventType() string
    OccurredAt() time.Time
    GetCompanionID() CompanionID  // ADD THIS
}
```

**BaseEvent (modify existing):**
```go
type BaseEvent struct {
    Type        string      `json:"type"`
    Timestamp   time.Time   `json:"occurredAt"`
    CompanionID CompanionID `json:"companionId"`
}
```

**Event naming convention:**
- Go type names: PascalCase past tense — `MessageReceivedEvent`, `ConversationCompletedEvent`
- String discriminators: `namespace.snake_case` — `conversation.message_received`
- Events describe what HAPPENED (past tense), not what should happen

**Event bus pattern:**
```go
// Publisher (subsystem use case):
bus.Publish(&conversation.MessageReceivedEvent{...})

// Subscriber (another subsystem's use case):
bus.Subscribe("conversation.message_received", func(e domain.Event) {
    msg := e.(*conversation.MessageReceivedEvent)
    // process...
})
```

### Repository Interface Pattern

**WithTx for atomic operations:**
```go
type ConversationRepository interface {
    WithTx(ctx context.Context, fn func(ConversationRepository) error) error
    CreateConversation(ctx context.Context, conv Conversation) error
    // ... other methods
}
```

- Use cases own transaction boundaries — they call `repo.WithTx()`
- Repository implementations receive the tx from context or closure
- This story defines INTERFACES ONLY — implementations are Story 1.4

### JSONB Payload Strategy

For entities with evolving schemas (KnowledgeNode payload, EmotionalState state, PersonalitySnapshot payload):
- Use `json.RawMessage` in Go structs — defer unmarshaling to use case layer
- JSONB in PostgreSQL — no structural migrations needed for payload changes
- Lax unmarshaling in V1: ignore unknown fields, backfill missing fields with defaults

### Naming Conventions (from Story 1.1, continue)

| Context | Convention | Example |
|---------|-----------|---------|
| Go packages | single lowercase word | `memory`, `conversation` |
| Go files | snake_case.go | `entities.go`, `repository.go` |
| Go exports | PascalCase | `ConversationRepository`, `MessageReceivedEvent` |
| Go unexported | camelCase | `handleEvent`, `getState` |
| JSON fields | camelCase with explicit tags | `json:"companionId"` |
| Event types (Go) | PascalCase past tense | `MessageReceivedEvent` |
| Event types (string) | namespace.snake_case | `conversation.message_received` |
| Error vars | ErrPascalCase | `ErrConversationNotFound` |

### What This Story Does NOT Include

- **No database migrations** — Story 1.4 creates tables and repository implementations
- **No use case logic** — Story 1.6 implements conversation use case
- **No HTTP/WebSocket handlers** — Story 1.5 creates WebSocket protocol
- **No LLM integration** — Story 1.3 creates LLM abstraction
- **No TUI changes** — Story 1.7 creates the TUI
- **No adapter layer code** — Only domain interfaces (ports) and infrastructure event bus

### Project Structure Notes

**New files to create:**
```
internal/domain/
├── types.go                          # MODIFY: add UUID types, MessageID, ConversationID, Role
├── events.go                         # MODIFY: add CompanionID to BaseEvent
├── errors.go                         # MODIFY: keep as-is (subsystem errors in subsystem packages)
├── companion.go                      # NEW: CompanionState entity + CompanionStateRepository interface
├── conversation/
│   ├── entities.go                   # NEW: Message, Conversation entities
│   ├── repository.go                 # NEW: ConversationRepository interface
│   ├── events.go                     # NEW: MessageReceivedEvent, ConversationCompletedEvent
│   └── errors.go                     # NEW: conversation-specific errors
├── memory/
│   ├── entities.go                   # NEW: KnowledgeNode, KnowledgeEdge
│   ├── repository.go                 # NEW: MemoryRepository interface (stub)
│   ├── provider.go                   # NEW: MemoryContextProvider interface
│   └── errors.go                     # NEW: memory-specific errors
├── personality/
│   ├── entities.go                   # NEW: PersonalitySnapshot
│   ├── repository.go                 # NEW: PersonalityRepository interface (stub)
│   ├── provider.go                   # NEW: PersonalityProvider interface
│   └── errors.go                     # NEW: personality-specific errors
├── emotional/
│   ├── entities.go                   # NEW: EmotionalState
│   ├── repository.go                 # NEW: EmotionalRepository interface (stub)
│   ├── provider.go                   # NEW: EmotionalStateProvider interface
│   └── errors.go                     # NEW: emotional-specific errors
├── lifecycle/
│   └── entities.go                   # NEW: stub entities (LifecycleTask, ProcessingResult)
└── initiation/
    └── entities.go                   # NEW: stub entities (InitiationEvent, UrgeState)

internal/infrastructure/eventbus/
├── bus.go                            # NEW: channel-based event dispatcher
└── bus_test.go                       # NEW: event bus unit tests
```

**Files to delete (doc.go stubs replaced by real files):**
- `internal/domain/conversation/doc.go`
- `internal/domain/memory/doc.go`
- `internal/domain/personality/doc.go`
- `internal/domain/emotional/doc.go`
- `internal/domain/lifecycle/doc.go`
- `internal/domain/initiation/doc.go`
- `internal/infrastructure/eventbus/doc.go`

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Clean Architecture] — Layer rules, dependency direction
- [Source: _bmad-output/planning-artifacts/architecture.md#Event System] — Event bus design, event types, channel dispatcher
- [Source: _bmad-output/planning-artifacts/architecture.md#Domain Entities] — Entity definitions, companion ID scoping
- [Source: _bmad-output/planning-artifacts/architecture.md#Repository Pattern] — WithTx, CRUD interfaces
- [Source: _bmad-output/planning-artifacts/architecture.md#Naming Conventions] — Go, JSON, event naming
- [Source: _bmad-output/planning-artifacts/architecture.md#Testing] — Three-tier testing, co-located unit tests
- [Source: _bmad-output/planning-artifacts/architecture.md#JSONB Strategy] — json.RawMessage for evolving payloads
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.2] — Acceptance criteria, user story
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#River Model] — Conversation continuity model influencing entity design
- [Source: _bmad-output/implementation-artifacts/1-1-project-scaffold-and-development-environment.md] — Existing code, patterns established, file list

### Previous Story Intelligence (from Story 1.1)

**Key learnings to build on:**
- Go Blueprint scaffolded the project but significant restructuring was done — respect the final structure, not Blueprint defaults
- `database/sql` was replaced with `pgxpool` — use pgx patterns for any future DB work
- Blueprint's WebSocket hub is broadcast-oriented — needs rewrite for 1:1 in Story 1.5, not this story
- Logger factory creates per-subsystem tagged loggers in composition root — event bus should accept `*slog.Logger` via constructor
- Config loader uses YAML + env var overrides — no config changes needed for this story
- `.air.toml` configured for Windows — file watching works for .go files
- All 6 subsystem directories already exist as stubs in domain/ and usecase/ — expand, don't recreate
- `go vet ./...` is available via `.claude/settings.local.json` — use it for verification

**Patterns established in Story 1.1 to follow:**
- Constructor injection (no DI framework, no init() wiring)
- `*slog.Logger` passed as dependency
- Explicit `json:"camelCase"` struct tags on all JSON-facing structs
- `context.Context` as first parameter on all public methods

## Library & Framework Versions

| Library | Version | Import Path | Notes |
|---------|---------|-------------|-------|
| **Go** | 1.26.0 | — | Module version in go.mod |
| **google/uuid** | latest | `github.com/google/uuid` | NEW dependency — add via `go get` |
| **pgx/v5** | v5.8.0 | `github.com/jackc/pgx/v5` | Already in go.mod (not used in this story directly) |

No other new dependencies needed. This story uses only Go stdlib + google/uuid in the domain layer, and Go stdlib + slog in the event bus.

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6

### Debug Log References

- Concurrent publish test initially had buffer overflow (1000 events > 256 buffer). Reduced to 200 events for thread-safety test; buffer overflow covered by dedicated test.
- `make` command not available on Windows — verified with `go build ./...` and `go test ./...` directly.

### Completion Notes List

- All 10 tasks completed successfully with 23 unit tests passing
- Domain layer has zero project imports (only stdlib + google/uuid) — verified via `go list`
- No cross-subsystem imports — each subsystem only imports parent `athema/internal/domain`
- Event bus supports async delivery, buffered channels (256), panic recovery, and graceful shutdown
- All entities have CompanionID scoping and explicit `json:"camelCase"` struct tags
- 7 doc.go stubs removed, replaced with real implementation files
- Full regression suite passes with no failures

### Change Log
| Change | Date | Reason |
|--------|------|--------|
| Story created | 2026-03-02 | Ultimate context engine analysis completed — comprehensive developer guide created |
| Story implementation completed | 2026-03-02 | All 10 tasks implemented: domain entities, repository interfaces, provider interfaces, events, errors, event bus, and unit tests |
| Code review fixes applied | 2026-03-02 | Fixed 8 issues: C1 (ID types JSON marshal as UUID strings), C2 (entity-specific ID types replacing incorrect MessageID usage), H1 (EventPublisher/EventSubscriber port interfaces), H2 (subsystem errors wrap domain errors), H3 (Subscribe after Close guard), M1 (subsystem config env var override), M2 (migrate.go error handling), M3 (sprint-status relative path). 9 new tests added. 32 total tests passing. |

### File List

**New files:**
- `internal/domain/companion.go` — CompanionState entity + CompanionStateRepository interface
- `internal/domain/conversation/entities.go` — Message, Conversation entities
- `internal/domain/conversation/repository.go` — ConversationRepository interface
- `internal/domain/conversation/events.go` — MessageReceivedEvent, ConversationCompletedEvent
- `internal/domain/conversation/errors.go` — conversation-specific sentinel errors
- `internal/domain/conversation/entities_test.go` — JSON marshaling tests
- `internal/domain/memory/entities.go` — KnowledgeNode, KnowledgeEdge entities
- `internal/domain/memory/repository.go` — MemoryRepository interface (stub)
- `internal/domain/memory/provider.go` — MemoryContextProvider interface
- `internal/domain/memory/events.go` — memory event type constants
- `internal/domain/memory/errors.go` — memory-specific sentinel errors
- `internal/domain/personality/entities.go` — PersonalitySnapshot entity
- `internal/domain/personality/repository.go` — PersonalityRepository interface (stub)
- `internal/domain/personality/provider.go` — PersonalityProvider interface
- `internal/domain/personality/events.go` — personality event type constants
- `internal/domain/personality/errors.go` — personality-specific sentinel errors
- `internal/domain/emotional/entities.go` — EmotionalState entity
- `internal/domain/emotional/repository.go` — EmotionalRepository interface (stub)
- `internal/domain/emotional/provider.go` — EmotionalStateProvider interface
- `internal/domain/emotional/events.go` — emotional event type constants
- `internal/domain/emotional/errors.go` — emotional-specific sentinel errors
- `internal/domain/lifecycle/entities.go` — stub entities (LifecycleTask, ProcessingResult)
- `internal/domain/lifecycle/events.go` — lifecycle event type constants
- `internal/domain/initiation/entities.go` — stub entities (InitiationEvent, UrgeState)
- `internal/domain/initiation/events.go` — initiation event type constants
- `internal/domain/types_test.go` — ID type creation and string conversion tests
- `internal/domain/events_test.go` — BaseEvent creation and interface tests
- `internal/infrastructure/eventbus/bus.go` — channel-based event dispatcher
- `internal/infrastructure/eventbus/bus_test.go` — comprehensive event bus tests (8 tests)

**Modified files:**
- `internal/domain/types.go` — refactored CompanionID to UUID-based, added MessageID, ConversationID, Role
- `internal/domain/events.go` — added CompanionID to BaseEvent, GetCompanionID() to Event interface
- `go.mod` — added github.com/google/uuid dependency
- `go.sum` — updated with google/uuid checksums

**Deleted files:**
- `internal/domain/conversation/doc.go` — replaced by real files
- `internal/domain/memory/doc.go` — replaced by real files
- `internal/domain/personality/doc.go` — replaced by real files
- `internal/domain/emotional/doc.go` — replaced by real files
- `internal/domain/lifecycle/doc.go` — replaced by real files
- `internal/domain/initiation/doc.go` — replaced by real files
- `internal/infrastructure/eventbus/doc.go` — replaced by bus.go
