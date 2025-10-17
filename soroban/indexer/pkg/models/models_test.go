package models

import (
	"testing"
	"time"

	"github.com/blockroma/soroban-indexer/pkg/models/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Auto-migrate all models
	// Note: Account/trustline/offer/data/claimable balance/liquidity pool tables removed
	// These classic Stellar ledger entries should be indexed via Horizon API instead
	err = db.AutoMigrate(
		&Event{},
		&Transaction{},
		&Cursor{},
		&TokenMetadata{},
		&TokenOperation{},
		&TokenBalance{},
		&ContractDataEntry{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestEventTableName(t *testing.T) {
	event := Event{}
	if event.TableName() != "events" {
		t.Errorf("TableName() = %v, want events", event.TableName())
	}
}

func TestUpsertEvent(t *testing.T) {
	db := setupTestDB(t)

	event := &Event{
		ID:                       "event-123",
		TxIndex:                  1,
		EventType:                "contract",
		Ledger:                   12345,
		LedgerClosedAt:           "2024-01-01T00:00:00Z",
		ContractID:               "contract-456",
		PagingToken:              "12345-1",
		Topic:                    `["transfer"]`,
		Value:                    `"1000000"`,
		InSuccessfulContractCall: true,
		LastModifiedLedgerSeq:    12345,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	err := UpsertEvent(db, event)
	if err != nil {
		t.Fatalf("UpsertEvent() error = %v", err)
	}

	// Verify it was inserted
	var retrieved Event
	result := db.First(&retrieved, "id = ?", "event-123")
	if result.Error != nil {
		t.Fatalf("Failed to retrieve event: %v", result.Error)
	}

	if retrieved.ID != event.ID {
		t.Errorf("ID = %v, want %v", retrieved.ID, event.ID)
	}
	if retrieved.ContractID != event.ContractID {
		t.Errorf("ContractID = %v, want %v", retrieved.ContractID, event.ContractID)
	}
	if retrieved.LastModifiedLedgerSeq != event.LastModifiedLedgerSeq {
		t.Errorf("LastModifiedLedgerSeq = %v, want %v", retrieved.LastModifiedLedgerSeq, event.LastModifiedLedgerSeq)
	}

	// Test upsert (update)
	event.Topic = `["mint"]`
	event.LastModifiedLedgerSeq = 12346
	err = UpsertEvent(db, event)
	if err != nil {
		t.Fatalf("UpsertEvent() update error = %v", err)
	}

	var updated Event
	db.First(&updated, "id = ?", "event-123")

	// Topic can be returned as string or *interface{} depending on database driver
	var topicStr string
	switch v := updated.Topic.(type) {
	case string:
		topicStr = v
	case *interface{}:
		if v != nil {
			if s, ok := (*v).(string); ok {
				topicStr = s
			}
		}
	default:
		t.Fatalf("Topic has unexpected type %T", updated.Topic)
	}

	if topicStr != `["mint"]` {
		t.Errorf("Topic after update = %v, want [\"mint\"]", topicStr)
	}
	if updated.LastModifiedLedgerSeq != 12346 {
		t.Errorf("LastModifiedLedgerSeq after update = %v, want 12346", updated.LastModifiedLedgerSeq)
	}
}

func TestUpsertTransaction(t *testing.T) {
	db := setupTestDB(t)

	ledger := uint32(12345)
	appOrder := int32(1)
	sourceAccount := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"
	fee := int32(100)
	feeCharged := int32(100)
	sequence := int64(987654321)
	feeBump := false
	createdAt := int64(1704067200)

	tx := &Transaction{
		ID:               "tx-hash-123",
		Status:           "SUCCESS",
		Ledger:           &ledger,
		LedgerCreatedAt:  &createdAt,
		ApplicationOrder: &appOrder,
		SourceAccount:    &sourceAccount,
		Fee:              &fee,
		FeeCharged:       &feeCharged,
		Sequence:         &sequence,
		FeeBump:          &feeBump,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := UpsertTransaction(db, tx)
	if err != nil {
		t.Fatalf("UpsertTransaction() error = %v", err)
	}

	var retrieved Transaction
	result := db.First(&retrieved, "id = ?", "tx-hash-123")
	if result.Error != nil {
		t.Fatalf("Failed to retrieve transaction: %v", result.Error)
	}

	if retrieved.Status != tx.Status {
		t.Errorf("Status = %v, want %v", retrieved.Status, tx.Status)
	}
	if retrieved.SourceAccount == nil || *retrieved.SourceAccount != *tx.SourceAccount {
		t.Errorf("SourceAccount = %v, want %v", retrieved.SourceAccount, tx.SourceAccount)
	}
	if retrieved.Ledger == nil || *retrieved.Ledger != *tx.Ledger {
		t.Errorf("Ledger = %v, want %v", retrieved.Ledger, tx.Ledger)
	}
	if retrieved.FeeBump == nil || *retrieved.FeeBump != *tx.FeeBump {
		t.Errorf("FeeBump = %v, want %v", retrieved.FeeBump, tx.FeeBump)
	}
}

func TestTransactionWithExtendedFields(t *testing.T) {
	db := setupTestDB(t)

	ledger := uint32(12345)
	appOrder := int32(1)
	sourceAccount := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"
	fee := int32(100)
	feeCharged := int32(100)
	sequence := int64(987654321)
	feeBump := true
	muxedId := int64(1234567890)

	// Extended fields
	memo := &util.TypeItem{
		Type:      "text",
		ItemValue: "test memo",
	}

	signatures := &util.Signatures{
		{Hint: "hint1", Signature: "sig1"},
		{Hint: "hint2", Signature: "sig2"},
	}

	minSeqNum := int64(100)
	preconditions := &util.Preconditions{
		MinSeqNum: &minSeqNum,
		TimeBounds: &util.Bonds{
			Min: 1000000,
			Max: 2000000,
		},
	}

	feeBumpInfo := &util.FeeBumpInfo{
		Fee:            200,
		SourceAccount:  &sourceAccount,
		MuxedAccountId: &muxedId,
	}

	tx := &Transaction{
		ID:               "tx-extended-123",
		Status:           "SUCCESS",
		Ledger:           &ledger,
		ApplicationOrder: &appOrder,
		SourceAccount:    &sourceAccount,
		Fee:              &fee,
		FeeCharged:       &feeCharged,
		Sequence:         &sequence,
		FeeBump:          &feeBump,
		MuxedAccountId:   &muxedId,
		Memo:             memo,
		Signatures:       signatures,
		Preconditions:    preconditions,
		FeeBumpInfo:      feeBumpInfo,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := UpsertTransaction(db, tx)
	if err != nil {
		t.Fatalf("UpsertTransaction() error = %v", err)
	}

	var retrieved Transaction
	result := db.First(&retrieved, "id = ?", "tx-extended-123")
	if result.Error != nil {
		t.Fatalf("Failed to retrieve transaction: %v", result.Error)
	}

	if retrieved.Memo == nil || retrieved.Memo.Type != "text" || retrieved.Memo.ItemValue != "test memo" {
		t.Errorf("Memo = %+v, want text/test memo", retrieved.Memo)
	}
	if retrieved.Signatures == nil || len(*retrieved.Signatures) != 2 {
		t.Errorf("Signatures = %+v, want 2 signatures", retrieved.Signatures)
	}
	if retrieved.Preconditions == nil || retrieved.Preconditions.MinSeqNum == nil {
		t.Errorf("Preconditions = %+v, want preconditions with MinSeqNum", retrieved.Preconditions)
	}
	if retrieved.FeeBumpInfo == nil || retrieved.FeeBumpInfo.Fee != 200 {
		t.Errorf("FeeBumpInfo = %+v, want fee 200", retrieved.FeeBumpInfo)
	}
	if retrieved.MuxedAccountId == nil || *retrieved.MuxedAccountId != muxedId {
		t.Errorf("MuxedAccountId = %v, want %v", retrieved.MuxedAccountId, muxedId)
	}
}

func TestCursor(t *testing.T) {
	db := setupTestDB(t)

	// Initial cursor should be 0
	cursor, err := GetCursor(db)
	if err != nil {
		t.Fatalf("GetCursor() error = %v", err)
	}
	if cursor != 0 {
		t.Errorf("Initial cursor = %v, want 0", cursor)
	}

	// Update cursor
	err = UpdateCursor(db, 12345)
	if err != nil {
		t.Fatalf("UpdateCursor() error = %v", err)
	}

	// Verify updated cursor
	cursor, err = GetCursor(db)
	if err != nil {
		t.Fatalf("GetCursor() after update error = %v", err)
	}
	if cursor != 12345 {
		t.Errorf("Updated cursor = %v, want 12345", cursor)
	}

	// Update again
	err = UpdateCursor(db, 67890)
	if err != nil {
		t.Fatalf("UpdateCursor() second update error = %v", err)
	}

	cursor, err = GetCursor(db)
	if err != nil {
		t.Fatalf("GetCursor() after second update error = %v", err)
	}
	if cursor != 67890 {
		t.Errorf("Second updated cursor = %v, want 67890", cursor)
	}
}

func TestUpsertTokenMetadata(t *testing.T) {
	db := setupTestDB(t)

	metadata := &TokenMetadata{
		ContractID:   "contract-123",
		AdminAddress: "admin-address",
		Decimal:      7,
		Name:         "Test Token",
		Symbol:       "TEST",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := UpsertTokenMetadata(db, metadata)
	if err != nil {
		t.Fatalf("UpsertTokenMetadata() error = %v", err)
	}

	var retrieved TokenMetadata
	result := db.First(&retrieved, "contract_id = ?", "contract-123")
	if result.Error != nil {
		t.Fatalf("Failed to retrieve metadata: %v", result.Error)
	}

	if retrieved.Name != metadata.Name {
		t.Errorf("Name = %v, want %v", retrieved.Name, metadata.Name)
	}
	if retrieved.Symbol != metadata.Symbol {
		t.Errorf("Symbol = %v, want %v", retrieved.Symbol, metadata.Symbol)
	}
	if retrieved.Decimal != metadata.Decimal {
		t.Errorf("Decimal = %v, want %v", retrieved.Decimal, metadata.Decimal)
	}
}

func TestUpsertTokenOperation(t *testing.T) {
	db := setupTestDB(t)

	to := "recipient-address"
	amount := util.Int128{}
	amount.SetString("1000000", 10)

	op := &TokenOperation{
		ID:             "op-123",
		Type:           "transfer",
		TxIndex:        1,
		Ledger:         12345,
		LedgerClosedAt: "2024-01-01T00:00:00Z",
		ContractID:     "contract-456",
		From:           "sender-address",
		To:             &to,
		Amount:         &amount,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := UpsertTokenOperation(db, op)
	if err != nil {
		t.Fatalf("UpsertTokenOperation() error = %v", err)
	}

	var retrieved TokenOperation
	result := db.First(&retrieved, "id = ?", "op-123")
	if result.Error != nil {
		t.Fatalf("Failed to retrieve token operation: %v", result.Error)
	}

	if retrieved.Type != op.Type {
		t.Errorf("Type = %v, want %v", retrieved.Type, op.Type)
	}
	if retrieved.From != op.From {
		t.Errorf("From = %v, want %v", retrieved.From, op.From)
	}
}

func TestUpsertTokenBalance(t *testing.T) {
	db := setupTestDB(t)

	balance := &TokenBalance{
		ContractID: "contract-789",
		Address:    "holder-address",
		Balance:    "5000000000",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := UpsertTokenBalance(db, balance)
	if err != nil {
		t.Fatalf("UpsertTokenBalance() error = %v", err)
	}

	var retrieved TokenBalance
	result := db.First(&retrieved, "contract_id = ? AND address = ?", "contract-789", "holder-address")
	if result.Error != nil {
		t.Fatalf("Failed to retrieve token balance: %v", result.Error)
	}

	if retrieved.Balance != balance.Balance {
		t.Errorf("Balance = %v, want %v", retrieved.Balance, balance.Balance)
	}
}

func TestUpsertContractDataEntry(t *testing.T) {
	db := setupTestDB(t)

	entry := &ContractDataEntry{
		KeyHash:    "hash-123",
		ContractID: "contract-abc",
		KeyXdr:     "key-xdr-data",
		Key:        []byte(`{"type":"instance"}`),
		Durability: "persistent",
		ValXdr:     "val-xdr-data",
		Val:        []byte(`{"value":1000}`),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := UpsertContractDataEntry(db, entry)
	if err != nil {
		t.Fatalf("UpsertContractDataEntry() error = %v", err)
	}

	var retrieved ContractDataEntry
	result := db.First(&retrieved, "key_hash = ?", "hash-123")
	if result.Error != nil {
		t.Fatalf("Failed to retrieve contract data entry: %v", result.Error)
	}

	if retrieved.ContractID != entry.ContractID {
		t.Errorf("ContractID = %v, want %v", retrieved.ContractID, entry.ContractID)
	}
	if retrieved.Durability != entry.Durability {
		t.Errorf("Durability = %v, want %v", retrieved.Durability, entry.Durability)
	}
}

// TestUpsertAccountEntry removed - AccountEntry model deleted
// These classic Stellar ledger entries should be indexed via Horizon API instead
/*
func TestUpsertAccountEntry(t *testing.T) {
	// REMOVED: AccountEntry model no longer supported
	// Use Horizon API for account entries
}
*/

func TestInt128(t *testing.T) {
	var i128 util.Int128

	// Test SetString
	i128.SetString("1000000000000", 10)
	if i128.String() != "1000000000000" {
		t.Errorf("Int128 string = %v, want 1000000000000", i128.String())
	}

	// Test Value (for database storage)
	val, err := i128.Value()
	if err != nil {
		t.Errorf("Int128.Value() error = %v", err)
	}
	if val != "1000000000000" {
		t.Errorf("Int128.Value() = %v, want 1000000000000", val)
	}

	// Test Scan (for database retrieval)
	var i128_2 util.Int128
	err = i128_2.Scan("5000000000000")
	if err != nil {
		t.Errorf("Int128.Scan() error = %v", err)
	}
	if i128_2.String() != "5000000000000" {
		t.Errorf("Int128 after scan = %v, want 5000000000000", i128_2.String())
	}
}
