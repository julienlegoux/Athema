.PHONY: all build test test-race test-integration test-e2e run run-server run-tui migrate-up migrate-down clean docker-up docker-down

all: build test

# Build the server binary
build:
	@echo "Building athema-server..."
	@go build -o bin/athema-server cmd/server/main.go

# Run unit tests
test:
	@echo "Running unit tests..."
	@go test ./... -v

# Run unit tests with race detector (requires CGO — matches CI)
test-race:
	@echo "Running unit tests with race detector..."
	@CGO_ENABLED=1 go test -race ./... -v

# Run integration tests (requires running PostgreSQL)
test-integration:
	@echo "Running integration tests..."
	@go test -tags integration ./... -v

# Run end-to-end tests
test-e2e:
	@echo "Running e2e tests..."
	@go test -tags e2e ./... -v

# Run full stack via Docker Compose
run:
	@docker compose up --build

# Run server with Air hot reload
run-server:
	@air

# Run TUI client
run-tui:
	@go run cmd/tui/main.go

# Run database migrations up
migrate-up:
	@go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate -path migrations -database "$${ATHEMA_DB_DSN}" up

# Run database migrations down
migrate-down:
	@go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate -path migrations -database "$${ATHEMA_DB_DSN}" down

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/ tmp/

# Start Docker services
docker-up:
	@docker compose up -d

# Stop Docker services
docker-down:
	@docker compose down
