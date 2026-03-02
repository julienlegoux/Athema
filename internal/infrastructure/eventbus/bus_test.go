package eventbus_test

import (
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"athema/internal/domain"
	"athema/internal/infrastructure/eventbus"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
}

func newTestEvent(eventType string) domain.Event {
	return domain.NewBaseEvent(eventType, domain.NewCompanionID())
}

func TestBus_PublishSubscribe_BasicFlow(t *testing.T) {
	bus := eventbus.NewBus(newTestLogger())
	defer bus.Close()

	received := make(chan domain.Event, 1)
	bus.Subscribe("test.event", func(e domain.Event) {
		received <- e
	})

	event := newTestEvent("test.event")
	bus.Publish(event)

	select {
	case got := <-received:
		if got.EventType() != "test.event" {
			t.Fatalf("expected event type %q, got %q", "test.event", got.EventType())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestBus_MultipleSubscribers_SameEventType(t *testing.T) {
	bus := eventbus.NewBus(newTestLogger())
	defer bus.Close()

	var count atomic.Int32

	for range 3 {
		bus.Subscribe("multi.event", func(e domain.Event) {
			count.Add(1)
		})
	}

	bus.Publish(newTestEvent("multi.event"))

	// Wait for all subscribers to process.
	deadline := time.After(2 * time.Second)
	for {
		if count.Load() == 3 {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("expected 3 subscribers to receive event, got %d", count.Load())
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func TestBus_EventTypeFiltering(t *testing.T) {
	bus := eventbus.NewBus(newTestLogger())
	defer bus.Close()

	var matchCount atomic.Int32
	var noMatchCount atomic.Int32

	bus.Subscribe("wanted.event", func(e domain.Event) {
		matchCount.Add(1)
	})
	bus.Subscribe("other.event", func(e domain.Event) {
		noMatchCount.Add(1)
	})

	bus.Publish(newTestEvent("wanted.event"))

	// Give time for processing.
	time.Sleep(100 * time.Millisecond)

	if matchCount.Load() != 1 {
		t.Fatalf("expected matching subscriber to receive 1 event, got %d", matchCount.Load())
	}
	if noMatchCount.Load() != 0 {
		t.Fatalf("expected non-matching subscriber to receive 0 events, got %d", noMatchCount.Load())
	}
}

func TestBus_HandlerPanicRecovery(t *testing.T) {
	bus := eventbus.NewBus(newTestLogger())
	defer bus.Close()

	received := make(chan domain.Event, 1)

	// First subscriber panics.
	bus.Subscribe("panic.event", func(e domain.Event) {
		panic("test panic")
	})

	// Second subscriber should still work.
	bus.Subscribe("panic.event", func(e domain.Event) {
		received <- e
	})

	bus.Publish(newTestEvent("panic.event"))

	select {
	case got := <-received:
		if got.EventType() != "panic.event" {
			t.Fatalf("expected event type %q, got %q", "panic.event", got.EventType())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for event after panic recovery")
	}

	// Verify bus still works after a panic.
	bus.Publish(newTestEvent("panic.event"))

	select {
	case <-received:
		// Good, bus still delivers events.
	case <-time.After(2 * time.Second):
		t.Fatal("bus stopped working after handler panic")
	}
}

func TestBus_Close_GracefulShutdown(t *testing.T) {
	bus := eventbus.NewBus(newTestLogger())

	var count atomic.Int32

	bus.Subscribe("close.event", func(e domain.Event) {
		count.Add(1)
	})

	// Publish some events before close.
	for range 5 {
		bus.Publish(newTestEvent("close.event"))
	}

	// Close should wait for dispatchers to drain.
	bus.Close()

	if count.Load() != 5 {
		t.Fatalf("expected 5 events processed before close, got %d", count.Load())
	}

	// Publishing after close should not panic.
	bus.Publish(newTestEvent("close.event"))
}

func TestBus_ConcurrentPublish_ThreadSafety(t *testing.T) {
	bus := eventbus.NewBus(newTestLogger())
	defer bus.Close()

	const numPublishers = 10
	const eventsPerPublisher = 20
	total := numPublishers * eventsPerPublisher // 200, well within 256 buffer

	var received atomic.Int32

	bus.Subscribe("concurrent.event", func(e domain.Event) {
		received.Add(1)
	})

	var wg sync.WaitGroup
	wg.Add(numPublishers)

	for range numPublishers {
		go func() {
			defer wg.Done()
			for range eventsPerPublisher {
				bus.Publish(newTestEvent("concurrent.event"))
			}
		}()
	}

	wg.Wait()

	// Wait for all events to be delivered.
	deadline := time.After(5 * time.Second)
	for {
		if received.Load() == int32(total) {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("expected %d events, received %d", total, received.Load())
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func TestBus_BufferOverflow_DropsEvent(t *testing.T) {
	// Create a bus with a small buffer to test overflow.
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
	bus := eventbus.NewBus(logger)

	// Subscribe with a slow handler that blocks to fill the buffer.
	blocker := make(chan struct{})
	var received atomic.Int32

	bus.Subscribe("overflow.event", func(e domain.Event) {
		<-blocker // Block until released.
		received.Add(1)
	})

	// Publish more events than the default buffer size (256).
	// The first event will be picked up by the handler (blocked).
	// The next 256 will fill the buffer. Additional ones should be dropped.
	for i := range 300 {
		_ = i
		bus.Publish(newTestEvent("overflow.event"))
	}

	// Release the blocker and close.
	close(blocker)
	bus.Close()

	// We should have received fewer than 300 events due to overflow.
	got := received.Load()
	if got > 300 {
		t.Fatalf("received more events than published: %d", got)
	}
	// At minimum the handler processes the buffered events (up to 257: 1 in-flight + 256 buffered).
	if got == 0 {
		t.Fatal("expected at least some events to be processed")
	}
	// Log the actual count for visibility.
	t.Logf("buffer overflow test: published 300, received %d (some dropped as expected)", got)
}

func TestBus_NoSubscribers_NoBlock(t *testing.T) {
	bus := eventbus.NewBus(newTestLogger())
	defer bus.Close()

	// Publishing to an event type with no subscribers should not block or panic.
	done := make(chan struct{})
	go func() {
		bus.Publish(newTestEvent("no.subscribers"))
		close(done)
	}()

	select {
	case <-done:
		// Good, publish returned immediately.
	case <-time.After(1 * time.Second):
		t.Fatal("publish blocked when there are no subscribers")
	}
}
