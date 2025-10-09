#!/bin/bash
set -e

echo "🚀 Starting Stellar RPC (Mainnet)..."

cd "$(dirname "$0")/.."

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Copying from .env.example..."
    cp .env.example .env
    echo "⚠️  Please edit .env with your configuration before starting!"
    exit 1
fi

# Load environment variables
source .env

# Create data directories if they don't exist
mkdir -p data/{stellar-rpc,captive-core,logs}

# Pull latest image
echo "📦 Pulling latest Stellar RPC image..."
docker compose pull

# Start services
echo "🔄 Starting services..."
docker compose up -d

# Wait for health check
echo "⏳ Waiting for RPC to be healthy..."
for i in {1..30}; do
    if docker compose exec -T stellar-rpc curl -sf http://localhost:8000/health > /dev/null 2>&1; then
        echo "✅ Stellar RPC is healthy!"
        echo ""
        echo "🌐 RPC Endpoint: http://localhost:${RPC_PORT:-8000}"
        echo "📊 Admin Endpoint: http://localhost:${ADMIN_PORT:-6061}"
        echo ""
        echo "View logs: docker compose logs -f stellar-rpc"
        exit 0
    fi
    echo "   Attempt $i/30..."
    sleep 10
done

echo "❌ Health check timeout. Check logs with: docker compose logs stellar-rpc"
exit 1
