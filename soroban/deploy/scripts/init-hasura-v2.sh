#!/bin/sh

# Initialize Hasura - Complete setup including database source
# This script ensures the database is connected and all tables are tracked

set -e

# Load environment variables from .env file if it exists
if [ -f .env ]; then
    echo "📄 Loading configuration from .env file..."
    # Export variables from .env (handles comments and empty lines)
    export $(grep -v '^#' .env | grep -v '^$' | xargs)
fi

HASURA_URL="${HASURA_URL:-http://localhost:8081}"
HASURA_ADMIN_SECRET="${HASURA_ADMIN_SECRET:-}"
MAX_RETRIES=30
RETRY_DELAY=2

# Build auth header if admin secret is set
if [ -n "$HASURA_ADMIN_SECRET" ]; then
    AUTH_HEADER="x-hasura-admin-secret: $HASURA_ADMIN_SECRET"
    echo "✅ Using admin secret for authentication"
else
    AUTH_HEADER=""
    echo "ℹ️  No admin secret configured (dev mode)"
fi

echo ""
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
echo "Step 1: Ensuring database source is connected..."

# Get database URL from environment (same as Hasura container)
DB_URL="${HASURA_GRAPHQL_DATABASE_URL:-postgresql://stellar:password@postgres:5432/stellar_indexer}"

# Try to add the database source (will fail if it already exists, which is fine)
if [ -n "$AUTH_HEADER" ]; then
    curl -s -X POST "${HASURA_URL}/v1/metadata" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -d "{
            \"type\": \"pg_add_source\",
            \"args\": {
                \"name\": \"default\",
                \"configuration\": {
                    \"connection_info\": {
                        \"database_url\": \"$DB_URL\",
                        \"pool_settings\": {
                            \"max_connections\": 50,
                            \"idle_timeout\": 180,
                            \"retries\": 1
                        }
                    }
                }
            }
        }" > /dev/null 2>&1 && echo "  ✅ Database source added" || echo "  ℹ️  Database source already exists"
else
    curl -s -X POST "${HASURA_URL}/v1/metadata" \
        -H "Content-Type: application/json" \
        -d "{
            \"type\": \"pg_add_source\",
            \"args\": {
                \"name\": \"default\",
                \"configuration\": {
                    \"connection_info\": {
                        \"database_url\": \"$DB_URL\",
                        \"pool_settings\": {
                            \"max_connections\": 50,
                            \"idle_timeout\": 180,
                            \"retries\": 1
                        }
                    }
                }
            }
        }" > /dev/null 2>&1 && echo "  ✅ Database source added" || echo "  ℹ️  Database source already exists"
fi

echo ""
echo "Step 2: Reloading metadata to ensure fresh state..."

if [ -n "$AUTH_HEADER" ]; then
    curl -s -X POST "${HASURA_URL}/v1/metadata" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -d '{"type": "reload_metadata", "args": {}}' > /dev/null 2>&1
else
    curl -s -X POST "${HASURA_URL}/v1/metadata" \
        -H "Content-Type: application/json" \
        -d '{"type": "reload_metadata", "args": {}}' > /dev/null 2>&1
fi
echo "  ✅ Metadata reloaded"

echo ""
echo "Step 3: Tracking all database tables..."

# Function to track a table
track_table() {
    table=$1
    echo "  Tracking: $table"

    if [ -n "$AUTH_HEADER" ]; then
        RESPONSE=$(curl -s -X POST "${HASURA_URL}/v1/metadata" \
            -H "Content-Type: application/json" \
            -H "$AUTH_HEADER" \
            -d "{
                \"type\": \"pg_track_table\",
                \"args\": {
                    \"source\": \"default\",
                    \"table\": {
                        \"schema\": \"public\",
                        \"name\": \"$table\"
                    }
                }
            }")
    else
        RESPONSE=$(curl -s -X POST "${HASURA_URL}/v1/metadata" \
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
            }")
    fi

    if echo "$RESPONSE" | grep -q "error"; then
        echo "    ℹ️  Already tracked or table doesn't exist"
    else
        echo "    ✅ Tracked successfully"
    fi
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
echo "Step 4: Setting up permissions for 'anonymous' role..."

# Function to grant permissions
grant_permissions() {
    table=$1
    echo "  Permissions: $table"

    if [ -n "$AUTH_HEADER" ]; then
        RESPONSE=$(curl -s -X POST "${HASURA_URL}/v1/metadata" \
            -H "Content-Type: application/json" \
            -H "$AUTH_HEADER" \
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
            }")
    else
        RESPONSE=$(curl -s -X POST "${HASURA_URL}/v1/metadata" \
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
            }")
    fi

    if echo "$RESPONSE" | grep -q "already exists"; then
        echo "    ℹ️  Permission already exists"
    elif echo "$RESPONSE" | grep -q "error"; then
        echo "    ⚠️  Error: $RESPONSE"
    else
        echo "    ✅ Permission granted"
    fi
}

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
echo "Step 5: Verifying setup with test query..."

TEST_RESPONSE=$(curl -s -X POST "${HASURA_URL}/v1/graphql" \
    -H "Content-Type: application/json" \
    -d '{
        "query": "{ __schema { queryType { name } } }"
    }')

if echo "$TEST_RESPONSE" | grep -q "queryType"; then
    echo "  ✅ GraphQL endpoint is working!"
else
    echo "  ⚠️  Warning: GraphQL test failed"
    echo "  Response: $TEST_RESPONSE"
fi

echo ""
echo "═══════════════════════════════════════"
echo "✅ Hasura initialization complete!"
echo "═══════════════════════════════════════"
echo ""
echo "📊 Access Points:"
echo "   Console:  ${HASURA_URL}/console"
echo "   GraphQL:  ${HASURA_URL}/v1/graphql"
echo "   Health:   ${HASURA_URL}/healthz"
echo ""
echo "🧪 Test Query:"
echo "   curl -X POST ${HASURA_URL}/v1/graphql \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"query\": \"{ cursor { id last_ledger } }\"}'"
echo ""
