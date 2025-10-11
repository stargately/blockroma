# Hasura GraphQL Engine Setup

This directory contains Hasura metadata configuration for exposing the Stellar indexer database via GraphQL API.

## Quick Start

Hasura is automatically started with the docker-compose stack:

```bash
cd deploy
docker compose up -d
```

Access the Hasura Console at: **http://localhost:8081**

## Features

- **Auto-tracked Tables**: All 13 indexed tables are pre-configured
- **GraphQL API**: Query any indexed data via GraphQL
- **GraphQL Subscriptions**: Real-time updates via WebSocket
- **Anonymous Access**: Public read-only access enabled by default
- **Admin Console**: Web UI for exploring schema and testing queries

## Exposed Tables

All tables from the indexer are available via GraphQL:

### Smart Contract Data
- **events** - Contract events with decoded topics/values
- **transactions** - Full transaction details
- **contract_data_entries** - Contract storage data

### Token Data
- **token_operations** - Token transfers, mints, burns
- **token_metadata** - Token name, symbol, decimals
- **token_balances** - Token holder balances

### Ledger State
- **account_entries** - Stellar account data
- **trust_line_entries** - Asset trustlines
- **offer_entries** - DEX offers
- **liquidity_pool_entries** - AMM pools
- **claimable_balance_entries** - Claimable balances
- **data_entries** - Account data entries

### Indexer State
- **cursor** - Current indexing position

## Example Queries

### Get Latest Events
```graphql
query GetLatestEvents {
  events(order_by: {ledger: desc}, limit: 10) {
    id
    contract_id
    type
    ledger
    ledger_closed_at
    topic
    value
  }
}
```

### Get Token Transfers
```graphql
query GetTokenTransfers($contract_id: String!) {
  token_operations(
    where: {
      contract_id: {_eq: $contract_id}
      operation_type: {_eq: "transfer"}
    }
    order_by: {ledger: desc}
    limit: 20
  ) {
    id
    from_address
    to_address
    amount
    ledger_closed_at
  }
}
```

### Get Account Balances
```graphql
query GetAccountBalances($address: String!) {
  token_balances(where: {address: {_eq: $address}}) {
    contract_id
    balance
    token_metadata {
      name
      symbol
      decimals
    }
  }
}
```

### Subscribe to New Events (Real-time)
```graphql
subscription NewEvents {
  events(order_by: {ledger: desc}, limit: 1) {
    id
    contract_id
    type
    ledger
    topic
    value
  }
}
```

### Aggregate Queries
```graphql
query TokenStats($contract_id: String!) {
  token_operations_aggregate(
    where: {contract_id: {_eq: $contract_id}}
  ) {
    aggregate {
      count
      sum {
        amount
      }
    }
  }
}
```

## Configuration

### Environment Variables

Set in `deploy/.env`:

```bash
# Hasura port (default: 8081)
HASURA_PORT=8081

# Admin secret for protected endpoints (optional, recommended for production)
HASURA_ADMIN_SECRET=your_secret_here

# Default role for unauthenticated users (default: anonymous)
HASURA_UNAUTHORIZED_ROLE=anonymous
```

### Access Control

**Default Setup (Development)**:
- Anonymous users: Read-only access to all tables
- No admin secret: Console accessible without authentication

**Production Recommendations**:
1. Set strong `HASURA_ADMIN_SECRET` in `.env`
2. Access console with `x-hasura-admin-secret` header
3. Configure fine-grained permissions per table/role as needed

## Metadata Management

Hasura metadata is stored in `./metadata/` and version-controlled.

### Export Metadata
```bash
# From inside Hasura container
docker exec -it stellar-hasura hasura-cli metadata export
```

### Apply Metadata
```bash
# Metadata is auto-applied on container startup from mounted volume
docker compose restart hasura
```

### Manual Tracking

If you add new tables to the indexer, track them via Console:
1. Go to http://localhost:8081/console/data
2. Click "Track" next to the new table
3. Export metadata to save changes

## GraphQL Endpoints

Once running, Hasura exposes:

- **GraphQL API**: `http://localhost:8081/v1/graphql`
- **GraphQL Console**: `http://localhost:8081/console`
- **Health Check**: `http://localhost:8081/healthz`

## Integration Examples

### cURL
```bash
curl -X POST http://localhost:8081/v1/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ events(limit: 5) { id contract_id ledger } }"
  }'
```

### JavaScript/TypeScript
```typescript
const response = await fetch('http://localhost:8081/v1/graphql', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    query: `
      query GetEvents($limit: Int!) {
        events(limit: $limit, order_by: {ledger: desc}) {
          id
          contract_id
          ledger
        }
      }
    `,
    variables: { limit: 10 }
  })
});

const data = await response.json();
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

## Troubleshooting

### Hasura won't start
```bash
# Check logs
docker logs stellar-hasura

# Verify PostgreSQL is healthy
docker compose ps postgres
```

### Tables not showing up
```bash
# Restart Hasura to re-apply metadata
docker compose restart hasura

# Check metadata volume is mounted correctly
docker inspect stellar-hasura | grep -A 5 Mounts
```

### Permission denied errors
- Check `HASURA_ADMIN_SECRET` is set correctly if protected
- Verify anonymous role permissions in metadata files
- Try with admin secret header: `x-hasura-admin-secret: your_secret`

## Advanced Features

### Custom Views
Create SQL views for complex queries:
```sql
CREATE VIEW token_transfer_summary AS
SELECT
  contract_id,
  COUNT(*) as transfer_count,
  SUM(amount) as total_volume
FROM token_operations
WHERE operation_type = 'transfer'
GROUP BY contract_id;
```

Then track the view in Hasura Console.

### Computed Fields
Add virtual fields by creating SQL functions and exposing them via Hasura.

### Actions & Events
Configure Hasura Actions for custom business logic or Event Triggers for webhooks on data changes.

## Resources

- [Hasura Documentation](https://hasura.io/docs/latest/index/)
- [GraphQL Queries](https://hasura.io/docs/latest/queries/postgres/index/)
- [Subscriptions](https://hasura.io/docs/latest/subscriptions/postgres/index/)
- [Authorization](https://hasura.io/docs/latest/auth/authorization/index/)
