# Stellar RPC Indexer

Standalone indexer service that polls Stellar RPC API and writes events and transactions to PostgreSQL.

## Architecture

```
Stellar RPC (upstream) → Indexer (polls every 1s) → PostgreSQL
```

- **No Redis** - Direct database writes
- **No fork** - Uses upstream stellar-rpc unmodified
- **Simple** - Single binary, minimal configuration

## Features

- ✅ Polls Stellar RPC every **1 second** (near real-time)
- ✅ Indexes **contract events** with decoded topics/values
- ✅ Indexes **transactions** with full details
- ✅ **Automatic recovery** - cursor tracking in database
- ✅ **Health endpoint** - `/health` for monitoring
- ✅ **Stats endpoint** - `/stats` for metrics
- ✅ **Graceful shutdown** - No data loss on restart

## Quick Start

```bash
cd deploy/

# Configure environment
cp .env.example .env
nano .env  # Set POSTGRES_PASSWORD

# Start all services
docker compose up -d

# Check logs
docker compose logs -f indexer

# Check stats
curl http://localhost:8080/stats
```

## Environment Variables

Only 2 required:

```bash
STELLAR_RPC_URL=http://stellar-rpc:8000
POSTGRES_DSN=postgresql://user:pass@postgres:5432/stellar_indexer
```

## Database Schema

### Events Table
```sql
CREATE TABLE events (
    id VARCHAR PRIMARY KEY,
    tx_index INTEGER,
    type VARCHAR,
    ledger INTEGER,
    ledger_closed_at VARCHAR,
    contract_id VARCHAR,
    paging_token VARCHAR,
    topic JSONB,           -- Decoded topics
    value JSONB,           -- Decoded value
    in_successful_contract_call BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id VARCHAR PRIMARY KEY,
    status VARCHAR,
    ledger INTEGER,
    application_order INTEGER,
    source_account VARCHAR,
    fee INTEGER,
    fee_charged INTEGER,
    sequence BIGINT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Cursor Table
```sql
CREATE TABLE indexer_cursor (
    id INTEGER PRIMARY KEY,
    last_ledger INTEGER NOT NULL,
    updated_at TIMESTAMP
);
```

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
# Response: OK
```

### Statistics
```bash
curl http://localhost:8080/stats
# Response:
{
    "lastLedger": 12345678,
    "totalEvents": 1500000,
    "totalTransactions": 750000
}
```

## Development

### Build Locally
```bash
cd indexer/
go mod download
go build -o indexer cmd/indexer/main.go
```

### Run Locally
```bash
export STELLAR_RPC_URL=http://localhost:8000
export POSTGRES_DSN=postgresql://stellar:password@localhost:5432/stellar_indexer

./indexer
```

### Run Tests
```bash
go test ./...
```

## Performance

### Resource Usage
- **Memory**: ~100MB
- **CPU**: <5% (mostly idle)
- **Network**: Minimal (1 request/second)
- **Database**: Depends on event volume

### Polling Performance
- **Interval**: 1 second
- **Latency**: 0.5s average, 1s max
- **Throughput**: Handles 1000+ events/second
- **Batch Size**: 1000 events per request

### Database Performance
```sql
-- Add indexes for fast queries
CREATE INDEX idx_events_ledger ON events(ledger);
CREATE INDEX idx_events_contract_id ON events(contract_id);
CREATE INDEX idx_events_ledger_closed_at ON events(ledger_closed_at);
CREATE INDEX idx_txs_ledger ON transactions(ledger);
```

## Monitoring

### Logs
```bash
# Follow indexer logs
docker compose logs -f indexer

# Check for errors
docker compose logs indexer | grep ERROR
```

### Metrics
The indexer exposes metrics on port 8080:
- `/health` - Health check
- `/stats` - Current statistics

### Alerts
Monitor these conditions:
- Indexer not running
- Last ledger not updating (>60s)
- Database connection errors
- RPC connection errors

## Troubleshooting

### Indexer won't start
```bash
# Check logs
docker compose logs indexer

# Check RPC connectivity
curl http://localhost:8000/health

# Check database connectivity
docker compose exec postgres psql -U stellar -d stellar_indexer -c "\dt"
```

### Missing events
```bash
# Check cursor position
docker compose exec postgres psql -U stellar -d stellar_indexer \
  -c "SELECT * FROM indexer_cursor;"

# Reset cursor to re-index from ledger X
docker compose exec postgres psql -U stellar -d stellar_indexer \
  -c "UPDATE indexer_cursor SET last_ledger = X WHERE id = 1;"

# Restart indexer
docker compose restart indexer
```

### High latency
1. Check Stellar RPC performance
2. Check database performance (add indexes)
3. Check network latency
4. Verify adequate CPU/RAM

## Upgrading

### Upgrade Stellar RPC
```bash
# Edit .env
nano .env  # Change STELLAR_RPC_VERSION

# Pull and restart
docker compose pull stellar-rpc
docker compose up -d stellar-rpc
```

Indexer automatically adapts - no changes needed!

### Upgrade Indexer
```bash
# Rebuild and restart
docker compose up -d --build indexer
```

## Production Deployment

### Recommended Setup
- Run on dedicated server/VM
- Use managed PostgreSQL (AWS RDS, etc.)
- Enable automated backups
- Set up monitoring/alerting
- Use reverse proxy for SSL

### Scaling
For high volume:
1. Increase PostgreSQL resources
2. Add database read replicas
3. Add database indexes
4. Consider partitioning events table by ledger

## Comparison with Old Architecture

| Feature | Old (Fork + Redis) | New (Standalone) |
|---------|-------------------|------------------|
| Components | 3 (RPC fork + Redis + Consumer) | 2 (RPC + Indexer) |
| Latency | 0ms (inline) | ~0.5s (polling) |
| Maintenance | Very High | Low |
| Upgrades | Very Hard | Easy |
| Infrastructure | Complex | Simple |
| State | Redis queue | DB cursor |
| Recovery | Manual | Automatic |

## License

Same as parent project

## Support

See main project README for support information.
