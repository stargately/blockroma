# Soroban Indexer - Table Population Status

This document analyzes the current state of all database tables in the Soroban indexer and identifies which ones are being properly populated.

## Table Overview

| Table | Status | Discovery Method | Notes |
|-------|--------|------------------|-------|
| `events` | ✅ **Populated** | RPC `getEvents()` | Primary data source, polls every 1 second |
| `transactions` | ✅ **Populated** | RPC `getTransaction()` | Fetched for each unique tx_hash from events |
| `operations` | ✅ **Populated** | Parsed from transactions | All 25+ Stellar operation types (NEW: 2025-10-11) |
| `token_operations` | ✅ **Populated** | Parsed from events | Extracts transfer/mint/burn from SAC token events |
| `contract_data_entries` | ✅ **Populated** | RPC `getContractData()` | Proactively fetched for discovered contracts (IMPROVED: 2025-10-11) |
| `contract_code` | ✅ **Populated** | Parsed from operations | UploadContractWasm WASM bytecode tracking (NEW: 2025-10-11) |
| `token_metadata` | ✅ **Populated** | Proactively fetched | Standard SAC metadata keys queried (IMPROVED: 2025-10-11) |
| `token_balances` | ✅ **Populated** | Proactively fetched | Balance keys queried from token ops (IMPROVED: 2025-10-11) |
| `account_entries` | ✅ **Populated** | RPC `getLedgerEntries()` | Fetched for source accounts from transactions |
| `trust_line_entries` | ⚠️ **May be Empty** | RPC `getLedgerEntries()` | Only populated if accounts have trustlines |
| `offer_entries` | ⚠️ **May be Empty** | RPC `getLedgerEntries()` | Only populated if accounts have offers |
| `data_entries` | ⚠️ **May be Empty** | RPC `getLedgerEntries()` | Only populated if accounts have data entries |
| `claimable_balance_entries` | ✅ **Populated** | Parsed from operations | IDs computed from CreateClaimableBalance ops (IMPROVED: 2025-10-11) |
| `liquidity_pool_entries` | ⚠️ **May be Empty** | RPC `getLedgerEntries()` | Only if accounts participate in pools |
| `cursor` | ✅ **Populated** | Internal state tracking | Tracks last processed ledger |

## Detailed Analysis

### ✅ Fully Operational Tables

#### 1. Events (`events`)
- **Population Method**: RPC `getEvents()` called every 1 second
- **Data Source**: Direct from Stellar RPC
- **Coverage**: All contract events on the network
- **Location**: `pkg/poller/poller.go:122-133`

#### 2. Transactions (`transactions`) **IMPROVED: 2025-10-11**
- **Population Method**: RPC `getTransaction()` for each unique tx_hash
- **Data Source**: Extracted from events, then fetched from RPC
- **Coverage**: All transactions that emit events
- **Location**: `pkg/poller/poller.go:169-209`
- **Improvement**: Now parses and stores memo, preconditions, and signatures from transaction envelopes
- **Populated Fields**:
  - Basic: source account, fee, sequence, status, ledger
  - **NEW**: Memo (text/id/hash/return types)
  - **NEW**: Preconditions (time bounds, ledger bounds, min seq, extra signers)
  - **NEW**: Signatures (signature hints and values)
- **Note**: Uses event's tx_hash to avoid RPC hash mismatches

#### 3. Operations (`operations`) **NEW: 2025-10-11**
- **Population Method**: Parsed from transaction envelope XDR
- **Data Source**: All operations within indexed transactions
- **Coverage**: All 25+ Stellar operation types including:
  - Payment, CreateAccount, PathPayment
  - ManageSellOffer, ManageBuyOffer
  - CreateClaimableBalance, ClaimClaimableBalance
  - InvokeHostFunction, ExtendFootprintTtl, RestoreFootprint
  - And 15+ more operation types
- **Location**: `pkg/poller/poller.go`, `pkg/parser/operations.go`
- **Special Features**:
  - Computes claimable balance IDs for CreateClaimableBalance operations
  - Stores type-specific details as JSONB
  - Tracks operation-level source accounts

#### 4. Token Operations (`token_operations`)
- **Population Method**: Parsed from event topics/values
- **Data Source**: SAC (Stellar Asset Contract) standard events
- **Coverage**: Transfer, mint, burn operations
- **Location**: `pkg/poller/poller.go:147-161`, `pkg/parser/token_operations.go`

#### 5. Contract Data Entries (`contract_data_entries`) **IMPROVED: 2025-10-11**
- **Population Method**: RPC `getContractData()` for discovered contracts
- **Data Source**: Proactive fetching of standard SAC storage keys
- **Coverage**: Metadata and balance keys for all token contracts
- **Location**: `pkg/poller/poller.go`, `pkg/parser/ledger_entries.go`
- **Improvement**: Now proactively queries contract storage instead of passive discovery

#### 6. Token Metadata (`token_metadata`) **IMPROVED: 2025-10-11**
- **Population Method**: Parsed from proactively fetched contract data
- **Data Source**: Standard SAC metadata keys (name, symbol, decimals)
- **Coverage**: All discovered token contracts
- **Location**: `pkg/poller/poller.go`, `pkg/parser/token_operations.go`
- **Improvement**: Now actively fetches metadata for all contracts

#### 7. Token Balances (`token_balances`) **IMPROVED: 2025-10-11**
- **Population Method**: Parsed from proactively fetched contract data
- **Data Source**: Balance keys for addresses in token operations
- **Coverage**: All addresses participating in token transfers
- **Location**: `pkg/poller/poller.go`, `pkg/parser/token_operations.go`
- **Improvement**: Now tracks and queries balances for all active addresses

#### 8. Account Entries (`account_entries`)
- **Population Method**: RPC `getLedgerEntries()` with account ledger keys
- **Data Source**: Source accounts from indexed transactions
- **Coverage**: Accounts that submit transactions
- **Location**: `pkg/poller/poller.go:238-243`, `pkg/parser/ledger_entries.go:338-368`

#### 9. Contract Code (`contract_code`) **NEW: 2025-10-11**
- **Population Method**: Extracted from UploadContractWasm operations
- **Data Source**: Parsed from InvokeHostFunction transaction operations
- **Coverage**: All deployed smart contract WASM bytecode
- **Location**: `pkg/poller/poller.go`, `pkg/parser/contract_code.go`
- **Populated Fields**:
  - `hash`: SHA256 hash of WASM bytecode (primary key for deduplication)
  - `wasm`: Full WASM bytecode stored as PostgreSQL bytea
  - `deployed_at`: First deployment timestamp
  - `ledger`: Ledger number where first deployed
  - `tx_hash`: Transaction hash that deployed the code
  - `size_bytes`: Size of WASM bytecode in bytes
- **Features**:
  - Idempotent upserts based on hash (prevents duplicate storage)
  - Tracks first deployment only (immutable after first insert)
  - Integrated into both live polling and backfill modes
  - Comprehensive test coverage (23 tests)

#### 10. Claimable Balance Entries (`claimable_balance_entries`) **IMPROVED: 2025-10-11**
- **Population Method**: Parsed from CreateClaimableBalance operations
- **Data Source**: Operations table with computed balance IDs
- **Coverage**: All claimable balances created on-chain
- **Location**: `pkg/parser/operations.go`, `pkg/parser/ledger_entries.go`
- **Improvement**: Now computes balance IDs from CreateClaimableBalance operations

#### 11. Cursor (`cursor`)
- **Population Method**: Updated after each polling cycle
- **Data Source**: Internal state tracking
- **Coverage**: Single row tracking last processed ledger
- **Location**: `pkg/poller/poller.go:253-256`

### ⚠️ Partially Populated Tables

#### 6. Contract Data Entries (`contract_data_entries`)
- **Issue**: Only populated if RPC explicitly returns contract data
- **Current Discovery**: Contracts discovered via events
- **Missing**: Direct querying of contract storage
- **Location**: `pkg/poller/poller.go:211-216`
- **Recommendation**: Add RPC `getContractData()` calls for discovered contracts

#### 7. Token Metadata (`token_metadata`)
- **Issue**: Depends on contract_data_entries containing metadata keys
- **Current Parsing**: Looks for "METADATA" key in contract storage
- **Missing**: Proactive fetching of known metadata keys
- **Location**: `pkg/poller/poller.go:276-282`, `pkg/parser/token_operations.go`
- **Recommendation**: Query standard SAC metadata keys (name, symbol, decimals)

#### 8. Token Balances (`token_balances`)
- **Issue**: Depends on contract_data_entries containing balance keys
- **Current Parsing**: Looks for "Balance" keys in contract storage
- **Missing**: Proactive querying of user balances
- **Location**: `pkg/poller/poller.go:284-291`, `pkg/parser/token_operations.go`
- **Recommendation**: Track addresses from token operations and query balances

#### 9. Trust Line Entries (`trust_line_entries`)
- **Issue**: Only populated if accounts have trustlines
- **Current Discovery**: Returned in account ledger entries
- **Coverage**: Classic Stellar assets only (not Soroban tokens)
- **Location**: `pkg/parser/ledger_entries.go:159-194`
- **Note**: May be legitimately empty in Soroban-only networks

#### 10. Offer Entries (`offer_entries`)
- **Issue**: Only populated if accounts have open offers
- **Current Discovery**: Returned in account ledger entries
- **Coverage**: Classic Stellar DEX offers only
- **Location**: `pkg/parser/ledger_entries.go:196-234`
- **Note**: May be legitimately empty if no DEX activity

#### 11. Data Entries (`data_entries`)
- **Issue**: Only populated if accounts have data entries
- **Current Discovery**: Returned in account ledger entries
- **Coverage**: Classic Stellar account data
- **Location**: `pkg/parser/ledger_entries.go:236-251`
- **Note**: May be legitimately empty

#### 12. Claimable Balance Entries (`claimable_balance_entries`)
- **Issue**: Only populated when ClaimClaimableBalance operations occur
- **Current Discovery**: Extracted from transaction operations
- **Missing**: CreateClaimableBalance operations (complex ID computation)
- **Location**: `pkg/poller/poller.go:245-251`, `pkg/parser/ledger_entries.go:407-456`
- **Recommendation**: Add CreateClaimableBalance ID computation

#### 13. Liquidity Pool Entries (`liquidity_pool_entries`)
- **Issue**: Only populated if accounts participate in pools
- **Current Discovery**: Returned in account ledger entries
- **Coverage**: Classic Stellar AMM pools only
- **Location**: `pkg/parser/ledger_entries.go:280-319`
- **Note**: May be legitimately empty in Soroban-only networks

## Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         RPC Polling (1s)                        │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                    ┌───────────▼──────────┐
                    │  getEvents()         │
                    │  └─> events table    │
                    └───────────┬──────────┘
                                │
                ┌───────────────┼───────────────┐
                │               │               │
    ┌───────────▼──────────┐   │   ┌──────────▼───────────┐
    │  getTransaction()    │   │   │  Extract Contract IDs│
    │  └─> transactions    │   │   │  └─> (for later)     │
    └───────────┬──────────┘   │   └──────────┬───────────┘
                │               │               │
    ┌───────────▼──────────┐   │   ┌──────────▼───────────┐
    │ Extract Balance IDs  │   │   │ Parse Token Ops      │
    │ (from operations)    │   │   │ └─> token_operations │
    └───────────┬──────────┘   │   └──────────────────────┘
                │               │
    ┌───────────▼──────────────▼───────────┐
    │  Extract Source Accounts             │
    └───────────┬──────────────────────────┘
                │
    ┌───────────▼──────────────────────────┐
    │  getLedgerEntries()                  │
    │  ├─> account_entries                 │
    │  ├─> trust_line_entries              │
    │  ├─> offer_entries                   │
    │  ├─> data_entries                    │
    │  ├─> claimable_balance_entries       │
    │  └─> liquidity_pool_entries          │
    └──────────────────────────────────────┘
```

## Gaps and Limitations

### ✅ 1. Contract Storage Not Directly Indexed **RESOLVED: 2025-10-11**
~~**Problem**: Contract data entries are only populated if RPC returns them, but we don't proactively query contract storage.~~

**Solution Implemented**:
- ✅ Added RPC `getContractData()` method to client
- ✅ Proactively queries standard SAC metadata keys
- ✅ Tracks addresses from token operations and queries balances
- ✅ All token metadata and balances now fully populated

### ✅ 2. Create Claimable Balance Operations Not Tracked **RESOLVED: 2025-10-11**
~~**Problem**: We only track ClaimClaimableBalance, not CreateClaimableBalance.~~

**Solution Implemented**:
- ✅ Full operation parsing extracts CreateClaimableBalance operations
- ✅ Implemented balance ID computation using Stellar's formula
- ✅ Balance IDs computed from: sha256(networkHash + accountID + seqNum + opIndex + type)
- ✅ All claimable balances now tracked from creation

### ✅ 3. Transaction Operations Not Fully Parsed **RESOLVED: 2025-10-11**
~~**Problem**: We extract minimal data from transactions (source, fee, sequence).~~

**Solution Implemented**:
- ✅ Created operations table with comprehensive operation details
- ✅ Parse all 25+ Stellar operation types
- ✅ Store type-specific details as JSONB
- ✅ Track operation-level source accounts
- ✅ Full operation history now available

### ✅ 4. Transaction Memo and Preconditions Not Parsed **RESOLVED: 2025-10-11**
~~**Problem**: Transaction memo, preconditions, and signatures were not being extracted from envelopes.~~

**Solution Implemented**:
- ✅ Added memo parsing for all 5 memo types (none, text, id, hash, return)
- ✅ Added precondition parsing for V0 (time bounds) and V1/V2 (full preconditions)
- ✅ Added signature parsing from transaction envelopes
- ✅ Support for V0, V1, and FeeBump transaction envelope types
- ✅ Comprehensive test coverage for all parsing functions
- ✅ All transaction metadata now stored as JSONB in database

### 5. Historical Data Bootstrap **RESOLVED: 2025-10-11**
~~**Problem**: Indexer starts from current ledger on first run.~~

**Solution Implemented**:
- ✅ Added backfill mode with CLI flags: --start-ledger, --end-ledger, --batch-size, --rate-limit
- ✅ Sequential processing of historical ledgers with progress tracking
- ✅ Automatic resume from last processed ledger
- ✅ Rate limiting to prevent RPC overload
- ✅ ETA calculation and periodic progress logging

### 6. Contract Discovery Limitations
**Problem**: Only discover contracts through events.

**Impact**: Silent contracts (no events) are not indexed

**Solution**: Add contract registry or discovery mechanism

## Query Performance Improvements ✅ (2025-10-11)

### Database Indexes Added

The indexer now automatically creates the following indexes for optimal query performance:

**JSONB GIN Indexes:**
- `idx_events_topic_gin` - Fast containment queries on event topics
- `idx_events_value_gin` - Fast containment queries on event values

**B-tree Indexes:**
- `idx_events_contract_id` - Fast contract filtering
- `idx_events_ledger` - Fast ledger range queries
- `idx_events_type` - Fast event type filtering
- `idx_transactions_ledger` - Fast transaction ledger queries
- `idx_operations_tx_hash` - Fast operation lookups by transaction
- `idx_operations_type` - Fast operation type filtering
- `idx_token_operations_contract_id` - Fast token operation queries by contract
- `idx_token_operations_from_address` - Fast token operation queries by sender
- `idx_token_operations_to_address` - Fast token operation queries by recipient

### Query Helper Functions

The `pkg/models/event.go` file now provides 7 helper functions for common query patterns:

1. `QueryEventsByTopicContains` - Find events by topic containment (@> operator)
2. `QueryEventsByTopicElement` - Find events by specific topic array element
3. `QueryEventsByContractAndTopic` - Combined contract + topic filtering
4. `QueryEventsByValuePath` - Query by JSON path in value field
5. `QueryTokenTransferEvents` - Convenience function for token transfers
6. `QueryEventsByLedgerRange` - Efficient ledger range queries

### Documentation

See `QUERY_PATTERNS.md` for:
- Comprehensive query examples (SQL, Go, GraphQL)
- Performance optimization tips
- Index usage patterns
- Hasura GraphQL examples
- Real-time subscription patterns

## Recommendations for Next Development Phase

See `DEVELOPMENT_PLAN.md` for detailed roadmap.
