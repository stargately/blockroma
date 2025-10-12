package worker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	tests := []struct {
		name       string
		maxWorkers int
		expected   int
	}{
		{"positive workers", 5, 5},
		{"zero workers", 0, 1},
		{"negative workers", -5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewPool(tt.maxWorkers)
			if pool.maxWorkers != tt.expected {
				t.Errorf("expected %d workers, got %d", tt.expected, pool.maxWorkers)
			}
		})
	}
}

func TestPool_ExecuteTasks(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(3)
	pool.Start(ctx)

	var counter int32
	taskCount := 10

	// Submit tasks that increment a counter
	for i := 0; i < taskCount; i++ {
		pool.Submit(func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
	}

	results := pool.Wait()

	// All tasks should complete successfully
	if len(results) != taskCount {
		t.Errorf("expected %d results, got %d", taskCount, len(results))
	}

	for i, result := range results {
		if result.Error != nil {
			t.Errorf("result %d has error: %v", i, result.Error)
		}
	}

	// Counter should be incremented taskCount times
	if counter != int32(taskCount) {
		t.Errorf("expected counter to be %d, got %d", taskCount, counter)
	}
}

func TestPool_ExecuteTasksWithErrors(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(2)
	pool.Start(ctx)

	expectedErr := errors.New("task error")
	taskCount := 5
	errorCount := 0

	for i := 0; i < taskCount; i++ {
		i := i // Capture loop variable
		pool.Submit(func(ctx context.Context) error {
			if i%2 == 0 {
				return expectedErr
			}
			return nil
		})
	}

	results := pool.Wait()

	for _, result := range results {
		if result.Error != nil {
			errorCount++
		}
	}

	// We expect 3 errors (indices 0, 2, 4)
	if errorCount != 3 {
		t.Errorf("expected 3 errors, got %d", errorCount)
	}
}

func TestPool_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := NewPool(2)
	pool.Start(ctx)

	var started int32
	var completed int32

	// Submit tasks that will be cancelled
	for i := 0; i < 5; i++ {
		pool.Submit(func(ctx context.Context) error {
			atomic.AddInt32(&started, 1)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				atomic.AddInt32(&completed, 1)
				return nil
			}
		})
	}

	// Cancel after a short delay
	time.Sleep(10 * time.Millisecond)
	cancel()

	results := pool.Wait()

	// Some tasks should have been cancelled
	t.Logf("Started: %d, Completed: %d, Results: %d", started, completed, len(results))
}

func TestExecute_ConvenienceMethod(t *testing.T) {
	ctx := context.Background()
	var counter int32

	tasks := []Task{
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		},
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		},
		func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return errors.New("error")
		},
	}

	results := Execute(ctx, 2, tasks)

	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	if counter != 3 {
		t.Errorf("expected counter to be 3, got %d", counter)
	}

	// One result should have an error
	errorCount := 0
	for _, result := range results {
		if result.Error != nil {
			errorCount++
		}
	}

	if errorCount != 1 {
		t.Errorf("expected 1 error, got %d", errorCount)
	}
}

func TestPool_Concurrency(t *testing.T) {
	ctx := context.Background()
	maxWorkers := 3
	pool := NewPool(maxWorkers)
	pool.Start(ctx)

	var concurrent int32
	var maxConcurrent int32
	var mu sync.Mutex

	taskCount := 20

	for i := 0; i < taskCount; i++ {
		pool.Submit(func(ctx context.Context) error {
			current := atomic.AddInt32(&concurrent, 1)

			mu.Lock()
			if current > maxConcurrent {
				maxConcurrent = current
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)
			atomic.AddInt32(&concurrent, -1)
			return nil
		})
	}

	pool.Wait()

	// Max concurrent should not exceed maxWorkers
	if maxConcurrent > int32(maxWorkers) {
		t.Errorf("max concurrent tasks (%d) exceeded max workers (%d)", maxConcurrent, maxWorkers)
	}

	t.Logf("Max concurrent tasks: %d (limit: %d)", maxConcurrent, maxWorkers)
}

func TestPool_EmptyTasks(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(2)
	pool.Start(ctx)

	results := pool.Wait()

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
