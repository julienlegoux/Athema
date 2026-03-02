package eventbus

import (
	"fmt"
	"log/slog"
	"sync"

	"athema/internal/domain"
)

const defaultBufferSize = 256

type subscriber struct {
	ch      chan domain.Event
	handler func(domain.Event)
	done    chan struct{}
}

// Bus is an in-process event bus that delivers events asynchronously
// to registered subscribers via buffered channels.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]*subscriber
	logger      *slog.Logger
	bufferSize  int
	closed      bool
}

// NewBus creates a new event bus with the given logger.
func NewBus(logger *slog.Logger) *Bus {
	return &Bus{
		subscribers: make(map[string][]*subscriber),
		logger:      logger,
		bufferSize:  defaultBufferSize,
	}
}

// Subscribe registers a handler for the given event type.
// The handler will be called asynchronously in its own goroutine.
func (b *Bus) Subscribe(eventType string, handler func(domain.Event)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	sub := &subscriber{
		ch:      make(chan domain.Event, b.bufferSize),
		handler: handler,
		done:    make(chan struct{}),
	}

	b.subscribers[eventType] = append(b.subscribers[eventType], sub)

	go b.dispatch(sub, eventType)
}

// Publish sends an event to all subscribers registered for its type.
// Publishing is non-blocking; events are buffered per subscriber.
func (b *Bus) Publish(event domain.Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return
	}

	subs := b.subscribers[event.EventType()]
	for _, sub := range subs {
		select {
		case sub.ch <- event:
		default:
			b.logger.Warn("event bus: subscriber buffer full, dropping event",
				"eventType", event.EventType(),
				"companionId", event.GetCompanionID().String(),
			)
		}
	}
}

// Close gracefully shuts down the event bus, draining all subscriber channels.
func (b *Bus) Close() {
	b.mu.Lock()
	b.closed = true
	for _, subs := range b.subscribers {
		for _, sub := range subs {
			close(sub.ch)
		}
	}
	b.mu.Unlock()

	// Wait for all dispatchers to finish processing remaining events.
	b.mu.RLock()
	for _, subs := range b.subscribers {
		for _, sub := range subs {
			<-sub.done
		}
	}
	b.mu.RUnlock()
}

// dispatch reads events from a subscriber's channel and calls its handler.
// It recovers from handler panics to prevent one bad handler from crashing the bus.
func (b *Bus) dispatch(sub *subscriber, eventType string) {
	defer close(sub.done)

	for event := range sub.ch {
		b.safeCall(sub.handler, event, eventType)
	}
}

// safeCall invokes a handler with panic recovery.
func (b *Bus) safeCall(handler func(domain.Event), event domain.Event, eventType string) {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error("event bus: handler panicked",
				"eventType", eventType,
				"panic", fmt.Sprintf("%v", r),
			)
		}
	}()

	handler(event)
}
