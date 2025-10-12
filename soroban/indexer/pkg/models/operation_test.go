package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestOperationTableName(t *testing.T) {
	op := Operation{}
	if op.TableName() != "operations" {
		t.Errorf("TableName() = %v, want operations", op.TableName())
	}
}

func TestUpsertOperation(t *testing.T) {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(&Operation{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create test operation
	op := &Operation{
		ID:               "test-tx-0",
		TxHash:           "test-tx",
		OperationIndex:   0,
		SourceAccount:    "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		OperationType:    "OPERATION_TYPE_PAYMENT",
		OperationDetails: []byte(`{"destination": "GABC...", "amount": 1000000}`),
	}

	// Insert
	if err := UpsertOperation(db, op); err != nil {
		t.Fatalf("UpsertOperation() error = %v", err)
	}

	// Verify
	var retrieved Operation
	if err := db.Where("id = ?", "test-tx-0").First(&retrieved).Error; err != nil {
		t.Fatalf("Failed to retrieve operation: %v", err)
	}

	if retrieved.TxHash != "test-tx" {
		t.Errorf("TxHash = %v, want test-tx", retrieved.TxHash)
	}

	if retrieved.OperationIndex != 0 {
		t.Errorf("OperationIndex = %v, want 0", retrieved.OperationIndex)
	}

	if retrieved.OperationType != "OPERATION_TYPE_PAYMENT" {
		t.Errorf("OperationType = %v, want OPERATION_TYPE_PAYMENT", retrieved.OperationType)
	}
}

func TestGetOperationsByTxHash(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&Operation{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Insert multiple operations for same transaction
	ops := []*Operation{
		{
			ID:             "test-tx-0",
			TxHash:         "test-tx",
			OperationIndex: 0,
			OperationType:  "OPERATION_TYPE_PAYMENT",
		},
		{
			ID:             "test-tx-1",
			TxHash:         "test-tx",
			OperationIndex: 1,
			OperationType:  "OPERATION_TYPE_BUMP_SEQUENCE",
		},
		{
			ID:             "other-tx-0",
			TxHash:         "other-tx",
			OperationIndex: 0,
			OperationType:  "OPERATION_TYPE_CREATE_ACCOUNT",
		},
	}

	for _, op := range ops {
		if err := UpsertOperation(db, op); err != nil {
			t.Fatalf("UpsertOperation() error = %v", err)
		}
	}

	// Retrieve operations for test-tx
	retrieved, err := GetOperationsByTxHash(db, "test-tx")
	if err != nil {
		t.Fatalf("GetOperationsByTxHash() error = %v", err)
	}

	if len(retrieved) != 2 {
		t.Fatalf("GetOperationsByTxHash() returned %d operations, want 2", len(retrieved))
	}

	// Check they're in order
	if retrieved[0].OperationIndex != 0 {
		t.Errorf("First operation index = %v, want 0", retrieved[0].OperationIndex)
	}

	if retrieved[1].OperationIndex != 1 {
		t.Errorf("Second operation index = %v, want 1", retrieved[1].OperationIndex)
	}
}

func TestGetOperationByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&Operation{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Insert operation
	op := &Operation{
		ID:             "test-tx-0",
		TxHash:         "test-tx",
		OperationIndex: 0,
		OperationType:  "OPERATION_TYPE_PAYMENT",
	}

	if err := UpsertOperation(db, op); err != nil {
		t.Fatalf("UpsertOperation() error = %v", err)
	}

	// Retrieve by ID
	retrieved, err := GetOperationByID(db, "test-tx-0")
	if err != nil {
		t.Fatalf("GetOperationByID() error = %v", err)
	}

	if retrieved.ID != "test-tx-0" {
		t.Errorf("ID = %v, want test-tx-0", retrieved.ID)
	}

	// Try to retrieve non-existent operation
	_, err = GetOperationByID(db, "non-existent")
	if err == nil {
		t.Error("GetOperationByID() expected error for non-existent operation")
	}
}

func TestUpsertOperation_Update(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&Operation{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Insert initial operation
	op := &Operation{
		ID:             "test-tx-0",
		TxHash:         "test-tx",
		OperationIndex: 0,
		OperationType:  "OPERATION_TYPE_PAYMENT",
		SourceAccount:  "GABC123",
	}

	if err := UpsertOperation(db, op); err != nil {
		t.Fatalf("UpsertOperation() insert error = %v", err)
	}

	// Update the operation
	op.SourceAccount = "GDEF456"
	if err := UpsertOperation(db, op); err != nil {
		t.Fatalf("UpsertOperation() update error = %v", err)
	}

	// Verify update
	retrieved, err := GetOperationByID(db, "test-tx-0")
	if err != nil {
		t.Fatalf("GetOperationByID() error = %v", err)
	}

	if retrieved.SourceAccount != "GDEF456" {
		t.Errorf("SourceAccount = %v, want GDEF456", retrieved.SourceAccount)
	}

	// Verify only one record exists
	var count int64
	db.Model(&Operation{}).Count(&count)
	if count != 1 {
		t.Errorf("Operation count = %v, want 1", count)
	}
}
