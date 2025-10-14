#!/bin/bash

# Check token_metadata table for data
echo "=== Checking token_metadata table ==="
docker exec stellar-db psql -U stellar -d stellar -c "
SELECT
    contract_id,
    name,
    symbol,
    decimal,
    admin_address,
    created_at
FROM token_metadata
ORDER BY created_at DESC
LIMIT 10;
"

echo ""
echo "=== Token metadata count ==="
docker exec stellar-db psql -U stellar -d stellar -c "
SELECT COUNT(*) as total_tokens FROM token_metadata;
"

echo ""
echo "=== Recent contract data entries ==="
docker exec stellar-db psql -U stellar -d stellar -c "
SELECT
    contract_id,
    key::text as key_preview,
    LEFT(val::text, 100) as val_preview,
    durability
FROM contract_data
WHERE key::text LIKE '%ScvLedgerKeyContractInstance%'
ORDER BY created_at DESC
LIMIT 5;
"

echo ""
echo "=== Sample of any contract data entries ==="
docker exec stellar-db psql -U stellar -d stellar -c "
SELECT
    COUNT(*) as total_entries,
    COUNT(DISTINCT contract_id) as unique_contracts
FROM contract_data;
"
