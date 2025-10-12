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
		&AccountEntry{},
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
			ID:             "tx1_0",
			TxHash:         "tx1",
			OperationIndex: 0,
			OperationType:  "payment",
		},
		{
			ID:             "tx1_1",
			TxHash:         "tx1",
			OperationIndex: 1,
			OperationType:  "create_account",
		},
		{
			ID:             "tx2_0",
			TxHash:         "tx2",
			OperationIndex: 0,
			OperationType:  "invoke_host_function",
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

func TestBatchUpsertAccountEntries(t *testing.T) {
	db := setupBatchTestDB(t)

	balance := int64(1000000)
	seqNum := int64(12345)

	// Ext needs to be properly JSON-encoded
	extJSON := []byte(`{"v":0}`)

	accounts := []*AccountEntry{
		{
			AccountID:     "account1",
			Balance:       balance,
			SeqNum:        seqNum,
			NumSubEntries: 0,
			Flags:         0,
			HomeDomain:    "",
			InflationDest: "",
			Ext:           extJSON,
		},
		{
			AccountID:     "account2",
			Balance:       balance,
			SeqNum:        seqNum,
			NumSubEntries: 0,
			Flags:         0,
			HomeDomain:    "",
			InflationDest: "",
			Ext:           extJSON,
		},
	}

	err := BatchUpsertAccountEntries(db, accounts)
	if err != nil {
		t.Fatalf("BatchUpsertAccountEntries failed: %v", err)
	}

	var count int64
	db.Model(&AccountEntry{}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 account entries, got %d", count)
	}
}

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
