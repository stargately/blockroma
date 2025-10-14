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

	// Log every call to verify function is being invoked
	fmt.Printf("[TRACE] ExtractContractDataFromMeta called for TX %s\n", txHash[:8])

	// Check metadata version and log what we actually have
	metaType := fmt.Sprintf("V%d", meta.V)
	isSoroban := false
	if meta.V4 != nil {
		metaType = "V4 (Soroban)"
		isSoroban = true
	} else if meta.V3 != nil {
		metaType = "V3 (Soroban)"
		isSoroban = true
	} else if meta.V2 != nil {
		metaType = "V2 (classic)"
	} else if meta.V1 != nil {
		metaType = "V1 (classic)"
	}
	fmt.Printf("[TRACE] TX %s: Metadata type=%s\n", txHash[:8], metaType)

	// Check if this is a Soroban transaction (V3 or V4)
	if !isSoroban {
		// Not a Soroban transaction
		fmt.Printf("[TRACE] TX %s: Not Soroban metadata (classic Stellar transaction)\n", txHash[:8])
		return entries, nil
	}

	fmt.Printf("[TRACE] TX %s: Has %s metadata (Soroban transaction)\n", txHash[:8], metaType)

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

	// Count total changes
	totalChanges := len(txChangesAfter)
	for _, op := range operations {
		if v3Op, ok := op.(xdr.OperationMeta); ok {
			totalChanges += len(v3Op.Changes)
		} else if v4Op, ok := op.(xdr.OperationMetaV2); ok {
			totalChanges += len(v4Op.Changes)
		}
	}

	// Log summary of what we found
	fmt.Printf("[TRACE] TX %s: TxChangesAfter=%d, Operations=%d, totalChanges=%d\n",
		txHash[:8], len(txChangesAfter), len(operations), totalChanges)

	// Early return if no changes at all
	if totalChanges == 0 {
		fmt.Printf("[TRACE] TX %s: No ledger changes, skipping\n", txHash[:8])
		return entries, nil
	}

	fmt.Printf("[DEBUG] TX %s: Processing %d ledger changes\n", txHash[:8], totalChanges)

	// Process ledger entry changes after the transaction
	// This includes created, updated, and restored contract data entries
	contractDataChanges := 0
	otherChanges := 0
	for _, change := range txChangesAfter {
		// Get the ledger entry from the change
		ledgerEntry, ok := getLedgerEntryFromChange(change)
		if !ok {
			continue
		}

		// Count entry types
		if ledgerEntry.Data.Type == xdr.LedgerEntryTypeContractData {
			contractDataChanges++
		} else {
			otherChanges++
		}

		entry, err := extractContractDataFromLedgerEntry(ledgerEntry)
		if err != nil {
			// Log error but continue processing
			continue
		}
		if entry != nil {
			entries = append(entries, entry)
		}
	}

	// Log what we found
	if contractDataChanges > 0 || otherChanges > 0 {
		fmt.Printf("[DEBUG] TX %s: TxChangesAfter has %d ContractData entries, %d other entries\n",
			txHash, contractDataChanges, otherChanges)
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

	// Final summary
	fmt.Printf("[TRACE] TX %s: Extracted %d contract data entries\n", txHash[:8], len(entries))
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
