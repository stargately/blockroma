package models

import (
	"encoding/json"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestContractDataEntryJSONBSerialization tests that interface{} values are properly serialized to JSON
func TestContractDataEntryJSONBSerialization(t *testing.T) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(&ContractDataEntry{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Helper to marshal values to JSONB
	mustMarshal := func(v interface{}) JSONB {
		bytes, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}
		return JSONB(bytes)
	}

	tests := []struct {
		name  string
		entry *ContractDataEntry
	}{
		{
			name: "Simple string key and value",
			entry: &ContractDataEntry{
				KeyHash:    "hash1",
				ContractID: "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM",
				Key:        mustMarshal("simple_key"),
				Val:        mustMarshal("simple_value"),
				Durability: "persistent",
			},
		},
		{
			name: "Map key and value",
			entry: &ContractDataEntry{
				KeyHash:    "hash2",
				ContractID: "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM",
				Key: mustMarshal(map[string]interface{}{
					"type": "balance",
					"user": "GABC123",
				}),
				Val: mustMarshal(map[string]interface{}{
					"amount":  "1000000",
					"enabled": true,
				}),
				Durability: "temporary",
			},
		},
		{
			name: "Nested structure",
			entry: &ContractDataEntry{
				KeyHash:    "hash3",
				ContractID: "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM",
				Key:        mustMarshal("ScvLedgerKeyContractInstance"),
				Val: mustMarshal(map[string]interface{}{
					"executable": map[string]interface{}{
						"type":     "Wasm",
						"wasmHash": "abc123",
					},
					"storage": []interface{}{
						map[string]interface{}{
							"key":   "name",
							"value": "My Token",
						},
						map[string]interface{}{
							"key":   "symbol",
							"value": "MTK",
						},
					},
				}),
				Durability: "persistent",
			},
		},
		{
			name: "Array value",
			entry: &ContractDataEntry{
				KeyHash:    "hash4",
				ContractID: "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM",
				Key:        mustMarshal("balances"),
				Val: mustMarshal([]interface{}{
					map[string]interface{}{"address": "G123", "balance": 100},
					map[string]interface{}{"address": "G456", "balance": 200},
				}),
				Durability: "persistent",
			},
		},
		{
			name: "Numeric values",
			entry: &ContractDataEntry{
				KeyHash:    "hash5",
				ContractID: "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM",
				Key:        mustMarshal("counter"),
				Val: mustMarshal(map[string]interface{}{
					"count":       42,
					"lastUpdated": 1234567890,
					"enabled":     true,
				}),
				Durability: "persistent",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Insert entry
			if err := UpsertContractDataEntry(db, tt.entry); err != nil {
				t.Errorf("Failed to upsert contract data entry: %v", err)
				return
			}

			// Retrieve entry
			var retrieved ContractDataEntry
			if err := db.Where("key_hash = ?", tt.entry.KeyHash).First(&retrieved).Error; err != nil {
				t.Errorf("Failed to retrieve contract data entry: %v", err)
				return
			}

			// Verify basic fields
			if retrieved.KeyHash != tt.entry.KeyHash {
				t.Errorf("KeyHash mismatch: got %v, want %v", retrieved.KeyHash, tt.entry.KeyHash)
			}
			if retrieved.ContractID != tt.entry.ContractID {
				t.Errorf("ContractID mismatch: got %v, want %v", retrieved.ContractID, tt.entry.ContractID)
			}
			if retrieved.Durability != tt.entry.Durability {
				t.Errorf("Durability mismatch: got %v, want %v", retrieved.Durability, tt.entry.Durability)
			}

			// Verify Key and Val are not nil
			if retrieved.Key == nil {
				t.Error("Retrieved Key is nil")
			}
			if retrieved.Val == nil {
				t.Error("Retrieved Val is nil")
			}
		})
	}
}

// TestContractDataEntryUpsert tests the upsert functionality
func TestContractDataEntryUpsert(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&ContractDataEntry{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Helper to marshal values to JSONB
	mustMarshal := func(v interface{}) JSONB {
		bytes, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}
		return JSONB(bytes)
	}

	// Initial insert
	entry := &ContractDataEntry{
		KeyHash:    "test_hash",
		ContractID: "CTEST",
		Key:        mustMarshal("key1"),
		Val:        mustMarshal("value1"),
		Durability: "persistent",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := UpsertContractDataEntry(db, entry); err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	// Verify insert
	var count int64
	db.Model(&ContractDataEntry{}).Where("key_hash = ?", "test_hash").Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 entry, got %d", count)
	}

	// Update with same key_hash
	time.Sleep(10 * time.Millisecond) // Ensure updated_at is different
	updatedEntry := &ContractDataEntry{
		KeyHash:    "test_hash",
		ContractID: "CTEST_UPDATED",
		Key:        mustMarshal("key2"),
		Val:        mustMarshal("value2"),
		Durability: "temporary",
	}

	if err := UpsertContractDataEntry(db, updatedEntry); err != nil {
		t.Fatalf("Failed to upsert: %v", err)
	}

	// Verify only one entry exists
	db.Model(&ContractDataEntry{}).Where("key_hash = ?", "test_hash").Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 entry after upsert, got %d", count)
	}

	// Verify fields were updated
	var retrieved ContractDataEntry
	db.Where("key_hash = ?", "test_hash").First(&retrieved)
	if retrieved.ContractID != "CTEST_UPDATED" {
		t.Errorf("ContractID not updated: got %v, want CTEST_UPDATED", retrieved.ContractID)
	}
	if retrieved.Durability != "temporary" {
		t.Errorf("Durability not updated: got %v, want temporary", retrieved.Durability)
	}
}

// TestContractDataEntryNilValues tests handling of nil values
func TestContractDataEntryNilValues(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&ContractDataEntry{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Entry with empty/nil JSONB fields
	entry := &ContractDataEntry{
		KeyHash:    "nil_test",
		ContractID: "CTEST",
		Key:        nil, // nil JSONB
		Val:        nil, // nil JSONB
		Durability: "persistent",
	}

	// Should insert without error
	if err := UpsertContractDataEntry(db, entry); err != nil {
		t.Errorf("Failed to insert entry with nil JSONB: %v", err)
	}

	// Retrieve and verify
	var retrieved ContractDataEntry
	if err := db.Where("key_hash = ?", "nil_test").First(&retrieved).Error; err != nil {
		t.Errorf("Failed to retrieve entry: %v", err)
	}
}

// TestContractDataEntryComplexStructures tests complex nested structures
func TestContractDataEntryComplexStructures(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&ContractDataEntry{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Helper to marshal values to JSONB
	mustMarshal := func(v interface{}) JSONB {
		bytes, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}
		return JSONB(bytes)
	}

	// Complex nested structure mimicking real Soroban data
	entry := &ContractDataEntry{
		KeyHash:    "complex_test",
		ContractID: "CTEST",
		Key:        mustMarshal("ScvLedgerKeyContractInstance"),
		Val: mustMarshal(map[string]interface{}{
			"executable": map[string]interface{}{
				"type": "ContractExecutableWasm",
				"wasm_hash": []interface{}{
					0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56, 0x78, 0x90,
				},
			},
			"storage": []interface{}{
				map[string]interface{}{
					"key": map[string]interface{}{
						"type":  "Symbol",
						"value": "name",
					},
					"value": map[string]interface{}{
						"type":  "String",
						"value": "My Token Name",
					},
				},
				map[string]interface{}{
					"key": map[string]interface{}{
						"type":  "Symbol",
						"value": "symbol",
					},
					"value": map[string]interface{}{
						"type":  "String",
						"value": "MTK",
					},
				},
				map[string]interface{}{
					"key": map[string]interface{}{
						"type":  "Symbol",
						"value": "decimal",
					},
					"value": map[string]interface{}{
						"type":  "U32",
						"value": 7,
					},
				},
			},
		}),
		Durability: "persistent",
	}

	// Insert
	if err := UpsertContractDataEntry(db, entry); err != nil {
		t.Errorf("Failed to insert complex structure: %v", err)
		return
	}

	// Retrieve
	var retrieved ContractDataEntry
	if err := db.Where("key_hash = ?", "complex_test").First(&retrieved).Error; err != nil {
		t.Errorf("Failed to retrieve complex structure: %v", err)
		return
	}

	// Unmarshal Val to verify structure
	var valMap map[string]interface{}
	if err := json.Unmarshal(retrieved.Val, &valMap); err != nil {
		t.Errorf("Failed to unmarshal Val: %v", err)
		return
	}

	// Verify nested structure exists
	if _, hasExecutable := valMap["executable"]; !hasExecutable {
		t.Error("Val missing 'executable' key")
	}
	if _, hasStorage := valMap["storage"]; !hasStorage {
		t.Error("Val missing 'storage' key")
	}
}

// BenchmarkContractDataEntryUpsert benchmarks the upsert performance
func BenchmarkContractDataEntryUpsert(b *testing.B) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&ContractDataEntry{}); err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	// Helper to marshal values to JSONB
	mustMarshal := func(v interface{}) JSONB {
		bytes, err := json.Marshal(v)
		if err != nil {
			b.Fatalf("Failed to marshal test data: %v", err)
		}
		return JSONB(bytes)
	}

	entry := &ContractDataEntry{
		KeyHash:    "bench_hash",
		ContractID: "CBENCH",
		Key: mustMarshal(map[string]interface{}{
			"type": "balance",
			"user": "G123",
		}),
		Val: mustMarshal(map[string]interface{}{
			"amount":  "1000000",
			"enabled": true,
		}),
		Durability: "persistent",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := UpsertContractDataEntry(db, entry); err != nil {
			b.Fatalf("Upsert failed: %v", err)
		}
	}
}
