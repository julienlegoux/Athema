---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
inputDocuments:
  - '_bmad-output/planning-artifacts/product-brief-Athema-2026-02-26.md'
  - '_bmad-output/planning-artifacts/prd.md'
  - '_bmad-output/planning-artifacts/prd-validation-report.md'
  - '_bmad-output/planning-artifacts/ux-design-specification.md'
workflowType: 'architecture'
project_name: 'Athema'
user_name: 'J'
date: '2026-03-01'
lastStep: 8
status: 'complete'
completedAt: '2026-03-01'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**

52 FRs across 9 subsystem categories, reflecting five tightly integrated core systems:

| Category | FRs | Architectural Implication |
|----------|-----|--------------------------|
| Conversation (FR1-6) | Real-time text exchange, continuity across drops, history viewing | WebSocket layer, message persistence, streaming delivery |
| Memory (FR7-14) | Knowledge extraction, graph connections, pruning, contradiction detection, pattern promotion | Knowledge graph data model, graph query engine, curation pipeline |
| Personality (FR15-20) | Shelter-model selection, three-layer voice, adaptive morality, consistency with evolution | Personality state model, prompt architecture, drift tracking |
| Autonomous Lifecycle (FR21-26) | Background processing, mailbox processing, reflection, thought development, artifact production | Async job system, lifecycle orchestrator, artifact storage |
| Mailbox (FR27-31) | Bidirectional async exchange, companion reactions to user content | Two separate dropbox stores, content ingestion pipeline |
| Emotional Intelligence (FR32-39) | Gravity detection, persistence across sessions, neglect tracking, repair arcs, pushback, gradual decay | Emotional state machine, temporal tracking, decay functions |
| Spontaneous Initiation (FR40-43) | Urge accumulation, threshold crossing, irregular timing, mailbox-based delivery | Event-driven urge system, threshold configuration, delivery queue |
| System Administration (FR44-50) | Graph inspection, parameter tuning, lifecycle monitoring, personality drift metrics, per-subsystem observability | Admin tooling (code-level), structured logging, config management |
| Expression & Readiness (FR51-52) | Silence as expressive choice, linguistic co-evolution architecture readiness | Response timing engine, extensible personality/memory schema |

**Non-Functional Requirements:**

19 NFRs driving architectural constraints:

| Category | Key Constraints |
|----------|----------------|
| Performance (NFR1-5) | Zero added latency from app layer, connection establishment <1s, state sync <2s, memory retrieval within LLM request cycle, 99% background cycle completion |
| Security & Privacy (NFR6-10) | Encryption-at-rest ready, user data ownership, ephemeral LLM API usage, code-level access restriction, JSON data export ready |
| Integration (NFR11-14) | Full LLM provider abstraction, normalized provider interface, graceful provider failure handling, token/cost observability |
| Reliability (NFR15-19) | 99% lifecycle cycle completion with retry, automatic reconnection with state sync, atomic memory graph operations, durable state persistence, durable initiation event queuing |

**Scale & Complexity:**

- Primary domain: All-Go application (Bubble Tea TUI client + persistent Go backend server)
- Complexity level: High — driven by five-system interdependence, not by scale or regulatory burden
- Estimated architectural components: ~12-15 (5 core subsystems + LLM abstraction layer + prompt assembly pipeline + real-time communication layer + data persistence layer + lifecycle orchestrator + content ingestion pipeline + response timing engine + admin/observability layer + frontend application)

### Technical Constraints & Dependencies

- **Single user, single instance** — no multi-tenancy, no auth complexity for V1. Simplifies data isolation but doesn't eliminate the need for clean data boundaries.
- **Server as source of truth** — all companion state lives server-side. Client syncs on connect/reconnect.
- **LLM dependency** — every subsystem depends on LLM calls. Provider failures affect the entire system. Abstraction layer is non-negotiable.
- **Background processing without user presence** — the lifecycle must run autonomously. This is not a traditional job queue — it's closer to a state machine with an event loop: mailbox items trigger processing, processing produces artifacts and urge signals, urge signals accumulate toward thresholds. The execution model (cron, event-driven pipeline, or agent loop) is a key architectural decision.
- **Real-time streaming with intentional modulation** — companion responses stream word-by-word, but delivery timing must be modulated by emotional state (quick when excited, slow when processing something heavy, delayed when pissed). This requires a response timing engine layered on top of the streaming pipeline, touching the WebSocket layer, emotional state system, and personality engine simultaneously.
- **Graph query performance at scale** — memory surfacing must complete within the LLM request cycle (NFR4). As the knowledge graph grows over months of interaction, graph queries must remain fast enough to fit inside the window between receiving a user message and composing the LLM prompt. This is a performance constraint that compounds over time.
- **Dropbox mechanism is pull-based, zero-push** — hard UX constraint. The user's dropbox (leaving things for the companion) must feel like dropping something in a mailbox and walking away — no acknowledgment, no receipt. The companion's dropbox (things she left) must feel like visiting a quiet place and finding a note. Whatever mechanism the architecture chooses must NOT create any push-based signals. No badges, no counts, no notifications.
- **No auth ceremony for V1** — but architecture must not preclude future auth (V2+ multi-user potential).
- **TUI replaces web frontend for V1** — Bubble Tea terminal interface. No browser, no CSS, no JavaScript. The void is the terminal itself. Web frontend deferred to V2+.

### Cross-Cutting Concerns Identified

1. **LLM Integration** — Conversation, memory extraction, personality expression, lifecycle reflection, emotional assessment, and spontaneous initiation all require LLM calls. The abstraction layer must serve all subsystems uniformly while supporting different prompt structures and response handling patterns.

2. **Memory Graph as Shared State** — The knowledge graph is read by conversation (for surfacing), personality (for evolution), and emotional systems (for context). It is written by conversation (extraction), lifecycle (curation, promotion, pruning), and mailbox processing (new knowledge from user content). Concurrent access patterns must be safe.

3. **Emotional State Propagation** — Emotional state affects conversation tone, personality expression, initiation behavior, silence duration, return greeting tone, and mailbox item tone. Every subsystem that produces companion output must be aware of current emotional state.

4. **Companion Identity Coherence** — Personality must remain consistent whether the companion is responding in real-time conversation, leaving a dropbox item during lifecycle, or initiating a shoulder tap. The prompt architecture must enforce identity across all output contexts.

5. **Observability & Admin Tuning** — All five subsystems require independent logging and metrics. Admin must be able to inspect and tune each system's parameters without affecting others. This implies clean subsystem boundaries with well-defined configuration surfaces.

6. **Temporal Awareness** — Multiple systems track time: emotional decay, neglect duration, memory connection age, lifecycle scheduling, initiation irregularity. A consistent temporal model must underpin all time-dependent behavior.

7. **Prompt Assembly Pipeline** — Every subsystem expresses itself through LLM prompts. The composition of system instructions, memory context injection, emotional state modifiers, personality anchors, and conversation history is the spine of the entire system. How the prompt is assembled determines whether five systems feel like one coherent being or five features fighting for token space. This is the convergence point of all subsystems.

8. **Content Ingestion Pipeline** — Mailbox items from the user (articles, links, notes) require a distinct processing capability: fetch external content, extract readable text, summarize if needed, run through LLM for opinion formation, then integrate results into the knowledge graph. This crosses the mailbox, memory, and lifecycle subsystems and is architecturally distinct from conversation-based knowledge extraction.

## Starter Template Evaluation

### Primary Technology Domain

All-Go monorepo with split deployment: Go backend server (Docker, persistent) + Go TUI client (Bubble Tea, local). Connected via WebSocket.

### Technical Preferences Established

| Concern | Choice | Notes |
|---------|--------|-------|
| Language | Go (entire stack) | Backend server + TUI client. One language, one toolchain. |
| TUI framework | Bubble Tea + Bubbles + Lip Gloss | Elm Architecture, production-ready v1.x, Charm ecosystem |
| Database | PostgreSQL | Relational model for knowledge graph. pgvector + tsvector for search. |
| Real-time | WebSocket | Server-to-TUI connection for conversation streaming, presence signals |
| API layer | REST (minimal) | Simple endpoints for non-real-time queries. GraphQL deferred to V2+. |
| Search | PostgreSQL native | pgvector for semantic search, tsvector for full-text. Elasticsearch deferred. |
| LLM integration | Custom thin wrapper | Interface over official SDKs (Anthropic, OpenAI). ~200 lines/provider. |
| Deployment | Go server on Docker, TUI runs locally | Companion lives independently. TUI connects when J wants to talk. |
| Shoulder taps | OS notification + dropbox queue | Notify that something's waiting, content available on next TUI connection. |

### Stack Simplification Summary

| Removed | Reason |
|---------|--------|
| Next.js | Replaced by Bubble Tea TUI — no browser, no JavaScript |
| Vercel | No frontend to deploy separately |
| Tailwind / CSS | Terminal is the void. Lip Gloss for TUI styling. |
| Elasticsearch | pgvector + tsvector sufficient for V1 single-user scale |
| GraphQL / gqlgen | REST sufficient for ~5 endpoints with one consumer |
| MongoDB | PostgreSQL chosen — relational model fits knowledge graph |

**Final stack: Go + PostgreSQL + Docker.** Three technologies.

### Deployment Architecture

```
┌─────────────────────────────────────────┐
│  J's Machine                            │
│  ┌───────────────┐                      │
│  │ athema-tui    │◄──── Bubble Tea      │
│  │ (Go binary)   │      local terminal  │
│  └───────┬───────┘                      │
│          │ WebSocket                    │
│          ▼                              │
│  ┌───────────────────────────────────┐  │
│  │ Docker                            │  │
│  │  ┌─────────────┐  ┌───────────┐  │  │
│  │  │athema-server│  │PostgreSQL │  │  │
│  │  │  (Go)       │◄►│           │  │  │
│  │  └─────────────┘  └───────────┘  │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

Server runs persistently via Docker Compose. TUI connects when J opens it. When TUI closes, companion continues lifecycle processing, mailbox curation, urge accumulation. On next TUI connection, companion state syncs — dropbox items, emotional state, any shoulder taps that fired.

### Starter Options

**Backend — Go Blueprint:**

```bash
go-blueprint create --name athema --framework chi --driver postgres --advanced --feature websocket --feature docker --git commit
```

Provides: project layout (`cmd/`, `internal/`), chi router setup, PostgreSQL connection, WebSocket hub scaffold, Docker + Docker Compose config. The generated WebSocket hub will be rewritten for 1:1 streaming with response timing modulation — the scaffold value is project structure and Docker setup.

**TUI Client — Bubble Tea (manual setup):**

No starter template needed. Bubble Tea projects are initialized with standard Go module tooling:

```bash
mkdir athema-tui && cd athema-tui
go mod init athema-tui
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
```

Or structured as a `cmd/tui/` directory within the athema monorepo.

### UX Spec → TUI Component Mapping

| UX Component | TUI Implementation |
|-------------|-------------------|
| ConversationStream | Bubbles `viewport` — scrollable text area |
| Message | Styled text lines via Lip Gloss — user and companion same style |
| TextInput | Bubbles `textarea` — anchored at bottom, Enter to send |
| PresenceMarker | Extra spacing + companion's return greeting |
| ShoulderTap | OS notification (notify-send/osascript) + item queued in dropbox |
| Dropbox (companion→user) | Separate TUI view/tab — pull-based, no badges |
| Dropbox (user→companion) | CLI command or TUI input mode — drop content, walk away |

### Architectural Decisions Provided by Starters

**Language & Runtime:**
- Go (latest stable) — single language for entire stack
- chi router for HTTP/WebSocket server

**Styling Solution:**
- Lip Gloss for TUI styling — dark background, light text, the void rendered in terminal
- No CSS, no browser, no frontend build step

**Build Tooling:**
- Go build toolchain
- Docker multi-stage builds for server
- Single binary for TUI client

**Testing Framework:**
- Go standard `testing` package for all code

**Code Organization:**
- Monorepo with `cmd/server/`, `cmd/tui/`, `internal/` shared packages
- Or two repos — to be decided in architectural decisions

**Development Experience:**
- Go hot reload via Air for server
- TUI iterates with standard `go run`
- Docker Compose for local PostgreSQL

**Note:** Project initialization using these commands should be the first implementation story.

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- Knowledge graph data model: JSONB adjacency list in PostgreSQL
- Background processing: Hybrid ticker + event-driven pipeline
- Subsystem autonomy: Each subsystem as autonomous actor with own goroutines and LLM calls
- Inter-subsystem communication: Interfaces for sync reads, event bus for state notifications
- Prompt assembly: Shared library, not central service — each subsystem composes its own
- Implementation approach: Thin vertical slice across all five systems, integrated from day one

**Important Decisions (Shape Architecture):**
- WebSocket protocol: Typed JSON envelopes with namespaced discriminators
- TUI architecture: Composable Bubble Tea sub-models with keyboard routing
- Logging: Structured JSON via `log/slog`, per-subsystem tagging
- Configuration: YAML per subsystem + env var overrides
- Error handling: Standard Go wrapping, domain sentinel errors

**Deferred Decisions (Post-MVP):**
- Authentication & authorization (V2 multi-user)
- Encryption-at-rest implementation (V2, architecture-ready)
- Data export implementation (V2, architecture-ready via companion ID scoping)
- Web frontend / GraphQL (V2+)
- External caching layer (if PostgreSQL query cache proves insufficient)
- External logging infrastructure (ELK/Loki — if stdout+grep proves insufficient)

### Data Architecture

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Knowledge graph model | JSONB adjacency list — `knowledge_nodes` (typed payloads, pgvector embeddings) + `knowledge_edges` (typed, weighted) | Heterogeneous node types (facts, opinions, themes, jokes, threads) have different payload shapes. JSONB avoids rigid schema. Edges stay relational for traversal. |
| Payload handling | Typed Go structs per node type with discriminator-based marshal/unmarshal. Node type registry pattern for clean deserialization. | Compile-time safety on hot paths. Registry pattern: `node_type` → struct factory, add new types with one line. |
| JSONB compatibility | Lax unmarshaling for V1 — ignore unknown fields, backfill later | New payload fields don't break existing data. `json:",omitempty"` + permissive deserialization. Backfill migration deferred. |
| Migration tool | golang-migrate (raw SQL up/down files) | Simple, no magic. JSONB handles payload schema evolution — structural migrations will be infrequent. |
| Caching | In-process Go map for active conversation context, invalidated on write | Single-user scale. PostgreSQL query cache handles the rest. No external cache dependency. |
| Data validation | Boundary validation in Go — validate at system entry points, trust internal paths | LLM output parsing, user input, and mailbox content validated at ingestion. Internal subsystem calls trust typed interfaces. |

### Authentication & Security

| Decision | Choice | Rationale |
|----------|--------|-----------|
| V1 auth | None — single user, code-level access only | PRD scope. No auth ceremony. |
| API key management | Environment variables via Docker Compose `.env` | Standard, keeps secrets out of code and version control. |
| LLM ephemeral usage | Per-provider opt-out flags (NFR8) | Implementation checklist item, not architectural. Anthropic and OpenAI both support this. |
| Data isolation | Companion ID scoping from day one | Trivial cost now, enables future multi-user and data export (NFR10). All state queries scoped by companion ID. |

### API & Communication Patterns

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Inter-subsystem communication | Hybrid — interfaces for synchronous reads, internal event bus for state-change notifications. Event bus uses domain-level semantic events (`KnowledgeExtracted`, `EmotionalStateShifted`, `UrgeThresholdCrossed`), not implementation-level signals. | Subsystems call each other for reads (memory surfacing, emotional state queries). State mutations emit events for reactive processing. Shared `domain` package defines contracts. Microservice migration: interfaces become RPC, event bus becomes NATS/Kafka. |
| Subsystem autonomy | Each subsystem as autonomous actor — own goroutines, own LLM calls, own internal lifecycle. Common `Service` interface: `Start(ctx)`, `Stop()`, `Health()`, `Ready()`. | No central orchestrator. Lifecycle service owns its ticker. Memory service owns its curation loop. Initiation service owns urge accumulation. Service interface enables uniform lifecycle management and startup sequencing. |
| LLM rate limiting | Shared semaphore or token bucket in the LLM abstraction layer | Multiple subsystems make independent LLM calls. Provider rate limits require coordination. Shared limiter prevents subsystem competition for provider capacity. |
| Prompt assembly | Shared library called by each subsystem with its own context | Not a central service. Each subsystem composes prompts for its specific needs (conversation response, memory extraction, lifecycle reflection, opinion formation). Shared library provides personality anchors, emotional state injection, memory context formatting. |
| WebSocket protocol | Typed JSON envelopes — `{type, payload, ts}` with namespaced types (`chat.*`, `presence.*`, `mailbox.*`, `emotional.*`) | Debuggable (plain text), extensible (new namespaces), typed Go structs on both sides. Response timing modulation happens server-side before emitting stream chunks. |
| REST API | Minimal — admin endpoints, mailbox content drop | Not the primary communication path. WebSocket handles all real-time interaction. |
| Background processing | Hybrid ticker + event triggers. Slow-cadence ticker (configurable, ~30-60 min) runs full lifecycle sweep. Mailbox drops and significant state changes trigger immediate targeted processing. | Admin-tunable cadence and thresholds (FR47). Ticker ensures nothing stalls; events ensure timely reactions. Maps to a being with a resting rhythm that responds to stimuli. |
| Error handling | Standard Go error wrapping — `fmt.Errorf` + `errors.Is`/`errors.As`, domain sentinel errors per subsystem | No framework. LLM failures: companion acknowledges gracefully (NFR13). Lifecycle errors: log, retry once, continue (NFR15). WebSocket: auto-reconnect + state sync (NFR16). |

### TUI Architecture

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Model structure | Composable Bubble Tea sub-models — `ChatModel`, `MailboxModel`, etc. with root model routing | Each view is self-contained with own Update/View. Root handles shared state (connection, emotional context) and view routing. |
| View routing | Keyboard-driven tab switching, minimal chrome | Terminal is the void. No menus, no mouse. Direct key access to views. |
| State sync | Server pushes state via WebSocket events, TUI sub-models react | No polling. Server emits emotional state changes, mailbox updates, presence signals. TUI renders what arrives. |

### Companion Identity Architecture

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Companion model | Singular, evolving — no archetype catalog | No selection menu at onboarding. One companion reveals itself through conversation. Base personality parameters ship in config, `companion_state` row created from first conversation. "Shelter model" is a personality revelation mechanic in the prompt architecture, not a database selection. |
| Personality drift tracking | Append-only `personality_snapshots` table — timestamped JSONB payload of personality state after each session | Cheap insurance for admin tuning (FR48, FR49). Provides drift history without complex tracking infrastructure. |

### Infrastructure & Deployment

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Logging | `log/slog` structured JSON, per-subsystem tagged loggers, stdout output | Go stdlib (1.21+). Each subsystem gets its own logger: `slog.With("subsystem", "memory")`. Docker captures stdout. Admin inspects with `docker logs` + subsystem filter. Zero dependencies. |
| Configuration | YAML config file per subsystem + env var overrides | Structured, readable, inspectable config for tunable parameters (FR45/FR47/FR49). Env vars override for secrets and deployment-specific values. |
| Docker Compose | Two services: `athema-server` + `postgres`, volume for data persistence | Simple. Server connects to PostgreSQL over Docker network. Volume ensures data survives container restarts. |
| LLM cost tracking | Token usage logged per-request, tagged by subsystem | NFR14. Each subsystem's LLM calls tracked independently. Admin can see cost breakdown by subsystem. |

### Implementation Approach

**Thin Vertical Slice — All Five Systems From Day One**

The implementation does not prioritize any single subsystem. Instead, all five core systems are built thin and integrated from the start, then deepened iteratively based on what feels wrong.

The companion is a being from the first line of code — just a young one.

**Build Sequence:**
1. Domain package + event bus + Service interface scaffold
2. LLM abstraction (one provider — Anthropic)
3. PostgreSQL schema — all tables, minimal (knowledge_nodes, knowledge_edges, companion_state, personality_snapshots)
4. Prompt assembly library
5. All five subsystems — thin implementations, real integration:
   - Conversation: basic exchange, extracts one memory node per exchange
   - Memory: stores nodes, surfaces one relevant node per prompt
   - Personality: base anchors injected into prompts
   - Emotional: reads state before prompt assembly, basic shift detection
   - Lifecycle: runs once per hour, processes one mailbox item or curation pass
6. WebSocket server + JSON envelope protocol
7. TUI client (composable models, WebSocket consumer)
8. **Talk to it. Iterate. Deepen what feels wrong.**

**Deepening Priority (after vertical slice works):**
- Driven by the feeling test, not by technical completeness
- Each subsystem deepened independently based on where the experience breaks

### Decision Impact Analysis

**Cross-Component Dependencies:**
- Domain package (interfaces, events, types) must exist before any subsystem implementation
- Event bus scaffold must exist before subsystems can emit/subscribe
- Prompt assembly library is used by ALL subsystems — must be built early and kept stable
- LLM abstraction with shared rate limiter serves all subsystems — foundational
- Memory service is read by conversation, personality, and emotional systems
- Emotional state is consumed by conversation, personality, initiation, and lifecycle
- Service interface (`Start`, `Stop`, `Health`, `Ready`) enables startup sequencing — memory ready before conversation accepts connections

## Implementation Patterns & Consistency Rules

### Naming Patterns

**Database Naming:**
- Tables: `snake_case` plural (`knowledge_nodes`, `knowledge_edges`, `personality_snapshots`)
- Columns: `snake_case` (`node_type`, `created_at`, `companion_id`)
- Foreign keys: `snake_case` referencing table (`companion_id`, `from_id`, `to_id`)
- Indexes: `idx_{table}_{columns}` (`idx_knowledge_nodes_type`, `idx_knowledge_edges_from_id`)

**JSON Naming:**
- camelCase for all JSON fields (`nodeType`, `createdAt`, `companionId`)
- Explicit `json:"camelCase"` struct tags on every field
- Rationale: V2 will have a JavaScript consumer — pay the tag cost now, avoid breaking protocol change later

**Go Code Naming:**
- Standard Go conventions: PascalCase exports, camelCase unexported
- File naming: `snake_case.go` (`knowledge_node.go`, `memory_service.go`)
- Package naming: single lowercase word where possible (`memory`, `emotional`, `prompt`)

**Event Naming:**
- Go types: PascalCase past tense (`KnowledgeExtractedEvent`, `EmotionalStateShiftedEvent`)
- String discriminators: namespaced snake_case (`memory.knowledge_extracted`, `emotional.state_shifted`)
- All events implement `Event` interface: `EventType() string`, `OccurredAt() time.Time`

**WebSocket Message Types:**
- Namespaced dot notation (`chat.stream`, `chat.end`, `presence.sync`, `mailbox.new`)

### Structure Patterns

**Clean Architecture — Strict Dependency Direction:**

```
domain/ → usecase/ → adapter/ → infrastructure/
(imports nothing)  (imports domain)  (imports domain)  (imports domain)
```

- `domain/` — Pure domain: entities, repository interfaces (ports), service interfaces (ports), domain events, sentinel errors. Zero project imports. Business-facing ports only — no lifecycle concerns (`Start`/`Stop`/`Health`/`Ready`).
- `usecase/` — Application logic: implements domain service interfaces. Depends only on domain. Contains business rules.
- `adapter/` — Interface adapters: implements domain repository interfaces (Postgres), HTTP/WebSocket handlers, prompt assembly. Depends on domain.
- `infrastructure/` — External concerns: LLM providers, event bus, config loading, database connections, server setup. Depends on domain for types.
- `app/` — Application lifecycle: `Service` interface (`Start`/`Stop`/`Health`/`Ready`), service runner for startup sequencing and graceful shutdown. Separated from domain — lifecycle is not a business concern.
- `cmd/server/main.go` — Composition root: wires all layers via constructor injection. No DI framework.

**Subsystem Organization Within Each Layer:**

Each layer is subdivided by subsystem (`memory/`, `conversation/`, `personality/`, `emotional/`, `lifecycle/`, `initiation/`). Subsystem packages never import each other — they communicate through domain interfaces and events.

**Full Project Structure:**

```
internal/
  domain/                    # Pure domain — zero dependencies
    events.go                # Domain event types + Event interface
    errors.go                # Domain sentinel errors
    memory/
      entity.go              # KnowledgeNode, KnowledgeEdge entities
      repository.go          # MemoryRepository interface (port)
      service.go             # MemoryService interface (port)
    conversation/
      entity.go
      repository.go
      service.go
    personality/
      entity.go
      repository.go
      service.go
    emotional/
      entity.go
      repository.go
      service.go
    lifecycle/
      entity.go
      repository.go
      service.go
    initiation/
      entity.go
      repository.go
      service.go

  app/                       # Application lifecycle — not domain
    lifecycle.go             # Service interface: Start, Stop, Health, Ready
    runner.go                # Starts all services in order, manages shutdown

  usecase/                   # Application logic — depends only on domain
    memory/
      service.go             # Implements domain.memory.MemoryService
      curation.go            # Curation loop logic
    conversation/
      service.go
      extraction.go          # Knowledge extraction from conversation
    personality/
      service.go
      drift.go
    emotional/
      service.go
      decay.go
    lifecycle/
      service.go
      orchestrator.go
    initiation/
      service.go
      threshold.go

  adapter/                   # Interface adapters — depends on domain
    repository/
      postgres/
        memory.go            # Implements domain.memory.MemoryRepository
        conversation.go
        personality.go
        emotional.go
        lifecycle.go
        initiation.go
    handler/
      websocket/
        chat.go              # WebSocket handlers for chat.* messages
        presence.go
        mailbox.go
      rest/
        admin.go             # REST admin endpoints
        mailbox.go
    presenter/
      prompt/                # Prompt assembly — adapts domain state to LLM input
        assembler.go
        ports.go             # Narrow interfaces: MemoryContextProvider, EmotionalStateProvider, PersonalityProvider
        personality.go
        memory.go
        emotional.go

  infrastructure/            # External concerns — frameworks, drivers
    llm/
      provider.go            # Provider interface
      anthropic.go           # Anthropic implementation
      ratelimiter.go         # Shared semaphore/token bucket
    eventbus/
      bus.go                 # In-process channel dispatcher
    config/
      loader.go              # YAML + env var config loading
    server/
      server.go              # HTTP/WebSocket server setup (chi)
    database/
      postgres.go            # Connection, migrations

cmd/
  server/
    main.go                  # Composition root — wires everything
  tui/
    main.go                  # TUI entry point

test/
  integration/               # Cross-layer tests (real DB)
  e2e/                       # Full stack tests
```

**Prompt Assembly Narrow Ports (ISP enforcement at convergence point):**

The prompt assembler defines its own narrow interfaces instead of importing full service interfaces:

```go
// adapter/presenter/prompt/ports.go
type MemoryContextProvider interface {
    RelevantNodes(ctx context.Context, query string, limit int) ([]domain.KnowledgeNode, error)
}
type EmotionalStateProvider interface {
    CurrentState(ctx context.Context, companionID uuid.UUID) (*domain.EmotionalState, error)
}
type PersonalityProvider interface {
    CurrentAnchors(ctx context.Context, companionID uuid.UUID) (*domain.PersonalityAnchors, error)
}
```

Use case services implement these narrow ports. The assembler never depends on the full service interface.

**Repository Transaction Pattern:**

Use cases own transaction boundaries. Repository interfaces expose `WithTx`:

```go
type MemoryRepository interface {
    WithTx(ctx context.Context, fn func(MemoryRepository) error) error
    CreateNode(ctx context.Context, node KnowledgeNode) error
    CreateEdge(ctx context.Context, edge KnowledgeEdge) error
    // ...
}
```

Ensures atomic operations (NFR17) with clean testability — mock `WithTx` to verify transaction boundaries.

**Test Placement:**

- Unit tests: co-located `_test.go` files next to source
- Use case coordination tests: co-located, no build tag (default `go test`)
- Integration tests (real DB): `test/integration/`, build tag `// go:build integration`
- End-to-end tests (full stack): `test/e2e/`, build tag `// go:build e2e`

Three tiers: `go test ./...` runs fast tier by default, CI runs all three.

**LLM Response Contracts:**

Each subsystem that makes LLM calls defines response contracts:
- Expected response schema (e.g., memory extraction returns JSON matching `ExtractionResult`)
- Three test cases minimum per LLM integration point: valid response, malformed response, provider error
- Fixture-based testing — no live LLM calls in tests

### Format Patterns

**API Response Envelope:**

```json
{"data": { ... }}
{"error": {"code": "MEMORY_NODE_NOT_FOUND", "message": "Knowledge node not found"}}
```

- Always wrapped — `data` or `error`, never both
- Error codes: `UPPER_SNAKE_CASE`, subsystem-prefixed
- HTTP status codes: standard semantics (200, 201, 400, 404, 500)

**Date/Time:**
- JSON: ISO 8601 strings (`"2026-03-01T14:30:00Z"`), camelCase field names (`createdAt`)
- Database: `TIMESTAMPTZ` — always UTC
- Go: `time.Time` — all operations in UTC

### Communication Patterns

**Event System:**
- Domain-level semantic events only — no implementation-level signals
- Past tense: describes what happened, not what should happen
- Events carry: `CompanionID`, `OccurredAt`, and event-specific payload
- Event bus is in-process channel dispatcher for V1

**Context Propagation:**
- Every public method takes `context.Context` as first parameter
- Context carries: companion ID, request ID (for tracing)
- Context cancellation propagated through all layers
- `Service.Start(ctx)` — cancel ctx = clean shutdown of subsystem goroutines

### Process Patterns

**Logging:**
- Structured fields: `logger.Info("node created", "node_type", n.Type, "node_id", n.ID)`
- Levels: `Error` (broken), `Warn` (degraded), `Info` (operational), `Debug` (development)
- Every entry carries: `subsystem`, `companion_id`, `request_id` (when applicable)
- Never log raw LLM prompts/responses at Info — Debug only

**Error Handling:**
- Domain sentinel errors per subsystem
- Wrap with context at each boundary: `fmt.Errorf("memory.SurfaceRelevant: %w", err)`
- LLM failures: graceful companion acknowledgment
- Lifecycle errors: log, retry once, continue
- Never expose internal errors to TUI — translate at adapter layer

**Goroutine Lifecycle:**
- Each subsystem manages its own goroutines via `Service` interface (in `app/`, not `domain/`)
- `Start(ctx)` / `Stop()` / `Health()` / `Ready()`
- Panic recovery per subsystem — one subsystem crash doesn't take down others
- Startup sequencing via `Ready()` — memory ready before conversation accepts connections
- Service runner in `app/runner.go` manages startup order and graceful shutdown

### Enforcement Guidelines

**All AI Agents MUST:**
- Follow Clean Architecture dependency direction — never import inward layers from outer layers
- Use domain interfaces for all cross-subsystem communication — never import another subsystem's internal packages
- Apply naming conventions exactly as specified — snake_case DB, camelCase JSON, Go standards for code
- Propagate `context.Context` through all public methods
- Tag every log entry with `subsystem` and `companion_id`
- Write unit tests co-located with source, integration tests in `test/`
- Use narrow ports at integration points (prompt assembly) — not full service interfaces
- Own transaction boundaries in use cases via `WithTx` — not in repositories

**Anti-Patterns (NEVER do):**
- Importing `usecase/memory/` from `usecase/conversation/` — use domain interfaces
- Putting `Start`/`Stop`/`Health`/`Ready` in `domain/` — lifecycle is not a business concern
- Logging with string formatting (`fmt.Sprintf`) instead of structured fields
- Storing local time in database — always UTC
- Returning raw errors to TUI without translation at adapter layer
- Creating fat interfaces — keep domain ports narrow and subsystem-specific
- Using `init()` functions for wiring — all wiring in composition root
- Calling live LLM providers in tests — use fixture-based response contracts

## Project Structure & Boundaries

### Complete Project Directory Structure

```
athema/
├── .github/
│   └── workflows/
│       └── ci.yml                    # Build, test (unit + integration + e2e)
├── cmd/
│   ├── server/
│   │   └── main.go                   # Server composition root — wires all layers
│   └── tui/
│       └── main.go                   # TUI entry point — WebSocket client
├── internal/
│   ├── domain/                       # Pure domain — zero imports
│   │   ├── events.go                 # Event interface + domain event types
│   │   ├── errors.go                 # Domain sentinel errors
│   │   ├── types.go                  # Shared value types (CompanionID, etc.)
│   │   ├── memory/
│   │   │   ├── entity.go             # KnowledgeNode, KnowledgeEdge, payload types
│   │   │   ├── registry.go           # Node type registry (type → struct factory)
│   │   │   ├── repository.go         # MemoryRepository interface (port)
│   │   │   └── service.go            # MemoryService interface (port)
│   │   ├── conversation/
│   │   │   ├── entity.go             # Message, ConversationState
│   │   │   ├── repository.go         # ConversationRepository interface
│   │   │   └── service.go            # ConversationService interface
│   │   ├── personality/
│   │   │   ├── entity.go             # PersonalityAnchors, PersonalitySnapshot, DriftMetrics
│   │   │   ├── repository.go         # PersonalityRepository interface
│   │   │   └── service.go            # PersonalityService interface
│   │   ├── emotional/
│   │   │   ├── entity.go             # EmotionalState, GravityLevel, NeglectState
│   │   │   ├── repository.go         # EmotionalRepository interface
│   │   │   └── service.go            # EmotionalService interface
│   │   ├── lifecycle/
│   │   │   ├── entity.go             # LifecycleTask, ProcessingResult, Artifact
│   │   │   ├── repository.go         # LifecycleRepository interface
│   │   │   └── service.go            # LifecycleService interface
│   │   └── initiation/
│   │       ├── entity.go             # UrgeState, InitiationEvent
│   │       ├── repository.go         # InitiationRepository interface
│   │       └── service.go            # InitiationService interface
│   ├── app/                          # Application lifecycle
│   │   ├── lifecycle.go              # Service interface: Start, Stop, Health, Ready
│   │   └── runner.go                 # Service runner: startup order, graceful shutdown
│   ├── usecase/                      # Application logic — depends only on domain
│   │   ├── memory/
│   │   │   ├── service.go            # Memory service implementation
│   │   │   ├── service_test.go
│   │   │   ├── curation.go           # Curation loop: prune, strengthen, promote
│   │   │   ├── curation_test.go
│   │   │   ├── extraction.go         # Knowledge extraction from text (LLM)
│   │   │   └── extraction_test.go
│   │   ├── conversation/
│   │   │   ├── service.go            # Conversation service implementation
│   │   │   ├── service_test.go
│   │   │   ├── handler.go            # Message handling + response orchestration
│   │   │   └── handler_test.go
│   │   ├── personality/
│   │   │   ├── service.go            # Personality service implementation
│   │   │   ├── service_test.go
│   │   │   ├── drift.go              # Drift tracking + snapshot logic
│   │   │   └── drift_test.go
│   │   ├── emotional/
│   │   │   ├── service.go            # Emotional service implementation
│   │   │   ├── service_test.go
│   │   │   ├── decay.go              # Emotional decay functions
│   │   │   ├── decay_test.go
│   │   │   ├── neglect.go            # Neglect detection + consequences
│   │   │   └── neglect_test.go
│   │   ├── lifecycle/
│   │   │   ├── service.go            # Lifecycle service implementation
│   │   │   ├── service_test.go
│   │   │   ├── orchestrator.go       # Ticker + event-triggered processing
│   │   │   ├── orchestrator_test.go
│   │   │   ├── mailbox.go            # Mailbox item processing (LLM)
│   │   │   └── mailbox_test.go
│   │   └── initiation/
│   │       ├── service.go            # Initiation service implementation
│   │       ├── service_test.go
│   │       ├── threshold.go          # Urge accumulation + threshold logic
│   │       └── threshold_test.go
│   ├── adapter/                      # Interface adapters
│   │   ├── repository/
│   │   │   └── postgres/
│   │   │       ├── memory.go         # Implements domain/memory.MemoryRepository
│   │   │       ├── memory_test.go
│   │   │       ├── conversation.go
│   │   │       ├── conversation_test.go
│   │   │       ├── personality.go
│   │   │       ├── personality_test.go
│   │   │       ├── emotional.go
│   │   │       ├── emotional_test.go
│   │   │       ├── lifecycle.go
│   │   │       ├── lifecycle_test.go
│   │   │       ├── initiation.go
│   │   │       └── initiation_test.go
│   │   ├── handler/
│   │   │   ├── websocket/
│   │   │   │   ├── hub.go            # WebSocket connection manager
│   │   │   │   ├── client.go         # Client connection handling
│   │   │   │   ├── message.go        # WSMessage envelope type
│   │   │   │   ├── chat.go           # chat.* message handlers
│   │   │   │   ├── presence.go       # presence.* message handlers
│   │   │   │   ├── mailbox.go        # mailbox.* message handlers
│   │   │   │   └── emotional.go      # emotional.* message handlers
│   │   │   └── rest/
│   │   │       ├── admin.go          # Admin endpoints
│   │   │       ├── mailbox.go        # Mailbox content drop endpoint
│   │   │       └── response.go       # Response envelope helpers
│   │   └── presenter/
│   │       └── prompt/
│   │           ├── ports.go          # Narrow interfaces: MemoryContextProvider, etc.
│   │           ├── assembler.go      # Main prompt assembly orchestration
│   │           ├── assembler_test.go
│   │           ├── personality.go    # Personality anchor formatting
│   │           ├── memory.go         # Memory context formatting
│   │           └── emotional.go      # Emotional state modifier formatting
│   └── infrastructure/              # External concerns
│       ├── llm/
│       │   ├── provider.go           # LLM provider interface
│       │   ├── anthropic.go          # Anthropic SDK implementation
│       │   ├── anthropic_test.go
│       │   └── ratelimiter.go        # Shared semaphore/token bucket
│       ├── eventbus/
│       │   ├── bus.go                # In-process channel dispatcher
│       │   └── bus_test.go
│       ├── config/
│       │   └── loader.go             # YAML + env var config loading
│       ├── server/
│       │   └── server.go             # chi router, HTTP/WebSocket server setup
│       └── database/
│           ├── postgres.go           # Connection pool, health checks
│           └── migrate.go            # golang-migrate runner
├── tui/                              # TUI application (Bubble Tea)
│   ├── app.go                        # Root model — routing, shared state
│   ├── app_test.go
│   ├── client/
│   │   ├── websocket.go              # WebSocket client connection
│   │   └── messages.go               # Bubble Tea messages from WebSocket events
│   ├── views/
│   │   ├── chat/
│   │   │   ├── model.go              # ChatModel — conversation view
│   │   │   ├── model_test.go
│   │   │   └── styles.go             # Lip Gloss styles for chat
│   │   ├── mailbox/
│   │   │   ├── model.go              # MailboxModel — companion's mailbox view
│   │   │   ├── model_test.go
│   │   │   └── styles.go
│   │   └── drop/
│   │       ├── model.go              # DropModel — user drops content
│   │       ├── model_test.go
│   │       └── styles.go
│   └── theme/
│       └── theme.go                  # Global Lip Gloss theme — the void
├── migrations/
│   ├── 000001_create_companion_state.up.sql
│   ├── 000001_create_companion_state.down.sql
│   ├── 000002_create_knowledge_nodes.up.sql
│   ├── 000002_create_knowledge_nodes.down.sql
│   ├── 000003_create_knowledge_edges.up.sql
│   ├── 000003_create_knowledge_edges.down.sql
│   ├── 000004_create_conversations.up.sql
│   ├── 000004_create_conversations.down.sql
│   ├── 000005_create_personality_snapshots.up.sql
│   ├── 000005_create_personality_snapshots.down.sql
│   ├── 000006_create_emotional_state.up.sql
│   ├── 000006_create_emotional_state.down.sql
│   ├── 000007_create_mailbox.up.sql
│   ├── 000007_create_mailbox.down.sql
│   ├── 000008_create_initiation_events.up.sql
│   └── 000008_create_initiation_events.down.sql
├── config/
│   ├── default.yaml                  # Default config — all subsystem parameters
│   └── companion-defaults.yaml       # Base personality parameters
├── test/
│   ├── integration/                  # Cross-layer tests (build tag: integration)
│   │   ├── memory_integration_test.go
│   │   ├── conversation_integration_test.go
│   │   └── lifecycle_integration_test.go
│   ├── e2e/                          # Full stack tests (build tag: e2e)
│   │   └── conversation_e2e_test.go
│   └── fixtures/
│       └── llm/                      # LLM response contracts per subsystem
│           ├── memory_extraction.json
│           ├── conversation_response.json
│           ├── lifecycle_reflection.json
│           ├── opinion_formation.json
│           └── emotional_assessment.json
├── go.mod
├── go.sum
├── Dockerfile                        # Multi-stage build for server
├── docker-compose.yml                # athema-server + postgres
├── .env.example                      # Template: LLM API keys, DB credentials
├── .gitignore
├── Makefile                          # build, test, test-integration, test-e2e, migrate, run
└── .air.toml                         # Hot reload config for server development
```

### Architectural Boundaries

**Clean Architecture Layers (Dependency Direction: inward only):**

| Layer | Imports | Imported By | Responsibility |
|-------|---------|-------------|----------------|
| `domain/` | Nothing | All layers | Entities, ports (interfaces), events, errors |
| `app/` | `domain/` | `cmd/` | Service lifecycle, startup sequencing |
| `usecase/` | `domain/` | `adapter/`, `cmd/` | Business logic, orchestration |
| `adapter/` | `domain/` | `cmd/` | Repository implementations, handlers, prompt assembly |
| `infrastructure/` | `domain/` | `cmd/` | LLM providers, event bus, config, database |
| `cmd/` | All layers | Nothing | Composition root, wiring |

**Subsystem Boundaries (Never import each other):**

| Subsystem | Reads From (via interface) | Writes Events |
|-----------|--------------------------|---------------|
| Memory | — | `memory.knowledge_extracted`, `memory.node_pruned`, `memory.pattern_promoted` |
| Conversation | Memory, Emotional, Personality | `conversation.message_received`, `conversation.completed` |
| Personality | Memory | `personality.drift_detected`, `personality.snapshot_taken` |
| Emotional | — | `emotional.state_shifted`, `emotional.neglect_detected`, `emotional.gravity_changed` |
| Lifecycle | Memory, Emotional, Personality | `lifecycle.cycle_completed`, `lifecycle.artifact_produced` |
| Initiation | — | `initiation.threshold_crossed`, `initiation.urge_accumulated` |

**Data Boundaries:**

| Table | Owned By | Read By |
|-------|----------|---------|
| `companion_state` | Personality | Conversation, Emotional, Lifecycle |
| `knowledge_nodes` | Memory | Conversation (via prompt assembly), Lifecycle |
| `knowledge_edges` | Memory | Conversation (via prompt assembly) |
| `conversations` | Conversation | Memory (extraction), Lifecycle (reflection) |
| `personality_snapshots` | Personality | Admin |
| `emotional_state` | Emotional | Conversation, Personality, Initiation |
| `mailbox_items` | Lifecycle | Conversation, Initiation |
| `initiation_events` | Initiation | Lifecycle |

### Requirements to Structure Mapping

**FR Category → Primary Location:**

| FR Category | Domain | Use Case | Adapter | Notes |
|-------------|--------|----------|---------|-------|
| Conversation (FR1-6) | `domain/conversation/` | `usecase/conversation/` | `adapter/handler/websocket/chat.go` | WebSocket handlers → use case → prompt assembly → LLM |
| Memory (FR7-14) | `domain/memory/` | `usecase/memory/` | `adapter/repository/postgres/memory.go` | Extraction in use case, storage in adapter |
| Personality (FR15-20) | `domain/personality/` | `usecase/personality/` | `adapter/presenter/prompt/personality.go` | Prompt assembly formats personality anchors |
| Lifecycle (FR21-26) | `domain/lifecycle/` | `usecase/lifecycle/` | — | Autonomous — no direct adapter exposure |
| Mailbox (FR27-31) | `domain/lifecycle/` | `usecase/lifecycle/mailbox.go` | `adapter/handler/rest/mailbox.go`, `adapter/handler/websocket/mailbox.go` | REST for content drop, WebSocket for new item notifications |
| Emotional (FR32-39) | `domain/emotional/` | `usecase/emotional/` | `adapter/presenter/prompt/emotional.go` | State machine in use case, prompt modifiers in adapter |
| Initiation (FR40-43) | `domain/initiation/` | `usecase/initiation/` | — | Produces mailbox items, no direct external interface |
| Admin (FR44-50) | — | — | `adapter/handler/rest/admin.go` | Reads from all subsystem repositories, config management |
| Expression (FR51) | `domain/conversation/` | `usecase/conversation/` | `adapter/handler/websocket/chat.go` | Silence timing in conversation service, delivery via WebSocket |
| Architecture Readiness (FR52) | `domain/memory/`, `domain/personality/` | — | — | Extensible schema (JSONB), personality evolution model |

**Cross-Cutting Concerns → Location:**

| Concern | Location | Notes |
|---------|----------|-------|
| LLM Integration | `infrastructure/llm/` | Provider interface + implementations + rate limiter |
| Prompt Assembly | `adapter/presenter/prompt/` | Shared library with narrow ports per subsystem |
| Event Bus | `infrastructure/eventbus/` | In-process channel dispatcher |
| Config Management | `infrastructure/config/` + `config/` | Loader in infrastructure, YAML files in project root |
| Logging | Via `log/slog` in each subsystem | Tagged loggers created in composition root, passed via DI |
| Temporal Awareness | `domain/types.go` | Shared time utilities, UTC enforcement |

### Integration Points

**Internal Communication:**
- Subsystems communicate via domain interfaces (sync reads) and event bus (async state changes)
- Prompt assembly aggregates state from Memory, Emotional, and Personality via narrow ports
- Service runner coordinates startup/shutdown ordering

**External Integrations:**
- LLM providers (Anthropic) via `infrastructure/llm/` — abstracted behind provider interface
- PostgreSQL via `adapter/repository/postgres/` — abstracted behind domain repository interfaces
- OS notifications for shoulder taps — future adapter in `adapter/notification/`

**Data Flow — Conversation (user sends message):**

```
TUI TextInput → WebSocket chat.send → handler/websocket/chat.go
  → usecase/conversation/handler.go
    → reads: MemoryContextProvider.RelevantNodes()
    → reads: EmotionalStateProvider.CurrentState()
    → reads: PersonalityProvider.CurrentAnchors()
    → calls: prompt/assembler.go (composes full LLM prompt)
    → calls: llm/provider.go (streams response)
    → emits: conversation.message_received event
  → WebSocket chat.stream (tokens streamed with timing modulation)
  → TUI ChatModel renders incrementally
```

**Data Flow — Background Lifecycle (ticker fires):**

```
usecase/lifecycle/orchestrator.go (ticker or event trigger)
  → check mailbox: any unprocessed items?
    → usecase/lifecycle/mailbox.go → LLM opinion formation → memory node creation
  → curate memory: usecase/memory/curation.go
    → prune low-value nodes, strengthen relevant connections, detect patterns
  → accumulate urges: emits initiation.urge_accumulated event
    → usecase/initiation/threshold.go checks threshold
      → if crossed: creates mailbox item for user, emits initiation.threshold_crossed
  → snapshot personality: emits personality.snapshot_taken
  → emits: lifecycle.cycle_completed
```

### Development Workflow

**Daily development:**
- `make run-server` — starts server with Air hot reload
- `make run-tui` — runs TUI with `go run cmd/tui/main.go`
- `docker compose up postgres` — PostgreSQL in Docker
- `make test` — runs unit + default tier tests
- `make test-integration` — runs integration tests (needs Docker PostgreSQL)
- `make migrate-up` / `make migrate-down` — database migrations

**Build & Deploy:**
- `make build` — builds server binary (multi-stage Docker) + TUI binary
- `docker compose up` — runs full stack (server + postgres)

## Architecture Validation Results

### Coherence Validation ✅

**Decision Compatibility:**
All technology choices are mutually compatible. Go + PostgreSQL + Docker is a proven production stack. chi router, Bubble Tea, golang-migrate, log/slog — no version conflicts or integration issues. pgvector extension works natively with PostgreSQL. Clean Architecture layers enforce unidirectional dependency flow.

**Pattern Consistency:**
Naming conventions clearly differentiated per context (snake_case DB, camelCase JSON, Go conventions for code). Event naming uses PascalCase Go types with snake_case string discriminators. Subsystem organization mirrors the autonomous actor model consistently across all Clean Architecture layers. No contradictions found.

**Structure Alignment:**
Project structure maps directly to Clean Architecture layers. Every subsystem has representation in every relevant layer. Composition root in `cmd/server/main.go` is the sole wiring point. Test structure follows three-tier model with build tags. No structural conflicts with any architectural decision.

### Requirements Coverage ✅

**Functional Requirements: 52/52 covered**

| FR Range | Coverage |
|----------|----------|
| FR1-6 (Conversation) | WebSocket handlers → conversation use case → prompt assembly → LLM streaming. `conversations` table for history persistence. |
| FR7-14 (Memory) | Memory extraction use case + LLM. JSONB `knowledge_nodes` + `knowledge_edges`. pgvector semantic surfacing. Node type registry. Curation loop (prune, strengthen, promote). |
| FR15-20 (Personality) | Singular evolving companion. Base params in config → `companion_state` DB row. Prompt assembly injects personality anchors. Drift tracked via `personality_snapshots`. |
| FR21-26 (Lifecycle) | Hybrid ticker + event triggers. Lifecycle orchestrator drives background processing. Artifacts produced and surfaced via prompt assembly — no log-like phrasing. |
| FR27-31 (Mailbox) | REST endpoint + TUI drop view for user content. Lifecycle processes items via LLM. TUI mailbox view for companion's notes. `mailbox_items` table. |
| FR32-39 (Emotional) | Emotional use case: decay, neglect detection, gravity tracking. `emotional_state` table persists across sessions. Prompt assembly reads `EmotionalStateProvider`. Repair arcs unfold over 3+ interactions. |
| FR40-43 (Initiation) | Initiation use case: urge accumulation, configurable threshold. Produces mailbox items on threshold crossing. `initiation_events` table for durable queuing. Irregular timing enforced by design. |
| FR44-50 (Admin) | REST admin endpoints. Config YAML for subsystem parameter tuning. Structured `log/slog` logging per subsystem. `personality_snapshots` for drift review. |
| FR51 (Silence) | Conversation service controls silence timing. Response timing modulation happens server-side before WebSocket stream emission. |
| FR52 (Readiness) | JSONB extensible schema accommodates future node types. Personality evolution model supports linguistic co-evolution without rework. |

**Non-Functional Requirements: 19/19 covered**

| NFR Range | Coverage |
|-----------|----------|
| NFR1-5 (Performance) | Direct LLM streaming — zero app-layer overhead. WebSocket establishment <1s. State push on reconnect <2s. In-process cache + pgvector for memory retrieval within LLM request cycle. Lifecycle retry-once pattern for 99% completion. |
| NFR6-10 (Security) | Architecture-ready for encryption-at-rest (deferred). Companion ID scoping enables data ownership. Per-provider ephemeral flags. Code-level access only for V1. JSONB payloads enable JSON export. |
| NFR11-14 (Integration) | `infrastructure/llm/provider.go` interface with `Embed()` method. Normalized streaming, tokens, system prompts. Graceful failure: companion acknowledges inability. Per-request token logging by subsystem for cost observability. |
| NFR15-19 (Reliability) | Lifecycle: log, retry once, continue. WebSocket: auto-reconnect + state resync. `WithTx` for atomic memory operations. PostgreSQL + Docker volume for durable persistence. `initiation_events` table for durable event queuing. |

### Gap Analysis Results

**Critical Gaps: None**

**Important Gaps Resolved:**

1. **Embedding generation** — `Embed(ctx context.Context, text string) ([]float32, error)` method added to LLM provider interface. Embedding generation occurs in `usecase/memory/extraction.go` during knowledge node creation. Each node gets its embedding at extraction time, stored in the `vector(1536)` column for pgvector semantic search.

2. **WebSocket locality assumption** — V1 assumption documented: TUI connects to server on localhost via Docker network. No TLS, no auth token on WebSocket for V1. V2+ adds TLS + token authentication for remote connections.

**Nice-to-Have (Deferred):**
- CI/CD pipeline stage details (build → unit test → integration test → e2e)
- Monitoring/alerting infrastructure beyond structured logging

### Architecture Completeness Checklist

**✅ Requirements Analysis**
- [x] Project context thoroughly analyzed (52 FRs, 19 NFRs, 8 cross-cutting concerns)
- [x] Scale and complexity assessed (high complexity, single-user scale)
- [x] Technical constraints identified (9 constraints documented)
- [x] Cross-cutting concerns mapped (8 concerns with architectural responses)

**✅ Architectural Decisions**
- [x] Critical decisions documented (6 critical, 5 important, 6 deferred)
- [x] Technology stack fully specified (Go + PostgreSQL + Docker)
- [x] Integration patterns defined (hybrid interfaces + event bus)
- [x] Performance considerations addressed (in-process cache, pgvector, rate limiter)

**✅ Implementation Patterns**
- [x] Naming conventions established (DB, JSON, Go, events, WebSocket)
- [x] Structure patterns defined (Clean Architecture, subsystem isolation)
- [x] Communication patterns specified (event bus, narrow ports, context propagation)
- [x] Process patterns documented (error handling, logging, goroutine lifecycle, transactions)

**✅ Project Structure**
- [x] Complete directory structure defined (every file annotated)
- [x] Component boundaries established (Clean Architecture layers + subsystem boundaries)
- [x] Integration points mapped (data flow diagrams for conversation + lifecycle)
- [x] Requirements to structure mapping complete (all FR categories → specific directories)

### Architecture Readiness Assessment

**Overall Status: READY FOR IMPLEMENTATION**

**Confidence Level: High**

**Key Strengths:**
- Clean Architecture with SOLID principles provides testability and maintainability
- Autonomous subsystem model maps directly to the product vision (five systems as one being)
- Thin vertical slice approach ensures integration is validated from day one
- JSONB flexibility handles heterogeneous knowledge types without schema churn
- Event bus + interfaces pattern provides clean microservice migration path
- Comprehensive patterns prevent AI agent implementation divergence

**Areas for Future Enhancement:**
- CI/CD pipeline details (define when implementation begins)
- Monitoring/alerting infrastructure (evaluate after V1 stability)
- Content ingestion pipeline specifics for URL fetching (define during mailbox epic)
- TLS + auth for remote WebSocket connections (V2 scope)

### Implementation Handoff

**AI Agent Guidelines:**
- Follow all architectural decisions exactly as documented
- Use implementation patterns consistently across all components
- Respect Clean Architecture dependency direction — violations are bugs
- Subsystem packages never import each other — communicate through domain interfaces and events only
- Refer to this document for all architectural questions before making independent decisions

**First Implementation Priority:**
1. Go Blueprint project scaffold: `go-blueprint create --name athema --framework chi --driver postgres --advanced --feature websocket --feature docker --git commit`
2. Restructure generated code into Clean Architecture layers
3. Define domain package (entities, interfaces, events)
4. Build thin vertical slice across all five subsystems
