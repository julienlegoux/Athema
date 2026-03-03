package memory

import (
	"fmt"

	"athema/internal/domain"
)

var (
	ErrNodeNotFound = fmt.Errorf("memory: node %w", domain.ErrNotFound)
	ErrEdgeNotFound = fmt.Errorf("memory: edge %w", domain.ErrNotFound)
)
