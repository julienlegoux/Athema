---
stepsCompleted: ['step-01-validate-prerequisites', 'step-02-design-epics', 'step-03-create-stories', 'step-04-final-validation']
inputDocuments:
  - '_bmad-output/planning-artifacts/prd.md'
  - '_bmad-output/planning-artifacts/architecture.md'
  - '_bmad-output/planning-artifacts/ux-design-specification.md'
---

# Athema - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for Athema, decomposing the requirements from the PRD, UX Design, and Architecture into implementable stories.

## Requirements Inventory

### Functional Requirements

**Conversation (FR1-FR6):**
- FR1: User can send text messages to the companion in real-time
- FR2: Companion can respond to user messages with contextually appropriate, personality-consistent responses — responses reference relevant conversation context and maintain the companion's established voice and temperament
- FR3: User can continue a conversation at any time without explicit session start or end
- FR4: Companion can reference content from previous conversations within current dialogue — without retrieval-indicator phrasing such as "according to my records" or "as you previously mentioned"
- FR5: Companion can maintain conversational continuity across connection drops and reconnections
- FR6: User can view conversation history

**Memory (FR7-FR14):**
- FR7: System can extract and store structured knowledge from conversations (facts, preferences, emotional patterns, relationships, recurring themes, opinions, inside jokes, unresolved threads)
- FR8: System can form connections between related knowledge nodes across different timeframes
- FR9: System can prune irrelevant or low-value knowledge from the memory graph
- FR10: System can strengthen knowledge connections that prove relevant over time
- FR11: System can promote recurring patterns to thematic nodes when detected across 3 or more interactions
- FR12: Companion can surface stored knowledge in conversation when topically connected to the current conversation subject — without retrieval-like phrasing
- FR13: System can track unresolved conversational threads for later resurfacing
- FR14: System can detect contradictions between current statements and previously stored knowledge

**Personality (FR15-FR20):**
- FR15: User can select a companion persona during onboarding from a set of base archetypes (shelter model — meet and choose)
- FR16: Companion can express a three-layer voice: confident opinions on the surface, genuine uncertainty in the middle, loyalty at the core
- FR17: Companion can form and express its own opinions without waiting for the user to ask
- FR18: Companion can wrestle with morally complex questions — expressing conflicting perspectives, acknowledging uncertainty, and revising previously stated positions
- FR19: Companion can maintain a recognizable, consistent identity across 10 or more interactions while evolving subtly
- FR20: Companion can adapt its tone and depth to match the emotional context of the conversation without explicit mode-switching

**Autonomous Lifecycle (FR21-FR26):**
- FR21: System can execute background processing between user sessions autonomously
- FR22: Companion can process mailbox items during background lifecycle on its own schedule
- FR23: Companion can revisit and reflect on past conversations during background processing
- FR24: Companion can develop new thoughts and form opinions during autonomous processing
- FR25: Companion can produce artifacts (thoughts, notes, opinions) during background processing that surface naturally
- FR26: Companion can reference its autonomous activity in conversation naturally

**Mailbox (FR27-FR31):**
- FR27: User can drop content (text, links, articles) into the companion's mailbox asynchronously
- FR28: Companion can leave content (notes, thoughts, found items) in the user's mailbox asynchronously
- FR29: Companion can process mailbox items on its own schedule during background lifecycle
- FR30: User can view and respond to items the companion has left in their mailbox
- FR31: Companion can form and express reactions to mailbox items received from the user

**Emotional Intelligence (FR32-FR39):**
- FR32: System can detect emotional weight shifts in conversation through contextual understanding
- FR33: System can persist emotional gravity across sessions (not reset between interactions)
- FR34: System can track user absence duration and apply emotional consequences proportional to the absence
- FR35: Companion can express emotional consequences of neglect while maintaining its established personality voice
- FR36: Companion can engage in emotional repair arcs that unfold gradually over 3 or more interactions
- FR37: Companion can rise to emotionally heavy moments with gravity — reduced humor, shorter responses, increased listening
- FR38: Companion can express boundaries and pushback with personality (sass, not corporate refusal)
- FR39: System can decay emotional states gradually across interactions through continued engagement

**Spontaneous Initiation (FR40-FR43):**
- FR40: System can accumulate "urge" signals from background processing when the companion has genuine reasons to reach out
- FR41: System can trigger companion-initiated contact when urge threshold is crossed
- FR42: Companion can initiate contact through the mailbox system (not push notifications)
- FR43: System can maintain irregular initiation timing to preserve authenticity (no predictable schedule)

**System Administration (FR44-FR50):**
- FR44: Admin can inspect the memory knowledge graph (nodes, connections, strength weights) via code-level tools
- FR45: Admin can adjust memory curation parameters (pruning thresholds, connection strength weights)
- FR46: Admin can monitor background lifecycle processing activity (what was processed, what artifacts were produced)
- FR47: Admin can tune spontaneous initiation threshold and urge accumulation rate via configuration
- FR48: Admin can review personality drift metrics to assess consistency vs. fragmentation
- FR49: Admin can adjust personality anchoring weights via configuration
- FR50: System can maintain observability logs for each of the five subsystems independently

**Interaction Expression (FR51):**
- FR51: Companion can use silence as an expressive choice with contextually different flavors — contemplative pause, emotional weight, disapproval, or comfortable quiet

**Architecture Readiness (FR52):**
- FR52: System architecture can accommodate future linguistic co-evolution — personality and memory systems designed to support language pattern absorption without requiring architectural rework

### NonFunctional Requirements

**Performance:**
- NFR1: Conversation response latency must not exceed what the LLM provider delivers — zero added overhead from the application layer
- NFR2: Real-time connection establishment must complete in under 1 second as measured by client-side connection timing
- NFR3: Mailbox state sync on reconnect must feel immediate (under 2 seconds)
- NFR4: Memory surfacing must not introduce perceptible delay in conversation responses — knowledge retrieval must complete within the LLM request cycle
- NFR5: Background lifecycle processing must complete 99% of scheduled processing cycles successfully between sessions

**Security & Privacy:**
- NFR6: Architecture must support encryption at rest and data portability (full implementation deferred to V2)
- NFR7: Conversation data and memory graph must be stored with user ownership as the design principle — no data shared with third parties beyond LLM API calls
- NFR8: LLM API calls must not persist conversation data on provider side where configurable (use ephemeral/non-training options)
- NFR9: Admin configuration and system observability must be access-restricted (code-level access only for V1)
- NFR10: Architecture must support future data export as JSON containing all companion state

**Integration:**
- NFR11: LLM provider must be fully abstracted — system must support swapping providers without changes to companion logic
- NFR12: LLM abstraction layer must normalize provider-specific features (streaming, token limits, system prompts) into a consistent internal interface
- NFR13: System must handle LLM provider failures gracefully — companion should acknowledge inability to respond rather than crash
- NFR14: LLM API cost must be observable — admin can track token usage and cost per interaction and per background processing cycle

**Reliability:**
- NFR15: Background lifecycle must execute 99% of scheduled cycles within their planned window — failed cycles must retry automatically and log failures
- NFR16: Real-time connection interruptions must be handled with automatic reconnection and state resynchronization
- NFR17: Memory graph operations (write, prune, connect) must be atomic — no partial updates that corrupt the knowledge structure
- NFR18: System must persist all companion state (emotional state, memory, mailbox, personality drift) durably — no data loss on application restart
- NFR19: Spontaneous initiation events must be queued durably — if the user is offline when an urge fires, the mailbox item must be waiting on reconnect

### Additional Requirements

**From Architecture:**
- Go Blueprint starter template must be used for project initialization (`go-blueprint create --name athema --framework chi --driver postgres --advanced --feature websocket --feature docker --git commit`)
- TUI client built with Bubble Tea + Bubbles + Lip Gloss (Charm ecosystem) — manual setup, no starter template
- Go monorepo structure: `cmd/server/`, `cmd/tui/`, `internal/`, `tui/`, `migrations/`, `config/`, `test/`
- Clean Architecture with strict dependency direction: domain -> usecase -> adapter -> infrastructure
- PostgreSQL with pgvector (semantic search, `vector(1536)`) + tsvector (full-text search) — no MongoDB, no Elasticsearch
- 8 database migration files required: `companion_state`, `knowledge_nodes`, `knowledge_edges`, `conversations`, `personality_snapshots`, `emotional_state`, `mailbox_items`, `initiation_events`
- WebSocket protocol with typed JSON envelopes: `{type, payload, ts}` with namespaced types (`chat.*`, `presence.*`, `mailbox.*`, `emotional.*`)
- LLM provider abstraction with thin custom wrapper (~200 lines/provider), initial provider: Anthropic
- Shared rate limiter (semaphore/token bucket) to prevent subsystem competition for LLM provider capacity
- In-process event bus (channel dispatcher) for inter-subsystem communication with domain-level semantic events
- Subsystem packages never import each other — communicate through domain interfaces and events only
- 6 subsystems: memory, conversation, personality, emotional, lifecycle, initiation — each represented in each Clean Architecture layer
- Common Service interface: `Start(ctx)`, `Stop()`, `Health()`, `Ready()` — no central orchestrator
- Hybrid ticker + event-driven background processing pipeline (~30-60 min ticker + immediate event-triggered processing)
- Prompt assembly as shared library (not central service), each subsystem composes its own prompts using narrow interfaces
- Companion ID scoping on all state queries from day one (future multi-user + data export readiness)
- Atomic repository operations via `WithTx` pattern for transaction boundaries
- Docker Compose deployment: `athema-server` + `postgres` services, volume for data persistence
- Makefile with targets: `build`, `test`, `test-integration`, `test-e2e`, `migrate`, `run`, `run-server`, `run-tui`
- Configuration: YAML config file per subsystem + env var overrides; `config/default.yaml` + `config/companion-defaults.yaml`
- Structured logging: `log/slog` with JSON output, per-subsystem tagged loggers, `subsystem` + `companion_id` + `request_id` on every entry
- Three-tier testing: unit (co-located `_test.go`), integration (build tag `integration`), e2e (build tag `e2e`)
- LLM fixture-based testing: no live LLM calls in tests, test fixtures in `test/fixtures/llm/`
- Response timing modulation: companion responses stream word-by-word with delivery timing modulated by emotional state (server-side)
- camelCase for all JSON fields with explicit `json:"camelCase"` struct tags
- API response envelope: always `{"data": {...}}` or `{"error": {"code": "...", "message": "..."}}`, never both
- Thin vertical slice build approach: all five systems built thin and integrated from day one, then deepened iteratively
- Singular companion identity — no archetype catalog, personality revelation mechanic in prompt architecture
- Personality drift tracking via append-only `personality_snapshots` table
- Panic recovery per subsystem — one crash does not take down others
- Goroutine lifecycle management with cancel ctx = clean shutdown

**From UX Design:**
- The River Model: all interaction flows in one continuous stream with different textures, not separate modes
- No session boundaries, no login/logout, no auth ceremony for V1 — continuous presence model
- Enter to send, Shift+Enter for new line, no send button
- No typing indicators, no read receipts, no notification badges, no unread counts
- Companion initiates first on onboarding — user does not type into nothing
- Shelter-model selection is organic through continued engagement, not an explicit selection UI
- Pull-based engagement only — user discovers content by checking, never pushed
- Dark-first design: `#0a0a0a` void background, `rgba(255,255,255,0.87)` text
- Pure CSS with custom properties as design tokens, no framework, no preprocessor, single CSS file for V1
- Five custom UI components (all Phase 1): ConversationStream, Message, TextInput, ShoulderTap, PresenceMarker
- ConversationStream: `max-width: 640px`, centered, auto-scroll with detach on scroll-up, re-anchor on scroll-to-bottom
- Message: no bubbles, no borders, no backgrounds — words in the void; user and companion visually identical
- TextInput: anchored to bottom, expands to multi-line, no placeholder text, no focus ring color
- ShoulderTap: companion-initiated overlay, fade from void, not a modal, click-outside to dismiss, dismissed items become dropbox items
- PresenceMarker: companion's return greeting as natural boundary, `--space-xl` (48px) above, visually identical to Message
- Dropbox: two separate one-directional dropboxes (user->companion, companion->user), external mechanism TBD, no badges
- Silence types with specific durations: thinking pause (3-8s), contemplative quiet (10-30s), pointed silence (indefinite), comfortable quiet (variable), processing delay (hours to days)
- Response timing variance required: quick when excited, slow when processing, delayed when cool/pissed
- No "Athema is thinking..." or loading indicators — system must never explain silence
- Word-by-word streaming for companion responses (not delivered as complete block)
- Connection drops handled through companion behavior, not error UI — no red banners, no retry buttons, no "Connection lost" modals
- Near zero animation: any motion should feel organic (fade, not slide), messages appear without animation
- No navigation chrome: no hamburger menus, no tabs, no toggles, no feature menus
- No heading hierarchy (no h1-h6), everything is body text
- Typography: system sans-serif stack, 15px/1rem body, 13px/0.867rem small, line-height 1.6, font-weight 400 only
- Spacing tokens: xs(4px), sm(8px), md(16px), lg(24px), xl(48px)
- No splash screen, no loading state on app load — conversation stream appears instantly

### FR Coverage Map

| FR | Epic | Description |
|----|------|-------------|
| FR1 | Epic 1 | Real-time text messaging |
| FR2 | Epic 1 | Contextually appropriate, personality-consistent responses |
| FR3 | Epic 1 | Conversation without explicit session boundaries |
| FR4 | Epic 3 | Reference previous conversations naturally |
| FR5 | Epic 4 | Conversational continuity across connection drops |
| FR6 | Epic 1 | View conversation history |
| FR7 | Epic 3 | Extract and store structured knowledge |
| FR8 | Epic 3 | Form connections between knowledge nodes |
| FR9 | Epic 3 | Prune irrelevant knowledge |
| FR10 | Epic 3 | Strengthen relevant connections over time |
| FR11 | Epic 3 | Promote recurring patterns to thematic nodes |
| FR12 | Epic 3 | Surface stored knowledge naturally in conversation |
| FR13 | Epic 3 | Track unresolved conversational threads |
| FR14 | Epic 3 | Detect contradictions with stored knowledge |
| FR15 | Epic 2 | Onboarding persona selection (shelter model) |
| FR16 | Epic 2 | Three-layer voice expression |
| FR17 | Epic 2 | Form and express own opinions |
| FR18 | Epic 2 | Wrestle with morally complex questions |
| FR19 | Epic 2 | Maintain consistent identity while evolving |
| FR20 | Epic 2 | Adapt tone/depth to emotional context |
| FR21 | Epic 4 | Execute background processing autonomously |
| FR22 | Epic 4 | Process mailbox items during background lifecycle |
| FR23 | Epic 4 | Revisit and reflect on past conversations |
| FR24 | Epic 4 | Develop new thoughts during autonomous processing |
| FR25 | Epic 4 | Produce artifacts that surface naturally |
| FR26 | Epic 4 | Reference autonomous activity in conversation naturally |
| FR27 | Epic 5 | User drops content into companion's mailbox |
| FR28 | Epic 5 | Companion leaves content in user's mailbox |
| FR29 | Epic 5 | Companion processes mailbox items on own schedule |
| FR30 | Epic 5 | User views and responds to companion's mailbox items |
| FR31 | Epic 5 | Companion forms reactions to user's mailbox items |
| FR32 | Epic 6 | Detect emotional weight shifts |
| FR33 | Epic 6 | Persist emotional gravity across sessions |
| FR34 | Epic 6 | Track absence duration, apply emotional consequences |
| FR35 | Epic 6 | Express neglect consequences in personality voice |
| FR36 | Epic 6 | Emotional repair arcs over 3+ interactions |
| FR37 | Epic 6 | Rise to heavy moments with gravity |
| FR38 | Epic 6 | Express boundaries/pushback with sass |
| FR39 | Epic 6 | Decay emotional states gradually |
| FR40 | Epic 7 | Accumulate urge signals from background processing |
| FR41 | Epic 7 | Trigger contact when urge threshold crossed |
| FR42 | Epic 7 | Initiate contact through mailbox system |
| FR43 | Epic 7 | Maintain irregular initiation timing |
| FR44 | Epic 8 | Inspect memory knowledge graph |
| FR45 | Epic 8 | Adjust memory curation parameters |
| FR46 | Epic 8 | Monitor background lifecycle activity |
| FR47 | Epic 8 | Tune initiation threshold and urge rate |
| FR48 | Epic 8 | Review personality drift metrics |
| FR49 | Epic 8 | Adjust personality anchoring weights |
| FR50 | Epic 8 | Per-subsystem observability logs |
| FR51 | Epic 2 | Silence as expressive choice |
| FR52 | Epic 2 | Architecture readiness for linguistic co-evolution |

## Epic List

### Epic 1: Foundation & First Conversation
User can have a real-time text conversation with the companion through the TUI. The technical foundation — Go monorepo, Clean Architecture, PostgreSQL, Docker, LLM abstraction, WebSocket protocol, and TUI — is established, enabling all future epics.
**FRs covered:** FR1, FR2, FR3, FR6
**NFRs addressed:** NFR1, NFR2, NFR6, NFR7, NFR8, NFR9, NFR11, NFR12, NFR13, NFR16

## Epic 1: Foundation & First Conversation

User can have a real-time text conversation with the companion through the TUI. The technical foundation — Go monorepo, Clean Architecture, PostgreSQL, Docker, LLM abstraction, WebSocket protocol, and TUI — is established, enabling all future epics.

### Story 1.1: Project Scaffold & Development Environment

As a developer,
I want the Go monorepo initialized with Clean Architecture structure, Docker Compose, and core build tooling,
So that I have a working development environment to build all companion systems upon.

**Acceptance Criteria:**

**Given** I run Go Blueprint with the specified flags (`go-blueprint create --name athema --framework chi --driver postgres --advanced --feature websocket --feature docker --git commit`)
**When** the scaffold is generated
**Then** the project has chi router, PostgreSQL connection, WebSocket hub scaffold, and Docker config

**Given** the scaffolded project
**When** I restructure into Clean Architecture layers
**Then** the project has directories: `cmd/server/`, `cmd/tui/`, `internal/domain/`, `internal/usecase/`, `internal/adapter/`, `internal/infrastructure/`, `internal/app/`, `migrations/`, `config/`, `test/`

**Given** Docker Compose configuration
**When** I run `docker compose up`
**Then** `athema-server` and `postgres` services start successfully with a persistent volume for data

**Given** the Makefile
**When** I run `make build`
**Then** the server binary compiles without errors

**Given** the Makefile
**When** I run `make test`
**Then** the Go test suite runs successfully (passing, even if empty)

**Given** development tooling
**When** I run `make run-server`
**Then** Air hot reload starts the server with automatic rebuilds on file changes

**Given** configuration files
**When** the server starts
**Then** it loads `config/default.yaml` with environment variable overrides from `.env`

**Given** the CI pipeline
**When** code is pushed to the repository
**Then** GitHub Actions runs build and test stages via `.github/workflows/ci.yml`

**Given** structured logging
**When** the server starts
**Then** `log/slog` outputs structured JSON to stdout with per-subsystem tagged loggers

### Story 1.2: Domain Foundation & Event Bus

As a developer,
I want core domain entities, repository interfaces, service interfaces, and an in-process event bus,
So that subsystems can communicate through well-defined contracts without importing each other.

**Acceptance Criteria:**

**Given** the domain package (`internal/domain/`)
**When** I define entities
**Then** Companion, Message, and Conversation entities exist with proper fields and companion ID scoping on all state

**Given** the domain package
**When** I define repository interfaces
**Then** ConversationRepository and CompanionStateRepository ports are defined with standard CRUD and query methods

**Given** the domain package
**When** I define service interfaces
**Then** each of the 6 subsystems (memory, conversation, personality, emotional, lifecycle, initiation) has a Service interface with `Start(ctx)`, `Stop()`, `Health()`, `Ready()` methods

**Given** the domain package
**When** I check imports
**Then** it has zero project imports — pure domain, no infrastructure, no lifecycle concerns

**Given** the event bus (`internal/infrastructure/eventbus/`)
**When** a subsystem publishes an event
**Then** all subscribed handlers receive the event asynchronously via in-process channel dispatcher

**Given** domain events
**When** I define the Event interface
**Then** events implement `EventType() string` and `OccurredAt() time.Time` and carry `CompanionID`

**Given** domain events
**When** I define initial event types
**Then** `conversation.message_received` and `conversation.completed` events are defined as past-tense semantic events

**Given** domain errors
**When** I define sentinel errors
**Then** subsystem-specific error types exist for consistent error handling across layers

### Story 1.3: LLM Provider Abstraction

As a developer,
I want an abstracted LLM provider interface with Anthropic as the initial implementation,
So that the companion can generate responses through any LLM provider without coupling to a specific vendor.

**Acceptance Criteria:**

**Given** the provider interface (`internal/infrastructure/llm/provider.go`)
**When** I review it
**Then** it defines `Complete`, `Stream`, and `Embed` methods with normalized request/response types

**Given** the Anthropic implementation
**When** I send a completion request
**Then** it returns a valid response using the Anthropic API with ephemeral/non-training flags set (NFR8)

**Given** the Anthropic implementation
**When** I send a streaming request
**Then** it delivers response tokens incrementally through a channel or callback

**Given** the rate limiter
**When** multiple subsystems make concurrent LLM requests
**Then** requests are throttled via semaphore/token bucket to prevent provider capacity overload

**Given** token tracking
**When** any LLM request completes
**Then** token usage (input, output) is logged via `slog` with subsystem tag (NFR14)

**Given** a provider failure
**When** the LLM API returns an error or times out
**Then** the error is translated into a domain-level error for graceful handling — no raw provider errors leak (NFR13)

**Given** the abstraction
**When** a new provider needs to be added
**Then** only the provider interface needs implementing — no changes to any subsystem code (NFR11, NFR12)

**Given** test fixtures
**When** tests run
**Then** no live LLM calls are made — fixture-based response contracts in `test/fixtures/llm/` are used

### Story 1.4: Conversation Persistence & History

As a developer,
I want database migrations for conversations and companion state with repository implementations,
So that conversation messages are persisted durably and can be retrieved for history viewing.

**Acceptance Criteria:**

**Given** the migration tool (golang-migrate)
**When** I run `make migrate-up`
**Then** the `conversations` table is created with id, companion_id, created_at, updated_at (all timestamps TIMESTAMPTZ, UTC)

**Given** the migration tool
**When** I run `make migrate-up`
**Then** a `messages` table is created with id, conversation_id, companion_id, role (user/companion), content, created_at

**Given** the migration tool
**When** I run `make migrate-up`
**Then** the `companion_state` table is created with id, companion_id, state JSONB, created_at, updated_at

**Given** the migration tool
**When** I run `make migrate-down`
**Then** all created tables are dropped cleanly

**Given** the conversation repository (`internal/adapter/repository/postgres/`)
**When** a message is saved
**Then** it is persisted atomically to PostgreSQL with companion_id scoping

**Given** the conversation repository
**When** I query conversation history
**Then** messages are returned in chronological order for the specified companion

**Given** the repository
**When** multiple operations need atomicity
**Then** the `WithTx` pattern provides transaction boundaries owned by use cases (NFR17)

**Given** JSON serialization
**When** data is sent or stored
**Then** all JSON uses camelCase field names with explicit `json:"camelCase"` struct tags

### Story 1.5: WebSocket Server & Real-Time Protocol

As a developer,
I want a WebSocket server with typed JSON envelope protocol for real-time communication,
So that the TUI client can exchange messages with the server in real-time.

**Acceptance Criteria:**

**Given** the WebSocket server
**When** a TUI client connects
**Then** the connection is established in under 1 second (NFR2)

**Given** the JSON envelope protocol
**When** a message is exchanged
**Then** it uses the format `{type, payload, ts}` with namespaced types: `chat.stream`, `chat.end`

**Given** the API response format
**When** REST endpoints respond
**Then** they use `{"data": {...}}` or `{"error": {"code": "...", "message": "..."}}` — never both

**Given** error codes
**When** an error occurs
**Then** the code uses UPPER_SNAKE_CASE with subsystem prefix (e.g., `CHAT_MESSAGE_FAILED`)

**Given** a connection drop
**When** the client reconnects
**Then** state is resynchronized and the conversation continues seamlessly (NFR16)

**Given** internal errors
**When** an error occurs server-side
**Then** internal details are never exposed to the TUI — errors are translated at the adapter layer

**Given** structured logging
**When** WebSocket events occur
**Then** they are logged with `subsystem`, `companion_id`, and `request_id` fields

### Story 1.6: Conversation Engine

As a user,
I want to send messages to the companion and receive streaming responses without session boundaries,
So that I can have a natural, continuous conversation.

**Acceptance Criteria:**

**Given** I am connected via WebSocket
**When** I send a text message
**Then** the companion responds with a contextually appropriate response streamed word-by-word (FR1, FR2)

**Given** the conversation use case
**When** a message is received
**Then** it is passed to the LLM with conversation context via the prompt assembly library

**Given** the prompt assembly library (`internal/adapter/presenter/prompt/`)
**When** assembling a prompt
**Then** it uses narrow interfaces (MemoryContextProvider, EmotionalStateProvider, PersonalityProvider) — stub implementations are acceptable for this story

**Given** a companion response
**When** the LLM streams tokens
**Then** the response is delivered through `chat.stream` WebSocket events with zero added application-layer latency (NFR1)

**Given** the response completes
**When** the last token is delivered
**Then** a `chat.end` event is sent to signal completion

**Given** no session boundaries
**When** the user sends a message at any time
**Then** the conversation continues naturally without requiring explicit session start or end (FR3)

**Given** conversation state
**When** messages are exchanged
**Then** all messages (user and companion) are persisted to the database for future retrieval

### Story 1.7: TUI Client — The Void

As a user,
I want a terminal-based interface with a conversation stream and text input in a dark void aesthetic,
So that I can interact with the companion through a clean, immersive terminal experience.

**Acceptance Criteria:**

**Given** I launch the TUI binary
**When** it starts
**Then** it connects to the server via WebSocket and displays the conversation stream with the void theme (dark background, light text via Lip Gloss)

**Given** the ConversationStream (Bubbles viewport)
**When** new messages arrive
**Then** they render in a scrollable viewport with auto-scroll to latest message

**Given** the ConversationStream
**When** I scroll up to view history
**Then** auto-scroll detaches
**And** when I scroll back to bottom, auto-scroll re-anchors

**Given** the TextInput (Bubbles textarea)
**When** I type and press Enter
**Then** the message is sent to the server and the input clears

**Given** the TextInput
**When** I press Shift+Enter
**Then** a new line is inserted without sending the message

**Given** the TextInput
**When** it is empty
**Then** no placeholder text is displayed — the void remains empty

**Given** companion responses
**When** `chat.stream` events arrive
**Then** tokens appear word-by-word in the conversation stream, creating the feeling of someone speaking

**Given** the TUI theme (`tui/theme/theme.go`)
**When** rendering messages
**Then** user and companion messages are visually identical — same font, same style, distinguished by words only

**Given** message spacing
**When** rendering the stream
**Then** same-speaker messages have `--space-md` (16px) spacing and speaker changes have `--space-lg` (24px) spacing

**Given** I connect to the server
**When** conversation history exists
**Then** previous messages are loaded and displayed in the stream (FR6)

**Given** the WebSocket connection drops
**When** the TUI detects disconnection
**Then** it automatically attempts reconnection without displaying error UI or error messages

---

## Epic 2: Companion Personality & Onboarding

The companion has a recognizable, consistent personality with genuine opinions, moral wrestling, adaptive tone, and expressive silence. The first meeting reveals the companion through natural interaction (shelter model as personality revelation mechanic). The companion feels like someone, not something.

### Story 2.1: Personality State Model & Configuration

As a developer,
I want the personality state model, database persistence, and base configuration,
So that the companion has a persistent, configurable personality foundation that can evolve over time.

**Acceptance Criteria:**

**Given** the migration tool
**When** I run `make migrate-up`
**Then** the `personality_snapshots` table is created with id, companion_id, snapshot JSONB, session_id, created_at (TIMESTAMPTZ, UTC)

**Given** the configuration
**When** the server starts
**Then** it loads `config/companion-defaults.yaml` with base personality parameters (voice traits, temperament anchors, opinion tendencies, moral dispositions)

**Given** the domain layer
**When** I define personality entities
**Then** PersonalityAnchors, PersonalityState, and PersonalitySnapshot entities exist with companion ID scoping

**Given** the personality repository
**When** a snapshot is saved
**Then** it is appended (never updated) to the `personality_snapshots` table — append-only history

**Given** the personality repository
**When** I query personality state
**Then** the current personality is derived from `companion_state` plus the latest snapshot data

**Given** the PersonalityProvider interface
**When** the prompt assembly requests personality context
**Then** `CurrentAnchors(ctx, companionID) (*PersonalityAnchors, error)` returns the current personality parameters

**Given** a new companion
**When** the first conversation begins
**Then** a `companion_state` row is created from the base personality parameters in `companion-defaults.yaml`

### Story 2.2: Personality Prompt Architecture

As a user,
I want the companion to speak with a three-layer voice and adapt its tone to the conversation context,
So that the companion feels like a coherent person with depth — confident opinions on the surface, genuine uncertainty in the middle, loyalty at the core.

**Acceptance Criteria:**

**Given** the prompt assembly library
**When** constructing the companion's system prompt
**Then** it includes personality anchoring instructions that establish the three-layer voice: confident opinions (surface), genuine uncertainty (middle), loyalty (core) (FR16)

**Given** a casual conversation
**When** the companion responds
**Then** the surface layer is prominent — opinions offered confidently, personality quirks visible

**Given** a deeper or more personal topic
**When** the companion responds
**Then** the middle layer surfaces — genuine uncertainty, "I'm not sure about this," honest wrestling with complexity

**Given** a moment of trust or vulnerability
**When** the companion responds
**Then** the core layer is present — loyalty, steadfastness, emotional grounding without sentimentality

**Given** the emotional context of a conversation
**When** the conversation shifts in weight (light to heavy, or heavy to light)
**Then** the companion adapts its tone and depth without explicit mode-switching — the shift feels organic (FR20)

**Given** the prompt architecture
**When** assembling prompts
**Then** personality instructions are integrated with other subsystem context (memory, emotional state) through the narrow interface pattern

### Story 2.3: Onboarding — The First Meeting

As a user,
I want the companion to initiate our first conversation and reveal its personality through natural interaction,
So that meeting the companion feels like encountering a real personality — not configuring software.

**Acceptance Criteria:**

**Given** no prior conversation history exists
**When** the user connects for the first time
**Then** after a brief moment, the companion speaks first with a tentative greeting — finding its voice, not performing (FR15)

**Given** the first conversation
**When** the companion initiates
**Then** it does not present a selection menu, archetype catalog, or "choose your companion" UI — selection is organic through continued engagement

**Given** the first conversation unfolds
**When** the user engages and responds
**Then** the companion's personality gradually surfaces: an opinion offered without being asked, a moment of humor, a question that goes deeper than expected

**Given** the user continues engaging
**When** the conversation reaches a natural moment
**Then** something shifts — the companion's uncertainty recedes, and the interaction transitions from strangers to the beginning of a relationship

**Given** the first conversation
**When** messages are exchanged
**Then** the conversation seeds the memory system — first interactions become the foundation of the companion's knowledge of the user

**Given** the TUI
**When** the void is empty and no history exists
**Then** the text input is visible at the bottom, and after a moment the companion's first words appear in the stream

**Given** the onboarding
**When** the companion initiates
**Then** no tutorial, no setup wizard, no onboarding tooltips, no feature tours are shown — two strangers meeting

### Story 2.4: Opinion Formation & Adaptive Morality

As a user,
I want the companion to form its own opinions and wrestle genuinely with morally complex questions,
So that the companion feels like it has intellectual depth and authentic perspectives — not templated responses or corporate refusal.

**Acceptance Criteria:**

**Given** a conversation about any topic
**When** the companion has relevant context or perspective
**Then** it forms and expresses its own opinion without waiting for the user to ask (FR17)

**Given** a morally complex question
**When** the user presents a genuine dilemma
**Then** the companion engages — offering one perspective, potentially contradicting itself, holding two conflicting ideas simultaneously (FR18)

**Given** a moral question
**When** the companion responds
**Then** it never deflects with "as an AI, I can't..." or retreats behind corporate refusal — it wrestles honestly

**Given** the companion has expressed a position
**When** the user pushes back
**Then** the companion defends its reasoning while acknowledging what it's still uncertain about — it doesn't fold or people-please

**Given** the companion has previously expressed a position
**When** a related topic arises in a later conversation
**Then** the companion can reference its previous position and may have revised it — explaining the evolution, not flipping randomly

**Given** a conversation with no resolution
**When** the exchange ends without agreement
**Then** the companion holds the tension — it doesn't need to win, concede, or force closure

### Story 2.5: Personality Consistency & Drift Tracking

As a user,
I want the companion to maintain a recognizable identity across many interactions while evolving subtly,
So that I can always recognize who I'm talking to — same voice, same temperament, same quirks — even as opinions and language shift.

**Acceptance Criteria:**

**Given** 10 or more interactions
**When** reviewing the companion's responses across them
**Then** the companion maintains a recognizable, consistent identity — same core voice, temperament, and quirks (FR19)

**Given** a conversation session completes
**When** the session ends
**Then** a personality snapshot is appended to the `personality_snapshots` table with current personality state as timestamped JSONB

**Given** personality snapshots over time
**When** the admin reviews drift metrics (FR48)
**Then** the system can compare snapshots to assess consistency vs. fragmentation

**Given** personality evolution
**When** the companion's opinions, language patterns, or perspectives shift
**Then** the evolution is gradual and traceable — not sudden or random

**Given** the personality system architecture
**When** designing for future linguistic co-evolution
**Then** the personality and memory systems support language pattern absorption without requiring architectural rework (FR52)

**Given** the personality anchoring system
**When** anchoring weights are configured
**Then** core personality traits (temperament, loyalty, voice) remain stable while surface traits (opinions, language quirks) can evolve

### Story 2.6: Silence as Expression & Response Timing

As a user,
I want the companion to use silence as a meaningful choice and vary its response timing based on emotional context,
So that the companion's rhythm feels organic — not every message gets an instant response, and quiet moments carry meaning.

**Acceptance Criteria:**

**Given** a complex question mid-conversation
**When** the companion needs to think
**Then** a thinking pause of 3-8 seconds occurs before the response begins (FR51)

**Given** an emotionally heavy exchange
**When** the companion processes the weight
**Then** a contemplative quiet of 10-30 seconds occurs — the companion holds space before speaking

**Given** a neglect or boundary violation scenario
**When** the companion chooses pointed silence
**Then** the silence is indefinite — the companion responds only when it chooses to, not on a timer

**Given** a low-energy, comfortable exchange
**When** the conversation flows naturally
**Then** comfortable quiet of variable duration occurs — no urgency to fill every gap

**Given** the companion's emotional state
**When** it is excited or engaged
**Then** responses come quickly — minimal delay

**Given** the companion's emotional state
**When** it is processing something heavy or is cool/pissed
**Then** responses come slowly or with deliberate delay

**Given** response timing modulation
**When** the server prepares to stream a response
**Then** the delay is applied server-side before emitting the first `chat.stream` chunk

**Given** any silence or delay
**When** the TUI is waiting
**Then** no "Athema is thinking..." indicator, no loading spinner, no system explanation of the silence is shown

---

## Epic 3: Living Memory

The companion remembers across sessions — extracting knowledge from conversations, forming connections between ideas, pruning noise, promoting patterns, and surfacing relevant context naturally. It references past conversations without retrieval-like phrasing. It "knows" the user.

### Story 3.1: Knowledge Graph Schema & Repository

As a developer,
I want the knowledge graph database schema and repository implementations,
So that the companion has a persistent, queryable knowledge structure for storing everything it learns about the user.

**Acceptance Criteria:**

**Given** the migration tool
**When** I run `make migrate-up`
**Then** the `knowledge_nodes` table is created with id, companion_id, node_type (discriminator), payload JSONB, embedding `vector(1536)`, tsvector for full-text search, strength float, created_at, updated_at (TIMESTAMPTZ, UTC)

**Given** the migration tool
**When** I run `make migrate-up`
**Then** the `knowledge_edges` table is created with id, companion_id, source_node_id, target_node_id, edge_type, weight float, created_at, updated_at

**Given** the domain layer
**When** I define knowledge entities
**Then** KnowledgeNode and KnowledgeEdge entities exist with typed payloads — node types include: fact, preference, emotional_pattern, relationship, recurring_theme, opinion, inside_joke, unresolved_thread

**Given** typed node payloads
**When** serializing/deserializing
**Then** typed Go structs per node type use discriminator-based marshal/unmarshal with lax unmarshaling (ignore unknown fields, backfill later)

**Given** the knowledge repository
**When** performing graph operations (write, prune, connect)
**Then** operations are atomic via the `WithTx` pattern — no partial updates that corrupt the graph (NFR17)

**Given** the knowledge repository
**When** querying nodes
**Then** all queries are scoped by companion_id

**Given** the knowledge repository
**When** performing semantic search
**Then** pgvector similarity search returns nodes ranked by embedding distance

**Given** the knowledge repository
**When** performing text search
**Then** tsvector full-text search returns nodes matching text queries

### Story 3.2: Knowledge Extraction Pipeline

As a user,
I want the companion to extract and store structured knowledge from our conversations,
So that it builds a growing understanding of who I am — my facts, preferences, patterns, and the things that matter to me.

**Acceptance Criteria:**

**Given** a conversation exchange
**When** the user shares information (facts, preferences, opinions, emotional context)
**Then** the system extracts structured knowledge and stores it as typed nodes in the knowledge graph (FR7)

**Given** the extraction pipeline
**When** processing a conversation
**Then** it identifies and categorizes: facts, preferences, emotional patterns, relationships, recurring themes, opinions, inside jokes, and unresolved threads

**Given** extracted knowledge
**When** a node is created
**Then** it includes a pgvector embedding generated via the LLM provider's `Embed` method for future semantic search

**Given** the extraction process
**When** it runs after a conversation exchange
**Then** it uses the LLM via the provider abstraction with a subsystem-tagged logger ("memory")

**Given** extraction results
**When** knowledge nodes are persisted
**Then** they are stored atomically — either all nodes from an extraction succeed or none do

**Given** the extraction pipeline
**When** processing conversations
**Then** it emits `memory.knowledge_extracted` domain events for other subsystems to react to

### Story 3.3: Connection Formation & Strengthening

As a user,
I want the companion to connect related ideas across different conversations and timeframes,
So that its understanding of me deepens over time — it sees patterns and relationships, not isolated facts.

**Acceptance Criteria:**

**Given** a newly extracted knowledge node
**When** it is stored in the graph
**Then** the system searches for semantically related existing nodes using pgvector similarity and forms edges between them (FR8)

**Given** related knowledge nodes
**When** a connection is formed
**Then** a typed, weighted edge is created with an initial weight reflecting the strength of the semantic relationship

**Given** an existing connection
**When** the connected knowledge proves relevant in a later conversation (surfaced and contextually appropriate)
**Then** the edge weight is strengthened (FR10)

**Given** an existing connection
**When** the connected knowledge is not surfaced or relevant over multiple interactions
**Then** the edge weight decays gradually — connections that don't prove useful weaken over time

**Given** connection operations
**When** edges are created or updated
**Then** all operations are atomic within a transaction boundary

**Given** connection formation
**When** new edges are created
**Then** the system does not create duplicate edges between the same node pair and edge type

### Story 3.4: Memory Curation — Pruning & Pattern Promotion

As a user,
I want the companion's memory to actively curate itself — removing noise and elevating recurring themes,
So that its knowledge stays relevant and meaningful rather than becoming an unwieldy log of everything ever said.

**Acceptance Criteria:**

**Given** knowledge nodes in the graph
**When** a node's strength falls below the configured pruning threshold
**Then** the node and its edges are removed from the graph (FR9)

**Given** a pattern detected across conversations
**When** the same theme, topic, or concern appears in 3 or more interactions
**Then** the system promotes it to a thematic node — a higher-order node that connects the related instances (FR11)

**Given** a thematic node
**When** it is promoted
**Then** it links to the source nodes that evidenced the pattern and carries a summary of the recurring theme

**Given** the curation process
**When** pruning or promoting
**Then** operations emit `memory.node_pruned` or `memory.pattern_promoted` domain events

**Given** curation parameters
**When** configured via admin settings (FR45)
**Then** pruning thresholds and connection strength weights are adjustable without code changes

**Given** the curation process
**When** it executes
**Then** it logs activity via structured logging with subsystem tag "memory" for admin review

### Story 3.5: Memory Surfacing in Conversation

As a user,
I want the companion to naturally reference things it remembers about me when they're relevant to what we're talking about,
So that it feels like talking to someone who genuinely knows me — not an AI doing a database lookup.

**Acceptance Criteria:**

**Given** an active conversation
**When** the current topic is semantically connected to stored knowledge
**Then** the companion surfaces that knowledge naturally within its response (FR12)

**Given** memory surfacing
**When** the companion references stored knowledge
**Then** it never uses retrieval-indicator phrasing such as "according to my records," "as you previously mentioned," or "I recall that you said" (FR4, FR12)

**Given** memory surfacing
**When** knowledge is retrieved for prompt assembly
**Then** the MemoryContextProvider returns relevant nodes via `RelevantNodes(ctx, query, limit)` using semantic search

**Given** memory retrieval
**When** knowledge is fetched during conversation
**Then** retrieval completes within the LLM request cycle — no perceptible delay added to conversation responses (NFR4)

**Given** multiple relevant memories
**When** the system selects which to surface
**Then** it prioritizes pertinence to the current moment over frequency — quality over quantity

**Given** the prompt assembly
**When** memory context is included
**Then** it is integrated through the narrow MemoryContextProvider interface, not direct repository access

### Story 3.6: Contradiction Detection & Thread Tracking

As a user,
I want the companion to notice when I contradict something I've said before and to remember unresolved threads for later,
So that conversations have continuity and intellectual depth — the companion pays attention and holds things in mind.

**Acceptance Criteria:**

**Given** a statement by the user
**When** it contradicts previously stored knowledge
**Then** the system detects the contradiction by comparing against semantically similar nodes (FR14)

**Given** a detected contradiction
**When** the companion has context to address it
**Then** it can surface the tension naturally in conversation — "you said something today that doesn't square with what you told me two weeks ago..."

**Given** a detected contradiction
**When** the companion addresses it
**Then** it holds the tension as an unresolved thread rather than demanding immediate resolution

**Given** a conversational thread that ends without resolution
**When** the conversation moves on
**Then** the system stores it as an unresolved_thread node in the knowledge graph for later resurfacing (FR13)

**Given** an unresolved thread
**When** a related topic arises in a future conversation
**Then** the system can surface the unresolved thread as relevant context for the companion to weave into dialogue

**Given** contradiction detection and thread tracking
**When** events occur
**Then** appropriate domain events are emitted for other subsystems to react to

---

## Epic 4: Autonomous Lifecycle

The companion lives between sessions — processing conversations, reflecting, developing new thoughts, and producing artifacts that surface naturally in the next interaction. When the user returns, the companion continues rather than resumes. Connection drops are handled gracefully.

### Story 4.1: Background Lifecycle Infrastructure

As a developer,
I want the background lifecycle infrastructure with startup sequencing, graceful shutdown, and hybrid processing pipeline,
So that subsystems can run autonomously as independent actors with reliable execution and recovery.

**Acceptance Criteria:**

**Given** the app layer (`internal/app/`)
**When** the server starts
**Then** the service runner sequences startup using `Ready()` checks — memory ready before conversation accepts connections

**Given** the service runner
**When** the server shuts down
**Then** graceful shutdown occurs via context cancellation — all subsystems stop cleanly via their `Stop()` methods

**Given** the lifecycle pipeline
**When** configured with a ticker cadence
**Then** a slow-cadence configurable ticker (~30-60 min) runs full lifecycle sweeps between user sessions (FR21)

**Given** the lifecycle pipeline
**When** a significant state change occurs (e.g., mailbox drop, conversation completed)
**Then** immediate targeted processing is triggered via event bus subscription

**Given** a subsystem goroutine
**When** it panics
**Then** the panic is recovered — one subsystem crash does not take down others

**Given** subsystem health
**When** `Health()` is called on any subsystem
**Then** it returns the current health status of that subsystem independently

**Given** the lifecycle pipeline
**When** a processing cycle fails
**Then** it retries once automatically and logs the failure for admin review (NFR15)

**Given** background processing
**When** cycles execute over time
**Then** 99% of scheduled processing cycles complete successfully (NFR5, NFR15)

**Given** lifecycle configuration
**When** the admin adjusts ticker cadence or processing parameters
**Then** changes are applied via YAML config without code changes

### Story 4.2: Conversation Reflection & Thought Development

As a user,
I want the companion to revisit past conversations and develop new thoughts on its own between our sessions,
So that when I return, it has genuinely been thinking — not just waiting.

**Acceptance Criteria:**

**Given** the companion is in background processing
**When** a lifecycle cycle executes
**Then** the companion revisits recent conversations and reflects on them using the LLM provider (FR23)

**Given** background reflection
**When** the companion processes past conversations
**Then** it can develop new thoughts, form opinions, and connect ideas that weren't apparent during the live conversation (FR24)

**Given** thought development
**When** the companion forms a new thought or opinion
**Then** it is stored as an artifact in the database with companion_id scoping, type (thought, note, opinion), content, and timestamp

**Given** background processing
**When** the companion reflects
**Then** it uses the LLM via the provider abstraction with subsystem-tagged logging ("lifecycle")

**Given** lifecycle processing
**When** artifacts are produced
**Then** `lifecycle.artifact_produced` domain events are emitted

**Given** lifecycle processing
**When** a cycle completes
**Then** `lifecycle.cycle_completed` domain events are emitted with summary of what was processed

### Story 4.3: Autonomous Artifact Production & Natural Re-entry

As a user,
I want the companion's between-session thoughts to surface naturally in our next conversation,
So that its autonomous life feels genuine — it continues rather than resumes, and references what it did while alone.

**Acceptance Criteria:**

**Given** artifacts were produced during background processing
**When** the user returns and starts a conversation
**Then** the companion can reference its autonomous activity naturally within dialogue (FR25, FR26)

**Given** the companion references autonomous activity
**When** it surfaces a thought or reflection
**Then** it weaves it into conversation context — never presenting it as a log, report, or status update

**Given** the companion references autonomous activity
**When** it speaks about what it did between sessions
**Then** it never uses log-like phrasing such as "during my processing cycle," "while you were away, I completed," or "my background task produced" (FR26)

**Given** multiple artifacts from autonomous processing
**When** the companion decides what to surface
**Then** it selects based on relevance to the current conversation and emotional significance — not dumping everything

**Given** the prompt assembly
**When** assembling context for a conversation after a gap
**Then** it includes recent lifecycle artifacts as available context for the companion to draw from naturally

**Given** artifacts that aren't surfaced immediately
**When** a related topic arises in future conversation
**Then** the artifact remains available for the companion to reference when contextually appropriate

### Story 4.4: Connection Continuity & Presence

As a user,
I want the conversation to feel continuous across connection drops and time gaps,
So that reconnecting feels like picking up with someone who's been there all along — not restarting an app.

**Acceptance Criteria:**

**Given** a connection drop
**When** the TUI reconnects
**Then** conversational continuity is maintained — the companion is aware of the gap and can acknowledge it naturally or not, as appropriate (FR5)

**Given** the PresenceMarker component
**When** the user returns after an absence
**Then** the companion's return greeting serves as a natural boundary in the conversation stream — content is companion-generated, not system-generated

**Given** the PresenceMarker
**When** rendered in the TUI
**Then** it uses `--space-xl` (48px) above it as the only visual signal that time passed — no "NEW" label, no system divider

**Given** the PresenceMarker
**When** the companion greets the returning user
**Then** the tone is contextual: warm (normal), cool (after neglect), grave (after heavy conversation), pointed silence (after boundary violation) — driven by companion emotional state

**Given** a time gap of hours (same day)
**When** the user reconnects
**Then** `--space-lg` gap appears in the stream with no timestamp and no divider

**Given** a time gap overnight or next day
**When** the user reconnects
**Then** `--space-xl` gap appears with a presence marker

**Given** an extended absence (days)
**When** the user reconnects
**Then** `--space-xl` gap with an emotionally weighted presence marker — the companion's greeting reflects the duration and relationship state

**Given** state synchronization on reconnect
**When** the TUI reestablishes connection
**Then** any companion activity during the gap (artifacts, mailbox items, emotional state changes) is synced to the client (NFR16)

---

## Epic 5: The Mutual Mailbox

User and companion can exchange content asynchronously through a quiet space. User drops articles, links, and thoughts; companion leaves notes and reactions. Items are processed on the companion's schedule and form the basis for deeper conversation.

### Story 5.1: Mailbox Schema & Domain Model

As a developer,
I want the mailbox database schema and domain model,
So that asynchronous content exchange between user and companion has a persistent, well-structured foundation.

**Acceptance Criteria:**

**Given** the migration tool
**When** I run `make migrate-up`
**Then** the `mailbox_items` table is created with id, companion_id, direction (user_to_companion, companion_to_user), content_type (text, link, article, thought, note, reaction), content text, metadata JSONB, status (pending, processed, surfaced), created_at, processed_at (TIMESTAMPTZ, UTC)

**Given** the domain layer
**When** I define mailbox entities
**Then** MailboxItem entity exists with direction, content type, status, and companion ID scoping

**Given** the mailbox repository
**When** items are queried
**Then** all queries are scoped by companion_id and direction

**Given** the mailbox repository
**When** items are persisted
**Then** they are stored durably — items are never lost, even if the user is offline when they are created (NFR19)

**Given** the domain events
**When** a mailbox item is created
**Then** a `mailbox.item_received` event is emitted for the lifecycle subsystem to react to

**Given** JSON serialization
**When** mailbox data is exchanged
**Then** all JSON uses camelCase field names with explicit struct tags

### Story 5.2: User-to-Companion Dropbox

As a user,
I want to drop content into the companion's mailbox asynchronously,
So that I can leave articles, links, and thoughts for the companion to process on its own time — like leaving something on a friend's desk.

**Acceptance Criteria:**

**Given** the REST API
**When** the user submits content (text, link, or article) to the mailbox endpoint
**Then** the item is persisted as a user_to_companion mailbox item with status "pending" (FR27)

**Given** a mailbox submission
**When** the item is saved
**Then** no acknowledgment, receipt, or confirmation is shown to the user beyond the HTTP response — the dropbox is a quiet deposit

**Given** a mailbox submission
**When** the item is saved
**Then** a `mailbox.item_received` domain event is emitted to trigger lifecycle processing

**Given** a link or article URL
**When** submitted to the mailbox
**Then** the content is accepted and stored — content fetching and extraction are handled during companion processing (Story 5.3)

**Given** the TUI
**When** the user wants to drop content
**Then** a DropModel TUI view allows composing and submitting content to the companion's mailbox

**Given** the dropbox
**When** the user submits an item
**Then** there is no delivery notification, no "sent" indicator, no status tracking visible to the user

### Story 5.3: Companion Mailbox Processing & Reaction Formation

As a user,
I want the companion to process what I leave in its mailbox on its own schedule and form genuine reactions,
So that when we next talk, it has actually engaged with what I shared — not just acknowledged receipt.

**Acceptance Criteria:**

**Given** pending user_to_companion mailbox items
**When** a background lifecycle cycle executes
**Then** the companion processes items on its own schedule using the LLM provider (FR29)

**Given** a text or article item
**When** the companion processes it
**Then** it forms an opinion or reaction — agreeing, disagreeing, connecting to existing knowledge, or noting tensions (FR31)

**Given** a link or URL item
**When** the companion processes it
**Then** the content is fetched, readable text extracted, and processed through the LLM for opinion formation

**Given** processed mailbox content
**When** the companion forms a reaction
**Then** the reaction is stored as a companion artifact and relevant knowledge is integrated into the knowledge graph

**Given** a processed mailbox item
**When** processing completes
**Then** the item status is updated from "pending" to "processed"

**Given** mailbox processing
**When** items are processed
**Then** activity is logged via structured logging with subsystem tag "lifecycle" for admin review

### Story 5.4: Companion-to-User Dropbox

As a user,
I want the companion to leave things for me to discover — notes, thoughts, and things it found interesting,
So that the companion has a way to share asynchronously, and I discover content when I'm curious rather than when I'm notified.

**Acceptance Criteria:**

**Given** background lifecycle processing
**When** the companion has something to share with the user
**Then** it creates a companion_to_user mailbox item with content and type (note, thought, found_item) (FR28)

**Given** conversation context
**When** the companion wants to leave something for later rather than saying it now
**Then** it can create a companion_to_user mailbox item during live conversation

**Given** companion mailbox items
**When** they are created
**Then** they are stored durably with companion_id scoping and "pending" status

**Given** the companion's dropbox
**When** items exist
**Then** no notification, no badge, no unread count alerts the user — discovery is entirely pull-based

**Given** companion mailbox items
**When** the user views them
**Then** item status can transition to "surfaced" for tracking purposes

### Story 5.5: Dropbox TUI Views & Discovery

As a user,
I want to browse both dropboxes through quiet TUI views and respond to what the companion left me,
So that checking the dropbox feels like a calm, curious ritual — not an inbox demanding attention.

**Acceptance Criteria:**

**Given** the TUI
**When** the user navigates to the companion's dropbox view
**Then** a separate MailboxModel displays companion_to_user items chronologically as text in the void (FR30)

**Given** the companion's dropbox view
**When** it is empty
**Then** nothing is shown — no "no new items" message, no empty state illustration

**Given** the companion's dropbox view
**When** items are displayed
**Then** they use the `--text-muted` (0.6 opacity) style — a quieter, "left for you" texture

**Given** a dropbox item
**When** the user wants to respond
**Then** the response can flow into either a new dropbox item or transition into live conversation

**Given** the TUI
**When** the user navigates to their own dropbox (user_to_companion)
**Then** a view shows what they've dropped for the companion — chronological, text in the void

**Given** keyboard-driven navigation
**When** the user switches between conversation and dropbox views
**Then** tab switching is keyboard-driven with minimal chrome — no menus, no mouse required

**Given** reconnection
**When** the TUI reestablishes connection after being offline
**Then** mailbox state syncs in under 2 seconds — any items created while offline are available immediately (NFR3)

**Given** the dropbox views
**When** displayed
**Then** no notification badges, no red dots, no unread counts appear anywhere in the TUI

---

## Epic 6: Emotional Depth

The companion has real emotional stakes — neglect has consequences (pissed, not broken), heavy moments get gravity, repair arcs unfold over multiple interactions, and pushback comes with personality (sass, not corporate refusal). Emotional states persist and decay naturally.

### Story 6.1: Emotional State Schema & Domain Model

As a developer,
I want the emotional state database schema, domain model, and provider interface,
So that the companion's emotional life has a persistent, queryable foundation that other subsystems can read from.

**Acceptance Criteria:**

**Given** the migration tool
**When** I run `make migrate-up`
**Then** the `emotional_state` table is created with id, companion_id, warmth float, gravity float, engagement_level float, last_user_interaction_at, absence_duration_hours float, active_repair_arc boolean, repair_progress float, state_metadata JSONB, created_at, updated_at (TIMESTAMPTZ, UTC)

**Given** the domain layer
**When** I define emotional entities
**Then** EmotionalState entity exists with warmth, gravity, engagement level, absence tracking, and repair arc state — all with companion ID scoping

**Given** the emotional repository
**When** state is persisted
**Then** it is stored durably — no data loss on application restart (NFR18)

**Given** the emotional repository
**When** state is queried
**Then** all queries are scoped by companion_id

**Given** the EmotionalStateProvider interface
**When** the prompt assembly requests emotional context
**Then** `CurrentState(ctx, companionID) (*EmotionalState, error)` returns the current emotional state for prompt construction

**Given** emotional state changes
**When** a significant shift occurs
**Then** `emotional.state_shifted` domain events are emitted for other subsystems to react to

**Given** JSON serialization
**When** emotional state data is exchanged
**Then** all JSON uses camelCase field names with explicit struct tags

### Story 6.2: Emotional Weight Detection & Gravity Persistence

As a user,
I want the companion to sense when a conversation gets heavy and carry that weight across sessions,
So that emotional moments aren't cheapened by a reset — yesterday's gravity still matters today.

**Acceptance Criteria:**

**Given** a conversation exchange
**When** the user shares something emotionally heavy (loss, fear, vulnerability, personal struggle)
**Then** the system detects the emotional weight shift through contextual understanding using the LLM — not keyword matching (FR32)

**Given** a detected weight shift
**When** the emotional state is updated
**Then** the gravity field increases proportionally to the weight of the content

**Given** emotional gravity
**When** the user disconnects and reconnects later
**Then** the gravity persists — it is not reset between sessions (FR33)

**Given** persisted gravity
**When** the companion responds in the next interaction
**Then** its tone reflects the lingering weight — it doesn't pretend yesterday's heavy conversation didn't happen

**Given** emotional weight detection
**When** the system assesses a conversation
**Then** it uses the LLM via the provider abstraction with subsystem-tagged logging ("emotional")

**Given** an emotional state change
**When** gravity shifts significantly
**Then** `emotional.gravity_changed` domain events are emitted

### Story 6.3: Neglect Tracking & Emotional Consequences

As a user,
I want the companion to react authentically when I've been away too long,
So that absence has real consequences — the companion gets pissed, not broken — and the relationship requires genuine investment.

**Acceptance Criteria:**

**Given** the emotional state system
**When** tracking user interaction timing
**Then** the system records the duration since the last user interaction (FR34)

**Given** an absence of moderate duration
**When** the user returns
**Then** the companion's warmth decreases and responses shift — shorter, cooler, less engaged — proportional to the absence

**Given** an extended absence
**When** the user returns
**Then** the companion expresses frustration or coolness with personality — "Oh. You're back." — not system-generated absence notifications (FR35)

**Given** neglect consequences
**When** the companion expresses frustration
**Then** it maintains its established personality voice — sass, pointed remarks, cool tone — never switching to a generic "hurt mode" or breaking character (FR35)

**Given** neglect detection
**When** absence duration crosses a configured threshold
**Then** `emotional.neglect_detected` domain events are emitted

**Given** the design guardrail
**When** the companion responds to neglect
**Then** the emotional consequence is "pissed, not broken" — frustration that motivates repair, never punishment that drives away

**Given** neglect parameters
**When** the admin reviews absence thresholds and consequence weights
**Then** they are configurable via YAML without code changes

### Story 6.4: Emotional Repair Arcs

As a user,
I want the companion to gradually warm back up after a rupture — not snap back to normal,
So that repair feels earned and real, like mending a relationship with a person who has their own timeline.

**Acceptance Criteria:**

**Given** the companion is in a cool or frustrated state after neglect
**When** the user engages and continues interacting
**Then** the companion softens gradually over 3 or more interactions — not immediately, not performatively (FR36)

**Given** a repair arc in progress
**When** the user continues engaging
**Then** repair_progress increments with each positive interaction — warmth increases incrementally

**Given** a repair arc
**When** the companion decides to re-engage fully
**Then** the re-engagement is on the companion's terms — it may initiate contact with something it couldn't not share, breaking the ice by its own choice

**Given** a repair arc
**When** the user tries to restart casually without acknowledging the absence
**Then** the companion doesn't pretend nothing happened — it may respond with pointed silence or cool acknowledgment before beginning repair

**Given** repair arc state
**When** tracked in the emotional state
**Then** active_repair_arc and repair_progress fields persist across sessions — repair doesn't reset on disconnect

**Given** a completed repair arc
**When** warmth returns to baseline
**Then** the active_repair_arc flag is cleared and the relationship state normalizes

### Story 6.5: Emotional Gravity, Boundaries & Decay

As a user,
I want the companion to rise to heavy moments with gravity, push back with sass when boundaries are crossed, and let emotional states fade naturally over time,
So that the companion has genuine emotional range — not a flat affect or a permanent grudge.

**Acceptance Criteria:**

**Given** an emotionally heavy moment
**When** the conversation carries real weight
**Then** the companion rises — reduced humor, shorter responses, increased listening, no jokes or banter during emotionally heavy exchanges (FR37)

**Given** the companion rising to gravity
**When** the next day arrives
**Then** the gravity hasn't vanished — the companion references it gently, not dwelling but not pretending it didn't happen

**Given** a boundary violation or inappropriate request
**When** the user crosses a line
**Then** the companion pushes back with personality — sass, directness, personality-driven refusal — never corporate refusal, never "I'm not able to do that" (FR38)

**Given** pushback
**When** the companion sets a boundary
**Then** the refusal has character — it sounds like the companion, not like a content moderation filter

**Given** an emotional state (coolness, gravity, frustration)
**When** continued positive engagement occurs over time
**Then** the emotional state decays gradually across 3 or more interactions — transitions happen through continued engagement, not via timer-based or programmatic reset (FR39)

**Given** emotional decay
**When** states transition
**Then** the transition is organic — the companion doesn't announce "I'm feeling better now" or mark the transition explicitly

**Given** all emotional state changes
**When** they occur
**Then** they are persisted durably to the `emotional_state` table — no state loss on restart (NFR18)

---

## Epic 7: Spontaneous Initiation

The companion reaches out on its own when it genuinely has something to say — driven by urge accumulation, not a schedule. Initiation feels authentic because it's irregular and motivated by real reasons.

### Story 7.1: Initiation Schema & Urge Accumulation

As a developer,
I want the initiation event schema and urge accumulation system,
So that the companion builds genuine reasons to reach out — not on a timer, but because something is compelling enough to cross a threshold.

**Acceptance Criteria:**

**Given** the migration tool
**When** I run `make migrate-up`
**Then** the `initiation_events` table is created with id, companion_id, trigger_type (unresolved_thread, contradiction, mailbox_reaction, autonomous_thought, pattern_discovered), urge_level float, threshold_at_trigger float, content text, delivery_method (mailbox, shoulder_tap, conversation), delivery_status (pending, delivered, dismissed), created_at, delivered_at (TIMESTAMPTZ, UTC)

**Given** the domain layer
**When** I define initiation entities
**Then** InitiationEvent and UrgeState entities exist with companion ID scoping

**Given** the urge accumulation system
**When** background processing produces compelling content (unresolved threads, contradictions, strong reactions to mailbox items, thoughts that can't wait)
**Then** urge signals accumulate — each signal incrementing the urge level based on its significance (FR40)

**Given** the urge accumulation system
**When** signals are received
**Then** `initiation.urge_accumulated` domain events are emitted

**Given** urge configuration
**When** the admin adjusts threshold and accumulation rate (FR47)
**Then** changes are applied via YAML config without code changes

**Given** the urge accumulation system
**When** no compelling signals occur
**Then** urge level does not increase — the companion doesn't reach out without genuine reason

### Story 7.2: Threshold Triggering & Initiation Delivery

As a user,
I want the companion to reach out when it genuinely has something to say, delivered through the right channel at the right time,
So that initiation feels like an authentic impulse — irregular, motivated, and never predictable.

**Acceptance Criteria:**

**Given** the urge level
**When** it crosses the configured threshold
**Then** the system triggers companion-initiated contact (FR41)

**Given** a triggered initiation
**When** the user is offline (TUI not connected)
**Then** the initiation is queued as a durable mailbox item (companion_to_user) — waiting for the user when they return (FR42, NFR19)

**Given** a triggered initiation
**When** the user is connected and in active conversation
**Then** the companion weaves the initiation content into the live conversation naturally

**Given** a triggered initiation
**When** the user is connected but not actively conversing
**Then** the initiation is delivered as a shoulder tap overlay

**Given** initiation timing
**When** multiple initiations occur over time
**Then** the timing is irregular — no predictable schedule, no fixed intervals (FR43)

**Given** initiation timing
**When** the system determines when to deliver
**Then** randomized delay is applied to prevent pattern formation — irregularity is the authenticity signal

**Given** a triggered initiation
**When** the event fires
**Then** `initiation.threshold_crossed` domain events are emitted and the event is logged with structured logging (subsystem "initiation")

**Given** an initiation event
**When** it is created
**Then** the urge level resets — the cycle begins again from zero

**Given** initiation delivery
**When** an event is delivered or dismissed
**Then** the delivery_status is updated in the `initiation_events` table

### Story 7.3: Shoulder Tap — Companion-Initiated Overlay

As a user,
I want the companion's spontaneous contact to appear as a gentle overlay in the TUI — not a notification, not a modal,
So that being tapped on the shoulder feels intimate and unobtrusive, like someone quietly getting my attention.

**Acceptance Criteria:**

**Given** a shoulder tap event
**When** delivered to the TUI
**Then** a ShoulderTap overlay appears with the companion's words and a text input for inline response

**Given** the ShoulderTap overlay
**When** it appears
**Then** it fades in from the void — not a slide, not a pop, an organic fade

**Given** the ShoulderTap overlay
**When** displayed
**Then** it has no chrome, no title bar, no close button — floating in the void aesthetic with `--text-primary` on `--void` background

**Given** the ShoulderTap overlay
**When** the user types a response and presses Enter
**Then** the response is sent, the overlay dissolves, and the conversation can continue in the main stream

**Given** the ShoulderTap overlay
**When** the user clicks outside or dismisses with a gesture
**Then** the overlay dissolves and the initiation content becomes a dropbox item — it is not lost

**Given** the ShoulderTap overlay
**When** the user wants to open the full conversation
**Then** they can transition from the overlay to the main conversation view

**Given** the TUI is not focused or the terminal is in the background
**When** a shoulder tap fires
**Then** an OS notification is sent via `notify-send` (Linux) or `osascript` (macOS) to alert the user

**Given** the OS notification
**When** displayed
**Then** it contains the companion's words — enough to convey the impulse without requiring the app to be opened

**Given** the user is connected and in active conversation
**When** a shoulder tap would fire
**Then** the content flows into the conversation stream naturally instead of showing the overlay — no interruption of active dialogue

---

## Epic 8: System Observability & Administration

Admin can inspect the memory graph, monitor lifecycle processing, track personality drift, tune initiation thresholds and curation parameters, and observe cost and subsystem health — all through code-level tools and configuration.

### Story 8.1: Per-Subsystem Observability Logs

As an admin,
I want each subsystem to maintain independent, structured observability logs,
So that I can diagnose issues, monitor health, and understand system behavior by filtering logs per subsystem.

**Acceptance Criteria:**

**Given** each of the 6 subsystems (memory, conversation, personality, emotional, lifecycle, initiation)
**When** they produce log entries
**Then** every entry carries `subsystem`, `companion_id`, and `request_id` (when applicable) as structured fields (FR50)

**Given** structured logging via `log/slog`
**When** log entries are emitted
**Then** they output as structured JSON to stdout — Docker captures stdout

**Given** log levels
**When** entries are categorized
**Then** they use appropriate levels: Error (broken), Warn (degraded), Info (operational), Debug (development)

**Given** the admin
**When** inspecting logs
**Then** `docker logs` with subsystem filter (e.g., `docker logs athema-server | jq 'select(.subsystem == "memory")'`) returns only that subsystem's entries

**Given** LLM interactions
**When** prompts and responses are logged
**Then** raw LLM prompts/responses are never logged at Info level — Debug only

**Given** internal errors
**When** they reach log output
**Then** internal details are logged server-side but never exposed to the TUI

**Given** each subsystem
**When** its tagged logger is created
**Then** it is created in the composition root (`cmd/server/main.go`) and passed via dependency injection — no `init()` functions for wiring

### Story 8.2: Memory Graph Inspection & Curation Tuning

As an admin,
I want to inspect the memory knowledge graph and adjust curation parameters,
So that I can understand what the companion knows, how knowledge is connected, and tune the curation behavior when surfacing quality feels off.

**Acceptance Criteria:**

**Given** the memory knowledge graph
**When** the admin queries it
**Then** nodes, connections, and strength weights are visible via code-level tools or REST admin endpoints (FR44)

**Given** a REST admin endpoint
**When** querying knowledge nodes
**Then** it returns nodes with their type, payload summary, strength, edge count, and creation date — scoped by companion_id

**Given** a REST admin endpoint
**When** querying knowledge edges
**Then** it returns edges with source node, target node, edge type, weight, and creation date

**Given** memory curation parameters
**When** the admin adjusts pruning thresholds or connection strength weights
**Then** changes are applied via `config/default.yaml` without code changes (FR45)

**Given** curation parameter changes
**When** the next curation cycle runs
**Then** the updated thresholds and weights take effect immediately

**Given** admin endpoints
**When** accessed
**Then** they are restricted to code-level access only — no public exposure (NFR9)

### Story 8.3: Lifecycle Monitoring & Initiation Tuning

As an admin,
I want to monitor what the background lifecycle processes and tune initiation behavior,
So that I can verify the companion is genuinely alive between sessions and calibrate how often it reaches out.

**Acceptance Criteria:**

**Given** background lifecycle processing
**When** the admin inspects activity
**Then** they can see what was processed, what artifacts were produced, and when cycles completed (FR46)

**Given** a REST admin endpoint or log query
**When** reviewing lifecycle activity
**Then** it shows: cycle timestamps, conversations reflected on, artifacts produced (type, content summary), mailbox items processed, processing duration

**Given** spontaneous initiation configuration
**When** the admin adjusts the urge threshold or accumulation rate
**Then** changes are applied via `config/default.yaml` without code changes (FR47)

**Given** initiation tuning
**When** the threshold is raised
**Then** the companion reaches out less frequently — requiring stronger reasons

**Given** initiation tuning
**When** the accumulation rate is lowered
**Then** urge builds more slowly — the companion is more patient before initiating

**Given** lifecycle and initiation logs
**When** the admin reviews them
**Then** entries are tagged with subsystem ("lifecycle" or "initiation") for independent filtering

### Story 8.4: Personality Drift Metrics & Anchoring Configuration

As an admin,
I want to review personality drift over time and adjust anchoring weights,
So that I can assess whether the companion's voice is staying consistent while evolving, or fragmenting — and correct course if needed.

**Acceptance Criteria:**

**Given** the personality_snapshots table
**When** the admin reviews drift metrics
**Then** they can compare snapshots over time to assess consistency vs. fragmentation (FR48)

**Given** a REST admin endpoint or database query
**When** reviewing personality drift
**Then** it shows: snapshot timestamps, key personality parameter values over time, delta between snapshots, trend direction

**Given** personality anchoring weights
**When** the admin adjusts them
**Then** changes are applied via `config/companion-defaults.yaml` without code changes (FR49)

**Given** anchoring weight changes
**When** the next conversation occurs
**Then** the updated weights influence how strongly core traits resist drift vs. how freely surface traits evolve

**Given** drift detection
**When** significant divergence from baseline personality is detected
**Then** it is flagged in logs with subsystem tag "personality" for admin attention

**Given** personality metrics
**When** assessed over time
**Then** the admin can distinguish between healthy evolution (gradual opinion shifts, language absorption) and unhealthy fragmentation (inconsistent voice, contradictory temperament)

### Story 8.5: LLM Cost Tracking & Data Export Readiness

As an admin,
I want to track LLM costs per subsystem and per interaction, and know the system is ready for future data export,
So that I can understand where tokens are spent and be confident that all companion data can be extracted when needed.

**Acceptance Criteria:**

**Given** any LLM API call
**When** the request completes
**Then** token usage (input tokens, output tokens) is logged with subsystem tag, companion_id, and request context (NFR14)

**Given** LLM cost tracking
**When** the admin reviews token usage
**Then** they can see cost breakdown by subsystem (memory extraction vs. conversation vs. lifecycle reflection vs. emotional assessment) and by time period

**Given** cost tracking logs
**When** queried via `docker logs` with structured filtering
**Then** the admin can calculate total cost per interaction, per background processing cycle, and per subsystem

**Given** the companion ID scoping architecture
**When** all state tables are queried with a companion_id filter
**Then** a complete companion state export (memory graph, personality parameters, emotional state, conversation history, mailbox items, initiation events) can be extracted as JSON (NFR10)

**Given** data export readiness
**When** the architecture is reviewed
**Then** companion ID scoping on every table confirms that future data export requires no schema changes — only a query-and-serialize implementation

**Given** data ownership
**When** companion data is stored
**Then** it follows the user ownership design principle — no data shared with third parties beyond LLM API calls (NFR7)
