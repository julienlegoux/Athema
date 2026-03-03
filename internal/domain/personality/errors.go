package personality

import (
	"fmt"

	"athema/internal/domain"
)

var (
	ErrSnapshotNotFound = fmt.Errorf("personality: snapshot %w", domain.ErrNotFound)
)
