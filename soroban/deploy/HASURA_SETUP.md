# Hasura GraphQL Setup Summary

This document summarizes the Hasura GraphQL Engine integration with the Stellar Soroban Indexer.

## What Was Added

### 1. Docker Service
Added Hasura GraphQL Engine v2.40.0 to `docker-compose.yml`:
- Runs on port 8081 (configurable via `HASURA_PORT`)
- Connects to the same PostgreSQL database as the indexer
- Pre-configured with metadata tracking all indexed tables
- Anonymous read-only access enabled by default

### 2. Metadata Configuration
Created complete Hasura metadata in `deploy/hasura/metadata/`:
- All 13 database tables are tracked and exposed via GraphQL
- Anonymous role with read permissions on all tables
- Aggregation queries enabled for analytics
- Version-controlled metadata that auto-applies on startup

**Tracked Tables:**
- `events` - Contract events with decoded topics/values
- `transactions` - Full transaction details
- `token_operations` - Token transfers, mints, burns
- `token_metadata` - Token info (name, symbol, decimals)
- `token_balances` - Token holder balances
- `account_entries` - Stellar account data
- `trust_line_entries` - Asset trustlines
- `offer_entries` - DEX offers
- `liquidity_pool_entries` - AMM pools
- `claimable_balance_entries` - Claimable balances
- `contract_data_entries` - Contract storage data
- `data_entries` - Account data entries
- `cursor` - Indexing position tracker

### 3. Environment Variables
Added to `.env.example`:
```bash
HASURA_PORT=8081                      # GraphQL API port
HASURA_ADMIN_SECRET=                  # Admin auth (optional)
HASURA_UNAUTHORIZED_ROLE=anonymous    # Default role for public access
```

### 4. Documentation
- `deploy/hasura/README.md` - Comprehensive Hasura usage guide with query examples
- Updated `deploy/README.md` with Hasura service information
- Updated main `README.md` with GraphQL API section
- Updated `CLAUDE.md` with Hasura integration details

### 5. Health Check Script
Created `deploy/scripts/check-hasura.sh` to verify:
- Hasura service health
- GraphQL endpoint connectivity
- Table tracking status
- Current data counts

## Quick Start

### 1. Start All Services
```bash
cd deploy
docker compose up -d
```

### 2. Access Hasura Console
Open http://localhost:8081/console in your browser

### 3. Try a Query
In the GraphQL tab of the console, run:
```graphql
query GetLatestEvents {
  events(order_by: {ledger: desc}, limit: 5) {
    id
    contract_id
    type
    ledger
    topic
    value
  }
}
```

### 4. Verify Setup
```bash
cd deploy
./scripts/check-hasura.sh
```

## GraphQL Features Available

### Queries
- **Simple queries**: Fetch data from any table
- **Filtering**: WHERE clauses with operators (_eq, _gt, _like, etc.)
- **Sorting**: ORDER BY on any column
- **Pagination**: LIMIT and OFFSET
- **Aggregations**: COUNT, SUM, AVG, MAX, MIN
- **Nested queries**: Join related data (when relationships configured)

### Subscriptions
Real-time updates via WebSocket:
```graphql
subscription WatchNewEvents {
  events(order_by: {ledger: desc}, limit: 1) {
    id
    contract_id
    ledger
  }
}
```

### Examples

**Get token transfers for an address:**
```graphql
query GetUserTransfers($address: String!) {
  token_operations(
    where: {
      _or: [
        {from_address: {_eq: $address}},
        {to_address: {_eq: $address}}
      ]
    }
    order_by: {ledger: desc}
  ) {
    operation_type
    from_address
    to_address
    amount
    contract_id
    ledger_closed_at
  }
}
```

**Get token metadata with holder count:**
```graphql
query GetTokenInfo($contract_id: String!) {
  token_metadata(where: {contract_id: {_eq: $contract_id}}) {
    name
    symbol
    decimals
    token_balances_aggregate {
      aggregate {
        count
      }
    }
  }
}
```

**Aggregate stats:**
```graphql
query GetStats {
  events_aggregate {
    aggregate {
      count
    }
  }
  transactions_aggregate {
    aggregate {
      count
    }
  }
  token_operations_aggregate {
    aggregate {
      count
      sum {
        amount
      }
    }
  }
}
```

## Security Configuration

### Development (Default)
- No admin secret required
- Console accessible to everyone
- Anonymous read-only access to all tables
- Suitable for local development and testing

### Production (Recommended)
1. Set strong `HASURA_ADMIN_SECRET` in `.env`:
   ```bash
   HASURA_ADMIN_SECRET=your-super-secret-key-here
   ```

2. Access console with admin secret:
   - Add header: `x-hasura-admin-secret: your-super-secret-key-here`

3. Configure fine-grained permissions as needed in metadata

4. Use reverse proxy (nginx/Caddy) for SSL

## Integration

### REST API (cURL)
```bash
curl -X POST http://localhost:8081/v1/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ events(limit: 5) { id contract_id } }"}'
```

### JavaScript/TypeScript
```typescript
const response = await fetch('http://localhost:8081/v1/graphql', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    query: `query GetEvents($limit: Int!) {
      events(limit: $limit) { id contract_id ledger }
    }`,
    variables: { limit: 10 }
  })
});
const { data } = await response.json();
```

### Python
```python
import requests

query = """
  query {
    events(limit: 5) {
      id
      contract_id
      ledger
    }
  }
"""

response = requests.post(
    'http://localhost:8081/v1/graphql',
    json={'query': query}
)
data = response.json()
```

## Maintenance

### Viewing Logs
```bash
docker logs -f stellar-hasura
```

### Restarting Hasura
```bash
docker compose restart hasura
```

### Exporting Metadata (after making changes in console)
```bash
docker exec -it stellar-hasura hasura-cli metadata export
```

### Adding New Tables
When you add new tables to the indexer:
1. Tables will appear in Hasura Console automatically
2. Click "Track" next to the new table
3. Configure permissions if needed
4. Export metadata to save changes

## Troubleshooting

### Hasura won't start
```bash
# Check logs
docker logs stellar-hasura

# Verify PostgreSQL is running
docker compose ps postgres
```

### Tables not showing
```bash
# Restart Hasura
docker compose restart hasura

# Check metadata is mounted
docker inspect stellar-hasura | grep -A 5 Mounts
```

### Permission errors
- Verify `HASURA_ADMIN_SECRET` matches in requests
- Check anonymous role permissions in metadata files
- For admin access, add header: `x-hasura-admin-secret: your-secret`

## Resources

- Hasura Console: http://localhost:8081/console
- GraphQL API: http://localhost:8081/v1/graphql
- [Hasura Documentation](https://hasura.io/docs/latest/index/)
- [GraphQL Queries Guide](https://hasura.io/docs/latest/queries/postgres/index/)
- [Subscriptions Guide](https://hasura.io/docs/latest/subscriptions/postgres/index/)

---

**Need Help?** See [deploy/hasura/README.md](hasura/README.md) for detailed query examples and advanced usage.
