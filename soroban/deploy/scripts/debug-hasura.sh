#!/bin/sh

# Debug Hasura configuration

set -e

HASURA_URL="${HASURA_URL:-http://localhost:8081}"

echo "=== Hasura Debug Information ==="
echo ""

echo "1. Checking Hasura Health:"
curl -s "${HASURA_URL}/healthz"
echo ""
echo ""

echo "2. Checking Database Connection:"
RESULT=$(curl -s -X POST "${HASURA_URL}/v1/metadata" \
    -H "Content-Type: application/json" \
    -d '{
        "type": "get_source_tables",
        "args": {
            "source": "default"
        }
    }')
echo "$RESULT" | python3 -m json.tool 2>/dev/null || echo "$RESULT"
echo ""

echo "3. Checking Tracked Tables:"
RESULT=$(curl -s -X POST "${HASURA_URL}/v1/metadata" \
    -H "Content-Type: application/json" \
    -d '{
        "type": "export_metadata",
        "args": {}
    }')
echo "$RESULT" | python3 -m json.tool 2>/dev/null | grep -A 5 "tables" | head -30 || echo "$RESULT"
echo ""

echo "4. Testing GraphQL Query (as anonymous):"
RESULT=$(curl -s -X POST "${HASURA_URL}/v1/graphql" \
    -H "Content-Type: application/json" \
    -d '{
        "query": "{ __schema { queryType { name } } }"
    }')
echo "$RESULT" | python3 -m json.tool 2>/dev/null || echo "$RESULT"
echo ""

echo "5. Testing Simple Table Query:"
RESULT=$(curl -s -X POST "${HASURA_URL}/v1/graphql" \
    -H "Content-Type: application/json" \
    -d '{
        "query": "{ cursor { id last_ledger } }"
    }')
echo "$RESULT" | python3 -m json.tool 2>/dev/null || echo "$RESULT"
echo ""

echo "=== End Debug ==="
