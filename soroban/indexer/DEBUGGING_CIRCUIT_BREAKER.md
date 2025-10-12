# Debugging Circuit Breaker Issues

## What is Happening

The error "circuit breaker is open" means:
- The RPC client has failed **5 consecutive times**
- The circuit breaker has opened to prevent overwhelming the RPC service
- The indexer will wait **30 seconds** before trying again

## Common Causes

### 1. RPC Service Not Running
The Stellar RPC might not be started or healthy.

**Check RPC status:**
```bash
# Start Docker if not running
# Then check RPC container
docker compose -f deploy/docker-compose.yml ps

# Check RPC logs
docker compose -f deploy/docker-compose.yml logs stellar-rpc | tail -50

# Test RPC directly
curl -X POST http://localhost:8000 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"getHealth"}'
```

### 2. RPC is Starting Up (First 10-15 minutes)
Stellar RPC with captive core takes 10-15 minutes to sync on first startup.

**Check if RPC is syncing:**
```bash
docker compose -f deploy/docker-compose.yml logs stellar-rpc | grep -i "sync\|health\|ready"
```

### 3. Wrong RPC URL Configuration
The indexer might be using the wrong URL.

**Check environment variables:**
```bash
# In deploy/.env
cat deploy/.env | grep STELLAR_RPC_URL

# Should be one of:
# - http://stellar-rpc:8000 (inside Docker)
# - http://localhost:8000 (outside Docker)
```

### 4. Network Issues
Docker network problems preventing communication.

**Test network connectivity:**
```bash
# From inside indexer container
docker compose -f deploy/docker-compose.yml exec stellar-indexer \
  wget -O- http://stellar-rpc:8000/health

# Or try from host
curl http://localhost:8000/health
```

### 5. RPC Rate Limiting or Timeout
Requests are timing out (>10 seconds) or being rate limited.

## Step-by-Step Debugging

### Step 1: Check Docker Status
```bash
# Start Docker Desktop or OrbStack
# Verify it's running
docker ps

# If not running, start services
cd deploy
docker compose up -d
```

### Step 2: Check RPC Health
```bash
# Wait for RPC to be healthy (may take 10-15 min on first start)
./deploy/check-rpc-health.sh

# Or manually check
curl http://localhost:8000/health
# Should return: {"status":"healthy"}
```

### Step 3: Check Indexer Logs
```bash
# Look for the actual error before circuit breaker opened
docker compose -f deploy/docker-compose.yml logs stellar-indexer | grep -E "error|Error|failed" | tail -20

# Look for RPC connection attempts
docker compose -f deploy/docker-compose.yml logs stellar-indexer | grep -E "rpc|RPC" | tail -20
```

### Step 4: Check Configuration
```bash
# Verify STELLAR_RPC_URL
docker compose -f deploy/docker-compose.yml exec stellar-indexer env | grep STELLAR_RPC_URL

# Should output: STELLAR_RPC_URL=http://stellar-rpc:8000
```

### Step 5: Test RPC Connectivity from Indexer
```bash
# Test from inside the indexer container
docker compose -f deploy/docker-compose.yml exec stellar-indexer \
  wget -qO- http://stellar-rpc:8000/health

# If this fails, there's a network issue
```

### Step 6: Restart with Fresh Logs
```bash
cd deploy

# Stop everything
docker compose down

# Start RPC first and wait for it to be healthy
docker compose up -d stellar-rpc
./check-rpc-health.sh

# Once healthy, start indexer
docker compose up -d stellar-indexer

# Watch logs
docker compose logs -f stellar-indexer
```

## Circuit Breaker Configuration

Located in `pkg/client/rpc.go:27`:

```go
circuitBreaker: worker.NewCircuitBreaker(
    5,              // Fail after 5 consecutive errors
    10*time.Second, // 10 second timeout per request
    30*time.Second, // Wait 30 seconds before retrying (half-open state)
)
```

To adjust these values, edit the NewCircuitBreaker parameters:
- **First parameter**: Number of failures before opening (default: 5)
- **Second parameter**: Request timeout (default: 10s)
- **Third parameter**: Reset timeout (default: 30s)

## Quick Fixes

### Fix 1: Increase Circuit Breaker Timeout
If RPC is slow, increase timeout:

```go
// pkg/client/rpc.go:27
circuitBreaker: worker.NewCircuitBreaker(
    5,
    30*time.Second, // Increased from 10s to 30s
    30*time.Second,
)
```

### Fix 2: Increase Failure Threshold
If RPC has intermittent failures:

```go
// pkg/client/rpc.go:27
circuitBreaker: worker.NewCircuitBreaker(
    10,             // Increased from 5 to 10
    10*time.Second,
    30*time.Second,
)
```

### Fix 3: Use External RPC (for testing)
If local RPC isn't working, use a public endpoint:

```bash
# In deploy/.env
STELLAR_RPC_URL=https://soroban-testnet.stellar.org

# Restart indexer
docker compose restart stellar-indexer
```

## Understanding the Logs

### Good Signs ✅
```
level=info msg="RPC health check passed"
level=info msg="Processing batch" events=50
level=info msg="Batch processed successfully"
```

### Bad Signs ❌
```
level=error msg="Poll failed" error="circuit breaker is open"
level=error msg="get latest ledger: circuit breaker is open"
level=warn msg="RPC request failed" error="context deadline exceeded"
level=error msg="dial tcp: connect: connection refused"
```

## Testing Account Entries Fix

Once the circuit breaker issue is resolved, you should see:

```bash
# These logs indicate the fix is working
docker compose logs -f stellar-indexer | grep -E "account|Account"

# Expected output:
# level=debug msg="Processing account ledger entries" accountCount=5
# level=info msg="Processed ledger entries" accounts=5 trustlines=2 offers=0
```

## Manual Testing

If you want to test the indexer without Docker:

```bash
cd indexer

# Set environment variables
export STELLAR_RPC_URL=https://soroban-testnet.stellar.org
export POSTGRES_DSN="host=localhost user=stellar password=stellar dbname=stellar port=5432 sslmode=disable"

# Start PostgreSQL separately if needed
# docker run -d -p 5432:5432 -e POSTGRES_USER=stellar -e POSTGRES_PASSWORD=stellar -e POSTGRES_DB=stellar postgres:15

# Run indexer
./build/indexer
```

## Need More Help?

1. **Enable debug logging**: Set log level to DEBUG in the indexer
2. **Check RPC logs**: `docker compose logs stellar-rpc -f`
3. **Verify PostgreSQL**: `docker compose logs postgres | tail -20`
4. **Network inspection**: `docker network inspect deploy_default`
5. **Port conflicts**: `lsof -i :8000` to check if port 8000 is in use

## Summary

The circuit breaker is **protecting** the system from a failing RPC service. The most likely causes are:

1. **RPC still syncing** (wait 10-15 min on first start)
2. **RPC not running** (start with `docker compose up -d`)
3. **Wrong URL configuration** (check deploy/.env)
4. **Network issues** (check docker network connectivity)

Follow the debugging steps above to identify and fix the root cause.
