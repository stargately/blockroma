# Hasura Troubleshooting Guide

## Issue: "no_queries_available" Error

**Symptom:**
When accessing Hasura console, you see:
```json
{
  "data": {
    "no_queries_available": "There are no queries available to the current role..."
  }
}
```

**Cause:** Tables haven't been tracked in Hasura yet.

**Solution:**
Run the initialization script to track all tables:
```bash
cd deploy
./scripts/init-hasura.sh
```

This will:
1. Wait for Hasura to be ready
2. Track all 13 database tables
3. Set up anonymous role permissions
4. Enable aggregation queries

Then refresh the Hasura console.

---

## Issue: Tables Don't Exist in Database

**Symptom:**
Init script shows errors about tables not existing.

**Cause:** The indexer hasn't created the tables yet (database is empty).

**Solution:**
1. Make sure the indexer is running:
   ```bash
   docker logs -f stellar-indexer
   ```

2. Wait for the indexer to process the first batch of events (this creates the tables)

3. Once tables are created, run the init script again:
   ```bash
   ./scripts/init-hasura.sh
   ```

---

## Issue: Permission Denied Errors

**Symptom:**
Queries fail with permission errors.

**Cause:** Anonymous role permissions not set up correctly.

**Solution:**
1. Re-run the init script:
   ```bash
   ./scripts/init-hasura.sh
   ```

2. Or manually grant permissions in Hasura Console:
   - Go to Data → [table_name] → Permissions
   - Add "anonymous" role
   - Select all columns
   - Set filter to `{}`
   - Enable "Allow aggregations"

---

## Issue: Hasura Won't Start

**Symptom:**
Docker shows hasura container is not running.

**Solution:**
1. Check logs:
   ```bash
   docker logs stellar-hasura
   ```

2. Verify PostgreSQL is running:
   ```bash
   docker ps | grep postgres
   ```

3. Check environment variables:
   ```bash
   docker compose config | grep HASURA
   ```

4. Restart the service:
   ```bash
   docker compose restart hasura
   ```

---

## Issue: Can't Connect to GraphQL Endpoint

**Symptom:**
HTTP requests to `http://localhost:8081/v1/graphql` fail.

**Solution:**
1. Verify Hasura is running:
   ```bash
   docker ps | grep hasura
   ```

2. Check health endpoint:
   ```bash
   curl http://localhost:8081/healthz
   ```

3. Verify port mapping:
   ```bash
   docker port stellar-hasura
   ```

4. Check if port 8081 is available:
   ```bash
   lsof -i :8081
   ```

---

## Issue: Metadata Not Loading

**Symptom:**
Tables show up but no permissions or relationships.

**Cause:** The metadata volume mount approach doesn't auto-apply metadata.

**Solution:**
Use the init script instead of relying on volume-mounted metadata:
```bash
./scripts/init-hasura.sh
```

This programmatically configures Hasura via the metadata API.

---

## Issue: Old Metadata Conflicts

**Symptom:**
Errors about existing permissions or tables already tracked.

**Solution:**
This is normal - the init script will skip items that already exist. If you need to reset Hasura completely:

1. Stop Hasura:
   ```bash
   docker compose stop hasura
   ```

2. Clear Hasura's internal state (stored in PostgreSQL):
   ```bash
   docker exec -it stellar-postgres psql -U stellar -d stellar_indexer -c "
   DROP SCHEMA IF EXISTS hdb_catalog CASCADE;
   "
   ```

3. Restart Hasura:
   ```bash
   docker compose up -d hasura
   ```

4. Re-run init script:
   ```bash
   ./scripts/init-hasura.sh
   ```

---

## Verifying Setup

After running the init script, verify everything is working:

1. **Check tables are tracked:**
   ```bash
   ./scripts/check-hasura.sh
   ```

2. **Test a simple query:**
   ```bash
   curl -X POST http://localhost:8081/v1/graphql \
     -H "Content-Type: application/json" \
     -d '{"query": "{ cursor { id last_ledger } }"}'
   ```

3. **Check console access:**
   Open http://localhost:8081/console and go to the "Data" tab.
   You should see all 13 tables listed.

---

## Manual Table Tracking (Alternative)

If the init script doesn't work, you can manually track tables in the console:

1. Open http://localhost:8081/console
2. Go to "Data" tab
3. Click "public" schema
4. For each untracked table, click "Track"
5. Go to table → Permissions
6. Add "anonymous" role with select permission

---

## Getting Help

If issues persist:

1. Check Hasura logs:
   ```bash
   docker logs stellar-hasura | tail -50
   ```

2. Check PostgreSQL connection:
   ```bash
   docker exec -it stellar-hasura env | grep HASURA_GRAPHQL_DATABASE_URL
   ```

3. Test database connectivity:
   ```bash
   docker exec -it stellar-postgres psql -U stellar -d stellar_indexer -c "\dt"
   ```

4. Consult Hasura documentation:
   https://hasura.io/docs/latest/index/

---

## Common Environment Variable Issues

Make sure these are set correctly in `.env`:

```bash
# Required
POSTGRES_PASSWORD=your_password_here

# Optional (with defaults)
HASURA_PORT=8081
HASURA_ADMIN_SECRET=               # Leave empty for dev mode
HASURA_UNAUTHORIZED_ROLE=anonymous
```

Verify they're loaded:
```bash
docker compose config | grep -A 5 hasura
```
