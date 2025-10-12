package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(5, 1*time.Second, 500*time.Millisecond)

	if cb.maxFailures != 5 {
		t.Errorf("expected maxFailures=5, got %d", cb.maxFailures)
	}

	if cb.state != StateClosed {
		t.Errorf("expected initial state=StateClosed, got %v", cb.state)
	}

	if cb.failures != 0 {
		t.Errorf("expected failures=0, got %d", cb.failures)
	}
}

func TestCircuitBreaker_SuccessfulCalls(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond, 500*time.Millisecond)
	ctx := context.Background()

	callCount := 0
	fn := func(ctx context.Context) error {
		callCount++
		return nil
	}

	// Make 5 successful calls
	for i := 0; i < 5; i++ {
		err := cb.Call(ctx, fn)
		if err != nil {
			t.Errorf("call %d failed: %v", i, err)
		}
	}

	if callCount != 5 {
		t.Errorf("expected 5 calls, got %d", callCount)
	}

	if cb.State() != StateClosed {
		t.Errorf("expected state=StateClosed, got %v", cb.State())
	}

	if cb.Failures() != 0 {
		t.Errorf("expected failures=0, got %d", cb.Failures())
	}
}

func TestCircuitBreaker_FailureOpensCircuit(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond, 500*time.Millisecond)
	ctx := context.Background()

	callCount := 0
	fn := func(ctx context.Context) error {
		callCount++
		return errors.New("error")
	}

	// Make calls until circuit opens
	for i := 0; i < 3; i++ {
		err := cb.Call(ctx, fn)
		if err == nil {
			t.Errorf("call %d should have failed", i)
		}
	}

	if cb.State() != StateOpen {
		t.Errorf("expected state=StateOpen after 3 failures, got %v", cb.State())
	}

	// Next call should be blocked
	err := cb.Call(ctx, fn)
	if err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}

	// Function should not have been called again
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestCircuitBreaker_HalfOpenState(t *testing.T) {
	resetTimeout := 100 * time.Millisecond
	cb := NewCircuitBreaker(2, resetTimeout, 500*time.Millisecond)
	ctx := context.Background()

	// Open the circuit with failures
	fn := func(ctx context.Context) error {
		return errors.New("error")
	}

	cb.Call(ctx, fn)
	cb.Call(ctx, fn)

	if cb.State() != StateOpen {
		t.Errorf("expected state=StateOpen, got %v", cb.State())
	}

	// Wait for reset timeout
	time.Sleep(resetTimeout + 50*time.Millisecond)

	// Next call should be allowed (half-open state)
	callCount := 0
	successFn := func(ctx context.Context) error {
		callCount++
		return nil
	}

	err := cb.Call(ctx, successFn)
	if err != nil {
		t.Errorf("expected successful call in half-open state, got %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Successful call should close the circuit
	if cb.State() != StateClosed {
		t.Errorf("expected state=StateClosed after successful half-open call, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	resetTimeout := 100 * time.Millisecond
	cb := NewCircuitBreaker(2, resetTimeout, 500*time.Millisecond)
	ctx := context.Background()

	// Open the circuit
	fn := func(ctx context.Context) error {
		return errors.New("error")
	}

	cb.Call(ctx, fn)
	cb.Call(ctx, fn)

	// Wait for reset timeout
	time.Sleep(resetTimeout + 50*time.Millisecond)

	// Call should fail in half-open state
	err := cb.Call(ctx, fn)
	if err == nil {
		t.Error("expected call to fail")
	}

	// Circuit should reopen
	if cb.State() != StateOpen {
		t.Errorf("expected state=StateOpen after failed half-open call, got %v", cb.State())
	}
}

func TestCircuitBreaker_Timeout(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond, 50*time.Millisecond)
	ctx := context.Background()

	// Function that takes longer than timeout
	fn := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return nil
		}
	}

	err := cb.Call(ctx, fn)
	if err == nil {
		t.Error("expected timeout error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond, 500*time.Millisecond)
	ctx := context.Background()

	// Open the circuit
	fn := func(ctx context.Context) error {
		return errors.New("error")
	}

	cb.Call(ctx, fn)
	cb.Call(ctx, fn)

	if cb.State() != StateOpen {
		t.Errorf("expected state=StateOpen, got %v", cb.State())
	}

	// Reset the circuit
	cb.Reset()

	if cb.State() != StateClosed {
		t.Errorf("expected state=StateClosed after reset, got %v", cb.State())
	}

	if cb.Failures() != 0 {
		t.Errorf("expected failures=0 after reset, got %d", cb.Failures())
	}
}

func TestCircuitBreaker_MixedResults(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond, 500*time.Millisecond)
	ctx := context.Background()

	successFn := func(ctx context.Context) error {
		return nil
	}

	errorFn := func(ctx context.Context) error {
		return errors.New("error")
	}

	// Mix of successes and failures
	cb.Call(ctx, errorFn)     // failure=1
	cb.Call(ctx, successFn)   // failure=0 (reset)
	cb.Call(ctx, errorFn)     // failure=1
	cb.Call(ctx, errorFn)     // failure=2
	cb.Call(ctx, successFn)   // failure=0 (reset)

	// Circuit should still be closed
	if cb.State() != StateClosed {
		t.Errorf("expected state=StateClosed, got %v", cb.State())
	}

	// Now hit max failures
	cb.Call(ctx, errorFn) // failure=1
	cb.Call(ctx, errorFn) // failure=2
	cb.Call(ctx, errorFn) // failure=3 -> OPEN

	if cb.State() != StateOpen {
		t.Errorf("expected state=StateOpen, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenMaxRequests(t *testing.T) {
	resetTimeout := 50 * time.Millisecond
	cb := NewCircuitBreaker(2, resetTimeout, 500*time.Millisecond)
	ctx := context.Background()

	// Open the circuit
	fn := func(ctx context.Context) error {
		return errors.New("error")
	}

	cb.Call(ctx, fn)
	cb.Call(ctx, fn)

	// Wait for reset timeout
	time.Sleep(resetTimeout + 20*time.Millisecond)

	// In half-open state, only halfOpenMax requests should be allowed
	// Use quick function to avoid slow test
	allowedCount := 0
	blockedCount := 0

	// Try to make more than halfOpenMax requests sequentially
	for i := 0; i < 10; i++ {
		err := cb.Call(ctx, func(ctx context.Context) error {
			// Fail the requests to stay in half-open or reopen
			if allowedCount < cb.halfOpenMax {
				return errors.New("error")
			}
			return errors.New("error")
		})
		if err == ErrCircuitOpen {
			blockedCount++
		} else {
			allowedCount++
		}
	}

	// Should allow up to halfOpenMax (3) requests before reopening or blocking
	if allowedCount > cb.halfOpenMax {
		t.Errorf("expected max %d allowed requests, got %d", cb.halfOpenMax, allowedCount)
	}

	if blockedCount == 0 {
		t.Error("expected some requests to be blocked")
	}

	t.Logf("Allowed: %d, Blocked: %d", allowedCount, blockedCount)
}
