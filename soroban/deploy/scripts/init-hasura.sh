#!/bin/sh

# Initialize Hasura by tracking all tables
# This script waits for Hasura to be ready, then tracks all database tables

set -e

HASURA_URL="${HASURA_URL:-http://localhost:8081}"
MAX_RETRIES=30
RETRY_DELAY=2

echo "Waiting for Hasura to be ready..."

# Wait for Hasura to be healthy
i=1
while [ $i -le $MAX_RETRIES ]; do
    if curl -s "${HASURA_URL}/healthz" | grep -q "OK"; then
        echo "✅ Hasura is ready!"
        break
    fi

    if [ $i -eq $MAX_RETRIES ]; then
        echo "❌ Hasura failed to start after ${MAX_RETRIES} attempts"
        exit 1
    fi

    echo "Waiting for Hasura... (attempt $i/$MAX_RETRIES)"
    sleep $RETRY_DELAY
    i=$((i + 1))
done

echo ""
echo "Tracking database tables..."

# Function to track a table
track_table() {
    table=$1
    echo "Tracking table: $table"

    curl -s -X POST "${HASURA_URL}/v1/metadata" \
        -H "Content-Type: application/json" \
        -d "{
            \"type\": \"pg_track_table\",
            \"args\": {
                \"source\": \"default\",
                \"table\": {
                    \"schema\": \"public\",
                    \"name\": \"$table\"
                }
            }
        }" > /dev/null 2>&1 || echo "  (table may already be tracked or doesn't exist yet)"
}

# Function to grant permissions
grant_permissions() {
    table=$1
    echo "Granting permissions on: $table"

    curl -s -X POST "${HASURA_URL}/v1/metadata" \
        -H "Content-Type: application/json" \
        -d "{
            \"type\": \"pg_create_select_permission\",
            \"args\": {
                \"source\": \"default\",
                \"table\": {
                    \"schema\": \"public\",
                    \"name\": \"$table\"
                },
                \"role\": \"anonymous\",
                \"permission\": {
                    \"columns\": \"*\",
                    \"filter\": {},
                    \"allow_aggregations\": true
                }
            }
        }" > /dev/null 2>&1 || echo "  (permission may already exist)"
}

# Track all tables
track_table "events"
track_table "transactions"
track_table "token_operations"
track_table "token_metadata"
track_table "token_balances"
track_table "account_entries"
track_table "trust_line_entries"
track_table "offer_entries"
track_table "liquidity_pool_entries"
track_table "claimable_balance_entries"
track_table "contract_data_entries"
track_table "data_entries"
track_table "cursor"

echo ""
echo "Setting up anonymous role permissions..."

# Grant permissions on all tables
grant_permissions "events"
grant_permissions "transactions"
grant_permissions "token_operations"
grant_permissions "token_metadata"
grant_permissions "token_balances"
grant_permissions "account_entries"
grant_permissions "trust_line_entries"
grant_permissions "offer_entries"
grant_permissions "liquidity_pool_entries"
grant_permissions "claimable_balance_entries"
grant_permissions "contract_data_entries"
grant_permissions "data_entries"
grant_permissions "cursor"

echo ""
echo "✅ Hasura initialization complete!"
echo ""
echo "Access Hasura Console: ${HASURA_URL}/console"
echo "GraphQL Endpoint: ${HASURA_URL}/v1/graphql"
