# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Blockroma Soroban Indexer is a standalone indexer for Stellar Soroban smart contract events and transactions. It polls the upstream Stellar RPC (v23.0.4) and writes data to PostgreSQL for analytics and querying.

**Key Architecture Decision**: This project uses the upstream Stellar RPC unmodified (no fork), with a standalone Go indexer that polls the RPC API. This eliminates fork maintenance overhead and simplifies upgrades.

**Module**: `github.com/blockroma/soroban-indexer`
**Go Version**: 1.23+

## Development Commands

### Build and Run
```bash
# Build the indexer binary
cd indexer && make build

# Run the indexer locally (requires env vars)
cd indexer && make run

# Run tests
cd indexer && make test

# Run tests with coverage
cd indexer && make test-coverage

# View coverage report in browser
cd indexer && make test-coverage-html

# Run specific package tests
cd indexer && make test-package PKG=./pkg/parser
```

### Docker Deployment
```bash
# Start all services (RPC, PostgreSQL, Indexer, Hasura)
cd deploy && docker compose up -d

# Check indexer logs
docker logs -f stellar-indexer

# Check Hasura logs
docker logs -f stellar-hasura

# Access Hasura Console
open http://localhost:8081

# Check RPC health
cd deploy && ./check-rpc-health.sh

# Stop all services
cd deploy && docker compose down
```

### Code Quality
```bash
cd indexer

# Format code
make fmt

# Run linter
make lint

# Run vet
make vet

# Run all checks (fmt, vet, test)
make check

# CI pipeline (deps, verify, checks, coverage)
make ci
```

## Architecture

### High-Level Data Flow

```
Stellar RPC (upstream) → Poller (1s interval) → Parser → PostgreSQL ← Hasura GraphQL
                                ↓                              ↑
                          Token Operations                     │
                          Contract Data                  (GraphQL API)
                          Ledger Entries
```

### Component Breakdown

**indexer/cmd/indexer/main.go**
- Entry point for the indexer service
- Initializes RPC client, database connection, and poller
- Starts HTTP server for health checks and metrics on port 8080
- Handles graceful shutdown via signals

**indexer/pkg/client/**
- `rpc.go`: JSON-RPC client for Stellar RPC
- Implements methods: `GetLatestLedger()`, `GetEvents()`, `GetTransaction()`, `GetLedgerEntries()`, `Health()`
- Uses standard JSON-RPC 2.0 protocol over HTTP

**indexer/pkg/poller/**
- `poller.go`: Main polling loop (1 second interval)
- Orchestrates the indexing workflow:
  1. Get cursor (last processed ledger) from database
  2. Fetch new events from RPC since cursor
  3. Parse events and extract token operations
  4. Fetch transactions for unique tx hashes
  5. Process contract data for discovered contracts
  6. Update cursor in database transaction
- Batch size: 1000 events per request
- All database writes happen in a single GORM transaction for consistency

**indexer/pkg/parser/**
- `parser.go`: Converts RPC responses to database models
  - `ParseEvent()`: Decodes XDR ScVal topics/values to JSON
  - `ParseTransaction()`: Extracts source account, fee, sequence from XDR envelope
- `token_operations.go`: Extracts token transfers/mints/burns from events
  - `ParseTokenOperation()`: Recognizes standard SAC token event patterns
  - `ParseTokenMetadata()`: Extracts name, symbol, decimals from contract storage
  - `ParseTokenBalance()`: Extracts holder balances from contract storage
- `contract_data_meta.go`: Parses contract data entries and contract code
  - Note: Classic Stellar ledger entries (accounts, trustlines, offers) are NOT indexed - use Horizon API instead

**indexer/pkg/models/**
- Database models using GORM (PostgreSQL)
- Key tables:
  - `events`: Contract events with decoded topics/values (JSONB)
  - `transactions`: Full transaction details including memos, signatures, preconditions
  - `operations`: Transaction operations
  - `token_operations`: Extracted token transfers/mints/burns
  - `token_metadata`: Token info (name, symbol, decimals)
  - `token_balances`: Token holder balances
  - `contract_data_entries`: Soroban contract data storage
  - `contract_code`: Soroban contract WASM code
  - `cursor`: Tracks last processed ledger for recovery
- All models have `UpsertX()` methods for idempotent writes
- Note: Classic Stellar ledger entries (accounts, trustlines, offers, etc.) are NOT indexed - use Horizon API instead

**indexer/pkg/db/**
- `db.go`: Database connection and migrations
- Auto-migrates all tables on startup

**deploy/hasura/**
- Hasura GraphQL Engine for querying indexed data
- Pre-configured metadata tracking all Soroban-specific database tables
- Provides GraphQL API with queries, mutations, and subscriptions
- Anonymous read-only access enabled by default
- Console UI at `http://localhost:8081` for schema exploration and testing
- Note: Classic Stellar ledger entries (accounts, trustlines, offers, etc.) are NOT included - use Horizon API for those

### Key Design Patterns

1. **Cursor-based Recovery**: The indexer tracks its position in the `cursor` table. On restart, it resumes from the last processed ledger.

2. **Transaction Batching**: All database writes for a polling cycle happen in a single GORM transaction. If any step fails, the entire batch rolls back and will be retried.

3. **Idempotent Upserts**: All database writes use upsert operations (INSERT ... ON CONFLICT DO UPDATE), so re-processing the same data is safe.

4. **XDR Decoding**: The parser uses Stellar's `go-xdr` package to decode base64 XDR into Go types, then converts to JSON for storage.

5. **Token Operation Detection**: The parser recognizes standard SAC (Stellar Asset Contract) event patterns by matching topic values against known patterns (e.g., `["transfer", Symbol("from"), ...]`).

## Configuration

The indexer requires only 2 environment variables:

- `STELLAR_RPC_URL`: URL of the Stellar RPC endpoint (default: `http://stellar-rpc:8000`)
- `POSTGRES_DSN`: PostgreSQL connection string (required)

Optional:
- `INDEXER_PORT`: HTTP server port for health/stats endpoints (default: `8080`)
- `HASURA_PORT`: Hasura GraphQL port (default: `8081`)
- `HASURA_ADMIN_SECRET`: Admin secret for Hasura (optional, recommended for production)
- `HASURA_UNAUTHORIZED_ROLE`: Default role for unauthenticated users (default: `anonymous`)

See `deploy/.env.example` for deployment configuration.

## Testing

The project has 36 tests across three categories:

1. **Parser Tests** (`pkg/parser/*_test.go`):
   - Event parsing and XDR decoding
   - Token operation extraction
   - Transaction parsing

2. **Client Tests** (`pkg/client/rpc_test.go`):
   - RPC method calls (uses `httptest.NewServer` for mocking)
   - JSON-RPC request/response handling
   - Error handling

3. **Model Tests** (`pkg/models/models_test.go`):
   - Database operations (uses in-memory SQLite)
   - Upsert logic
   - Custom types (e.g., Int128)

Tests use table-driven test patterns and in-memory databases (no external dependencies).

## GraphQL API with Hasura

The deployment includes Hasura GraphQL Engine that automatically exposes all indexed tables via GraphQL.

### Accessing Hasura

- **Console**: http://localhost:8081/console
- **GraphQL Endpoint**: http://localhost:8081/v1/graphql
- **Health Check**: http://localhost:8081/healthz

### Example Queries

**Get latest events:**
```graphql
query GetLatestEvents {
  events(order_by: {ledger: desc}, limit: 10) {
    id
    contract_id
    type
    ledger
    topic
    value
  }
}
```

**Get token transfers for a contract:**
```graphql
query GetTokenTransfers($contract_id: String!) {
  token_operations(
    where: {contract_id: {_eq: $contract_id}, operation_type: {_eq: "transfer"}}
    order_by: {ledger: desc}
  ) {
    from_address
    to_address
    amount
    ledger_closed_at
  }
}
```

**Real-time subscription to new events:**
```graphql
subscription NewEvents {
  events(order_by: {ledger: desc}, limit: 1) {
    id
    contract_id
    ledger
    topic
    value
  }
}
```

See `deploy/hasura/README.md` for more query examples and advanced usage.

### Metadata Management

Hasura metadata is stored in `deploy/hasura/metadata/` and includes:
- All Soroban-specific tables pre-tracked with anonymous read permissions
- Proper column configurations for each table
- Version-controlled metadata that auto-applies on startup

If you add new tables to the indexer:
1. Track them in Hasura Console (http://localhost:8081/console/data)
2. Export metadata to save changes
3. Commit metadata files to version control

**Important**: This indexer focuses on Soroban-specific data only. For classic Stellar ledger entries (accounts, trustlines, offers, data entries, claimable balances, liquidity pools), use the Horizon API instead.

## Common Development Scenarios

### Adding a New Event Parser

1. Add parsing logic to `pkg/parser/parser.go` or a new file in `pkg/parser/`
2. Create corresponding model in `pkg/models/` if needed
3. Add upsert method to model
4. Call parser and upsert in `pkg/poller/poller.go:poll()`
5. Add tests in `pkg/parser/*_test.go`

### Adding a New RPC Method

1. Add method to `pkg/client/rpc.go`
2. Define request/response types
3. Implement using `c.call()` helper
4. Add test in `pkg/client/rpc_test.go` with mock server

### Modifying Database Schema

1. Update model in `pkg/models/`
2. GORM auto-migrates on startup (no manual migrations)
3. For breaking changes, consider creating a new table and migrating data
4. If adding a new table, track it in Hasura Console and export metadata

### Debugging Indexing Issues

1. Check indexer logs: `docker logs -f stellar-indexer`
2. Check cursor position: `SELECT * FROM cursor;`
3. Check stats endpoint: `curl http://localhost:8080/stats`
4. Reset cursor to re-index: `UPDATE cursor SET last_ledger = X WHERE id = 1;`
5. For parser issues, add logging in `pkg/parser/` and rebuild

## Important Notes

- The indexer is designed to run continuously. Stopping and restarting is safe due to cursor-based recovery.
- The first startup of Stellar RPC takes 10-15 minutes as it syncs captive core. The indexer will wait for RPC to be healthy.
- All database writes are in transactions, so partial writes cannot occur.
- The poller processes events sequentially (not in parallel) to maintain order.
- Token operations are opportunistically extracted from events. Not all events will have token operations.
- Contract data processing (for metadata/balances) happens after events are stored, and failures don't block event indexing.
