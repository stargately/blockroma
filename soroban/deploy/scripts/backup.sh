#!/bin/bash
set -e

BACKUP_DIR="${BACKUP_DIR:-./backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="stellar-rpc-backup-${TIMESTAMP}"

echo "💾 Creating backup: ${BACKUP_NAME}"

cd "$(dirname "$0")/.."

# Create backup directory
mkdir -p "${BACKUP_DIR}"

# Stop the service to ensure consistent backup
echo "⏸️  Stopping Stellar RPC..."
docker compose stop stellar-rpc

# Create tarball of data directory
echo "📦 Compressing data..."
tar -czf "${BACKUP_DIR}/${BACKUP_NAME}.tar.gz" \
    -C data \
    stellar-rpc captive-core logs

# Restart service
echo "▶️  Restarting Stellar RPC..."
docker compose start stellar-rpc

# Get backup size
BACKUP_SIZE=$(du -h "${BACKUP_DIR}/${BACKUP_NAME}.tar.gz" | cut -f1)

echo "✅ Backup completed!"
echo "   Location: ${BACKUP_DIR}/${BACKUP_NAME}.tar.gz"
echo "   Size: ${BACKUP_SIZE}"

# Optional: Clean up old backups (keep last 7 days)
echo "🧹 Cleaning up old backups (keeping last 7 days)..."
find "${BACKUP_DIR}" -name "stellar-rpc-backup-*.tar.gz" -mtime +7 -delete

echo "✅ Done!"
