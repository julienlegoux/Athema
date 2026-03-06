package llm_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"athema/internal/infrastructure/llm"
)

func TestRateLimiterBasicAcquireRelease(t *testing.T) {
	rl := llm.NewRateLimiter(2)

	ctx := context.Background()
	if err := rl.Acquire(ctx); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if err := rl.Acquire(ctx); err != nil {
		t.Fatalf("second acquire failed: %v", err)
	}

	rl.Release()
	rl.Release()
}

func TestRateLimiterBlocksWhenExhausted(t *testing.T) {
	rl := llm.NewRateLimiter(1)
	ctx := context.Background()

	if err := rl.Acquire(ctx); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}

	// Second acquire should block
	done := make(chan struct{})
	go func() {
		ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		err := rl.Acquire(ctx2)
		if err == nil {
			t.Error("expected timeout error, got nil")
		}
		close(done)
	}()

	<-done
	rl.Release()
}

func TestRateLimiterContextCancellation(t *testing.T) {
	rl := llm.NewRateLimiter(1)
	ctx := context.Background()

	if err := rl.Acquire(ctx); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}

	ctx2, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := rl.Acquire(ctx2)
	if err == nil {
		t.Error("expected error from canceled context")
	}

	rl.Release()
}

func TestRateLimiterReleaseRestoresCapacity(t *testing.T) {
	rl := llm.NewRateLimiter(1)
	ctx := context.Background()

	if err := rl.Acquire(ctx); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	rl.Release()

	// Should be able to acquire again
	if err := rl.Acquire(ctx); err != nil {
		t.Fatalf("re-acquire after release failed: %v", err)
	}
	rl.Release()
}

func TestRateLimiterConcurrentAccess(t *testing.T) {
	maxConcurrent := 3
	rl := llm.NewRateLimiter(maxConcurrent)
	ctx := context.Background()

	var active atomic.Int32
	var maxSeen atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := rl.Acquire(ctx); err != nil {
				t.Errorf("acquire failed: %v", err)
				return
			}
			cur := active.Add(1)
			if cur > int32(maxConcurrent) {
				t.Errorf("concurrent count %d exceeds max %d", cur, maxConcurrent)
			}
			if cur > maxSeen.Load() {
				maxSeen.Store(cur)
			}
			time.Sleep(5 * time.Millisecond)
			active.Add(-1)
			rl.Release()
		}()
	}

	wg.Wait()

	if maxSeen.Load() == 0 {
		t.Error("expected some concurrent activity")
	}
}
