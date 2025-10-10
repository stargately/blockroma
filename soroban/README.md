# Blockroma Soroban Indexer

Standalone indexer for Stellar Soroban smart contract events and transactions. Polls Stellar RPC and writes data to PostgreSQL for analytics and querying.

## Architecture

```
┌─────────────────────────────┐
│  Upstream Stellar RPC       │
│  (v23.0.4 - Unmodified)     │
│  Port 8000 (JSON-RPC API)   │
└─────────────┬───────────────┘
              │ Poll every 1s
              ▼
┌─────────────────────────────┐
│  Standalone Indexer         │
│  (Go Binary)                │
└─────────────┬───────────────┘
              │ Write events/txs
              ▼
       ┌──────────────┐
       │ PostgreSQL   │◄────────┐
       └──────┬───────┘         │
              │                 │
              ▼                 │
       ┌──────────────┐         │
       │   Hasura     │─────────┘
       │  (GraphQL)   │
       └──────────────┘
```

## Features

✅ **No Fork Maintenance** - Uses upstream Stellar RPC (v23.0.4)
✅ **Direct PostgreSQL Writes** - No Redis, simpler architecture
✅ **GraphQL API** - Query indexed data via Hasura GraphQL Engine
✅ **Complete Transaction Metadata** - Stores full tx data including memos, signatures, preconditions
✅ **Token Operations** - Tracks SAC token transfers, mints, burns
✅ **Ledger State Tracking** - Indexes all ledger entries (accounts, trustlines, offers, etc.)
✅ **Easy Upgrades** - Just `docker pull` new RPC version
✅ **Comprehensive Tests** - 36 tests, full feature parity

## Quick Start

### Prerequisites

- Docker & Docker Compose
- 200GB+ disk space (for RPC captive core)
- 4GB+ RAM

### 1. Configure Environment

```bash
cd deploy
cp .env.example .env
# Edit .env with your PostgreSQL password
```

### 2. Start Services

```bash
docker compose up -d
```

This starts:
- **stellar-rpc** - Upstream Stellar RPC (port 8000)
- **postgres** - PostgreSQL database (port 5432)
- **indexer** - Event/transaction indexer (port 8080)
- **hasura** - GraphQL API (port 8081)

### 3. Check Status

```bash
# Check RPC health
./check-rpc-health.sh

# Check logs
docker logs -f stellar-indexer

# Check database
docker exec -it stellar-postgres psql -U stellar -d stellar_indexer -c "SELECT COUNT(*) FROM events;"

# Access Hasura Console
open http://localhost:8081
```

## Project Structure

```
soroban/
├── indexer/              # Go indexer service
│   ├── cmd/indexer/      # Main entry point
│   ├── pkg/
│   │   ├── client/       # RPC JSON-RPC client
│   │   ├── parser/       # Event/tx parser
│   │   ├── models/       # Database models
│   │   ├── poller/       # Polling logic
│   │   └── db/           # Database connection
│   ├── Dockerfile
│   ├── Makefile
│   └── README.md         # Development docs
└── deploy/               # Docker deployment
    ├── docker-compose.yml
    ├── config/           # RPC configuration
    ├── hasura/           # Hasura GraphQL metadata
    └── scripts/          # Helper scripts
```

## Development

See [indexer/README.md](indexer/README.md) for development instructions.

See [deploy/README.md](deploy/README.md) for deployment details.

## Database Tables

The indexer creates these tables:

- `events` - Contract events
- `transactions` - Transaction details
- `token_operations` - Token transfers/mints/burns
- `token_metadata` - Token info (name, symbol, decimals)
- `token_balances` - Token holder balances
- `account_entries` - Account data
- `trust_line_entries` - Trustlines
- `offer_entries` - DEX offers
- `liquidity_pool_entries` - Liquidity pools
- `claimable_balance_entries` - Claimable balances
- `contract_data_entries` - Contract storage
- `data_entries` - Account data entries
- `cursor` - Indexer sync position

## Module Information

**Module:** `github.com/blockroma/soroban-indexer`
**Go Version:** 1.23+

## Testing

```bash
cd indexer
make test              # Run all tests
make test-coverage     # Generate coverage report
make test-coverage-html # Open coverage in browser
```

**Test Results:** 36 tests, 33 passing, 3 skipped

## GraphQL API

The deployment includes Hasura GraphQL Engine for querying indexed data.

**Access the GraphQL Console**: http://localhost:8081

**Example Query:**
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

See [deploy/hasura/README.md](deploy/hasura/README.md) for more query examples and API usage.

## Documentation

- [TESTING.md](indexer/TESTING.md) - Testing guide
- [FEATURE_PARITY.md](indexer/FEATURE_PARITY.md) - Feature comparison with old indexer
- [CHANGES_SUMMARY.md](indexer/CHANGES_SUMMARY.md) - Implementation details
- [HEALTH_CHECK.md](deploy/HEALTH_CHECK.md) - RPC health check guide
- [deploy/hasura/README.md](deploy/hasura/README.md) - Hasura GraphQL usage

## Monitoring

The indexer exposes metrics on port 8080 (configurable via `INDEXER_PORT`).

### Health Check

```bash
curl http://localhost:8080/health
```

### Logs

```bash
# Indexer logs
docker logs -f stellar-indexer

# RPC logs
docker logs -f stellar-rpc-mainnet

# PostgreSQL logs
docker logs -f stellar-postgres

# Hasura logs
docker logs -f stellar-hasura
```

## Performance

- **Polling Interval:** 1 second
- **Batch Size:** 1000 events per request
- **Latency:** ~1 second behind RPC
- **Memory:** ~100MB (indexer) + ~50MB (Hasura)
- **CPU:** <5% idle

## Troubleshooting

### RPC takes long to start

RPC needs to sync captive core. First startup can take 10-15 minutes. Check logs:

```bash
docker logs -f stellar-rpc-mainnet
```

Wait for: `INFO Initializing transaction store...`

### Indexer connection errors

Make sure RPC is healthy:

```bash
./deploy/check-rpc-health.sh
```

### Database connection failed

Check PostgreSQL is running:

```bash
docker ps | grep postgres
docker logs stellar-postgres
```

## License

[Your License Here]

## Support

For issues and questions, please open an issue on the Blockroma repository.

---

**Built with ❤️ for the Stellar ecosystem**
