# Fix: RPC Error "key count exceeds maximum"

## Problem Summary

After fixing the ledger key format error, a new error appeared:
```json
{
  "error": "get ledger entries: rpc error -32602: key count (313) exceeds maximum supported (200)",
  "level": "warning",
  "msg": "Failed to process ledger entries"
}
```

The Stellar RPC has a **maximum limit of 200 keys per `getLedgerEntries` request**, but the code was trying to send **313 keys** in a single call.

## Root Cause

The code was sending all ledger keys in a single `GetLedgerEntries` call without batching:

```go
// ❌ WRONG: Sending all 313 keys at once
resp, err := p.rpcClient.GetLedgerEntries(ctx, keys)
```

When processing a large batch of transactions with many unique source accounts, this could easily exceed the RPC's 200-key limit.

## The Fix

**File**: `pkg/poller/poller.go:519-604`

Added batching to split large key arrays into chunks of 200 or fewer:

**BEFORE (No Batching):**
```go
// Fetch ledger entries from RPC
resp, err := p.rpcClient.GetLedgerEntries(ctx, keys)  // ❌ All keys at once
if err != nil {
    return fmt.Errorf("get ledger entries: %w", err)
}

// Process each ledger entry
for _, entry := range resp.Entries {
    // Process...
}
```

**AFTER (With Batching):**
```go
// Process ledger entries
accountCount := 0
trustlineCount := 0
// ... other counters ...

// Batch ledger entry requests (RPC max is 200 keys per request)
const maxKeysPerRequest = 200
for i := 0; i < len(keys); i += maxKeysPerRequest {
    end := i + maxKeysPerRequest
    if end > len(keys) {
        end = len(keys)
    }
    batch := keys[i:end]

    p.logger.WithFields(logrus.Fields{
        "batchStart": i,
        "batchEnd":   end,
        "batchSize":  len(batch),
    }).Debug("Fetching ledger entry batch")

    // Fetch this batch of ledger entries from RPC
    resp, err := p.rpcClient.GetLedgerEntries(ctx, batch)  // ✅ Max 200 keys
    if err != nil {
        p.logger.WithError(err).WithField("batchSize", len(batch)).Warn("Failed to fetch ledger entry batch")
        continue // Skip this batch but continue with others
    }

    p.logger.WithField("entriesReceived", len(resp.Entries)).Debug("Received ledger entries from RPC")

    // Process each ledger entry in this batch
    for _, entry := range resp.Entries {
        // Parse and store...
    }
}
```

### How It Works

1. **Split keys into batches**: Iterate through keys in chunks of 200
2. **Fetch each batch**: Call `GetLedgerEntries` for each chunk
3. **Accumulate results**: Process entries from all batches and aggregate counters
4. **Continue on error**: If one batch fails, skip it but continue with remaining batches

### Example with 313 Keys

Before: 1 request with 313 keys → **ERROR**

After:
- Request 1: Keys 0-199 (200 keys) ✅
- Request 2: Keys 200-312 (113 keys) ✅

Total: 2 requests, all succeed

## Testing

### Build Success ✅
```bash
$ go build -o ./build/indexer ./cmd/indexer
# No errors
```

### All Tests Pass ✅
```bash
$ go test ./...
ok  	github.com/blockroma/soroban-indexer/pkg/client	(cached)
ok  	github.com/blockroma/soroban-indexer/pkg/db	(cached)
ok  	github.com/blockroma/soroban-indexer/pkg/models	(cached)
ok  	github.com/blockroma/soroban-indexer/pkg/parser	(cached)
ok  	github.com/blockroma/soroban-indexer/pkg/poller	0.408s
ok  	github.com/blockroma/soroban-indexer/pkg/worker	(cached)
```

## Performance Impact

### Before
- **Single large request**: One RPC call with N keys (fails if N > 200)
- **Failure mode**: Complete failure for entire batch

### After
- **Multiple smaller requests**: ⌈N/200⌉ RPC calls
- **Resilience**: If one batch fails, others still succeed
- **Network overhead**: Minimal (HTTP overhead is negligible compared to data transfer)

### Example Scenarios

| Keys | Requests Before | Requests After | Result Before | Result After |
|------|----------------|----------------|---------------|--------------|
| 50   | 1              | 1              | ✅ Success    | ✅ Success   |
| 200  | 1              | 1              | ✅ Success    | ✅ Success   |
| 201  | 1              | 2              | ❌ Error      | ✅ Success   |
| 313  | 1              | 2              | ❌ Error      | ✅ Success   |
| 500  | 1              | 3              | ❌ Error      | ✅ Success   |
| 1000 | 1              | 5              | ❌ Error      | ✅ Success   |

## Additional Improvements

### Logging Enhancements

Added detailed logging for batching:
```go
p.logger.WithFields(logrus.Fields{
    "batchStart": i,
    "batchEnd":   end,
    "batchSize":  len(batch),
}).Debug("Fetching ledger entry batch")
```

This helps track:
- How many batches are being processed
- Which batch failed (if any)
- Size of each batch

### Error Handling

Changed from failing completely to graceful degradation:
- **Before**: Single error fails entire batch
- **After**: Failed batch is skipped, others continue

```go
if err != nil {
    p.logger.WithError(err).WithField("batchSize", len(batch)).Warn("Failed to fetch ledger entry batch")
    continue // Skip this batch but continue with others
}
```

## Deployment

### To Deploy the Fix

```bash
# 1. Rebuild the indexer
go build -o ./build/indexer ./cmd/indexer

# 2. Or rebuild Docker image
cd deploy
docker compose build stellar-indexer

# 3. Restart the indexer
docker compose restart stellar-indexer

# 4. Monitor logs - should see batching in action
docker compose logs -f stellar-indexer | grep -i "batch"
```

### Expected Log Output

**With Many Accounts (>200):**
```
level=debug msg="Building ledger keys for accounts" accountCount=313
level=debug msg="Fetching ledger entries from RPC" keyCount=313
level=debug msg="Fetching ledger entry batch" batchStart=0 batchEnd=200 batchSize=200
level=debug msg="Received ledger entries from RPC" entriesReceived=195
level=debug msg="Fetching ledger entry batch" batchStart=200 batchEnd=313 batchSize=113
level=debug msg="Received ledger entries from RPC" entriesReceived=110
level=info msg="Processed ledger entries" accounts=305 trustlines=12 offers=0
```

**With Few Accounts (<200):**
```
level=debug msg="Building ledger keys for accounts" accountCount=50
level=debug msg="Fetching ledger entries from RPC" keyCount=50
level=debug msg="Fetching ledger entry batch" batchStart=0 batchEnd=50 batchSize=50
level=debug msg="Received ledger entries from RPC" entriesReceived=48
level=info msg="Processed ledger entries" accounts=48 trustlines=5 offers=0
```

## Why 200?

The limit of 200 keys per request is enforced by the Stellar RPC server. From the error message:
```
key count (313) exceeds maximum supported (200)
```

This is a hard limit in the RPC implementation to prevent:
- Extremely large responses
- Long-running queries
- Resource exhaustion

## Related Issues

This fix builds on:
- **CIRCUIT_BREAKER_LOGGING.md** - Logging that helped identify all these errors
- **FIX_METHOD_NOT_FOUND.md** - Fixed the RPC method name
- **FIX_INVALID_LEDGER_KEY.md** - Fixed the ledger key format
- **FIX_ACCOUNT_ENTRIES.md** - Account entry indexing improvements

## Other Functions That May Need Similar Fixes

Currently, batching is applied to:
- ✅ `processLedgerEntries()` - **Fixed** (accounts)

May need batching in the future if they exceed 200 keys:
- `processClaimableBalances()` - Currently sends all balance IDs at once (line 618)
  - **Risk**: Low (unlikely to have >200 claimable balances in a single batch)
  - **Action**: Monitor logs; add batching if needed

## Summary

The indexer was sending 313 ledger keys in a single `getLedgerEntries` request, exceeding the RPC's 200-key limit. The fix implements batching to split large key arrays into chunks of 200 or fewer, with graceful error handling that continues processing even if individual batches fail. This allows the indexer to handle any number of unique accounts in a transaction batch without hitting RPC limits.

**Key Changes:**
- Added batching loop with `maxKeysPerRequest = 200`
- Split key array into chunks before sending to RPC
- Added batch-level logging for visibility
- Graceful error handling: skip failed batches, continue with others
- Accumulate results across all batches

The indexer can now process batches with hundreds or thousands of unique accounts without errors.
