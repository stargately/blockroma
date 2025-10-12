package worker

import (
	"context"
	"errors"
	"fmt"
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

// Logger interface for circuit breaker logging
type Logger interface {
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	Warn(msg string)
	Error(msg string)
	Info(msg string)
}

// FailureRecord tracks a recent failure
type FailureRecord struct {
	Error     string
	Timestamp time.Time
}

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

	// Enhanced logging
	logger         Logger
	recentErrors   []FailureRecord // Track recent errors for debugging
	maxErrorHistory int            // Max errors to keep in history
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration, requestTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:     maxFailures,
		resetTimeout:    resetTimeout,
		halfOpenMax:     3,
		requestTimeout:  requestTimeout,
		state:           StateClosed,
		maxErrorHistory: 10, // Keep last 10 errors for debugging
		recentErrors:    make([]FailureRecord, 0, 10),
	}
}

// SetLogger sets the logger for the circuit breaker
func (cb *CircuitBreaker) SetLogger(logger Logger) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.logger = logger
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

		// Track the error in history
		cb.addErrorToHistory(err)

		previousState := cb.state

		// If we've exceeded max failures, open the circuit
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			// Log when circuit opens
			if previousState != StateOpen && cb.logger != nil {
				cb.logCircuitOpened()
			}
		} else if cb.state == StateHalfOpen {
			// Failure in half-open state reopens the circuit
			cb.state = StateOpen
			if cb.logger != nil {
				cb.logger.WithError(err).WithFields(map[string]interface{}{
					"previousState": "half_open",
					"failures":      cb.failures,
				}).Warn("Circuit breaker reopened after half-open failure")
			}
		}
	} else {
		// Success
		if cb.state == StateHalfOpen {
			// Successful request in half-open state closes the circuit
			previousState := cb.state
			cb.state = StateClosed
			cb.failures = 0
			cb.halfOpenCount = 0
			if previousState == StateHalfOpen && cb.logger != nil {
				cb.logger.Info("Circuit breaker closed after successful half-open requests")
			}
		} else if cb.state == StateClosed {
			// Reset failure count on success
			cb.failures = 0
		}
	}
}

// addErrorToHistory adds an error to the recent errors list
func (cb *CircuitBreaker) addErrorToHistory(err error) {
	record := FailureRecord{
		Error:     err.Error(),
		Timestamp: time.Now(),
	}

	// Add to the front of the list
	cb.recentErrors = append([]FailureRecord{record}, cb.recentErrors...)

	// Trim to max history size
	if len(cb.recentErrors) > cb.maxErrorHistory {
		cb.recentErrors = cb.recentErrors[:cb.maxErrorHistory]
	}
}

// logCircuitOpened logs detailed information when the circuit opens
func (cb *CircuitBreaker) logCircuitOpened() {
	fields := map[string]interface{}{
		"failures":     cb.failures,
		"maxFailures":  cb.maxFailures,
		"resetTimeout": cb.resetTimeout.String(),
	}

	// Add recent errors to the log
	if len(cb.recentErrors) > 0 {
		recentErrorMsgs := make([]string, 0, len(cb.recentErrors))
		for i, record := range cb.recentErrors {
			if i >= 5 { // Only show last 5 errors in log
				break
			}
			recentErrorMsgs = append(recentErrorMsgs, fmt.Sprintf("[%s] %s",
				record.Timestamp.Format("15:04:05"), record.Error))
		}
		fields["recentErrors"] = recentErrorMsgs
	}

	cb.logger.WithFields(fields).Error("Circuit breaker opened due to consecutive failures")
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

// GetRecentErrors returns a copy of recent error history
func (cb *CircuitBreaker) GetRecentErrors() []FailureRecord {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	// Return a copy to prevent external modification
	errors := make([]FailureRecord, len(cb.recentErrors))
	copy(errors, cb.recentErrors)
	return errors
}
