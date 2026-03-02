package websocket

import (
	"log/slog"
	"sync"
)

// Hub manages active WebSocket client connections.
// Scaffold — full 1:1 streaming implementation in Story 1.5.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
	logger  *slog.Logger
}

// Client represents a single WebSocket connection to the hub.
type Client struct {
	Hub  *Hub
	Send chan []byte
}

// NewHub creates a new Hub.
func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
		logger:  logger,
	}
}

// Register adds a client to the hub.
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = struct{}{}
	h.logger.Info("client registered", "total", len(h.clients))
}

// Unregister removes a client from the hub and closes its send channel.
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.Send)
		h.logger.Info("client unregistered", "total", len(h.clients))
	}
}

// Count returns the number of connected clients.
func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
