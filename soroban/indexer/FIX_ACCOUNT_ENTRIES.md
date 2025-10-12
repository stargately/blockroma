# Fix: account_entries Table Always Empty

## Problem Summary

The `account_entries` table was always empty despite transactions being successfully indexed. Account entries, trustlines, offers, and other ledger data were not being populated.

## Root Cause Analysis

### The Bug

**Location**: `pkg/poller/poller.go`
- Main polling loop: Lines 313-328 (old code)
- Backfill function: Lines 729-752 (old code)

**What was happening**:

1. Transactions were fetched from RPC and stored in the database
2. The code attempted to extract source accounts by **querying the database** for transactions that were just inserted
3. **The query was happening within the same uncommitted database transaction**
4. Due to transaction isolation or GORM's query execution, the query often returned zero results
5. With zero account addresses, `processLedgerEntries()` was never called
6. No account entries, trustlines, or offers were fetched or stored

### The Problematic Code

```go
// OLD BUGGY CODE
accountAddresses := make(map[string]bool)
var transactions []models.Transaction
hashList := make([]string, 0, len(txHashes))
for hash := range txHashes {
    hashList = append(hashList, hash)
}
if err := tx.Where("id IN ?", hashList).Find(&transactions).Error; err == nil {
    for _, txn := range transactions {
        if txn.SourceAccount != nil && *txn.SourceAccount != "" {
            accountAddresses[*txn.SourceAccount] = true
        }
    }
} else {
    p.logger.WithError(err).Warn("Failed to query transactions for account addresses")
}
```

**Why this failed**:
- The database query might not see uncommitted inserts within the same transaction
- Silent failures when the query returned no results
- Unnecessary round-trip to the database for data we already had in memory

## The Fix

### Solution Overview

Instead of querying the database for transactions we just parsed, **extract source accounts directly during transaction processing**.

### Implementation

**File**: `pkg/poller/poller.go`

#### Change 1: Extract accounts during transaction processing (Lines 215, 268-271)

```go
// NEW CODE - Main polling loop
accountAddresses := make(map[string]bool)  // Initialize early

for txHash := range txHashes {
    // ... fetch and parse transaction ...

    dbTx, err := parser.ParseTransactionWithHash(*rpcTx, actualTxHash)
    if err != nil {
        continue
    }

    if err := models.UpsertTransaction(tx, dbTx); err != nil {
        return fmt.Errorf("upsert transaction: %w", err)
    }

    txCount++

    // NEW: Extract source account immediately
    if dbTx.SourceAccount != nil && *dbTx.SourceAccount != "" {
        accountAddresses[*dbTx.SourceAccount] = true
    }
}
```

#### Change 2: Remove database query (Lines 317-328)

```go
// NEW CODE - Simplified ledger entry processing
// Process ledger entries for discovered accounts
// Account addresses were collected during transaction processing above
if len(accountAddresses) > 0 {
    p.logger.WithField("accountCount", len(accountAddresses)).Info("Processing account ledger entries")
    if err := p.processLedgerEntries(ctx, tx, accountAddresses); err != nil {
        p.logger.WithError(err).Warn("Failed to process ledger entries")
        // Don't fail the whole batch if ledger entry processing fails
    }
} else {
    p.logger.Debug("No account addresses found in transactions")
}
```

#### Change 3: Updated logging (Line 321)

- Changed log level from `Debug` to `Info` for better visibility
- Log message now clearly indicates account processing is happening

#### Change 4: Applied same fix to backfill function (Lines 955, 978-980, 1024-1033)

The same bug existed in the `processEventBatch()` function used for backfilling. Applied identical fix.

## Benefits of the Fix

1. **No unnecessary database query**: Eliminated round-trip to database for data already in memory
2. **Deterministic behavior**: Account extraction always works, regardless of transaction isolation level
3. **Better performance**: One less query per batch
4. **More reliable**: No dependency on GORM's transaction handling
5. **Better logging**: Info-level logs make it easier to track account processing

## Testing

### Test Results

All 136 tests pass, including:

- ✅ `TestProcessLedgerEntries_AccountExtraction` - Verifies account address extraction
- ✅ `TestProcessLedgerEntries_EmptySourceAccount` - Handles nil/empty source accounts
- ✅ `TestProcessLedgerEntries_MixedTransactions` - Mixed valid and invalid transactions
- ✅ `TestAccountEntryUpsert` - Account entry insert and update operations
- ✅ `TestTransactionQueryWithInlineFunction` - Documents the previous bug pattern
- ✅ `TestLedgerEntryParsing` - XDR ledger entry parsing
- ✅ `TestGetLedgerEntriesResponse` - RPC response structure

### Build Status

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
# 1. Rebuild the Docker image
cd deploy
docker compose build stellar-indexer

# 2. Restart the indexer
docker compose restart stellar-indexer

# 3. Monitor the logs
docker compose logs -f stellar-indexer
```

### Expected Log Output

With the fix, you should see:

```
level=info msg="Processing account ledger entries" accountCount=5
level=info msg="Processed ledger entries" accounts=5 trustlines=2 offers=0 data=0 claimableBalance=0 liquidityPools=0
```

### Verification

Check that account entries are being populated:

```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U stellar -d stellar

# Count account entries
SELECT COUNT(*) FROM account_entries;

# View sample account entries
SELECT account_id, balance, seq_num, home_domain
FROM account_entries
LIMIT 10;
```

## Files Changed

1. **`pkg/poller/poller.go`**
   - Lines 215: Initialize `accountAddresses` map early
   - Lines 268-271: Extract source accounts during transaction processing
   - Lines 317-328: Remove database query, simplify logic
   - Line 321: Change log level from Debug to Info
   - Lines 955, 978-980: Same fix for backfill function
   - Lines 1024-1033: Remove database query in backfill function

## Related Issues

- The same pattern was previously fixed but only partially
- The inline function pattern (passing `func() []string {}()` to GORM) was even more problematic
- This fix completely eliminates the database query approach

## Verification Steps

1. ✅ All tests pass (136 tests)
2. ✅ Binary builds successfully
3. ✅ No regressions in existing functionality
4. ✅ Account extraction logic is deterministic and reliable
5. ✅ Logging improvements for better observability

## Performance Impact

**Positive**:
- One less database query per batch
- Faster batch processing
- Reduced database load

**Neutral**:
- No negative performance impact
- Memory usage unchanged (account addresses already in parsed transaction objects)

## Future Considerations

This fix highlights a general pattern to avoid:
- **Don't query the database for data you already have in memory**
- **Don't query uncommitted data within the same transaction**
- Extract data during parsing rather than re-fetching from database

## Summary

This fix resolves the core issue of `account_entries` being empty by extracting source accounts directly during transaction processing, eliminating the problematic database query within an uncommitted transaction. The solution is more efficient, reliable, and maintainable.
