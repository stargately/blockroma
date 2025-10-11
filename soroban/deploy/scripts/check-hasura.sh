#!/bin/bash

# Check Hasura GraphQL Engine health and connectivity

set -e

HASURA_URL="${HASURA_URL:-http://localhost:8081}"

echo "Checking Hasura GraphQL Engine..."
echo "URL: $HASURA_URL"
echo ""

# Check health endpoint
echo "1. Health Check:"
HEALTH=$(curl -s "${HASURA_URL}/healthz" || echo "FAILED")
if [ "$HEALTH" == "OK" ]; then
    echo "   ✅ Hasura is healthy"
else
    echo "   ❌ Hasura health check failed"
    echo "   Response: $HEALTH"
    exit 1
fi

# Check GraphQL endpoint with introspection query
echo ""
echo "2. GraphQL Endpoint Check:"
INTROSPECTION=$(curl -s -X POST "${HASURA_URL}/v1/graphql" \
    -H "Content-Type: application/json" \
    -d '{"query": "{ __schema { queryType { name } } }"}' || echo "FAILED")

if echo "$INTROSPECTION" | grep -q "queryType"; then
    echo "   ✅ GraphQL endpoint is responding"
else
    echo "   ❌ GraphQL endpoint check failed"
    echo "   Response: $INTROSPECTION"
    exit 1
fi

# Check if tables are tracked
echo ""
echo "3. Tracked Tables Check:"
TABLES_QUERY='{"query": "{ events_aggregate { aggregate { count } } transactions_aggregate { aggregate { count } } }"}'
TABLES=$(curl -s -X POST "${HASURA_URL}/v1/graphql" \
    -H "Content-Type: application/json" \
    -d "$TABLES_QUERY" || echo "FAILED")

if echo "$TABLES" | grep -q "aggregate"; then
    echo "   ✅ Tables are tracked and queryable"

    # Extract counts if available
    EVENTS_COUNT=$(echo "$TABLES" | grep -o '"count":[0-9]*' | head -1 | cut -d':' -f2)
    TXS_COUNT=$(echo "$TABLES" | grep -o '"count":[0-9]*' | tail -1 | cut -d':' -f2)

    if [ -n "$EVENTS_COUNT" ] && [ -n "$TXS_COUNT" ]; then
        echo ""
        echo "   📊 Current Data:"
        echo "      Events: $EVENTS_COUNT"
        echo "      Transactions: $TXS_COUNT"
    fi
else
    echo "   ❌ Tables query failed"
    echo "   Response: $TABLES"
    exit 1
fi

echo ""
echo "✅ All checks passed!"
echo ""
echo "Access Hasura Console: ${HASURA_URL}/console"
echo "GraphQL API Endpoint: ${HASURA_URL}/v1/graphql"
