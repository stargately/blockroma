package worker

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// ErrCircuitOpen is returned when the circuit breaker is open
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrTooManyRequests is returned when too many requests are in flight
	ErrTooManyRequests = errors.New("too many requests")
)

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern for RPC calls
type CircuitBreaker struct {
	mu sync.RWMutex

	maxFailures     int           // Number of failures before opening
	resetTimeout    time.Duration // Time to wait before attempting to close
	halfOpenMax     int           // Max requests in half-open state
	requestTimeout  time.Duration // Timeout for individual requests

	state           CircuitBreakerState
	failures        int
	lastFailureTime time.Time
	halfOpenCount   int
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration, requestTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:    maxFailures,
		resetTimeout:   resetTimeout,
		halfOpenMax:    3,
		requestTimeout: requestTimeout,
		state:          StateClosed,
	}
}

// Call executes a function with circuit breaker protection
func (cb *CircuitBreaker) Call(ctx context.Context, fn func(context.Context) error) error {
	// Check if circuit breaker allows the call
	if !cb.canAttempt() {
		return ErrCircuitOpen
	}

	// Create timeout context for the request
	reqCtx, cancel := context.WithTimeout(ctx, cb.requestTimeout)
	defer cancel()

	// Execute the function
	err := fn(reqCtx)

	// Record the result
	cb.recordResult(err)

	return err
}

// canAttempt checks if a request can be attempted
func (cb *CircuitBreaker) canAttempt() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Check if we should transition to half-open
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.halfOpenCount = 1 // Count this as the first half-open request
			return true
		}
		return false
	case StateHalfOpen:
		// Allow limited requests in half-open state
		if cb.halfOpenCount < cb.halfOpenMax {
			cb.halfOpenCount++
			return true
		}
		return false
	}

	return false
}

// recordResult records the result of a request
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailureTime = time.Now()

		// If we've exceeded max failures, open the circuit
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		} else if cb.state == StateHalfOpen {
			// Failure in half-open state reopens the circuit
			cb.state = StateOpen
		}
	} else {
		// Success
		if cb.state == StateHalfOpen {
			// Successful request in half-open state closes the circuit
			cb.state = StateClosed
			cb.failures = 0
			cb.halfOpenCount = 0
		} else if cb.state == StateClosed {
			// Reset failure count on success
			cb.failures = 0
		}
	}
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Failures returns the current failure count
func (cb *CircuitBreaker) Failures() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
	cb.halfOpenCount = 0
}
