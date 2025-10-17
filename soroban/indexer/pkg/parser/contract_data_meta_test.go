package parser

import (
	"encoding/base64"
	"testing"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

// TestExtractContractDataFromMeta tests basic contract data extraction
func TestExtractContractDataFromMeta(t *testing.T) {
	tests := []struct {
		name        string
		metaXdr     string
		wantEntries int
		wantErr     bool
	}{
		{
			name:        "empty metadata",
			metaXdr:     "",
			wantEntries: 0,
			wantErr:     true,
		},
		{
			name:        "invalid base64",
			metaXdr:     "not-valid-base64!!!",
			wantEntries: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := ExtractContractDataFromMeta("test-hash", tt.metaXdr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractContractDataFromMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(entries) != tt.wantEntries {
				t.Errorf("ExtractContractDataFromMeta() got %d entries, want %d", len(entries), tt.wantEntries)
			}
		})
	}
}

// TestExtractContractInstanceFromMeta tests contract instance extraction with real XDR structures
func TestExtractContractInstanceFromMeta(t *testing.T) {
	// Create a test contract ID
	contractIDBytes := [32]byte{}
	for i := range contractIDBytes {
		contractIDBytes[i] = byte(i)
	}
	contractID, _ := strkey.Encode(strkey.VersionByteContract, contractIDBytes[:])

	// Create a contract instance value with token metadata
	metadataMap := xdr.ScMap{
		{
			Key: xdr.ScVal{
				Type: xdr.ScValTypeScvString,
				Str:  scStringPtr("name"),
			},
			Val: xdr.ScVal{
				Type: xdr.ScValTypeScvString,
				Str:  scStringPtr("Test Token"),
			},
		},
		{
			Key: xdr.ScVal{
				Type: xdr.ScValTypeScvString,
				Str:  scStringPtr("symbol"),
			},
			Val: xdr.ScVal{
				Type: xdr.ScValTypeScvString,
				Str:  scStringPtr("TEST"),
			},
		},
		{
			Key: xdr.ScVal{
				Type: xdr.ScValTypeScvString,
				Str:  scStringPtr("decimal"),
			},
			Val: xdr.ScVal{
				Type: xdr.ScValTypeScvU32,
				U32:  uint32Ptr(7),
			},
		},
	}

	adminAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

	contractInstanceStorage := xdr.ScMap{
		{
			Key: xdr.ScVal{
				Type: xdr.ScValTypeScvSymbol,
				Sym:  scSymbolPtr("METADATA"),
			},
			Val: func() xdr.ScVal {
				mapPtr := &metadataMap
				return xdr.ScVal{
					Type: xdr.ScValTypeScvMap,
					Map:  &mapPtr,
				}
			}(),
		},
		{
			Key: xdr.ScVal{
				Type: xdr.ScValTypeScvSymbol,
				Sym:  scSymbolPtr("Admin"),
			},
			Val: xdr.ScVal{
				Type:    xdr.ScValTypeScvAddress,
				Address: &xdr.ScAddress{Type: xdr.ScAddressTypeScAddressTypeAccount, AccountId: &adminAccount},
			},
		},
	}

	contractInstance := xdr.ScVal{
		Type: xdr.ScValTypeScvContractInstance,
		Instance: &xdr.ScContractInstance{
			Executable: xdr.ContractExecutable{Type: 1}, // CONTRACT_EXECUTABLE_STELLAR_ASSET
			Storage:    &contractInstanceStorage,
		},
	}

	// Create contract instance key (special ledger key type)
	contractInstanceKey := xdr.ScVal{
		Type: xdr.ScValTypeScvLedgerKeyContractInstance,
	}

	// Create a contract data entry with the instance
	contractIDPtr := (*xdr.ContractId)(&contractIDBytes)
	contractAddress := xdr.ScAddress{
		Type:       xdr.ScAddressTypeScAddressTypeContract,
		ContractId: contractIDPtr,
	}

	contractData := &xdr.ContractDataEntry{
		Contract: contractAddress,
		Key:      contractInstanceKey,
		Val:      contractInstance,
		Durability: xdr.ContractDataDurabilityPersistent,
	}

	ledgerEntry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:         xdr.LedgerEntryTypeContractData,
			ContractData: contractData,
		},
	}

	// Create a ledger entry change
	ledgerChange := xdr.LedgerEntryChange{
		Type:    xdr.LedgerEntryChangeTypeLedgerEntryCreated,
		Created: &ledgerEntry,
	}

	// Create transaction metadata V3 with this change
	transactionMeta := xdr.TransactionMeta{
		V: 3,
		V3: &xdr.TransactionMetaV3{
			TxChangesAfter: xdr.LedgerEntryChanges{ledgerChange},
			Operations:     []xdr.OperationMeta{},
		},
	}

	// Encode to base64
	metaXdrBytes, err := transactionMeta.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal transaction meta: %v", err)
	}
	metaXdr := base64.StdEncoding.EncodeToString(metaXdrBytes)

	// Test extraction
	t.Run("extract contract instance metadata", func(t *testing.T) {
		metadataMap, err := ExtractContractInstanceFromMeta("test-tx-hash", metaXdr)
		if err != nil {
			t.Fatalf("ExtractContractInstanceFromMeta() error = %v", err)
		}

		if len(metadataMap) != 1 {
			t.Fatalf("Expected 1 metadata entry, got %d", len(metadataMap))
		}

		metadata, ok := metadataMap[contractID]
		if !ok {
			t.Fatalf("Expected metadata for contract %s, got keys: %v", contractID, getMapKeys(metadataMap))
		}

		if metadata.ContractID != contractID {
			t.Errorf("ContractID = %v, want %v", metadata.ContractID, contractID)
		}
		if metadata.Name != "Test Token" {
			t.Errorf("Name = %v, want Test Token", metadata.Name)
		}
		if metadata.Symbol != "TEST" {
			t.Errorf("Symbol = %v, want TEST", metadata.Symbol)
		}
		if metadata.Decimal != 7 {
			t.Errorf("Decimal = %v, want 7", metadata.Decimal)
		}
		if metadata.AdminAddress != "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H" {
			t.Errorf("AdminAddress = %v, want GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", metadata.AdminAddress)
		}
	})
}

// TestExtractContractInstanceFromMeta_NoInstance tests that non-instance contract data is ignored
func TestExtractContractInstanceFromMeta_NoInstance(t *testing.T) {
	// Create a regular contract data entry (not an instance)
	contractIDBytes := [32]byte{}
	for i := range contractIDBytes {
		contractIDBytes[i] = byte(i)
	}

	// Create a regular storage key (not instance key)
	regularKey := xdr.ScVal{
		Type: xdr.ScValTypeScvSymbol,
		Sym:  scSymbolPtr("Balance"),
	}

	regularValue := xdr.ScVal{
		Type: xdr.ScValTypeScvU128,
		U128: &xdr.UInt128Parts{Lo: 1000000, Hi: 0},
	}

	contractIDPtr := (*xdr.ContractId)(&contractIDBytes)
	contractAddress := xdr.ScAddress{
		Type:       xdr.ScAddressTypeScAddressTypeContract,
		ContractId: contractIDPtr,
	}

	contractData := &xdr.ContractDataEntry{
		Contract:   contractAddress,
		Key:        regularKey,
		Val:        regularValue,
		Durability: xdr.ContractDataDurabilityPersistent,
	}

	ledgerEntry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:         xdr.LedgerEntryTypeContractData,
			ContractData: contractData,
		},
	}

	ledgerChange := xdr.LedgerEntryChange{
		Type:    xdr.LedgerEntryChangeTypeLedgerEntryCreated,
		Created: &ledgerEntry,
	}

	transactionMeta := xdr.TransactionMeta{
		V: 3,
		V3: &xdr.TransactionMetaV3{
			TxChangesAfter: xdr.LedgerEntryChanges{ledgerChange},
			Operations:     []xdr.OperationMeta{},
		},
	}

	metaXdrBytes, _ := transactionMeta.MarshalBinary()
	metaXdr := base64.StdEncoding.EncodeToString(metaXdrBytes)

	t.Run("regular contract data should not extract metadata", func(t *testing.T) {
		metadataMap, err := ExtractContractInstanceFromMeta("test-tx-hash", metaXdr)
		if err != nil {
			t.Fatalf("ExtractContractInstanceFromMeta() error = %v", err)
		}

		if len(metadataMap) != 0 {
			t.Errorf("Expected 0 metadata entries for regular contract data, got %d", len(metadataMap))
		}
	})
}

// TestExtractContractInstanceFromMeta_V4 tests extraction from V4 transaction metadata
func TestExtractContractInstanceFromMeta_V4(t *testing.T) {
	// Create similar test data but with V4 metadata
	contractIDBytes := [32]byte{}
	for i := range contractIDBytes {
		contractIDBytes[i] = byte(i + 10) // Different ID for V4 test
	}
	contractID, _ := strkey.Encode(strkey.VersionByteContract, contractIDBytes[:])

	metadataMap := xdr.ScMap{
		{
			Key: xdr.ScVal{Type: xdr.ScValTypeScvString, Str: scStringPtr("name")},
			Val: xdr.ScVal{Type: xdr.ScValTypeScvString, Str: scStringPtr("V4 Token")},
		},
		{
			Key: xdr.ScVal{Type: xdr.ScValTypeScvString, Str: scStringPtr("symbol")},
			Val: xdr.ScVal{Type: xdr.ScValTypeScvString, Str: scStringPtr("V4T")},
		},
		{
			Key: xdr.ScVal{Type: xdr.ScValTypeScvString, Str: scStringPtr("decimal")},
			Val: xdr.ScVal{Type: xdr.ScValTypeScvU32, U32: uint32Ptr(8)},
		},
	}

	contractInstanceStorage := xdr.ScMap{
		{
			Key: xdr.ScVal{Type: xdr.ScValTypeScvSymbol, Sym: scSymbolPtr("METADATA")},
			Val: func() xdr.ScVal {
				mapPtr := &metadataMap
				return xdr.ScVal{Type: xdr.ScValTypeScvMap, Map: &mapPtr}
			}(),
		},
	}

	contractInstance := xdr.ScVal{
		Type: xdr.ScValTypeScvContractInstance,
		Instance: &xdr.ScContractInstance{
			Executable: xdr.ContractExecutable{Type: 1}, // CONTRACT_EXECUTABLE_STELLAR_ASSET
			Storage:    &contractInstanceStorage,
		},
	}

	contractInstanceKey := xdr.ScVal{Type: xdr.ScValTypeScvLedgerKeyContractInstance}

	contractIDPtr := (*xdr.ContractId)(&contractIDBytes)
	contractAddress := xdr.ScAddress{
		Type:       xdr.ScAddressTypeScAddressTypeContract,
		ContractId: contractIDPtr,
	}

	contractData := &xdr.ContractDataEntry{
		Contract:   contractAddress,
		Key:        contractInstanceKey,
		Val:        contractInstance,
		Durability: xdr.ContractDataDurabilityPersistent,
	}

	ledgerEntry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:         xdr.LedgerEntryTypeContractData,
			ContractData: contractData,
		},
	}

	ledgerChange := xdr.LedgerEntryChange{
		Type:    xdr.LedgerEntryChangeTypeLedgerEntryCreated,
		Created: &ledgerEntry,
	}

	// Use V4 transaction metadata
	transactionMeta := xdr.TransactionMeta{
		V: 4,
		V4: &xdr.TransactionMetaV4{
			TxChangesAfter: xdr.LedgerEntryChanges{ledgerChange},
			Operations:     []xdr.OperationMetaV2{},
		},
	}

	metaXdrBytes, _ := transactionMeta.MarshalBinary()
	metaXdr := base64.StdEncoding.EncodeToString(metaXdrBytes)

	t.Run("extract from V4 metadata", func(t *testing.T) {
		metadataMap, err := ExtractContractInstanceFromMeta("test-tx-hash-v4", metaXdr)
		if err != nil {
			t.Fatalf("ExtractContractInstanceFromMeta() error = %v", err)
		}

		if len(metadataMap) != 1 {
			t.Fatalf("Expected 1 metadata entry, got %d", len(metadataMap))
		}

		metadata, ok := metadataMap[contractID]
		if !ok {
			t.Fatalf("Expected metadata for contract %s", contractID)
		}

		if metadata.Name != "V4 Token" {
			t.Errorf("Name = %v, want V4 Token", metadata.Name)
		}
		if metadata.Symbol != "V4T" {
			t.Errorf("Symbol = %v, want V4T", metadata.Symbol)
		}
		if metadata.Decimal != 8 {
			t.Errorf("Decimal = %v, want 8", metadata.Decimal)
		}
	})
}

// TestExtractContractInstanceFromMeta_NonSoroban tests that non-Soroban transactions are handled
func TestExtractContractInstanceFromMeta_NonSoroban(t *testing.T) {
	// Create a V2 transaction metadata (classic Stellar, not Soroban)
	transactionMeta := xdr.TransactionMeta{
		V:  2,
		V2: &xdr.TransactionMetaV2{},
	}

	metaXdrBytes, _ := transactionMeta.MarshalBinary()
	metaXdr := base64.StdEncoding.EncodeToString(metaXdrBytes)

	t.Run("non-Soroban transaction should return empty map", func(t *testing.T) {
		metadataMap, err := ExtractContractInstanceFromMeta("test-tx-hash", metaXdr)
		if err != nil {
			t.Fatalf("ExtractContractInstanceFromMeta() error = %v", err)
		}

		if len(metadataMap) != 0 {
			t.Errorf("Expected 0 metadata entries for non-Soroban tx, got %d", len(metadataMap))
		}
	})
}

// Helper functions

func scStringPtr(s string) *xdr.ScString {
	v := xdr.ScString(s)
	return &v
}

func uint32Ptr(u uint32) *xdr.Uint32 {
	v := xdr.Uint32(u)
	return &v
}

func scSymbolPtr(s string) *xdr.ScSymbol {
	v := xdr.ScSymbol(s)
	return &v
}

func getMapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
