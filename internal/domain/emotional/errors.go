package emotional

import (
	"fmt"

	"athema/internal/domain"
)

var (
	ErrStateNotFound = fmt.Errorf("emotional: state %w", domain.ErrNotFound)
)
