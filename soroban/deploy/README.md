# Stellar RPC Indexer Deployment

This directory contains the production deployment configuration for the Stellar RPC indexer on mainnet.

## Quick Start

### 1. Create Environment File

Copy the example environment file and update with your values:

```bash
cp .env.example .env
nano .env  # or use your preferred editor
```

**IMPORTANT**: Change the `POSTGRES_PASSWORD` to a strong password!

### 2. Start Services

```bash
docker compose up -d
```

This will start:
- **stellar-rpc**: Mainnet RPC node with captive core
- **postgres**: PostgreSQL database for indexed data
- **indexer**: Standalone indexer polling RPC every 1 second
- **hasura**: GraphQL API for querying indexed data

### 3. Check Status

```bash
# View logs
docker compose logs -f

# Check individual service logs
docker compose logs -f stellar-rpc
docker compose logs -f indexer
docker compose logs -f postgres

# Check indexer stats
curl http://localhost:8080/stats

# Check indexer health
curl http://localhost:8080/health

# Access Hasura Console
open http://localhost:8081
```

## Configuration

All configuration is done via the `.env` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `POSTGRES_PASSWORD` | *(required)* | PostgreSQL password |
| `POSTGRES_DB` | `stellar_indexer` | Database name |
| `POSTGRES_USER` | `stellar` | Database user |
| `POSTGRES_PORT` | `5432` | PostgreSQL port |
| `STELLAR_RPC_VERSION` | `23.0.4` | Stellar RPC Docker image version |
| `RPC_PORT` | `8000` | RPC HTTP port |
| `ADMIN_PORT` | `6061` | RPC admin port |
| `INDEXER_PORT` | `8080` | Indexer HTTP port |
| `HASURA_PORT` | `8081` | Hasura GraphQL port |
| `HASURA_ADMIN_SECRET` | *(optional)* | Hasura admin authentication secret |
| `HASURA_UNAUTHORIZED_ROLE` | `anonymous` | Default role for unauthenticated users |

## Data Directories

All data is persisted in:
- `./data/stellar-rpc` - RPC SQLite database
- `./data/captive-core` - Captive core storage
- `./data/logs` - Log files
- PostgreSQL volume (managed by Docker)

## Services

### Stellar RPC (Port 8000)

The mainnet RPC node with embedded captive core.

**Config files:**
- `./config/stellar-rpc.toml` - RPC configuration
- `./config/stellar-core.cfg` - Captive core configuration

**Endpoints:**
- `http://localhost:8000` - JSON-RPC API
- `http://localhost:8000/health` - Health check
- `http://localhost:6061` - Admin/metrics

### PostgreSQL (Port 5432)

Database storing all indexed data.

**Connection string:**
```
postgresql://stellar:${POSTGRES_PASSWORD}@localhost:5432/stellar_indexer
```

**Tables created automatically:**
- `events` - Contract events
- `transactions` - Soroban transactions
- `token_operations` - Token transfers, mints, burns, etc.
- `token_metadata` - Token info (name, symbol, decimal)
- `token_balances` - Token holder balances
- `contract_data_entries` - Contract storage
- `account_entries` - Stellar accounts
- `trust_line_entries` - Asset trust lines
- `offer_entries` - DEX offers
- `liquidity_pool_entries` - AMM pools
- `claimable_balance_entries` - Claimable balances
- `data_entries` - Account data
- `cursor` - Indexing progress

### Indexer (Port 8080)

Standalone indexer polling RPC every 1 second.

**Endpoints:**
- `http://localhost:8080/health` - Health check
- `http://localhost:8080/stats` - Statistics (last ledger, counts)

**Polling behavior:**
- Fetches events every 1 second
- Extracts token operations from events
- Fetches transactions by hash
- Processes contract data for metadata/balances
- Updates cursor after each batch

### Hasura GraphQL Engine (Port 8081)

GraphQL API for querying all indexed data.

**Endpoints:**
- `http://localhost:8081/console` - Admin console (web UI)
- `http://localhost:8081/v1/graphql` - GraphQL API endpoint
- `http://localhost:8081/healthz` - Health check

**Features:**
- Query all indexed tables via GraphQL
- Real-time subscriptions for live data
- Aggregations, filtering, sorting
- Anonymous read-only access by default

See [hasura/README.md](hasura/README.md) for query examples and usage.

## Maintenance

### Stop Services

```bash
docker compose down
```

### Restart Services

```bash
docker compose restart
```

### Rebuild Indexer

```bash
docker compose up -d --build indexer
```

### View Resource Usage

```bash
docker stats
```

### Clear All Data (DANGEROUS)

```bash
docker compose down -v
rm -rf ./data/*
```

## Troubleshooting

### Check if RPC is syncing

```bash
docker compose logs -f stellar-rpc | grep "ledger"
```

### Check indexer progress

```bash
curl http://localhost:8080/stats
```

### Connect to PostgreSQL

```bash
docker compose exec postgres psql -U stellar -d stellar_indexer
```

### View indexer errors

```bash
docker compose logs indexer | grep -i error
```

### RPC not healthy

The RPC service takes 2-5 minutes to start (captive core initialization). Wait for the health check to pass before the indexer starts.

## Production Recommendations

1. **Use strong passwords**: Generate secure `POSTGRES_PASSWORD`
2. **Enable backups**: Set up PostgreSQL backups
3. **Monitor disk space**: Stellar RPC and PostgreSQL can grow large
4. **Set up monitoring**: Use Prometheus/Grafana for metrics
5. **Enable SSL**: Use reverse proxy (nginx/Caddy) with SSL
6. **Firewall**: Only expose necessary ports
7. **Regular updates**: Keep Docker images up to date

## Architecture

```
┌──────────────────┐
│   Stellar Core   │
│   (Captive)      │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐      ┌──────────────────┐
│   Stellar RPC    │◄─────┤     Indexer      │
│   (Port 8000)    │      │   (Port 8080)    │
│   SQLite (7d)    │      └────────┬─────────┘
└──────────────────┘               │
                                   ▼
                          ┌──────────────────┐
                          │   PostgreSQL     │◄──┐
                          │   (Port 5432)    │   │
                          │  All ledger data │   │
                          └──────────────────┘   │
                                   │             │
                                   ▼             │
                          ┌──────────────────┐   │
                          │  Hasura GraphQL  │───┘
                          │   (Port 8081)    │
                          │  Query Interface │
                          └──────────────────┘
```

- **Captive Core**: Embedded stellar-core managed by RPC
- **Stellar RPC**: 7-day SQLite retention, serves JSON-RPC
- **Indexer**: Polls every 1s, writes to PostgreSQL
- **PostgreSQL**: Long-term storage for all data
- **Hasura**: GraphQL API for querying indexed data
