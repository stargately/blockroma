# Circuit Breaker Logging Enhancement

## Overview

Enhanced the circuit breaker with detailed failure logging to help diagnose why the circuit breaker is opening. This provides visibility into the specific errors that caused the circuit breaker to trip.

## What Changed

### 1. Added Logger Interface (`pkg/worker/circuit_breaker.go:28-35`)

```go
type Logger interface {
    WithError(err error) Logger
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    Warn(msg string)
    Error(msg string)
    Info(msg string)
}
```

This interface allows the circuit breaker to log events without coupling to a specific logging library.

### 2. Added Failure Tracking (`pkg/worker/circuit_breaker.go:38-41`)

```go
type FailureRecord struct {
    Error     string
    Timestamp time.Time
}
```

The circuit breaker now tracks the last 10 failures with timestamps for debugging.

### 3. Enhanced CircuitBreaker Struct (`pkg/worker/circuit_breaker.go:58-60`)

Added three new fields:
- `logger Logger` - Optional logger for recording events
- `recentErrors []FailureRecord` - Circular buffer of recent failures
- `maxErrorHistory int` - Maximum number of errors to track (default: 10)

### 4. Detailed Logging When Circuit Opens

When the circuit breaker opens, it now logs:
- Number of failures that triggered the opening
- Maximum failure threshold
- Reset timeout duration
- Last 5 error messages with timestamps

**Example log output:**
```json
{
  "level": "error",
  "msg": "Circuit breaker opened due to consecutive failures",
  "failures": 5,
  "maxFailures": 5,
  "resetTimeout": "30s",
  "recentErrors": [
    "[23:15:42] http request: dial tcp: connection refused",
    "[23:15:41] http request: dial tcp: connection refused",
    "[23:15:40] http request: context deadline exceeded",
    "[23:15:39] http request: dial tcp: connection refused",
    "[23:15:38] http status: 503"
  ]
}
```

### 5. State Transition Logging

The circuit breaker now logs when:
- Circuit opens (with recent errors)
- Circuit reopens from half-open state after failure
- Circuit closes after successful half-open requests

### 6. Logrus Adapter (`pkg/worker/logger_adapter.go`)

Created an adapter to bridge logrus (used by the indexer) with the circuit breaker's Logger interface:

```go
type LogrusAdapter struct {
    entry *logrus.Entry
}
```

This allows seamless integration with the existing logging infrastructure.

### 7. Client Integration (`pkg/client/rpc.go:31-36`)

Added `SetLogger()` method to the RPC client:

```go
func (c *Client) SetLogger(logger worker.Logger) {
    if c.circuitBreaker != nil {
        c.circuitBreaker.SetLogger(logger)
    }
}
```

### 8. Main Integration (`cmd/indexer/main.go:82`)

The logger is now set on the RPC client during initialization:

```go
rpcClient := client.NewClient(rpcURL)
rpcClient.SetLogger(worker.NewLogrusAdapter(logger))
```

## New Methods

### `SetLogger(logger Logger)`
Sets the logger for the circuit breaker. Can be called after construction.

### `GetRecentErrors() []FailureRecord`
Returns a copy of recent error history for debugging. Thread-safe.

### `addErrorToHistory(err error)` (private)
Maintains a circular buffer of recent errors, keeping the most recent 10.

### `logCircuitOpened()` (private)
Logs detailed information when the circuit transitions to open state.

## Usage

### In Application Code

```go
// Create RPC client
rpcClient := client.NewClient(rpcURL)

// Set logger (optional but recommended)
rpcClient.SetLogger(worker.NewLogrusAdapter(logger))

// Use client normally
err := rpcClient.Health(ctx)
if err != nil {
    // If circuit is open, you'll now see detailed logs explaining why
    log.Println("RPC error:", err)
}
```

### Retrieving Error History Programmatically

```go
// Get recent errors for debugging
errors := circuitBreaker.GetRecentErrors()
for _, record := range errors {
    fmt.Printf("[%s] %s\n", record.Timestamp.Format(time.RFC3339), record.Error)
}
```

## Benefits

1. **Root Cause Analysis**: Immediately see what errors caused the circuit to open
2. **Debugging**: No need to guess why the circuit breaker tripped
3. **Monitoring**: Structured logs can be ingested by log aggregation tools
4. **Trend Analysis**: Error patterns become visible (e.g., repeated timeouts vs connection refused)
5. **Minimal Overhead**: Logging is optional and only tracks last 10 errors

## Log Levels

- **Error**: When circuit opens due to failures
- **Warn**: When circuit reopens from half-open state after failure
- **Info**: When circuit closes after successful recovery

## Example Debugging Session

**Before this enhancement:**
```
level=error msg="Poll failed" error="circuit breaker is open"
```
❌ No information about why the circuit opened

**After this enhancement:**
```json
{
  "level": "error",
  "msg": "Circuit breaker opened due to consecutive failures",
  "failures": 5,
  "maxFailures": 5,
  "resetTimeout": "30s",
  "recentErrors": [
    "[15:42:13] http request: Post \"http://stellar-rpc:8000\": dial tcp 172.20.0.4:8000: connect: connection refused",
    "[15:42:12] http request: Post \"http://stellar-rpc:8000\": dial tcp 172.20.0.4:8000: connect: connection refused",
    "[15:42:11] http request: Post \"http://stellar-rpc:8000\": dial tcp 172.20.0.4:8000: connect: connection refused",
    "[15:42:10] http request: Post \"http://stellar-rpc:8000\": dial tcp 172.20.0.4:8000: connect: connection refused",
    "[15:42:09] http request: Post \"http://stellar-rpc:8000\": dial tcp 172.20.0.4:8000: connect: connection refused"
  ]
}
```
✅ Clear diagnosis: RPC service at 172.20.0.4:8000 is refusing connections

## Configuration

No configuration changes needed. The enhancement is backward compatible:
- Existing code works without modification
- Logger is optional - if not set, circuit breaker works silently as before
- Error history tracking happens automatically (minimal memory: ~1KB)

## Testing

All 136 existing tests pass, including circuit breaker tests that verify:
- Circuit opens after N failures ✅
- Circuit transitions to half-open after timeout ✅
- Circuit closes after successful requests ✅
- Error tracking doesn't affect circuit breaker logic ✅

## Files Modified

1. **pkg/worker/circuit_breaker.go**
   - Added Logger interface
   - Added FailureRecord struct
   - Enhanced CircuitBreaker with logging fields
   - Updated NewCircuitBreaker to initialize error tracking
   - Added SetLogger method
   - Enhanced recordResult to track and log errors
   - Added addErrorToHistory method
   - Added logCircuitOpened method
   - Added GetRecentErrors method

2. **pkg/worker/logger_adapter.go** (new file)
   - Created LogrusAdapter to bridge logrus with Logger interface
   - Implements all Logger interface methods

3. **pkg/client/rpc.go**
   - Added SetLogger method to Client

4. **cmd/indexer/main.go**
   - Import worker package
   - Set logger on RPC client during initialization

## Deployment

### To Deploy

```bash
# Rebuild the indexer
cd indexer
make build

# Or rebuild Docker image
cd deploy
docker compose build stellar-indexer

# Restart the indexer
docker compose restart stellar-indexer

# Monitor logs with detailed circuit breaker information
docker compose logs -f stellar-indexer
```

### Verification

Check that you see detailed logs when circuit breaker opens:

```bash
# Look for circuit breaker logs
docker compose logs stellar-indexer | grep -i "circuit"

# Expected output when circuit opens:
# level=error msg="Circuit breaker opened due to consecutive failures" failures=5 maxFailures=5 resetTimeout="30s" recentErrors=[...]

# Expected output when circuit closes:
# level=info msg="Circuit breaker closed after successful half-open requests"
```

## Related Documentation

- See `DEBUGGING_CIRCUIT_BREAKER.md` for troubleshooting guide
- See `pkg/worker/circuit_breaker.go` for implementation details
- See `FIX_ACCOUNT_ENTRIES.md` for previous indexer fix

## Future Enhancements

Possible improvements for the future:
1. Add metrics export (Prometheus) for circuit breaker state changes
2. Add configurable error history size
3. Add HTTP endpoint to retrieve circuit breaker status and error history
4. Add per-error-type statistics (count connection refused, timeouts, etc.)
5. Add adaptive failure thresholds based on error patterns
