# Event Query Patterns

This document provides examples of common query patterns for the Soroban indexer, with a focus on efficient JSONB queries using GIN indexes.

## Table of Contents

- [Index Overview](#index-overview)
- [Helper Functions](#helper-functions)
- [Common Query Patterns](#common-query-patterns)
- [Performance Tips](#performance-tips)
- [GraphQL Examples](#graphql-examples)

## Index Overview

The indexer creates the following indexes for optimal query performance:

### JSONB GIN Indexes
- `idx_events_topic_gin` - GIN index on `events.topic` JSONB column
- `idx_events_value_gin` - GIN index on `events.value` JSONB column

### B-tree Indexes
- `idx_events_contract_id` - Regular index on `events.contract_id`
- `idx_events_ledger` - Regular index on `events.ledger`
- `idx_events_type` - Regular index on `events.type`
- `idx_transactions_ledger` - Regular index on `transactions.ledger`
- `idx_operations_tx_hash` - Regular index on `operations.tx_hash`
- `idx_operations_type` - Regular index on `operations.operation_type`
- `idx_token_operations_contract_id` - Regular index on `token_operations.contract_id`
- `idx_token_operations_from_address` - Regular index on `token_operations.from_address`
- `idx_token_operations_to_address` - Regular index on `token_operations.to_address`

## Helper Functions

The `pkg/models/event.go` file provides several helper functions for common query patterns:

### QueryEventsByTopicContains
Finds events where the topic JSONB contains a specific value.

```go
events, err := models.QueryEventsByTopicContains(db, `["transfer"]`)
```

### QueryEventsByTopicElement
Finds events where a specific element of the topic array matches a value.

```go
// Find events where topic[0] = "transfer"
events, err := models.QueryEventsByTopicElement(db, 0, "transfer")
```

### QueryEventsByContractAndTopic
Combines contract ID filtering with topic pattern matching (uses both indexes).

```go
events, err := models.QueryEventsByContractAndTopic(db, contractID, `["transfer"]`)
```

### QueryEventsByValuePath
Queries events by a specific path in the value JSONB.

```go
// Find events where value.amount = "1000"
events, err := models.QueryEventsByValuePath(db, "{amount}", `"1000"`)
```

### QueryTokenTransferEvents
Convenience function for finding token transfer events with pagination.

```go
events, err := models.QueryTokenTransferEvents(db, contractID, limit, offset)
```

### QueryEventsByLedgerRange
Queries events within a ledger range (uses ledger index).

```go
events, err := models.QueryEventsByLedgerRange(db, startLedger, endLedger)
```

## Common Query Patterns

### 1. Find All Transfer Events for a Token

**Go Code:**
```go
events, err := models.QueryTokenTransferEvents(db, contractID, 100, 0)
```

**Raw SQL:**
```sql
SELECT * FROM events
WHERE contract_id = 'CONTRACT_ID'
  AND topic @> '["transfer"]'
ORDER BY ledger DESC
LIMIT 100;
```

**GraphQL (Hasura):**
```graphql
query GetTokenTransfers($contractId: String!) {
  events(
    where: {
      contract_id: {_eq: $contractId}
      topic: {_contains: ["transfer"]}
    }
    order_by: {ledger: desc}
    limit: 100
  ) {
    id
    ledger
    topic
    value
    ledger_closed_at
  }
}
```

### 2. Find Events by Topic Array Element

**Go Code:**
```go
// Find all "mint" events (where topic[0] = "mint")
events, err := models.QueryEventsByTopicElement(db, 0, "mint")
```

**Raw SQL:**
```sql
SELECT * FROM events
WHERE topic ->> 0 = 'mint'
ORDER BY ledger DESC;
```

**GraphQL (Hasura):**
```graphql
query GetMintEvents {
  events(
    where: {
      topic: {_cast: {String: {_regex: "^\\[\"mint\""} } }
    }
    order_by: {ledger: desc}
  ) {
    id
    ledger
    topic
    value
  }
}
```

### 3. Find Events by Contract and Multiple Topic Conditions

**Raw SQL:**
```sql
-- Find transfer events from a specific address
SELECT * FROM events
WHERE contract_id = 'CONTRACT_ID'
  AND topic @> '["transfer"]'
  AND topic::text LIKE '%"from_address"%'
ORDER BY ledger DESC
LIMIT 100;
```

**GraphQL (Hasura):**
```graphql
query GetTransfersFromAddress($contractId: String!, $address: String!) {
  events(
    where: {
      contract_id: {_eq: $contractId}
      topic: {_contains: ["transfer", $address]}
    }
    order_by: {ledger: desc}
    limit: 100
  ) {
    id
    ledger
    topic
    value
    ledger_closed_at
  }
}
```

### 4. Find Events by Value JSON Path

**Go Code:**
```go
// Find events where value contains a specific amount
events, err := models.QueryEventsByValuePath(db, "{amount}", `"1000000"`)
```

**Raw SQL:**
```sql
SELECT * FROM events
WHERE value #> '{amount}' = '"1000000"'
ORDER BY ledger DESC;
```

### 5. Complex Multi-Condition Queries

**Raw SQL:**
```sql
-- Find recent high-value transfers
SELECT *
FROM events
WHERE contract_id = 'CONTRACT_ID'
  AND topic @> '["transfer"]'
  AND (value->>'amount')::bigint > 1000000
  AND ledger > 100000
ORDER BY ledger DESC
LIMIT 50;
```

**GraphQL (Hasura):**
```graphql
query GetHighValueTransfers($contractId: String!, $minLedger: Int!) {
  events(
    where: {
      contract_id: {_eq: $contractId}
      topic: {_contains: ["transfer"]}
      ledger: {_gt: $minLedger}
      value: {_cast: {String: {_regex: "\"amount\":[^,}]*[0-9]{7,}"}}}
    }
    order_by: {ledger: desc}
    limit: 50
  ) {
    id
    ledger
    topic
    value
    ledger_closed_at
  }
}
```

### 6. Find Events by Type and Ledger Range

**Go Code:**
```go
var events []models.Event
err := db.Where("type = ? AND ledger >= ? AND ledger <= ?", "contract", startLedger, endLedger).
    Order("ledger ASC").
    Find(&events).Error
```

**Raw SQL:**
```sql
SELECT * FROM events
WHERE type = 'contract'
  AND ledger BETWEEN 100000 AND 200000
ORDER BY ledger ASC;
```

**GraphQL (Hasura):**
```graphql
query GetEventsByTypeAndLedgerRange($startLedger: Int!, $endLedger: Int!) {
  events(
    where: {
      type: {_eq: "contract"}
      ledger: {_gte: $startLedger, _lte: $endLedger}
    }
    order_by: {ledger: asc}
  ) {
    id
    type
    ledger
    contract_id
    topic
    value
  }
}
```

### 7. Aggregate Queries

**Raw SQL:**
```sql
-- Count events by contract
SELECT contract_id, COUNT(*) as event_count
FROM events
WHERE ledger > 100000
GROUP BY contract_id
ORDER BY event_count DESC
LIMIT 10;
```

**GraphQL (Hasura):**
```graphql
query GetEventCountsByContract($minLedger: Int!) {
  events_aggregate(
    where: {ledger: {_gt: $minLedger}}
  ) {
    nodes {
      contract_id
    }
    aggregate {
      count
    }
  }
}
```

### 8. Join Events with Transactions

**Raw SQL:**
```sql
SELECT e.*, t.source_account, t.fee
FROM events e
JOIN transactions t ON e.tx_hash = t.hash
WHERE e.contract_id = 'CONTRACT_ID'
  AND t.ledger > 100000
ORDER BY e.ledger DESC
LIMIT 100;
```

**GraphQL (Hasura):**
```graphql
query GetEventsWithTransactions($contractId: String!, $minLedger: Int!) {
  events(
    where: {
      contract_id: {_eq: $contractId}
      transaction: {ledger: {_gt: $minLedger}}
    }
    order_by: {ledger: desc}
    limit: 100
  ) {
    id
    ledger
    topic
    value
    transaction {
      hash
      source_account
      fee
      ledger_closed_at
    }
  }
}
```

## Performance Tips

### 1. Use GIN Indexes for JSONB Queries

GIN indexes are specifically designed for JSONB and array operations. The following operators use GIN indexes:

- `@>` (contains)
- `?` (key exists)
- `?|` (any key exists)
- `?&` (all keys exist)

**Good (uses GIN index):**
```sql
WHERE topic @> '["transfer"]'
```

**Avoid (doesn't use GIN index):**
```sql
WHERE topic::text LIKE '%transfer%'
```

### 2. Combine Indexes for Multi-Condition Queries

When filtering by both regular columns and JSONB, put the regular column filter first:

**Good:**
```sql
WHERE contract_id = 'CONTRACT_ID'  -- Uses B-tree index first
  AND topic @> '["transfer"]'      -- Then uses GIN index
```

### 3. Use JSONB Operators Correctly

Different operators have different performance characteristics:

- `@>` (contains) - Fast with GIN index
- `->` (get JSON object) - Moderate performance
- `->>` (get JSON text) - Moderate performance
- `#>` (get JSON at path) - Moderate performance
- `#>>` (get JSON text at path) - Moderate performance

### 4. Limit Result Sets

Always use `LIMIT` and `OFFSET` for pagination:

```sql
SELECT * FROM events
WHERE contract_id = 'CONTRACT_ID'
ORDER BY ledger DESC
LIMIT 100 OFFSET 0;
```

### 5. Use Covering Indexes When Possible

If you only need specific columns, the query planner may use index-only scans:

```sql
-- If you only need id and ledger, this can use index-only scan
SELECT id, ledger FROM events
WHERE ledger > 100000;
```

### 6. Analyze Query Performance

Use `EXPLAIN ANALYZE` to understand query performance:

```sql
EXPLAIN ANALYZE
SELECT * FROM events
WHERE contract_id = 'CONTRACT_ID'
  AND topic @> '["transfer"]'
LIMIT 100;
```

Look for:
- **Index Scan** - Good (using index)
- **Seq Scan** - Bad (full table scan, add index)
- **Bitmap Index Scan** - Good (combining multiple indexes)

## GraphQL Examples

### Hasura Relationships

Set up relationships in Hasura for efficient joins:

**events → transaction:**
```json
{
  "name": "transaction",
  "using": {
    "foreign_key_constraint_on": "tx_hash"
  }
}
```

**events → token_operations:**
```json
{
  "name": "token_operations",
  "using": {
    "manual_configuration": {
      "remote_table": "token_operations",
      "column_mapping": {
        "id": "event_id"
      }
    }
  }
}
```

### Real-time Subscriptions

```graphql
subscription WatchNewTransfers($contractId: String!) {
  events(
    where: {
      contract_id: {_eq: $contractId}
      topic: {_contains: ["transfer"]}
    }
    order_by: {ledger: desc}
    limit: 10
  ) {
    id
    ledger
    topic
    value
    ledger_closed_at
  }
}
```

### Aggregations

```graphql
query GetTokenStats($contractId: String!) {
  events_aggregate(
    where: {
      contract_id: {_eq: $contractId}
      topic: {_contains: ["transfer"]}
    }
  ) {
    aggregate {
      count
      max {
        ledger
      }
      min {
        ledger
      }
    }
  }
}
```

## Advanced Patterns

### Full-Text Search on JSONB

```sql
-- Search for events containing a specific address in topic or value
SELECT * FROM events
WHERE (topic::text || value::text) ILIKE '%ADDRESS%'
LIMIT 100;
```

### JSON Array Expansion

```sql
-- Expand topic array elements
SELECT id, ledger, jsonb_array_elements_text(topic) as topic_element
FROM events
WHERE contract_id = 'CONTRACT_ID'
LIMIT 100;
```

### Complex JSONB Path Queries

```sql
-- Find events with nested JSON structures
SELECT * FROM events
WHERE value @> '{"transfer": {"from": "ADDRESS"}}'
ORDER BY ledger DESC;
```

## References

- [PostgreSQL JSONB Documentation](https://www.postgresql.org/docs/current/datatype-json.html)
- [PostgreSQL GIN Indexes](https://www.postgresql.org/docs/current/gin-intro.html)
- [Hasura GraphQL Documentation](https://hasura.io/docs/latest/index/)
- [JSONB Performance Tips](https://www.postgresql.org/docs/current/datatype-json.html#JSON-INDEXING)

---

**Last Updated**: 2025-10-11
**Author**: Claude Code
