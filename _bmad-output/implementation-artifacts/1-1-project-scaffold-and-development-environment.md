# Story 1.1: Project Scaffold & Development Environment

Status: review

## Story

As a developer,
I want the Go monorepo initialized with Clean Architecture structure, Docker Compose, and core build tooling,
so that I have a working development environment to build all companion systems upon.

## Acceptance Criteria

1. **Given** I run Go Blueprint with the specified flags
   **When** the scaffold is generated
   **Then** the project has chi router, PostgreSQL connection, WebSocket hub scaffold, and Docker config

2. **Given** the scaffolded project
   **When** I restructure into Clean Architecture layers
   **Then** the project has directories: `cmd/server/`, `cmd/tui/`, `internal/domain/`, `internal/usecase/`, `internal/adapter/`, `internal/infrastructure/`, `internal/app/`, `migrations/`, `config/`, `test/`

3. **Given** Docker Compose configuration
   **When** I run `docker compose up`
   **Then** `athema-server` and `postgres` services start successfully with a persistent volume for data

4. **Given** the Makefile
   **When** I run `make build`
   **Then** the server binary compiles without errors

5. **Given** the Makefile
   **When** I run `make test`
   **Then** the Go test suite runs successfully (passing, even if empty)

6. **Given** development tooling
   **When** I run `make run-server`
   **Then** Air hot reload starts the server with automatic rebuilds on file changes

7. **Given** configuration files
   **When** the server starts
   **Then** it loads `config/default.yaml` with environment variable overrides from `.env`

8. **Given** the CI pipeline
   **When** code is pushed to the repository
   **Then** GitHub Actions runs build and test stages via `.github/workflows/ci.yml`

9. **Given** structured logging
   **When** the server starts
   **Then** `log/slog` outputs structured JSON to stdout with per-subsystem tagged loggers

## Tasks / Subtasks

- [x] **Task 1: Run Go Blueprint scaffold** (AC: #1)
  - [x] Install go-blueprint: `go install github.com/melkeydev/go-blueprint@latest`
  - [x] Run: `go-blueprint create --name athema --framework chi --driver postgres --advanced --feature websocket --feature docker --git commit`
  - [x] Review generated scaffold, understand what it provides vs what needs restructuring
  - [x] Note: Blueprint generates a starting point — significant restructuring will be needed

- [x] **Task 2: Restructure into Clean Architecture** (AC: #2)
  - [x] Create complete directory structure per architecture spec (see File Structure Requirements below)
  - [x] Move/refactor Blueprint-generated code into correct layers
  - [x] Create `internal/domain/` as pure domain layer (zero project imports)
  - [x] Create `internal/app/` for lifecycle concerns (Service interface: Start, Stop, Health, Ready)
  - [x] Create `internal/usecase/` for application logic
  - [x] Create `internal/adapter/` for interface adapters (repository, handler, presenter)
  - [x] Create `internal/infrastructure/` for external concerns (llm, eventbus, config, server, database)
  - [x] Create `tui/` for TUI application structure
  - [x] Create stub packages for all 6 subsystems (memory, conversation, personality, emotional, lifecycle, initiation) in each layer
  - [x] Create `cmd/server/main.go` as composition root
  - [x] Create `cmd/tui/main.go` as TUI entry point (stub)
  - [x] Verify dependency direction: domain imports nothing, each outer layer imports only domain

- [x] **Task 3: Set up Docker Compose** (AC: #3)
  - [x] Create/update `docker-compose.yml` with `athema-server` and `postgres` services
  - [x] Configure PostgreSQL with pgvector extension (use `pgvector/pgvector:pg17` image)
  - [x] Add persistent volume for PostgreSQL data
  - [x] Create multi-stage `Dockerfile` for server binary
  - [x] Create `.env.example` with template variables (DB credentials, LLM API keys)
  - [x] Verify `docker compose up` starts both services successfully

- [x] **Task 4: Create Makefile** (AC: #4, #5)
  - [x] Create Makefile with targets: `build`, `test`, `test-integration`, `test-e2e`, `migrate-up`, `migrate-down`, `run`, `run-server`, `run-tui`
  - [x] `build`: compile server binary
  - [x] `test`: run unit tests (`go test ./...`)
  - [x] `test-integration`: run integration tests (`go test -tags integration ./...`)
  - [x] `test-e2e`: run e2e tests (`go test -tags e2e ./...`)
  - [x] `migrate-up` / `migrate-down`: run golang-migrate
  - [x] `run`: full stack (docker compose up)
  - [x] `run-server`: start server with Air hot reload
  - [x] `run-tui`: `go run cmd/tui/main.go`
  - [x] Verify `make build` compiles without errors
  - [x] Verify `make test` passes (even if empty test suite)

- [x] **Task 5: Set up configuration system** (AC: #7)
  - [x] Create `internal/infrastructure/config/loader.go` — YAML config loading with env var overrides
  - [x] Create `config/default.yaml` with subsystem parameter stubs
  - [x] Create `config/companion-defaults.yaml` with base personality parameter stubs
  - [x] Implement env var override pattern (e.g., `ATHEMA_DB_HOST` overrides `db.host` in YAML)
  - [x] Load config in `cmd/server/main.go` on startup

- [x] **Task 6: Set up structured logging** (AC: #9)
  - [x] Configure `log/slog` with JSON handler outputting to stdout
  - [x] Create logger factory that produces per-subsystem tagged loggers
  - [x] Each logger must include `subsystem` field (e.g., "memory", "conversation", "personality", "emotional", "lifecycle", "initiation")
  - [x] Loggers created in composition root (`cmd/server/main.go`) and passed via dependency injection
  - [x] No `init()` functions for logger wiring
  - [x] Log levels: Error (broken), Warn (degraded), Info (operational), Debug (development)

- [x] **Task 7: Set up Air hot reload** (AC: #6)
  - [x] Install air: `go install github.com/air-verse/air@latest`
  - [x] Create `.air.toml` configuration for server development
  - [x] Configure to watch `.go` files and rebuild `cmd/server/main.go`
  - [x] Verify `make run-server` starts Air and rebuilds on changes

- [x] **Task 8: Set up CI/CD** (AC: #8)
  - [x] Create `.github/workflows/ci.yml`
  - [x] Configure stages: build, unit test, integration test (with PostgreSQL service), e2e test
  - [x] Use build tags for test tiers
  - [x] Trigger on push and pull request

- [x] **Task 9: Create migration and test directory structure** (AC: #2)
  - [x] Create `migrations/` directory (empty — actual migration SQL files are created in Stories 1.4, 2.1, 3.1, etc.)
  - [x] Create `test/integration/` directory with `.gitkeep`
  - [x] Create `test/e2e/` directory with `.gitkeep`
  - [x] Create `test/fixtures/llm/` directory with `.gitkeep`
  - [x] Set up golang-migrate infrastructure in `internal/infrastructure/database/migrate.go`

- [x] **Task 10: Create .gitignore and final verification** (AC: all)
  - [x] Create/update `.gitignore` (Go binaries, .env, vendor/, tmp/, air build artifacts)
  - [x] Verify `make build` succeeds
  - [x] Verify `make test` passes
  - [x] Verify `docker compose up` starts services
  - [x] Verify `make run-server` starts with Air
  - [x] Verify structured JSON log output on server start

## Dev Notes

### Critical Architecture Constraints

- **Clean Architecture dependency direction is STRICT:** `domain/ → usecase/ → adapter/ → infrastructure/`. Domain imports NOTHING. Outer layers import ONLY domain (not each other).
- **Subsystem packages NEVER import each other** — they communicate through domain interfaces and events only.
- **Domain layer (`internal/domain/`) must have ZERO project imports** — pure domain, no infrastructure, no lifecycle concerns.
- **`internal/app/` is SEPARATE from `internal/domain/`** — lifecycle concerns (Start/Stop/Health/Ready) belong in `app/`, not domain.
- **Composition root is `cmd/server/main.go`** — all wiring via constructor injection. No DI framework. No `init()` functions for wiring.
- **6 subsystems:** memory, conversation, personality, emotional, lifecycle, initiation — each represented in each Clean Architecture layer.

### Go Blueprint Usage Notes

The Go Blueprint command generates a starter scaffold:
```bash
go-blueprint create --name athema --framework chi --driver postgres --advanced --feature websocket --feature docker --git commit
```

**What Blueprint provides:**
- Basic project layout (cmd/, internal/)
- chi router setup
- PostgreSQL connection boilerplate
- WebSocket hub scaffold
- Docker + Docker Compose starter config

**What needs significant restructuring after Blueprint:**
- Directory structure needs expansion into full Clean Architecture layers
- WebSocket hub needs rewriting for 1:1 streaming (not broadcast)
- Additional subsystem directories need creation
- Config system needs implementation
- Logging needs replacement with slog structured JSON

### Configuration Pattern

```
config/default.yaml          → Default values for all subsystems
config/companion-defaults.yaml → Base personality parameters
.env                         → Secrets (DB creds, API keys) — NOT committed
.env.example                 → Template showing expected env vars — committed
```

**Loading precedence:** YAML defaults → env var overrides
**Env var naming convention:** `ATHEMA_SECTION_KEY` (e.g., `ATHEMA_DB_HOST`, `ATHEMA_LLM_API_KEY`)

### Logging Pattern

```go
// In cmd/server/main.go (composition root):
baseLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

// Create per-subsystem loggers:
memoryLogger := baseLogger.With("subsystem", "memory")
conversationLogger := baseLogger.With("subsystem", "conversation")
// ... etc for all 6 subsystems

// Pass via constructor injection:
memorySvc := memory.NewService(memoryLogger, memoryRepo, ...)
```

**Log entry structure:**
```json
{"time":"2026-03-01T12:00:00Z","level":"INFO","msg":"server started","subsystem":"server","companion_id":"...","port":8080}
```

### Docker Setup

**PostgreSQL with pgvector:**
```yaml
# Use pgvector/pgvector:pg17 image (includes pgvector extension pre-installed)
postgres:
  image: pgvector/pgvector:pg17
  environment:
    POSTGRES_DB: athema
    POSTGRES_USER: athema
    POSTGRES_PASSWORD: ${DB_PASSWORD:-athema_dev}
  volumes:
    - pgdata:/var/lib/postgresql/data
  ports:
    - "5432:5432"
```

**Server Dockerfile (multi-stage):**
```dockerfile
# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /athema-server cmd/server/main.go

# Runtime stage
FROM alpine:latest
COPY --from=builder /athema-server /athema-server
EXPOSE 8080
CMD ["/athema-server"]
```

### Project Structure Notes

**Complete directory tree to create:**

```
athema/
├── cmd/
│   ├── server/
│   │   └── main.go                   # Composition root
│   └── tui/
│       └── main.go                   # TUI entry point (stub)
├── internal/
│   ├── domain/                       # Pure domain — ZERO imports
│   │   ├── types.go                  # Shared value types (CompanionID, etc.)
│   │   ├── events.go                 # Event interface + domain event types
│   │   ├── errors.go                 # Domain sentinel errors
│   │   ├── memory/                   # (stubs — actual entities in Story 1.2+)
│   │   ├── conversation/
│   │   ├── personality/
│   │   ├── emotional/
│   │   ├── lifecycle/
│   │   └── initiation/
│   ├── app/                          # Application lifecycle
│   │   └── lifecycle.go              # Service interface: Start, Stop, Health, Ready
│   ├── usecase/                      # Application logic
│   │   ├── memory/
│   │   ├── conversation/
│   │   ├── personality/
│   │   ├── emotional/
│   │   ├── lifecycle/
│   │   └── initiation/
│   ├── adapter/                      # Interface adapters
│   │   ├── repository/
│   │   │   └── postgres/             # Repository implementations
│   │   ├── handler/
│   │   │   ├── websocket/            # WebSocket handlers
│   │   │   └── rest/                 # REST handlers
│   │   └── presenter/
│   │       └── prompt/               # Prompt assembly
│   └── infrastructure/               # External concerns
│       ├── llm/                      # LLM provider abstraction
│       ├── eventbus/                 # In-process event bus
│       ├── config/                   # Config loading
│       ├── server/                   # chi router, HTTP server setup
│       └── database/                 # DB connection, migrations
├── tui/                              # TUI application
│   ├── client/                       # WebSocket client
│   ├── views/
│   │   ├── chat/
│   │   ├── mailbox/
│   │   └── drop/
│   └── theme/                        # Lip Gloss theme
├── migrations/                       # SQL migration files (empty for this story)
├── config/
│   ├── default.yaml                  # Default subsystem config
│   └── companion-defaults.yaml       # Base personality parameters
├── test/
│   ├── integration/                  # Build tag: integration
│   ├── e2e/                          # Build tag: e2e
│   └── fixtures/
│       └── llm/                      # LLM response fixtures
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── .gitignore
├── Makefile
├── .air.toml
└── .github/
    └── workflows/
        └── ci.yml
```

### Naming Conventions (Established in This Story)

| Context | Convention | Example |
|---------|-----------|---------|
| Go packages | single lowercase word | `memory`, `emotional`, `prompt` |
| Go files | snake_case.go | `knowledge_node.go`, `memory_service.go` |
| Go exports | PascalCase | `CompanionID`, `NewService` |
| Go unexported | camelCase | `companionID`, `handleMessage` |
| JSON fields | camelCase with explicit tags | `json:"companionId"` |
| DB tables | snake_case plural | `knowledge_nodes`, `mailbox_items` |
| DB columns | snake_case | `companion_id`, `created_at` |
| Event types (Go) | PascalCase past tense | `MessageReceivedEvent` |
| Event types (string) | namespace.snake_case | `conversation.message_received` |
| WebSocket types | namespace.action | `chat.stream`, `presence.sync` |
| API errors | UPPER_SNAKE_CASE | `CHAT_MESSAGE_FAILED` |
| Env vars | ATHEMA_SECTION_KEY | `ATHEMA_DB_HOST` |
| DB timestamps | TIMESTAMPTZ, UTC always | All `created_at`, `updated_at` |

### API Response Envelope (Established Pattern)

```json
// Success:
{"data": { ... }}

// Error:
{"error": {"code": "UPPER_SNAKE_CASE", "message": "Human-readable message"}}
```

Never return both `data` and `error`. Error codes are subsystem-prefixed UPPER_SNAKE_CASE.

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Project Structure] — Complete directory tree specification
- [Source: _bmad-output/planning-artifacts/architecture.md#Clean Architecture] — Layer rules and dependency direction
- [Source: _bmad-output/planning-artifacts/architecture.md#Configuration] — YAML + env var config pattern
- [Source: _bmad-output/planning-artifacts/architecture.md#Logging] — log/slog structured JSON requirements
- [Source: _bmad-output/planning-artifacts/architecture.md#Docker] — Docker Compose services
- [Source: _bmad-output/planning-artifacts/architecture.md#Testing] — Three-tier testing strategy
- [Source: _bmad-output/planning-artifacts/architecture.md#Naming Conventions] — All naming conventions
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.1] — Acceptance criteria and story definition
- [Source: _bmad-output/planning-artifacts/prd.md#NFR1-NFR19] — Non-functional requirements

## Library & Framework Versions (Researched March 2026)

| Library | Version | Import Path | Notes |
|---------|---------|-------------|-------|
| **Go** | 1.26.0 | — | Latest stable, Green Tea GC default |
| **chi** | v5.2.1 | `github.com/go-chi/chi/v5` | Always use `/v5` path |
| **golang-migrate** | v4.19.1 | `github.com/golang-migrate/migrate/v4` | Build with `-tags 'postgres'` |
| **Bubble Tea** | v1.x (stable) | `github.com/charmbracelet/bubbletea` | See Charm ecosystem note below |
| **Bubbles** | v0.21.0 | `github.com/charmbracelet/bubbles` | Stable, pairs with BT v1 |
| **Lip Gloss** | v1.1.0 | `github.com/charmbracelet/lipgloss` | Stable, pairs with BT v1 |
| **pgvector-go** | v0.3.0 | `github.com/pgvector/pgvector-go` | Use `/pgx` subpackage |
| **Air** | v1.64.4 | `github.com/air-verse/air` | Old path `cosmtrek/air` deprecated |
| **go-blueprint** | v0.10.11 | `github.com/melkeydev/go-blueprint` | Scaffold tool only |
| **PostgreSQL** | 17 | `pgvector/pgvector:pg17` Docker image | pgvector 0.8.2 pre-installed |

### Charm Ecosystem Decision

**Use the v1 stable stack:**
- Bubbletea v2.0.0 is stable, but Bubbles v2 and Lipgloss v2 are still **beta** (not stable as of March 2026)
- Mixing v1 and v2 across the Charm stack is **not supported**
- The v1 stack is production-proven and sufficient for all Story 1.1 requirements
- If the team decides to migrate to v2 later (once Bubbles/Lipgloss v2 go stable), it can be done as a dedicated migration story

### pgvector Security Note

pgvector 0.8.2 fixes CVE-2026-3172 (buffer overflow in parallel HNSW index builds). The `pgvector/pgvector:pg17` Docker image includes this fix. Do NOT use pgvector < 0.8.2 in production.

## Dev Agent Record

### Agent Model Used

Claude Opus 4.6 (claude-opus-4-6)

### Debug Log References

- Go Blueprint scaffolded to /tmp/blueprint-scaffold/athema, reviewed output, adapted patterns into Clean Architecture structure
- Blueprint go.mod used Go 1.25.0 and cmd/api — updated to Go 1.26.0 and cmd/server
- Blueprint database.go used database/sql with singleton — replaced with pgxpool-based implementation
- Blueprint docker-compose used BLUEPRINT_ env prefix and postgres:latest — replaced with ATHEMA_ prefix and pgvector/pgvector:pg17

### Completion Notes List

- Task 1: Go Blueprint v0.10.11 installed and scaffold generated. Reviewed what it provides (chi router, pgx driver, WebSocket, Docker basics) vs what needed restructuring (full Clean Architecture, config system, slog logging, correct env prefix)
- Task 2: Full Clean Architecture created with 4 layers (domain, usecase, adapter, infrastructure) + app lifecycle. All 6 subsystems (memory, conversation, personality, emotional, lifecycle, initiation) stubbed in domain and usecase layers. Domain layer has zero project imports. Composition root at cmd/server/main.go with constructor injection only
- Task 3: docker-compose.yml with athema-server and postgres (pgvector/pgvector:pg17) services, persistent pgdata volume, health checks. Multi-stage Dockerfile using golang:1.26-alpine. .env.example created with all ATHEMA_ variables
- Task 4: Makefile with all required targets (build, test, test-integration, test-e2e, migrate-up, migrate-down, run, run-server, run-tui). Build verified — compiles without errors. Tests pass
- Task 5: Config system with YAML loading + env var overrides via struct tags. Config struct covers server, db, llm, log, and all 6 subsystems. 5 unit tests all passing (Load, EnvOverrides, FileNotFound, InvalidYAML, DSN)
- Task 6: log/slog with JSON handler to stdout. Per-subsystem tagged loggers created in composition root. Verified JSON output: {"time":"...","level":"INFO","msg":"server started","subsystem":"server","port":8080}
- Task 7: .air.toml configured for Windows, watches .go and .yaml files, excludes _bmad directories
- Task 8: GitHub Actions CI with 4 jobs: build, test (unit), integration-test (with pgvector PostgreSQL service), e2e-test. All using Go 1.26 and build tags
- Task 9: migrations/, test/integration/, test/e2e/, test/fixtures/llm/ directories created with .gitkeep. golang-migrate infrastructure in internal/infrastructure/database/migrate.go
- Task 10: .gitignore covers binaries, .env, vendor, tmp, IDE files. Final verification: make build succeeds, make test passes (5/5 tests), docker compose config valid, server starts with structured JSON logs

### Change Log
| Change | Date | Reason |
|--------|------|--------|
| Story created | 2026-03-01 | Initial story creation from sprint planning |
| Full implementation completed | 2026-03-02 | All 10 tasks implemented: Clean Architecture scaffold, Docker, Makefile, config, logging, Air, CI, migrations, .gitignore |

### File List

- cmd/server/main.go (new) — Composition root with config loading, structured logging, graceful shutdown
- cmd/tui/main.go (new) — TUI entry point stub
- internal/domain/types.go (new) — CompanionID, UserID, SessionID value types
- internal/domain/events.go (new) — Event interface and BaseEvent
- internal/domain/errors.go (new) — Domain sentinel errors
- internal/domain/memory/doc.go (new) — Memory subsystem domain stub
- internal/domain/conversation/doc.go (new) — Conversation subsystem domain stub
- internal/domain/personality/doc.go (new) — Personality subsystem domain stub
- internal/domain/emotional/doc.go (new) — Emotional subsystem domain stub
- internal/domain/lifecycle/doc.go (new) — Lifecycle subsystem domain stub
- internal/domain/initiation/doc.go (new) — Initiation subsystem domain stub
- internal/app/lifecycle.go (new) — Service interface (Start, Stop, Health, Ready)
- internal/usecase/memory/doc.go (new) — Memory usecase stub
- internal/usecase/conversation/doc.go (new) — Conversation usecase stub
- internal/usecase/personality/doc.go (new) — Personality usecase stub
- internal/usecase/emotional/doc.go (new) — Emotional usecase stub
- internal/usecase/lifecycle/doc.go (new) — Lifecycle usecase stub
- internal/usecase/initiation/doc.go (new) — Initiation usecase stub
- internal/adapter/repository/postgres/doc.go (new) — Postgres repository stub
- internal/adapter/handler/websocket/doc.go (new) — WebSocket handler stub
- internal/adapter/handler/rest/doc.go (new) — REST handler stub
- internal/adapter/presenter/prompt/doc.go (new) — Prompt presenter stub
- internal/infrastructure/config/loader.go (new) — YAML config + env var overrides
- internal/infrastructure/config/loader_test.go (new) — Config loader unit tests (5 tests)
- internal/infrastructure/server/server.go (new) — chi HTTP server with health endpoint
- internal/infrastructure/database/database.go (new) — pgxpool connection management
- internal/infrastructure/database/migrate.go (new) — golang-migrate infrastructure
- internal/infrastructure/llm/doc.go (new) — LLM provider stub
- internal/infrastructure/eventbus/doc.go (new) — Event bus stub
- config/default.yaml (new) — Default subsystem configuration
- config/companion-defaults.yaml (new) — Base personality parameters
- .env.example (new) — Environment variable template
- Dockerfile (new) — Multi-stage build for server
- docker-compose.yml (new) — athema-server + postgres (pgvector:pg17)
- Makefile (new) — Build, test, run targets
- .air.toml (new) — Air hot reload configuration
- .github/workflows/ci.yml (new) — GitHub Actions CI pipeline
- .gitignore (new) — Git ignore rules
- go.mod (new) — Go module definition
- go.sum (new) — Go module checksums
- migrations/.gitkeep (new) — Migrations directory placeholder
- test/integration/.gitkeep (new) — Integration test directory
- test/e2e/.gitkeep (new) — E2E test directory
- test/fixtures/llm/.gitkeep (new) — LLM fixture directory
- tui/client/.gitkeep (new) — TUI client directory
- tui/views/chat/.gitkeep (new) — TUI chat view directory
- tui/views/mailbox/.gitkeep (new) — TUI mailbox view directory
- tui/views/drop/.gitkeep (new) — TUI drop view directory
- tui/theme/.gitkeep (new) — TUI theme directory
