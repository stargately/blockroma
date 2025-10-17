package models

import (
	"testing"
	"time"

	"github.com/blockroma/soroban-indexer/pkg/models/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupBatchTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Migrate all tables
	if err := db.AutoMigrate(
		&Event{},
		&Transaction{},
		&Operation{},
		&TokenOperation{},
		&ContractCode{},
	); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestBatchUpsertEvents(t *testing.T) {
	db := setupBatchTestDB(t)

	events := []*Event{
		{
			ID:                       "event1",
			EventType:                "contract",
			Ledger:                   100,
			ContractID:               "contract1",
			LastModifiedLedgerSeq:    100,
		},
		{
			ID:                       "event2",
			EventType:                "contract",
			Ledger:                   101,
			ContractID:               "contract2",
			LastModifiedLedgerSeq:    101,
		},
		{
			ID:                       "event3",
			EventType:                "contract",
			Ledger:                   102,
			ContractID:               "contract3",
			LastModifiedLedgerSeq:    102,
		},
	}

	// Test batch insert
	err := BatchUpsertEvents(db, events)
	if err != nil {
		t.Fatalf("BatchUpsertEvents failed: %v", err)
	}

	// Verify all events were inserted
	var count int64
	db.Model(&Event{}).Count(&count)
	if count != 3 {
		t.Errorf("Expected 3 events, got %d", count)
	}

	// Test batch update
	events[0].Ledger = 200
	events[1].Ledger = 201

	err = BatchUpsertEvents(db, events[:2])
	if err != nil {
		t.Fatalf("BatchUpsertEvents update failed: %v", err)
	}

	// Verify updates
	var event1 Event
	db.Where("id = ?", "event1").First(&event1)
	if event1.Ledger != 200 {
		t.Errorf("Expected ledger 200, got %d", event1.Ledger)
	}
}

func TestBatchUpsertEvents_EmptySlice(t *testing.T) {
	db := setupBatchTestDB(t)

	events := []*Event{}
	err := BatchUpsertEvents(db, events)
	if err != nil {
		t.Fatalf("BatchUpsertEvents with empty slice failed: %v", err)
	}
}

func TestBatchUpsertEvents_CustomBatchSize(t *testing.T) {
	db := setupBatchTestDB(t)

	// Create 250 events
	events := make([]*Event, 250)
	for i := 0; i < 250; i++ {
		events[i] = &Event{
			ID:                       string(rune('A' + i%26)) + string(rune('0' + i)),
			EventType:                "contract",
			Ledger:                   int32(100 + i),
			LastModifiedLedgerSeq:    uint32(100 + i),
		}
	}

	// Use custom batch size
	config := BatchConfig{BatchSize: 50}
	err := BatchUpsertEvents(db, events, config)
	if err != nil {
		t.Fatalf("BatchUpsertEvents with custom batch size failed: %v", err)
	}

	// Verify all events were inserted
	var count int64
	db.Model(&Event{}).Count(&count)
	if count != 250 {
		t.Errorf("Expected 250 events, got %d", count)
	}
}

func TestBatchUpsertTransactions(t *testing.T) {
	db := setupBatchTestDB(t)

	fee := int32(100)
	feeCharged := int32(100)
	ledger1 := uint32(100)
	ledger2 := uint32(101)

	transactions := []*Transaction{
		{
			ID:         "tx1",
			Fee:        &fee,
			FeeCharged: &feeCharged,
			Ledger:     &ledger1,
			Status:     "SUCCESS",
		},
		{
			ID:         "tx2",
			Fee:        &fee,
			FeeCharged: &feeCharged,
			Ledger:     &ledger2,
			Status:     "SUCCESS",
		},
	}

	err := BatchUpsertTransactions(db, transactions)
	if err != nil {
		t.Fatalf("BatchUpsertTransactions failed: %v", err)
	}

	var count int64
	db.Model(&Transaction{}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 transactions, got %d", count)
	}
}

func TestBatchUpsertOperations(t *testing.T) {
	db := setupBatchTestDB(t)

	operations := []*Operation{
		{
			ID:               "tx1_0",
			TxHash:           "tx1",
			OperationIndex:   0,
			OperationType:    "payment",
			OperationDetails: []byte(`{"destination":"GABC","amount":"100"}`),
		},
		{
			ID:               "tx1_1",
			TxHash:           "tx1",
			OperationIndex:   1,
			OperationType:    "create_account",
			OperationDetails: []byte(`{"destination":"GDEF","starting_balance":"10"}`),
		},
		{
			ID:               "tx2_0",
			TxHash:           "tx2",
			OperationIndex:   0,
			OperationType:    "invoke_host_function",
			OperationDetails: []byte(`{"function":"transfer"}`),
		},
	}

	err := BatchUpsertOperations(db, operations)
	if err != nil {
		t.Fatalf("BatchUpsertOperations failed: %v", err)
	}

	var count int64
	db.Model(&Operation{}).Count(&count)
	if count != 3 {
		t.Errorf("Expected 3 operations, got %d", count)
	}

	// Test update
	operations[0].OperationDetails = []byte(`{"destination":"GXYZ","amount":"200"}`)
	err = BatchUpsertOperations(db, operations[:1])
	if err != nil {
		t.Fatalf("BatchUpsertOperations update failed: %v", err)
	}

	var op Operation
	db.Where("id = ?", "tx1_0").First(&op)
	if string(op.OperationDetails) != `{"destination":"GXYZ","amount":"200"}` {
		t.Errorf("Expected updated details, got %s", string(op.OperationDetails))
	}
}

func TestBatchUpsertTokenOperations(t *testing.T) {
	db := setupBatchTestDB(t)

	amount1 := &util.Int128{}
	amount2 := &util.Int128{}

	tokenOps := []*TokenOperation{
		{
			ID:         "event1",
			ContractID: "token1",
			Type:       "transfer",
			From:       "addr1",
			To:         strPtr("addr2"),
			Amount:     amount1,
			Ledger:     100,
		},
		{
			ID:         "event2",
			ContractID: "token1",
			Type:       "mint",
			From:       "addr1",
			To:         strPtr("addr2"),
			Amount:     amount2,
			Ledger:     101,
		},
	}

	err := BatchUpsertTokenOperations(db, tokenOps)
	if err != nil {
		t.Fatalf("BatchUpsertTokenOperations failed: %v", err)
	}

	var count int64
	db.Model(&TokenOperation{}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 token operations, got %d", count)
	}
}

func TestBatchUpsertContractCode(t *testing.T) {
	db := setupBatchTestDB(t)

	codes := []*ContractCode{
		{
			Hash:       "hash1",
			Wasm:       []byte{0x00, 0x01, 0x02},
			DeployedAt: time.Now(),
			Ledger:     100,
			TxHash:     "tx1",
			SizeBytes:  3,
		},
		{
			Hash:       "hash2",
			Wasm:       []byte{0x03, 0x04, 0x05},
			DeployedAt: time.Now(),
			Ledger:     101,
			TxHash:     "tx2",
			SizeBytes:  3,
		},
	}

	err := BatchUpsertContractCode(db, codes)
	if err != nil {
		t.Fatalf("BatchUpsertContractCode failed: %v", err)
	}

	var count int64
	db.Model(&ContractCode{}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 contract codes, got %d", count)
	}

	// Test idempotency - inserting same codes again should not increase count
	err = BatchUpsertContractCode(db, codes)
	if err != nil {
		t.Fatalf("BatchUpsertContractCode idempotency test failed: %v", err)
	}

	db.Model(&ContractCode{}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 contract codes after re-insert, got %d", count)
	}
}

// TestBatchUpsertAccountEntries removed - AccountEntry model deleted
// These classic Stellar ledger entries should be indexed via Horizon API instead
/*
func TestBatchUpsertAccountEntries(t *testing.T) {
	// REMOVED: AccountEntry model no longer supported
	// Use Horizon API for account entries
}
*/

func TestDefaultBatchConfig(t *testing.T) {
	config := DefaultBatchConfig()
	if config.BatchSize != 100 {
		t.Errorf("Expected default batch size of 100, got %d", config.BatchSize)
	}
}

func BenchmarkBatchUpsertEvents_100(b *testing.B) {
	db := setupBenchDB(b)
	events := make([]*Event, 100)
	for i := 0; i < 100; i++ {
		events[i] = &Event{
			ID:                    string(rune('A'+i%26)) + string(rune('0'+i)),
			EventType:             "contract",
			Ledger:                int32(100 + i),
			LastModifiedLedgerSeq: uint32(100 + i),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BatchUpsertEvents(db, events)
	}
}

func BenchmarkSequentialUpsertEvents_100(b *testing.B) {
	db := setupBenchDB(b)
	events := make([]*Event, 100)
	for i := 0; i < 100; i++ {
		events[i] = &Event{
			ID:                    string(rune('A'+i%26)) + string(rune('0'+i)),
			EventType:             "contract",
			Ledger:                int32(100 + i),
			LastModifiedLedgerSeq: uint32(100 + i),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, event := range events {
			UpsertEvent(db, event)
		}
	}
}

func setupBenchDB(b *testing.B) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&Event{}); err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func strPtr(s string) *string {
	return &s
}

// TestBatchFunctionsWithinTransaction tests that all batch functions work correctly within a transaction
// This is critical because we removed the transaction wrappers from batch functions
func TestBatchFunctionsWithinTransaction(t *testing.T) {
	db := setupBatchTestDB(t)

	// Test all batch functions within a single transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// Test BatchUpsertEvents
		events := []*Event{
			{
				ID:                    "event_tx1",
				EventType:             "contract",
				Ledger:                100,
				LastModifiedLedgerSeq: 100,
			},
		}
		if err := BatchUpsertEvents(tx, events); err != nil {
			return err
		}

		// Test BatchUpsertTransactions
		fee := int32(100)
		ledger := uint32(100)
		transactions := []*Transaction{
			{
				ID:     "tx_tx1",
				Fee:    &fee,
				Ledger: &ledger,
				Status: "SUCCESS",
			},
		}
		if err := BatchUpsertTransactions(tx, transactions); err != nil {
			return err
		}

		// Test BatchUpsertOperations
		operations := []*Operation{
			{
				ID:               "op_tx1",
				TxHash:           "tx_tx1",
				OperationIndex:   0,
				OperationType:    "payment",
				OperationDetails: []byte(`{}`),
			},
		}
		if err := BatchUpsertOperations(tx, operations); err != nil {
			return err
		}

		// Test BatchUpsertTokenOperations
		amount := &util.Int128{}
		tokenOps := []*TokenOperation{
			{
				ID:         "tokenop_tx1",
				ContractID: "token1",
				Type:       "transfer",
				From:       "addr1",
				Amount:     amount,
				Ledger:     100,
			},
		}
		if err := BatchUpsertTokenOperations(tx, tokenOps); err != nil {
			return err
		}

		// Test BatchUpsertContractCode
		codes := []*ContractCode{
			{
				Hash:       "hash_tx1",
				Wasm:       []byte{0x00},
				DeployedAt: time.Now(),
				Ledger:     100,
				TxHash:     "tx_tx1",
				SizeBytes:  1,
			},
		}
		if err := BatchUpsertContractCode(tx, codes); err != nil {
			return err
		}

		// BatchUpsertAccountEntries removed - not supported

		return nil
	})

	if err != nil {
		t.Fatalf("Transaction with batch functions failed: %v", err)
	}

	// Verify all records were inserted
	var eventCount, txCount, opCount, tokenOpCount, codeCount int64
	db.Model(&Event{}).Count(&eventCount)
	db.Model(&Transaction{}).Count(&txCount)
	db.Model(&Operation{}).Count(&opCount)
	db.Model(&TokenOperation{}).Count(&tokenOpCount)
	db.Model(&ContractCode{}).Count(&codeCount)

	if eventCount != 1 {
		t.Errorf("Expected 1 event, got %d", eventCount)
	}
	if txCount != 1 {
		t.Errorf("Expected 1 transaction, got %d", txCount)
	}
	if opCount != 1 {
		t.Errorf("Expected 1 operation, got %d", opCount)
	}
	if tokenOpCount != 1 {
		t.Errorf("Expected 1 token operation, got %d", tokenOpCount)
	}
	if codeCount != 1 {
		t.Errorf("Expected 1 contract code, got %d", codeCount)
	}
}

// TestBatchFunctionsTransactionRollback tests that rollback works correctly with batch functions
func TestBatchFunctionsTransactionRollback(t *testing.T) {
	db := setupBatchTestDB(t)

	// Attempt transaction that will be rolled back
	err := db.Transaction(func(tx *gorm.DB) error {
		events := []*Event{
			{
				ID:                    "event_rollback",
				EventType:             "contract",
				Ledger:                100,
				LastModifiedLedgerSeq: 100,
			},
		}
		if err := BatchUpsertEvents(tx, events); err != nil {
			return err
		}

		operations := []*Operation{
			{
				ID:               "op_rollback",
				TxHash:           "tx_rollback",
				OperationIndex:   0,
				OperationType:    "payment",
				OperationDetails: []byte(`{}`),
			},
		}
		if err := BatchUpsertOperations(tx, operations); err != nil {
			return err
		}

		// Force rollback
		return gorm.ErrInvalidTransaction
	})

	if err == nil {
		t.Error("Expected transaction to fail")
	}

	// Verify nothing was inserted
	var eventCount, opCount int64
	db.Model(&Event{}).Count(&eventCount)
	db.Model(&Operation{}).Count(&opCount)

	if eventCount != 0 {
		t.Errorf("Expected 0 events after rollback, got %d", eventCount)
	}
	if opCount != 0 {
		t.Errorf("Expected 0 operations after rollback, got %d", opCount)
	}
}

// TestBatchFunctionsNestedTransaction tests that batch functions work in nested transaction contexts
// This test ensures our fix doesn't break if GORM creates savepoints
func TestBatchFunctionsNestedTransaction(t *testing.T) {
	db := setupBatchTestDB(t)

	// Outer transaction
	err := db.Transaction(func(tx1 *gorm.DB) error {
		events1 := []*Event{
			{
				ID:                    "event_outer",
				EventType:             "contract",
				Ledger:                100,
				LastModifiedLedgerSeq: 100,
			},
		}
		if err := BatchUpsertEvents(tx1, events1); err != nil {
			return err
		}

		// Inner transaction (GORM will create a savepoint)
		return tx1.Transaction(func(tx2 *gorm.DB) error {
			events2 := []*Event{
				{
					ID:                    "event_inner",
					EventType:             "contract",
					Ledger:                101,
					LastModifiedLedgerSeq: 101,
				},
			}
			return BatchUpsertEvents(tx2, events2)
		})
	})

	if err != nil {
		t.Fatalf("Nested transaction failed: %v", err)
	}

	// Verify both events were inserted
	var count int64
	db.Model(&Event{}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 events, got %d", count)
	}
}

// TestBatchFunctionsLargeBatchWithinTransaction tests large batches within a transaction
func TestBatchFunctionsLargeBatchWithinTransaction(t *testing.T) {
	db := setupBatchTestDB(t)

	// Create 500 operations (will be split into multiple batches)
	operations := make([]*Operation, 500)
	for i := 0; i < 500; i++ {
		operations[i] = &Operation{
			ID:               string(rune('A'+i%26)) + string(rune('0'+i)),
			TxHash:           "tx_large",
			OperationIndex:   int32(i),
			OperationType:    "payment",
			OperationDetails: []byte(`{}`),
		}
	}

	// Insert within transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		return BatchUpsertOperations(tx, operations)
	})

	if err != nil {
		t.Fatalf("Large batch within transaction failed: %v", err)
	}

	// Verify all operations were inserted
	var count int64
	db.Model(&Operation{}).Count(&count)
	if count != 500 {
		t.Errorf("Expected 500 operations, got %d", count)
	}
}
