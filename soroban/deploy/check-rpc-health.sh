#!/bin/bash
# Script to check stellar-rpc health status

RPC_URL="${1:-http://localhost:8000}"

echo "Checking stellar-rpc health at: $RPC_URL"
echo ""

# Make JSON-RPC getHealth request
RESPONSE=$(curl -s -X POST "$RPC_URL" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"getHealth"}')

echo "Response:"
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
echo ""

# Check if healthy
if echo "$RESPONSE" | grep -q '"status":"healthy"'; then
  echo "✅ Status: HEALTHY"
  exit 0
else
  echo "❌ Status: UNHEALTHY or ERROR"
  exit 1
fi
