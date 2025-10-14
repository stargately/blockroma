#!/bin/bash

# Local development script
# Runs indexer locally connecting to remote services at 85.239.232.74

set -e

echo "=== Starting Local Indexer ==="
echo ""

# Load environment variables from .env.local
export STELLAR_RPC_URL="http://85.239.232.74:8000"
export POSTGRES_DSN="postgresql://stellar:StellaBiliBilirIndexer2024!ChangeMe@85.239.232.74:5432/stellar_indexer?sslmode=disable"
export INDEXER_PORT="8080"
export LOG_LEVEL="debug"

echo "Environment:"
echo "  STELLAR_RPC_URL=$STELLAR_RPC_URL"
echo "  POSTGRES_DSN=$POSTGRES_DSN"
echo "  INDEXER_PORT=$INDEXER_PORT"
echo "  LOG_LEVEL=$LOG_LEVEL"
echo ""

echo "=== Starting Indexer ==="
echo "Press Ctrl+C to stop"
echo ""

cd cmd/indexer
go run *.go
