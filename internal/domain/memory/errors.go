package memory

import "errors"

var (
	ErrNodeNotFound = errors.New("memory.node_not_found")
	ErrEdgeNotFound = errors.New("memory.edge_not_found")
)
