package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupContractCodeTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the ContractCode table
	err = db.AutoMigrate(&ContractCode{})
	require.NoError(t, err)

	return db
}

func TestContractCodeTableName(t *testing.T) {
	cc := ContractCode{}
	assert.Equal(t, "contract_code", cc.TableName())
}

func TestUpsertContractCode_Insert(t *testing.T) {
	db := setupContractCodeTestDB(t)

	wasm := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	code := &ContractCode{
		Hash:       "abc123def456",
		Wasm:       wasm,
		DeployedAt: time.Now().UTC(),
		Ledger:     12345,
		TxHash:     "tx-hash-123",
		SizeBytes:  len(wasm),
	}

	err := UpsertContractCode(db, code)
	require.NoError(t, err)

	// Verify the code was inserted
	retrieved, err := GetContractCodeByHash(db, "abc123def456")
	require.NoError(t, err)
	assert.Equal(t, code.Hash, retrieved.Hash)
	assert.Equal(t, code.Wasm, retrieved.Wasm)
	assert.Equal(t, code.Ledger, retrieved.Ledger)
	assert.Equal(t, code.TxHash, retrieved.TxHash)
	assert.Equal(t, code.SizeBytes, retrieved.SizeBytes)
}

func TestUpsertContractCode_Idempotent(t *testing.T) {
	db := setupContractCodeTestDB(t)

	wasm := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	code1 := &ContractCode{
		Hash:       "abc123",
		Wasm:       wasm,
		DeployedAt: time.Now().UTC(),
		Ledger:     12345,
		TxHash:     "tx-hash-1",
		SizeBytes:  len(wasm),
	}

	// First insert
	err := UpsertContractCode(db, code1)
	require.NoError(t, err)

	// Second insert with same hash but different metadata
	code2 := &ContractCode{
		Hash:       "abc123", // Same hash
		Wasm:       wasm,
		DeployedAt: time.Now().UTC().Add(time.Hour),
		Ledger:     12346, // Different ledger
		TxHash:     "tx-hash-2", // Different tx
		SizeBytes:  len(wasm),
	}

	err = UpsertContractCode(db, code2)
	require.NoError(t, err)

	// Verify that the original code is unchanged
	retrieved, err := GetContractCodeByHash(db, "abc123")
	require.NoError(t, err)
	assert.Equal(t, code1.Hash, retrieved.Hash)
	assert.Equal(t, code1.Ledger, retrieved.Ledger)
	assert.Equal(t, code1.TxHash, retrieved.TxHash)
	// Should still have the first deployment metadata
	assert.NotEqual(t, code2.Ledger, retrieved.Ledger)
	assert.NotEqual(t, code2.TxHash, retrieved.TxHash)
}

func TestUpsertContractCode_MultipleDifferentCodes(t *testing.T) {
	db := setupContractCodeTestDB(t)

	wasm1 := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	wasm2 := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x01}

	code1 := &ContractCode{
		Hash:       "hash1",
		Wasm:       wasm1,
		DeployedAt: time.Now().UTC(),
		Ledger:     100,
		TxHash:     "tx1",
		SizeBytes:  len(wasm1),
	}

	code2 := &ContractCode{
		Hash:       "hash2",
		Wasm:       wasm2,
		DeployedAt: time.Now().UTC(),
		Ledger:     200,
		TxHash:     "tx2",
		SizeBytes:  len(wasm2),
	}

	// Insert both codes
	err := UpsertContractCode(db, code1)
	require.NoError(t, err)

	err = UpsertContractCode(db, code2)
	require.NoError(t, err)

	// Verify both codes are stored
	retrieved1, err := GetContractCodeByHash(db, "hash1")
	require.NoError(t, err)
	assert.Equal(t, code1.Hash, retrieved1.Hash)

	retrieved2, err := GetContractCodeByHash(db, "hash2")
	require.NoError(t, err)
	assert.Equal(t, code2.Hash, retrieved2.Hash)

	// Verify they are different
	assert.NotEqual(t, retrieved1.Wasm, retrieved2.Wasm)
}

func TestGetContractCodeByHash_NotFound(t *testing.T) {
	db := setupContractCodeTestDB(t)

	code, err := GetContractCodeByHash(db, "non-existent-hash")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Nil(t, code)
}

func TestGetAllContractCodes(t *testing.T) {
	db := setupContractCodeTestDB(t)

	// Insert multiple codes
	baseTime := time.Now().UTC()
	for i := 0; i < 5; i++ {
		code := &ContractCode{
			Hash:       string(rune('a' + i)) + "hash",
			Wasm:       []byte{byte(i)},
			DeployedAt: baseTime.Add(time.Duration(i) * time.Minute),
			Ledger:     uint32(100 + i),
			TxHash:     "tx" + string(rune('a'+i)),
			SizeBytes:  1,
		}
		err := UpsertContractCode(db, code)
		require.NoError(t, err)
	}

	// Get all codes with limit
	codes, err := GetAllContractCodes(db, 3, 0)
	require.NoError(t, err)
	assert.Len(t, codes, 3)

	// Verify they are ordered by deployed_at DESC
	assert.True(t, codes[0].DeployedAt.After(codes[1].DeployedAt) || codes[0].DeployedAt.Equal(codes[1].DeployedAt))
	assert.True(t, codes[1].DeployedAt.After(codes[2].DeployedAt) || codes[1].DeployedAt.Equal(codes[2].DeployedAt))
}

func TestGetAllContractCodes_WithOffset(t *testing.T) {
	db := setupContractCodeTestDB(t)

	// Insert multiple codes
	baseTime := time.Now().UTC()
	for i := 0; i < 5; i++ {
		code := &ContractCode{
			Hash:       string(rune('a' + i)) + "hash",
			Wasm:       []byte{byte(i)},
			DeployedAt: baseTime.Add(time.Duration(i) * time.Minute),
			Ledger:     uint32(100 + i),
			TxHash:     "tx" + string(rune('a'+i)),
			SizeBytes:  1,
		}
		err := UpsertContractCode(db, code)
		require.NoError(t, err)
	}

	// Get codes with offset
	codes, err := GetAllContractCodes(db, 2, 2)
	require.NoError(t, err)
	assert.Len(t, codes, 2)
}

func TestGetAllContractCodes_Empty(t *testing.T) {
	db := setupContractCodeTestDB(t)

	codes, err := GetAllContractCodes(db, 10, 0)
	require.NoError(t, err)
	assert.Empty(t, codes)
}

func TestContractCode_BinaryWasm(t *testing.T) {
	db := setupContractCodeTestDB(t)

	// Test with realistic WASM magic number and some random bytes
	wasm := []byte{
		0x00, 0x61, 0x73, 0x6d, // WASM magic number
		0x01, 0x00, 0x00, 0x00, // Version
		0x01, 0x05, 0x01, 0x60, 0x00, 0x01, 0x7f, // Some WASM sections
		0xff, 0xfe, 0xfd, 0xfc, // Binary data
	}

	code := &ContractCode{
		Hash:       "binary-wasm-test",
		Wasm:       wasm,
		DeployedAt: time.Now().UTC(),
		Ledger:     12345,
		TxHash:     "tx-binary",
		SizeBytes:  len(wasm),
	}

	err := UpsertContractCode(db, code)
	require.NoError(t, err)

	// Retrieve and verify binary data is intact
	retrieved, err := GetContractCodeByHash(db, "binary-wasm-test")
	require.NoError(t, err)
	assert.Equal(t, wasm, retrieved.Wasm)
	assert.Equal(t, len(wasm), retrieved.SizeBytes)
}

func TestContractCode_LargeWasm(t *testing.T) {
	db := setupContractCodeTestDB(t)

	// Create a large WASM (1MB)
	largeWasm := make([]byte, 1024*1024)
	for i := range largeWasm {
		largeWasm[i] = byte(i % 256)
	}

	code := &ContractCode{
		Hash:       "large-wasm-test",
		Wasm:       largeWasm,
		DeployedAt: time.Now().UTC(),
		Ledger:     99999,
		TxHash:     "tx-large",
		SizeBytes:  len(largeWasm),
	}

	err := UpsertContractCode(db, code)
	require.NoError(t, err)

	// Retrieve and verify
	retrieved, err := GetContractCodeByHash(db, "large-wasm-test")
	require.NoError(t, err)
	assert.Equal(t, largeWasm, retrieved.Wasm)
	assert.Equal(t, 1024*1024, retrieved.SizeBytes)
}

func TestContractCode_CreatedUpdatedTimestamps(t *testing.T) {
	db := setupContractCodeTestDB(t)

	code := &ContractCode{
		Hash:       "timestamp-test",
		Wasm:       []byte{0x00, 0x61, 0x73, 0x6d},
		DeployedAt: time.Now().UTC(),
		Ledger:     100,
		TxHash:     "tx-timestamp",
		SizeBytes:  4,
	}

	beforeInsert := time.Now().UTC()
	err := UpsertContractCode(db, code)
	require.NoError(t, err)
	afterInsert := time.Now().UTC()

	retrieved, err := GetContractCodeByHash(db, "timestamp-test")
	require.NoError(t, err)

	// Verify CreatedAt and UpdatedAt were set automatically
	assert.True(t, retrieved.CreatedAt.After(beforeInsert) || retrieved.CreatedAt.Equal(beforeInsert))
	assert.True(t, retrieved.CreatedAt.Before(afterInsert) || retrieved.CreatedAt.Equal(afterInsert))
	assert.True(t, retrieved.UpdatedAt.After(beforeInsert) || retrieved.UpdatedAt.Equal(beforeInsert))
	assert.True(t, retrieved.UpdatedAt.Before(afterInsert) || retrieved.UpdatedAt.Equal(afterInsert))
}
