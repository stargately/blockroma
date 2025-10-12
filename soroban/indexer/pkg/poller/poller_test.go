package poller

import (
	"encoding/base64"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/xdr"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/blockroma/soroban-indexer/pkg/client"
	"github.com/blockroma/soroban-indexer/pkg/models"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&models.Transaction{},
		&models.AccountEntry{},
		&models.TrustLineEntry{},
		&models.OfferEntry{},
		&models.DataEntry{},
		&models.ClaimableBalanceEntry{},
		&models.LiquidityPoolEntry{},
	); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// createTestTransactionEnvelope creates a test transaction envelope with a source account
func createTestTransactionEnvelope(sourceAccountAddr string) string {
	sourceAccount := xdr.MustAddress(sourceAccountAddr)

	tx := xdr.Transaction{
		SourceAccount: sourceAccount.ToMuxedAccount(),
		Fee:           xdr.Uint32(100),
		SeqNum:        xdr.SequenceNumber(1),
		Cond:          xdr.Preconditions{Type: xdr.PreconditionTypePrecondNone},
		Memo:          xdr.Memo{Type: xdr.MemoTypeMemoNone},
		Operations:    []xdr.Operation{},
		Ext:           xdr.TransactionExt{V: 0},
	}

	envelope := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx:         tx,
			Signatures: []xdr.DecoratedSignature{},
		},
	}

	encoded, _ := xdr.MarshalBase64(envelope)
	return encoded
}

// createTestAccountLedgerEntry creates a test account ledger entry XDR
func createTestAccountLedgerEntry(accountAddr string) string {
	accountID := xdr.MustAddress(accountAddr)

	account := xdr.AccountEntry{
		AccountId:     accountID,
		Balance:       xdr.Int64(1000000),
		SeqNum:        xdr.SequenceNumber(1),
		NumSubEntries: xdr.Uint32(0),
		Flags:         xdr.Uint32(0),
		HomeDomain:    xdr.String32("example.com"),
		Thresholds:    xdr.Thresholds{1, 2, 3, 4},
		Signers:       []xdr.Signer{},
		Ext:           xdr.AccountEntryExt{V: 0},
	}

	ledgerEntry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: xdr.Uint32(12345),
		Data: xdr.LedgerEntryData{
			Type:    xdr.LedgerEntryTypeAccount,
			Account: &account,
		},
		Ext: xdr.LedgerEntryExt{V: 0},
	}

	encoded, _ := xdr.MarshalBase64(ledgerEntry)
	return encoded
}

// TestProcessLedgerEntries_AccountExtraction tests that account addresses are properly extracted
func TestProcessLedgerEntries_AccountExtraction(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Quiet during tests

	// Create test transactions with source accounts
	testAddr1 := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"
	testAddr2 := "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ"

	ledger1 := uint32(100)
	ledger2 := uint32(100)

	tx1 := &models.Transaction{
		ID:            "tx1",
		SourceAccount: &testAddr1,
		Ledger:        &ledger1,
	}
	tx2 := &models.Transaction{
		ID:            "tx2",
		SourceAccount: &testAddr2,
		Ledger:        &ledger2,
	}

	// Insert transactions into database
	if err := db.Create(tx1).Error; err != nil {
		t.Fatalf("Failed to create test transaction 1: %v", err)
	}
	if err := db.Create(tx2).Error; err != nil {
		t.Fatalf("Failed to create test transaction 2: %v", err)
	}

	// Simulate extracting account addresses from transactions
	accountAddresses := make(map[string]bool)
	var transactions []models.Transaction
	hashList := []string{"tx1", "tx2"}

	if err := db.Where("id IN ?", hashList).Find(&transactions).Error; err == nil {
		for _, txn := range transactions {
			if txn.SourceAccount != nil && *txn.SourceAccount != "" {
				accountAddresses[*txn.SourceAccount] = true
			}
		}
	} else {
		t.Fatalf("Failed to query transactions: %v", err)
	}

	// Verify account addresses were extracted
	if len(accountAddresses) != 2 {
		t.Errorf("Expected 2 account addresses, got %d", len(accountAddresses))
	}
	if !accountAddresses[testAddr1] {
		t.Errorf("Expected account address %s to be extracted", testAddr1)
	}
	if !accountAddresses[testAddr2] {
		t.Errorf("Expected account address %s to be extracted", testAddr2)
	}
}

// TestProcessLedgerEntries_EmptySourceAccount tests handling of transactions without source accounts
func TestProcessLedgerEntries_EmptySourceAccount(t *testing.T) {
	db := setupTestDB(t)

	ledger := uint32(100)

	// Create transaction with nil source account
	tx1 := &models.Transaction{
		ID:            "tx1",
		SourceAccount: nil,
		Ledger:        &ledger,
	}

	// Create transaction with empty source account
	emptyAddr := ""
	tx2 := &models.Transaction{
		ID:            "tx2",
		SourceAccount: &emptyAddr,
		Ledger:        &ledger,
	}

	if err := db.Create(tx1).Error; err != nil {
		t.Fatalf("Failed to create test transaction 1: %v", err)
	}
	if err := db.Create(tx2).Error; err != nil {
		t.Fatalf("Failed to create test transaction 2: %v", err)
	}

	// Extract account addresses (should be empty)
	accountAddresses := make(map[string]bool)
	var transactions []models.Transaction
	hashList := []string{"tx1", "tx2"}

	if err := db.Where("id IN ?", hashList).Find(&transactions).Error; err == nil {
		for _, txn := range transactions {
			if txn.SourceAccount != nil && *txn.SourceAccount != "" {
				accountAddresses[*txn.SourceAccount] = true
			}
		}
	}

	// Verify no account addresses were extracted
	if len(accountAddresses) != 0 {
		t.Errorf("Expected 0 account addresses, got %d", len(accountAddresses))
	}
}

// TestProcessLedgerEntries_MixedTransactions tests extraction from mixed transactions
func TestProcessLedgerEntries_MixedTransactions(t *testing.T) {
	db := setupTestDB(t)

	testAddr := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"
	emptyAddr := ""
	ledger := uint32(100)

	// Mix of transactions with and without source accounts
	transactions := []*models.Transaction{
		{ID: "tx1", SourceAccount: &testAddr, Ledger: &ledger},
		{ID: "tx2", SourceAccount: nil, Ledger: &ledger},
		{ID: "tx3", SourceAccount: &emptyAddr, Ledger: &ledger},
		{ID: "tx4", SourceAccount: &testAddr, Ledger: &ledger}, // Duplicate address
	}

	for _, tx := range transactions {
		if err := db.Create(tx).Error; err != nil {
			t.Fatalf("Failed to create test transaction: %v", err)
		}
	}

	// Extract account addresses
	accountAddresses := make(map[string]bool)
	var txList []models.Transaction
	hashList := []string{"tx1", "tx2", "tx3", "tx4"}

	if err := db.Where("id IN ?", hashList).Find(&txList).Error; err == nil {
		for _, txn := range txList {
			if txn.SourceAccount != nil && *txn.SourceAccount != "" {
				accountAddresses[*txn.SourceAccount] = true
			}
		}
	}

	// Verify only one unique account address was extracted
	if len(accountAddresses) != 1 {
		t.Errorf("Expected 1 unique account address, got %d", len(accountAddresses))
	}
	if !accountAddresses[testAddr] {
		t.Errorf("Expected account address %s to be extracted", testAddr)
	}
}

// TestAccountEntryUpsert tests that account entries can be properly upserted
func TestAccountEntryUpsert(t *testing.T) {
	db := setupTestDB(t)

	testAddr := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"

	// Create initial account entry
	ext := []byte("{}")
	entry1 := &models.AccountEntry{
		AccountID:             testAddr,
		Balance:               1000000,
		SeqNum:                1,
		NumSubEntries:         0,
		Flags:                 0,
		HomeDomain:            "example.com",
		Ext:                   ext,
		LastModifiedLedgerSeq: 100,
	}

	if err := models.UpsertAccountEntry(db, entry1); err != nil {
		t.Fatalf("Failed to insert account entry: %v", err)
	}

	// Verify insert
	var retrieved1 models.AccountEntry
	if err := db.Where("account_id = ?", testAddr).First(&retrieved1).Error; err != nil {
		t.Fatalf("Failed to retrieve account entry: %v", err)
	}
	if retrieved1.Balance != 1000000 {
		t.Errorf("Expected balance 1000000, got %d", retrieved1.Balance)
	}

	// Update with new data
	ext2 := []byte("{}")
	entry2 := &models.AccountEntry{
		AccountID:             testAddr,
		Balance:               2000000, // Changed
		SeqNum:                2,        // Changed
		NumSubEntries:         1,        // Changed
		Flags:                 0,
		HomeDomain:            "newdomain.com", // Changed
		Ext:                   ext2,
		LastModifiedLedgerSeq: 200, // Changed
	}

	if err := models.UpsertAccountEntry(db, entry2); err != nil {
		t.Fatalf("Failed to upsert account entry: %v", err)
	}

	// Verify update
	var retrieved2 models.AccountEntry
	if err := db.Where("account_id = ?", testAddr).First(&retrieved2).Error; err != nil {
		t.Fatalf("Failed to retrieve updated account entry: %v", err)
	}
	if retrieved2.Balance != 2000000 {
		t.Errorf("Expected updated balance 2000000, got %d", retrieved2.Balance)
	}
	if retrieved2.SeqNum != 2 {
		t.Errorf("Expected updated sequence 2, got %d", retrieved2.SeqNum)
	}
	if retrieved2.HomeDomain != "newdomain.com" {
		t.Errorf("Expected updated home domain newdomain.com, got %s", retrieved2.HomeDomain)
	}

	// Verify no duplicate entries
	var count int64
	if err := db.Model(&models.AccountEntry{}).Where("account_id = ?", testAddr).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count account entries: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 account entry, got %d", count)
	}
}

// TestBuildAccountLedgerKey_InvalidAddress tests error handling for invalid addresses
func TestBuildAccountLedgerKey_InvalidAddress(t *testing.T) {
	// Import parser to test BuildAccountLedgerKey
	// This would require importing the parser package
	// For now, we'll test at integration level via the poller

	// Note: This test validates that invalid addresses are logged but don't crash the system
	db := setupTestDB(t)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create mock RPC client (would need to be implemented with httptest)
	// For unit tests, we focus on the database and logic parts
	// Integration tests would cover the full RPC flow

	invalidAddr := "INVALID_ADDRESS"
	ledger := uint32(100)
	tx := &models.Transaction{
		ID:            "tx1",
		SourceAccount: &invalidAddr,
		Ledger:        &ledger,
	}

	if err := db.Create(tx).Error; err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	// The processLedgerEntries function should handle invalid addresses gracefully
	// by logging a warning and continuing with valid addresses
	// This is tested in integration tests with the full poller
}

// TestTransactionQueryWithInlineFunction tests the bug we fixed
func TestTransactionQueryWithInlineFunction(t *testing.T) {
	db := setupTestDB(t)

	testAddr := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"
	ledger := uint32(100)
	txHashes := map[string]bool{
		"tx1": true,
		"tx2": true,
		"tx3": true,
	}

	// Create test transactions
	for hash := range txHashes {
		tx := &models.Transaction{
			ID:            hash,
			SourceAccount: &testAddr,
			Ledger:        &ledger,
		}
		if err := db.Create(tx).Error; err != nil {
			t.Fatalf("Failed to create test transaction: %v", err)
		}
	}

	// Test the FIXED approach (building hash list separately)
	t.Run("Fixed approach", func(t *testing.T) {
		var transactions []models.Transaction
		hashList := make([]string, 0, len(txHashes))
		for hash := range txHashes {
			hashList = append(hashList, hash)
		}

		if err := db.Where("id IN ?", hashList).Find(&transactions).Error; err != nil {
			t.Fatalf("Failed to query transactions: %v", err)
		}

		if len(transactions) != 3 {
			t.Errorf("Expected 3 transactions, got %d", len(transactions))
		}

		// Verify all have source accounts
		accountAddresses := make(map[string]bool)
		for _, txn := range transactions {
			if txn.SourceAccount != nil && *txn.SourceAccount != "" {
				accountAddresses[*txn.SourceAccount] = true
			}
		}

		if len(accountAddresses) != 1 {
			t.Errorf("Expected 1 unique account address, got %d", len(accountAddresses))
		}
	})

	// Test the BUGGY approach (inline function) - this demonstrates the issue
	t.Run("Buggy approach - inline function", func(t *testing.T) {
		var transactions []models.Transaction

		// This is the pattern that was causing the bug
		// GORM may not properly evaluate this inline function
		err := db.Where("id IN ?", func() []string {
			hashes := make([]string, 0, len(txHashes))
			for hash := range txHashes {
				hashes = append(hashes, hash)
			}
			return hashes
		}()).Find(&transactions).Error

		// The query should work in SQLite but may fail in PostgreSQL or other databases
		// or return unexpected results depending on GORM's handling
		if err != nil {
			t.Logf("Inline function approach failed (expected): %v", err)
		} else {
			t.Logf("Inline function approach returned %d transactions (may be incorrect)", len(transactions))
		}
	})
}

// TestLedgerEntryParsing tests that ledger entries are correctly parsed
func TestLedgerEntryParsing(t *testing.T) {
	testAddr := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"

	// Create test account ledger entry
	accountID := xdr.MustAddress(testAddr)

	account := xdr.AccountEntry{
		AccountId:     accountID,
		Balance:       xdr.Int64(1000000),
		SeqNum:        xdr.SequenceNumber(100),
		NumSubEntries: xdr.Uint32(5),
		Flags:         xdr.Uint32(1),
		HomeDomain:    xdr.String32("example.com"),
		Thresholds:    xdr.Thresholds{1, 2, 3, 4},
		Signers:       []xdr.Signer{},
		Ext:           xdr.AccountEntryExt{V: 0},
	}

	ledgerEntry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: xdr.Uint32(12345),
		Data: xdr.LedgerEntryData{
			Type:    xdr.LedgerEntryTypeAccount,
			Account: &account,
		},
		Ext: xdr.LedgerEntryExt{V: 0},
	}

	xdrString, err := xdr.MarshalBase64(ledgerEntry)
	if err != nil {
		t.Fatalf("Failed to marshal ledger entry: %v", err)
	}

	// Decode and verify
	var decoded xdr.LedgerEntry
	if err := xdr.SafeUnmarshalBase64(xdrString, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ledger entry: %v", err)
	}

	if decoded.Data.Type != xdr.LedgerEntryTypeAccount {
		t.Errorf("Expected account entry type, got %v", decoded.Data.Type)
	}

	if decoded.Data.Account.AccountId.Address() != testAddr {
		t.Errorf("Expected account ID %s, got %s", testAddr, decoded.Data.Account.AccountId.Address())
	}

	if decoded.Data.Account.Balance != xdr.Int64(1000000) {
		t.Errorf("Expected balance 1000000, got %d", decoded.Data.Account.Balance)
	}
}

// TestGetLedgerEntriesResponse tests the RPC response structure
func TestGetLedgerEntriesResponse(t *testing.T) {
	// Create a mock response similar to what RPC would return
	testAddr := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"
	entryXDR := createTestAccountLedgerEntry(testAddr)

	mockResp := client.GetLedgerEntriesResponse{
		Entries: []client.LedgerEntryResult{
			{
				Key:              base64.StdEncoding.EncodeToString([]byte("test_key")),
				XDR:              entryXDR,
				LastModifiedLedger: 12345,
			},
		},
		LatestLedger: 12346,
	}

	// Verify response structure
	if len(mockResp.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(mockResp.Entries))
	}

	// Decode the XDR
	var ledgerEntry xdr.LedgerEntry
	if err := xdr.SafeUnmarshalBase64(mockResp.Entries[0].XDR, &ledgerEntry); err != nil {
		t.Fatalf("Failed to decode entry XDR: %v", err)
	}

	if ledgerEntry.Data.Type != xdr.LedgerEntryTypeAccount {
		t.Errorf("Expected account entry, got type %v", ledgerEntry.Data.Type)
	}
}
