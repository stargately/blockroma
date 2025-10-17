package parser

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"

	"github.com/blockroma/soroban-indexer/pkg/models"
)

// ExtractContractDataFromMeta extracts contract data entries from transaction metadata
// This is the passive indexing approach - it extracts contract storage changes from successful transactions
func ExtractContractDataFromMeta(txHash string, metaXdr string) ([]*models.ContractDataEntry, error) {
	// Decode metadata XDR
	data, err := base64.StdEncoding.DecodeString(metaXdr)
	if err != nil {
		return nil, fmt.Errorf("decode meta xdr: %w", err)
	}

	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("unmarshal transaction meta: %w", err)
	}

	var entries []*models.ContractDataEntry

	// Check if this is a Soroban transaction (V3 or V4)
	isSoroban := meta.V4 != nil || meta.V3 != nil
	if !isSoroban {
		// Not a Soroban transaction (classic Stellar)
		return entries, nil
	}

	// Get the changes from either V3 or V4
	var txChangesAfter xdr.LedgerEntryChanges
	var operations []interface{} // Can be []OperationMeta or []OperationMetaV2

	if meta.V4 != nil {
		txChangesAfter = meta.V4.TxChangesAfter
		for _, op := range meta.V4.Operations {
			operations = append(operations, op)
		}
	} else if meta.V3 != nil {
		txChangesAfter = meta.V3.TxChangesAfter
		for _, op := range meta.V3.Operations {
			operations = append(operations, op)
		}
	}

	// Process ledger entry changes after the transaction
	// This includes created, updated, and restored contract data entries
	for _, change := range txChangesAfter {
		ledgerEntry, ok := getLedgerEntryFromChange(change)
		if !ok {
			continue
		}

		entry, err := extractContractDataFromLedgerEntry(ledgerEntry)
		if err != nil {
			continue
		}
		if entry != nil {
			entries = append(entries, entry)
		}
	}

	// Also process operation-level changes
	for _, op := range operations {
		var changes xdr.LedgerEntryChanges
		if v3Op, ok := op.(xdr.OperationMeta); ok {
			changes = v3Op.Changes
		} else if v4Op, ok := op.(xdr.OperationMetaV2); ok {
			changes = v4Op.Changes
		}

		for _, change := range changes {
			ledgerEntry, ok := getLedgerEntryFromChange(change)
			if !ok {
				continue
			}

			entry, err := extractContractDataFromLedgerEntry(ledgerEntry)
			if err != nil {
				continue
			}
			if entry != nil {
				entries = append(entries, entry)
			}
		}
	}

	return entries, nil
}

// getLedgerEntryFromChange extracts a ledger entry from a ledger entry change
func getLedgerEntryFromChange(change xdr.LedgerEntryChange) (xdr.LedgerEntry, bool) {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
		if created, ok := change.GetCreated(); ok {
			return created, true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
		if updated, ok := change.GetUpdated(); ok {
			return updated, true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryRestored:
		if restored, ok := change.GetRestored(); ok {
			return restored, true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryState:
		// State changes represent the state before the transaction
		// We want the new state (created/updated), not the old state
		return xdr.LedgerEntry{}, false
	case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
		// Removed entries are not useful for indexing
		return xdr.LedgerEntry{}, false
	}
	return xdr.LedgerEntry{}, false
}

// extractContractDataFromLedgerEntry extracts contract data from a ledger entry
func extractContractDataFromLedgerEntry(entry xdr.LedgerEntry) (*models.ContractDataEntry, error) {
	// Only process contract data entries
	if entry.Data.Type != xdr.LedgerEntryTypeContractData {
		return nil, nil
	}

	contractData := entry.Data.ContractData
	if contractData == nil {
		return nil, fmt.Errorf("contract data is nil")
	}

	// Extract contract ID from the contract address
	if contractData.Contract.Type != xdr.ScAddressTypeScAddressTypeContract {
		// Not a contract address (might be an account address)
		return nil, nil
	}

	contractIDHash := contractData.Contract.ContractId
	if contractIDHash == nil {
		return nil, fmt.Errorf("contract ID is nil")
	}

	// Encode contract ID to C... format
	contractID, err := strkey.Encode(strkey.VersionByteContract, (*contractIDHash)[:])
	if err != nil {
		return nil, fmt.Errorf("encode contract ID: %w", err)
	}

	// Convert key and value to interface{} first
	keyInterface := ScValToInterface(contractData.Key)
	valInterface := ScValToInterface(contractData.Val)

	// Marshal to JSON bytes for JSONB storage
	keyBytes, err := json.Marshal(keyInterface)
	if err != nil {
		return nil, fmt.Errorf("marshal key to JSON: %w", err)
	}
	valBytes, err := json.Marshal(valInterface)
	if err != nil {
		return nil, fmt.Errorf("marshal val to JSON: %w", err)
	}

	// Create ledger key hash for unique identification
	ledgerKey, err := entry.LedgerKey()
	if err != nil {
		return nil, fmt.Errorf("get ledger key: %w", err)
	}

	bin, err := ledgerKey.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal ledger key: %w", err)
	}

	keyHash := sha256.Sum256(bin)
	hexKey := hex.EncodeToString(keyHash[:])

	// Marshal key and value to base64 XDR
	keyXDR, err := xdr.MarshalBase64(contractData.Key)
	if err != nil {
		return nil, fmt.Errorf("marshal key xdr: %w", err)
	}

	valXDR, err := xdr.MarshalBase64(contractData.Val)
	if err != nil {
		return nil, fmt.Errorf("marshal value xdr: %w", err)
	}

	// Determine durability
	durability := "persistent"
	if contractData.Durability == xdr.ContractDataDurabilityTemporary {
		durability = "temporary"
	}

	return &models.ContractDataEntry{
		KeyHash:    hexKey,
		ContractID: contractID,
		Key:        models.JSONB(keyBytes),
		KeyXdr:     keyXDR,
		Val:        models.JSONB(valBytes),
		ValXdr:     valXDR,
		Durability: durability,
	}, nil
}

// ExtractContractInstanceFromMeta extracts contract instance data from transaction metadata
// This specifically looks for contract deployment transactions and extracts the contract instance,
// which contains token metadata (name, symbol, decimals)
func ExtractContractInstanceFromMeta(txHash string, metaXdr string) (map[string]*models.TokenMetadata, error) {
	// Decode metadata XDR
	data, err := base64.StdEncoding.DecodeString(metaXdr)
	if err != nil {
		return nil, fmt.Errorf("decode meta xdr: %w", err)
	}

	var meta xdr.TransactionMeta
	if err := xdr.SafeUnmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("unmarshal transaction meta: %w", err)
	}

	metadataMap := make(map[string]*models.TokenMetadata)

	// Check if this is a Soroban transaction (V3 or V4)
	isSoroban := meta.V4 != nil || meta.V3 != nil
	if !isSoroban {
		return metadataMap, nil
	}

	// Get the changes from either V3 or V4
	var txChangesAfter xdr.LedgerEntryChanges
	var operations []interface{}

	if meta.V4 != nil {
		txChangesAfter = meta.V4.TxChangesAfter
		for _, op := range meta.V4.Operations {
			operations = append(operations, op)
		}
	} else if meta.V3 != nil {
		txChangesAfter = meta.V3.TxChangesAfter
		for _, op := range meta.V3.Operations {
			operations = append(operations, op)
		}
	}

	// Process ledger entry changes after the transaction
	for _, change := range txChangesAfter {
		ledgerEntry, ok := getLedgerEntryFromChange(change)
		if !ok {
			continue
		}

		metadata := extractTokenMetadataFromLedgerEntry(ledgerEntry)
		if metadata != nil {
			metadataMap[metadata.ContractID] = metadata
		}
	}

	// Also process operation-level changes
	for _, op := range operations {
		var changes xdr.LedgerEntryChanges
		if v3Op, ok := op.(xdr.OperationMeta); ok {
			changes = v3Op.Changes
		} else if v4Op, ok := op.(xdr.OperationMetaV2); ok {
			changes = v4Op.Changes
		}

		for _, change := range changes {
			ledgerEntry, ok := getLedgerEntryFromChange(change)
			if !ok {
				continue
			}

			metadata := extractTokenMetadataFromLedgerEntry(ledgerEntry)
			if metadata != nil {
				metadataMap[metadata.ContractID] = metadata
			}
		}
	}

	return metadataMap, nil
}

// extractTokenMetadataFromLedgerEntry extracts token metadata from a contract data ledger entry
// if the entry is a contract instance with token metadata
func extractTokenMetadataFromLedgerEntry(entry xdr.LedgerEntry) *models.TokenMetadata {
	// Only process contract data entries
	if entry.Data.Type != xdr.LedgerEntryTypeContractData {
		return nil
	}

	contractData := entry.Data.ContractData
	if contractData == nil {
		return nil
	}

	// Extract contract ID from the contract address
	if contractData.Contract.Type != xdr.ScAddressTypeScAddressTypeContract {
		return nil
	}

	contractIDHash := contractData.Contract.ContractId
	if contractIDHash == nil {
		return nil
	}

	// Encode contract ID to C... format
	contractID, err := strkey.Encode(strkey.VersionByteContract, (*contractIDHash)[:])
	if err != nil {
		return nil
	}

	// Check if the key indicates this is a contract instance
	// Contract instance keys are special - they're stored as ScvLedgerKeyContractInstance
	keyInterface := ScValToInterface(contractData.Key)
	valInterface := ScValToInterface(contractData.Val)

	// Try to parse token metadata from this entry
	metadata := ParseTokenMetadata(contractID, keyInterface, valInterface)
	return metadata
}
