# Stellar RPC Health Check Guide

## Important: Health Endpoint is JSON-RPC, not REST

The stellar-rpc health check is **NOT** a REST endpoint like `/health`.  
It's a **JSON-RPC method** called `getHealth`.

## How to Check Health

### From Host Machine

```bash
curl -X POST http://localhost:8000 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"getHealth"}'
```

### From Inside Container

```bash
docker exec stellar-rpc-mainnet wget -q -O- \
  --post-data='{"jsonrpc":"2.0","id":1,"method":"getHealth"}' \
  --header='Content-Type: application/json' \
  http://localhost:8000
```

### Using Helper Script

```bash
./check-rpc-health.sh
# Or specify URL:
./check-rpc-health.sh http://localhost:8000
```

## Response Format

### Healthy Response

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "status": "healthy",
    "latestLedger": 12345678,
    "oldestLedger": 12245678,
    "ledgerRetentionWindow": 100800
  }
}
```

### Unhealthy Response (Not Initialized)

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32603,
    "message": "data stores are not initialized"
  }
}
```

### Unhealthy Response (High Latency)

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32603,
    "message": "latency (30s) since last known ledger closed is too high (>10s)"
  }
}
```

## Health Check Logic

The RPC is considered healthy when:

1. ✅ Data stores are initialized (at least 1 ledger ingested)
2. ✅ Ledger latency < MAX_HEALTHY_LEDGER_LATENCY (default: 30s in config)

## Typical Startup Timeline

| Time | Status | What's Happening |
|------|--------|------------------|
| 0-2 min | ❌ Error: Connection refused | Container starting |
| 2-5 min | ❌ Error: Data stores not initialized | Captive core starting |
| 5-10 min | ❌ Error: Data stores not initialized | Syncing first ledgers |
| 10+ min | ✅ Healthy | Fully operational |

## Docker Health Check

The docker-compose health check runs this command every 30 seconds:

```bash
wget -q -O- \
  --post-data='{"jsonrpc":"2.0","id":1,"method":"getHealth"}' \
  --header='Content-Type: application/json' \
  http://localhost:8000 | grep -q '"status":"healthy"'
```

It waits up to 5 minutes (start_period: 300s) before marking as unhealthy.

## Troubleshooting

### Health check keeps failing

1. Check if RPC is actually running:
   ```bash
   docker ps | grep stellar-rpc
   ```

2. Check RPC logs:
   ```bash
   docker logs -f stellar-rpc-mainnet
   ```

3. Try manual health check:
   ```bash
   ./check-rpc-health.sh
   ```

4. Check if port is accessible:
   ```bash
   nc -zv localhost 8000
   ```

### Using simplified docker-compose

If health checks are problematic, use the simplified version:

```bash
docker compose -f docker-compose.simple.yml up -d
```

This version:
- No health check dependency for indexer
- Indexer will retry connections automatically
- Simpler and more reliable

## Configuration

Health check settings in `config/stellar-rpc.toml`:

```toml
# Maximum acceptable latency for health check
MAX_HEALTHY_LEDGER_LATENCY = "30s"

# Ledger retention (affects oldest ledger in health response)
HISTORY_RETENTION_WINDOW = 100800  # 7 days
```
