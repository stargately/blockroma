# Fix: RPC Error -32601 "method not found"

## Problem Summary

The indexer was failing with:
```json
{
  "level": "error",
  "msg": "Circuit breaker opened due to consecutive failures",
  "failures": 5,
  "maxFailures": 5,
  "recentErrors": [
    "[22:22:24] rpc error -32601: method not found",
    "[22:22:24] rpc error -32601: method not found",
    "[22:22:24] rpc error -32601: method not found",
    "[22:22:24] rpc error -32601: method not found",
    "[22:22:24] rpc error -32601: method not found"
  ],
  "resetTimeout": "10s"
}
```

Error code **-32601** means "Method not found" according to JSON-RPC 2.0 specification.

## Root Cause

The code was calling `getContractData`, which **does not exist** in the Stellar RPC API.

### Official Stellar RPC Methods (2024-2025)

According to https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods:

1. ✅ getEvents
2. ✅ getFeeStats
3. ✅ getHealth
4. ✅ getLatestLedger
5. ✅ **getLedgerEntries** ← Use this instead
6. ✅ getLedgers
7. ✅ getNetwork
8. ✅ getTransaction
9. ✅ getTransactions
10. ✅ getVersionInfo
11. ✅ sendTransaction
12. ✅ simulateTransaction

❌ **getContractData** does not exist (was never part of the official API)

## The Fix

### Changed Files

**1. pkg/client/rpc.go (lines 288-303)**

**BEFORE:**
```go
// ContractDataRequest parameters for getContractData
type ContractDataRequest struct {
	ContractID string `json:"contractId"`
	Key        string `json:"key"`
	Durability string `json:"durability"` // "persistent" or "temporary"
}

// ContractDataResponse response from getContractData
type ContractDataResponse struct {
	XDR                string `json:"xdr"`
	LastModifiedLedger uint32 `json:"lastModifiedLedgerSeq,omitempty"`
	LiveUntilLedgerSeq uint32 `json:"liveUntilLedgerSeq,omitempty"`
}

// GetContractData fetches contract data for a specific key
func (c *Client) GetContractData(ctx context.Context, contractID, key, durability string) (*ContractDataResponse, error) {
	req := ContractDataRequest{
		ContractID: contractID,
		Key:        key,
		Durability: durability,
	}

	var result ContractDataResponse
	if err := c.call(ctx, "getContractData", req, &result); err != nil {  // ❌ This method doesn't exist!
		return nil, err
	}

	return &result, nil
}
```

**AFTER:**
```go
// GetContractData fetches contract data for a specific key using getLedgerEntries
// This is a convenience wrapper around GetLedgerEntries for fetching contract data
func (c *Client) GetContractData(ctx context.Context, contractID, key, durability string) (*LedgerEntryResult, error) {
	// Use getLedgerEntries with the provided key
	// The key should already be a base64-encoded ledger key from BuildContractDataKey
	resp, err := c.GetLedgerEntries(ctx, []string{key})  // ✅ Use official API method
	if err != nil {
		return nil, err
	}

	if len(resp.Entries) == 0 {
		return nil, fmt.Errorf("contract data not found")
	}

	return &resp.Entries[0], nil
}
```

**2. pkg/client/rpc_test.go (lines 370-427)**

Updated test to:
- Expect `getLedgerEntries` method instead of `getContractData`
- Mock `GetLedgerEntriesResponse` instead of `ContractDataResponse`
- Verify the wrapper works correctly

### Key Changes

1. **Removed non-existent RPC method**: `getContractData` → `getLedgerEntries`
2. **Changed return type**: `*ContractDataResponse` → `*LedgerEntryResult`
3. **Simplified implementation**: Now wraps `GetLedgerEntries` instead of making its own RPC call
4. **No changes to callers**: The function signature remains compatible (same parameters)

### How It Works

The `GetContractData` method now:
1. Takes the same parameters (contractID, key, durability)
2. Calls the official `getLedgerEntries` API method with the provided key
3. Returns the first ledger entry result
4. Maintains backward compatibility with existing code in `pkg/poller/poller.go:432`

## Why This Happened

The `getContractData` method was likely:
- Part of an early/beta version of Stellar RPC
- Documentation that was outdated
- A proposed method that was never implemented
- Confusion with Horizon API methods

The official Stellar RPC only supports `getLedgerEntries` for fetching contract data by constructing the appropriate ledger key.

## Testing

### All Tests Pass ✅

```bash
$ go test ./pkg/client/...
PASS
ok  	github.com/blockroma/soroban-indexer/pkg/client	2.351s
```

All 12 client tests pass, including:
- ✅ TestClient_GetContractData (now uses getLedgerEntries)
- ✅ TestClient_GetLedgerEntries
- ✅ All other RPC method tests

### Build Success ✅

```bash
$ make build
Building indexer...
go build -o ./build/indexer ./cmd/indexer
Binary built: ./build/indexer
✓ Build successful
```

## Deployment

### To Deploy the Fix

```bash
# 1. Rebuild the indexer
cd indexer
make build

# 2. Or rebuild Docker image
cd deploy
docker compose build stellar-indexer

# 3. Restart the indexer
docker compose restart stellar-indexer

# 4. Monitor logs - should no longer see "method not found"
docker compose logs -f stellar-indexer
```

### Expected Behavior After Fix

**BEFORE (with bug):**
```json
{
  "level": "error",
  "msg": "Circuit breaker opened due to consecutive failures",
  "recentErrors": ["rpc error -32601: method not found", ...]
}
```

**AFTER (fixed):**
```json
{
  "level": "info",
  "msg": "Processing batch",
  "events": 50,
  "transactions": 12,
  "contracts": 3
}
```

The circuit breaker should remain closed and indexing should proceed normally.

## Related Issues

This fix complements:
- **CIRCUIT_BREAKER_LOGGING.md** - Enhanced logging that helped identify the specific error
- **FIX_ACCOUNT_ENTRIES.md** - Previous fix for account data indexing
- **DEBUGGING_CIRCUIT_BREAKER.md** - Troubleshooting guide for circuit breaker issues

## Verification Steps

1. ✅ All tests pass (12/12 in pkg/client)
2. ✅ Binary builds successfully
3. ✅ No compilation errors
4. ✅ Backward compatible with existing code
5. ✅ Uses only official Stellar RPC methods

## Summary

The indexer was calling a non-existent RPC method (`getContractData`), causing the RPC server to return "method not found" errors. The circuit breaker logging enhancement (from CIRCUIT_BREAKER_LOGGING.md) made this immediately visible. The fix replaces the non-existent method call with the official `getLedgerEntries` API, which is the correct way to fetch contract data from Stellar RPC.

**Key Lesson**: Always verify RPC methods against the official documentation at https://developers.stellar.org/docs/data/apis/rpc/api-reference/methods
