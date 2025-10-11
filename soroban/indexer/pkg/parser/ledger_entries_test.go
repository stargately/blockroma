package parser

import (
	"encoding/base64"
	"testing"

	"github.com/blockroma/soroban-indexer/pkg/models"
	"github.com/stellar/go/xdr"
)

func TestBuildAccountLedgerKey(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "valid account address",
			address: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			wantErr: false,
		},
		{
			name:    "another valid address",
			address: "GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWHF",
			wantErr: false,
		},
		{
			name:    "invalid address - wrong format",
			address: "INVALID_ADDRESS",
			wantErr: true,
		},
		{
			name:    "empty address",
			address: "",
			wantErr: true,
		},
		{
			name:    "contract address instead of account",
			address: "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABSC4",
			wantErr: true, // Should fail because it's a contract, not an account
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := BuildAccountLedgerKey(tt.address)

			if tt.wantErr {
				if err == nil {
					t.Errorf("BuildAccountLedgerKey() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("BuildAccountLedgerKey() error = %v", err)
			}

			if key == "" {
				t.Error("BuildAccountLedgerKey() returned empty key")
			}

			// Verify the key is valid base64
			decoded, err := base64.StdEncoding.DecodeString(key)
			if err != nil {
				t.Errorf("BuildAccountLedgerKey() returned invalid base64: %v", err)
			}

			// Verify the key can be unmarshaled as a LedgerKey
			var ledgerKey xdr.LedgerKey
			if err := xdr.SafeUnmarshal(decoded, &ledgerKey); err != nil {
				t.Errorf("BuildAccountLedgerKey() returned invalid XDR: %v", err)
			}

			// Verify it's an account type
			if ledgerKey.Type != xdr.LedgerEntryTypeAccount {
				t.Errorf("BuildAccountLedgerKey() wrong type = %v, want Account", ledgerKey.Type)
			}

			// Verify the account ID matches
			if ledgerKey.Account != nil {
				accountAddr := ledgerKey.Account.AccountId.Address()
				if accountAddr != tt.address {
					t.Errorf("BuildAccountLedgerKey() account address = %v, want %v", accountAddr, tt.address)
				}
			} else {
				t.Error("BuildAccountLedgerKey() returned nil Account")
			}
		})
	}
}

func TestBuildAccountLedgerKey_RoundTrip(t *testing.T) {
	// Test that we can build a key, decode it, and get back the same address
	originalAddress := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"

	// Build the key
	key, err := BuildAccountLedgerKey(originalAddress)
	if err != nil {
		t.Fatalf("BuildAccountLedgerKey() error = %v", err)
	}

	// Decode the key
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Fatalf("Failed to decode key: %v", err)
	}

	// Unmarshal as LedgerKey
	var ledgerKey xdr.LedgerKey
	if err := xdr.SafeUnmarshal(decoded, &ledgerKey); err != nil {
		t.Fatalf("Failed to unmarshal ledger key: %v", err)
	}

	// Extract the account ID
	if ledgerKey.Account == nil {
		t.Fatal("Ledger key Account is nil")
	}

	recoveredAddress := ledgerKey.Account.AccountId.Address()
	if recoveredAddress != originalAddress {
		t.Errorf("Round trip failed: got %v, want %v", recoveredAddress, originalAddress)
	}
}

func TestParseAccountEntry(t *testing.T) {
	// Create a mock account ledger entry
	accountID := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

	account := xdr.AccountEntry{
		AccountId:     accountID,
		Balance:       10000000000, // 1000 XLM
		SeqNum:        123456,
		NumSubEntries: 5,
		Flags:         1,
		HomeDomain:    xdr.String32("example.com"),
		Thresholds:    xdr.Thresholds{1, 2, 3, 4},
	}

	entry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: 12345,
		Data: xdr.LedgerEntryData{
			Type:    xdr.LedgerEntryTypeAccount,
			Account: &account,
		},
	}

	// Parse the entry
	result := ParseAccountEntry(entry)

	// Verify the parsed result
	if result == nil {
		t.Fatal("ParseAccountEntry() returned nil")
	}

	if result.AccountID != "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H" {
		t.Errorf("AccountID = %v, want GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", result.AccountID)
	}

	if result.Balance != 10000000000 {
		t.Errorf("Balance = %v, want 10000000000", result.Balance)
	}

	if result.SeqNum != 123456 {
		t.Errorf("SeqNum = %v, want 123456", result.SeqNum)
	}

	if result.NumSubEntries != 5 {
		t.Errorf("NumSubEntries = %v, want 5", result.NumSubEntries)
	}

	if result.Flags != 1 {
		t.Errorf("Flags = %v, want 1", result.Flags)
	}

	if result.HomeDomain != "example.com" {
		t.Errorf("HomeDomain = %v, want example.com", result.HomeDomain)
	}

	if result.LastModifiedLedgerSeq != 12345 {
		t.Errorf("LastModifiedLedgerSeq = %v, want 12345", result.LastModifiedLedgerSeq)
	}

	if len(result.Thresholds) != 4 {
		t.Errorf("Thresholds length = %v, want 4", len(result.Thresholds))
	}
}

func TestParseAccountEntry_WithSigners(t *testing.T) {
	// Create account with signers
	accountID := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	signerKey := xdr.MustSigner("GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWHF")

	account := xdr.AccountEntry{
		AccountId:     accountID,
		Balance:       1000000000,
		SeqNum:        100,
		NumSubEntries: 1,
		Signers: []xdr.Signer{
			{
				Key:    signerKey,
				Weight: 10,
			},
		},
	}

	entry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: 1000,
		Data: xdr.LedgerEntryData{
			Type:    xdr.LedgerEntryTypeAccount,
			Account: &account,
		},
	}

	result := ParseAccountEntry(entry)

	if result == nil {
		t.Fatal("ParseAccountEntry() returned nil")
	}

	// Check that signers is not nil/empty JSON
	if result.Signers == nil {
		t.Error("Signers should not be nil")
	}

	signersJSON, ok := result.Signers.([]byte)
	if !ok {
		t.Errorf("Signers type = %T, want []byte", result.Signers)
	}

	if len(signersJSON) < 3 { // At least "[]" or more
		t.Errorf("Signers JSON too short: %s", string(signersJSON))
	}

	if string(signersJSON) == "[]" {
		t.Error("Signers JSON should not be empty array when signers exist")
	}
}

func TestParseLedgerEntry_InvalidXDR(t *testing.T) {
	// Test with invalid XDR
	_, err := ParseLedgerEntry("invalid-base64-xdr")
	if err == nil {
		t.Error("ParseLedgerEntry() expected error for invalid XDR")
	}
}

func TestParseLedgerEntry_AccountType(t *testing.T) {
	// Create a valid account entry
	accountID := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

	account := xdr.AccountEntry{
		AccountId:     accountID,
		Balance:       5000000000,
		SeqNum:        999,
		NumSubEntries: 2,
		Flags:         0,
		HomeDomain:    xdr.String32("test.com"),
		Thresholds:    xdr.Thresholds{1, 1, 1, 1},
	}

	entry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: 5000,
		Data: xdr.LedgerEntryData{
			Type:    xdr.LedgerEntryTypeAccount,
			Account: &account,
		},
	}

	// Marshal to XDR
	xdrBytes, err := entry.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal entry: %v", err)
	}

	xdrString := base64.StdEncoding.EncodeToString(xdrBytes)

	// Parse it
	results, err := ParseLedgerEntry(xdrString)
	if err != nil {
		t.Fatalf("ParseLedgerEntry() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("ParseLedgerEntry() returned %d results, want 1", len(results))
	}

	// Check it's an AccountEntry
	accountEntry, ok := results[0].(*models.AccountEntry)
	if !ok {
		t.Fatalf("ParseLedgerEntry() returned type %T, want *models.AccountEntry", results[0])
	}

	if accountEntry.Balance != 5000000000 {
		t.Errorf("Balance = %v, want 5000000000", accountEntry.Balance)
	}

	if accountEntry.SeqNum != 999 {
		t.Errorf("SeqNum = %v, want 999", accountEntry.SeqNum)
	}
}

// Note: We need to import models package to test AccountEntry type
// Since this would create circular dependency, we keep the test simple
// and just check that the function returns the right interface

func TestGetSponsoringID(t *testing.T) {
	// Test with entry that has no sponsor
	entry := &xdr.LedgerEntry{
		Ext: xdr.LedgerEntryExt{
			V: 0,
		},
	}

	sponsor := getSponsoringID(entry)
	if sponsor != "" {
		t.Errorf("getSponsoringID() = %v, want empty string for no sponsor", sponsor)
	}
}

func TestPriceToString(t *testing.T) {
	tests := []struct {
		name  string
		price xdr.Price
		want  string
	}{
		{
			name:  "simple fraction",
			price: xdr.Price{N: 1, D: 2},
			want:  "1/2",
		},
		{
			name:  "whole number",
			price: xdr.Price{N: 5, D: 1},
			want:  "5/1",
		},
		{
			name:  "zero denominator",
			price: xdr.Price{N: 10, D: 0},
			want:  "0",
		},
		{
			name:  "large numbers",
			price: xdr.Price{N: 1000000, D: 999999},
			want:  "1000000/999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := priceToString(tt.price)
			if got != tt.want {
				t.Errorf("priceToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildClaimableBalanceLedgerKey(t *testing.T) {
	tests := []struct {
		name      string
		balanceID string
		wantErr   bool
	}{
		{
			name:      "valid balance ID",
			balanceID: "000000000000000000000000000000000000000000000000000000000000dead",
			wantErr:   false,
		},
		{
			name:      "another valid balance ID",
			balanceID: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			wantErr:   false,
		},
		{
			name:      "invalid hex",
			balanceID: "invalid-hex",
			wantErr:   true,
		},
		{
			name:      "empty balance ID",
			balanceID: "",
			wantErr:   true,
		},
		{
			name:      "wrong length - too short",
			balanceID: "dead",
			wantErr:   true,
		},
		{
			name:      "wrong length - too long",
			balanceID: "000000000000000000000000000000000000000000000000000000000000deadbeef",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := BuildClaimableBalanceLedgerKey(tt.balanceID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("BuildClaimableBalanceLedgerKey() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("BuildClaimableBalanceLedgerKey() error = %v", err)
			}

			if key == "" {
				t.Error("BuildClaimableBalanceLedgerKey() returned empty key")
			}

			// Verify the key is valid base64
			decoded, err := base64.StdEncoding.DecodeString(key)
			if err != nil {
				t.Errorf("BuildClaimableBalanceLedgerKey() returned invalid base64: %v", err)
			}

			// Verify the key can be unmarshaled as a LedgerKey
			var ledgerKey xdr.LedgerKey
			if err := xdr.SafeUnmarshal(decoded, &ledgerKey); err != nil {
				t.Errorf("BuildClaimableBalanceLedgerKey() returned invalid XDR: %v", err)
			}

			// Verify it's a claimable balance type
			if ledgerKey.Type != xdr.LedgerEntryTypeClaimableBalance {
				t.Errorf("BuildClaimableBalanceLedgerKey() wrong type = %v, want ClaimableBalance", ledgerKey.Type)
			}
		})
	}
}

func TestBuildClaimableBalanceLedgerKey_RoundTrip(t *testing.T) {
	// Test that we can build a key, decode it, and get back the same balance ID
	originalBalanceID := "000000000000000000000000000000000000000000000000000000000000dead"

	// Build the key
	key, err := BuildClaimableBalanceLedgerKey(originalBalanceID)
	if err != nil {
		t.Fatalf("BuildClaimableBalanceLedgerKey() error = %v", err)
	}

	// Decode the key
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Fatalf("Failed to decode key: %v", err)
	}

	// Unmarshal as LedgerKey
	var ledgerKey xdr.LedgerKey
	if err := xdr.SafeUnmarshal(decoded, &ledgerKey); err != nil {
		t.Fatalf("Failed to unmarshal ledger key: %v", err)
	}

	// Extract the balance ID
	if ledgerKey.ClaimableBalance == nil {
		t.Fatal("Ledger key ClaimableBalance is nil")
	}

	recoveredBalanceID := ledgerKey.ClaimableBalance.BalanceId.V0.HexString()
	if recoveredBalanceID != originalBalanceID {
		t.Errorf("Round trip failed: got %v, want %v", recoveredBalanceID, originalBalanceID)
	}
}

func TestExtractClaimableBalanceIDs(t *testing.T) {
	// Create a transaction with a ClaimClaimableBalance operation
	sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

	// Create a claimable balance ID
	var hash xdr.Hash
	hashBytes, _ := base64.StdEncoding.DecodeString("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA3q0=")
	copy(hash[:], hashBytes)

	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &hash,
	}

	claimOp := xdr.ClaimClaimableBalanceOp{
		BalanceId: balanceID,
	}

	operation := xdr.Operation{
		SourceAccount: nil,
		Body: xdr.OperationBody{
			Type:                     xdr.OperationTypeClaimClaimableBalance,
			ClaimClaimableBalanceOp: &claimOp,
		},
	}

	tx := xdr.Transaction{
		SourceAccount: sourceAccount.ToMuxedAccount(),
		Fee:           100,
		SeqNum:        123456,
		Operations:    []xdr.Operation{operation},
	}

	envelope := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: tx,
		},
	}

	// Marshal to XDR
	xdrBytes, err := envelope.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal envelope: %v", err)
	}

	xdrString := base64.StdEncoding.EncodeToString(xdrBytes)

	// Extract balance IDs
	balanceIDs, err := ExtractClaimableBalanceIDs(xdrString)
	if err != nil {
		t.Fatalf("ExtractClaimableBalanceIDs() error = %v", err)
	}

	if len(balanceIDs) != 1 {
		t.Fatalf("ExtractClaimableBalanceIDs() returned %d IDs, want 1", len(balanceIDs))
	}

	expectedID := "000000000000000000000000000000000000000000000000000000000000dead"
	if balanceIDs[0] != expectedID {
		t.Errorf("ExtractClaimableBalanceIDs() returned %v, want %v", balanceIDs[0], expectedID)
	}
}

func TestExtractClaimableBalanceIDs_NoOperations(t *testing.T) {
	// Create a transaction with no operations
	sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

	tx := xdr.Transaction{
		SourceAccount: sourceAccount.ToMuxedAccount(),
		Fee:           100,
		SeqNum:        123456,
		Operations:    []xdr.Operation{},
	}

	envelope := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: tx,
		},
	}

	xdrBytes, err := envelope.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal envelope: %v", err)
	}

	xdrString := base64.StdEncoding.EncodeToString(xdrBytes)

	balanceIDs, err := ExtractClaimableBalanceIDs(xdrString)
	if err != nil {
		t.Fatalf("ExtractClaimableBalanceIDs() error = %v", err)
	}

	if len(balanceIDs) != 0 {
		t.Errorf("ExtractClaimableBalanceIDs() returned %d IDs, want 0", len(balanceIDs))
	}
}
