package app

import "context"

// Service defines the lifecycle contract for all application services.
// Every subsystem service must implement this interface.
type Service interface {
	// Start initializes the service and begins processing.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the service.
	Stop(ctx context.Context) error

	// Health returns nil if the service is healthy, or an error describing the issue.
	Health(ctx context.Context) error

	// Ready returns nil if the service is ready to accept work, or an error if not.
	Ready(ctx context.Context) error
}
