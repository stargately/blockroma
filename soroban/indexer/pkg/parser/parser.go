package parser

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
	"github.com/blockroma/soroban-indexer/pkg/client"
	"github.com/blockroma/soroban-indexer/pkg/models"
	"github.com/blockroma/soroban-indexer/pkg/models/util"
)

// ParseEvent converts RPC event to database model
func ParseEvent(event client.Event) (*models.Event, error) {
	// Parse topics
	topics := make([]interface{}, 0, len(event.Topic))
	for _, topicXdr := range event.Topic {
		topic, err := parseScVal(topicXdr)
		if err != nil {
			// Keep raw XDR if parse fails
			topics = append(topics, topicXdr)
		} else {
			topics = append(topics, topic)
		}
	}

	// Parse value
	value, err := parseScVal(event.Value)
	if err != nil {
		// Keep raw XDR if parse fails
		value = event.Value
	}

	// Convert topics to JSON
	topicJSON, _ := json.Marshal(topics)
	valueJSON, _ := json.Marshal(value)

	// Extract tx index from event ID
	// ID format: "0000123456-0000000001"
	txIndex := int32(0)
	parts := strings.Split(event.ID, "-")
	if len(parts) == 2 {
		var ledger, index uint64
		fmt.Sscanf(parts[0], "%d", &ledger)
		fmt.Sscanf(parts[1], "%d", &index)
		txIndex = int32(index)
	}

	return &models.Event{
		ID:                       event.ID,
		TxIndex:                  txIndex,
		EventType:                event.Type,
		Ledger:                   int32(event.Ledger),
		LedgerClosedAt:           event.LedgerClosedAt,
		ContractID:               event.ContractID,
		PagingToken:              event.PagingToken,
		Topic:                    string(topicJSON),
		Value:                    string(valueJSON),
		InSuccessfulContractCall: event.InSuccessfulContractCall,
		LastModifiedLedgerSeq:    event.Ledger,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}, nil
}

// ComputeTransactionHash computes the transaction hash from the envelope XDR
// This is the proper way to get the transaction hash, as the RPC may return empty hashes
func ComputeTransactionHash(envelopeXDR string, networkPassphrase string) (string, error) {
	// Decode the envelope
	envelope, err := decodeEnvelope(envelopeXDR)
	if err != nil {
		return "", fmt.Errorf("decode envelope: %w", err)
	}

	// Hash the transaction envelope with the network passphrase
	hash, err := network.HashTransactionInEnvelope(*envelope, networkPassphrase)
	if err != nil {
		return "", fmt.Errorf("hash transaction: %w", err)
	}

	return hex.EncodeToString(hash[:]), nil
}

// ParseTransaction converts RPC transaction to database model
// Uses the transaction hash from RPC response
func ParseTransaction(tx client.Transaction) (*models.Transaction, error) {
	return ParseTransactionWithHash(tx, tx.Hash)
}

// ParseTransactionWithHash converts RPC transaction to database model using provided hash
// This is useful when the RPC response hash might be unreliable
func ParseTransactionWithHash(tx client.Transaction, txHash string) (*models.Transaction, error) {
	// Decode envelope to extract source account and fee
	envelope, err := decodeEnvelope(tx.EnvelopeXdr)
	if err != nil {
		return nil, fmt.Errorf("decode envelope: %w", err)
	}

	sourceAccount := ""
	fee := int32(0)
	sequence := int64(0)
	var memo *util.TypeItem
	var preconditions *util.Preconditions
	var signatures *util.Signatures

	// Extract details from envelope
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		if v0, ok := envelope.GetV0(); ok {
			// V0 uses Ed25519 public key, convert to address
			sourceAccount, _ = strkey.Encode(strkey.VersionByteAccountID, v0.Tx.SourceAccountEd25519[:])
			fee = int32(v0.Tx.Fee)
			sequence = int64(v0.Tx.SeqNum)

			// Parse memo
			memo = parseMemo(v0.Tx.Memo)

			// Parse preconditions (V0 only has time bounds)
			preconditions = parsePreconditionsV0(v0.Tx.TimeBounds)

			// Parse signatures
			signatures = parseSignatures(v0.Signatures)
		}
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		if v1, ok := envelope.GetV1(); ok {
			sourceAccount = v1.Tx.SourceAccount.ToAccountId().Address()
			fee = int32(v1.Tx.Fee)
			sequence = int64(v1.Tx.SeqNum)

			// Parse memo
			memo = parseMemo(v1.Tx.Memo)

			// Parse preconditions (V1 has full preconditions support)
			preconditions = parsePreconditionsV1(v1.Tx.Cond)

			// Parse signatures
			signatures = parseSignatures(v1.Signatures)
		}
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		if fb, ok := envelope.GetFeeBump(); ok {
			sourceAccount = fb.Tx.FeeSource.ToAccountId().Address()
			fee = int32(fb.Tx.Fee)

			// Fee bump doesn't have sequence, memo, or preconditions
			// Get from inner tx
			if innerV1, ok := fb.Tx.InnerTx.GetV1(); ok {
				sequence = int64(innerV1.Tx.SeqNum)
				memo = parseMemo(innerV1.Tx.Memo)
				preconditions = parsePreconditionsV1(innerV1.Tx.Cond)
			}

			// Parse signatures from fee bump envelope
			signatures = parseSignatures(fb.Signatures)
		}
	}

	// Decode result to get fee charged
	result, err := decodeResult(tx.ResultXdr)
	feeCharged := int32(0)
	if err == nil {
		feeCharged = int32(result.FeeCharged)
	}

	ledger := tx.Ledger
	appOrder := tx.ApplicationOrder
	ledgerCreatedAt := tx.LedgerCloseTime

	return &models.Transaction{
		ID:               txHash,  // Use the provided hash parameter, not tx.Hash
		Status:           tx.Status,
		Ledger:           &ledger,
		LedgerCreatedAt:  &ledgerCreatedAt,
		ApplicationOrder: &appOrder,
		SourceAccount:    &sourceAccount,
		Fee:              &fee,
		FeeCharged:       &feeCharged,
		Sequence:         &sequence,
		Memo:             memo,
		Preconditions:    preconditions,
		Signatures:       signatures,
		CreatedAt:        time.Unix(tx.LedgerCloseTime, 0),
		UpdatedAt:        time.Now(),
	}, nil
}

// parseScVal decodes base64 XDR ScVal to Go interface
func parseScVal(xdrStr string) (interface{}, error) {
	if xdrStr == "" {
		return nil, nil
	}

	data, err := base64.StdEncoding.DecodeString(xdrStr)
	if err != nil {
		return nil, err
	}

	var scVal xdr.ScVal
	if err := xdr.SafeUnmarshal(data, &scVal); err != nil {
		return nil, err
	}

	return ScValToInterface(scVal), nil
}

// ScValToInterface converts XDR ScVal to Go interface
func ScValToInterface(val xdr.ScVal) interface{} {
	switch val.Type {
	case xdr.ScValTypeScvBool:
		return val.MustB()
	case xdr.ScValTypeScvVoid:
		return nil
	case xdr.ScValTypeScvU32:
		return val.MustU32()
	case xdr.ScValTypeScvI32:
		return val.MustI32()
	case xdr.ScValTypeScvU64:
		return val.MustU64()
	case xdr.ScValTypeScvI64:
		return val.MustI64()
	case xdr.ScValTypeScvU128:
		parts := val.MustU128()
		// Convert to big integer string for consistency with soroban-rpc-indexer
		hi := uint64(parts.Hi)
		lo := uint64(parts.Lo)
		result := (hi << 64) | lo
		return fmt.Sprintf("%d", result)
	case xdr.ScValTypeScvI128:
		parts := val.MustI128()
		// Convert to big integer string for consistency with soroban-rpc-indexer
		hi := int64(parts.Hi)
		lo := uint64(parts.Lo)
		// Note: This is simplified; proper int128 handling needs big.Int for large values
		if hi >= 0 {
			result := (uint64(hi) << 64) | lo
			return fmt.Sprintf("%d", result)
		} else {
			// For negative numbers, use two's complement
			result := (uint64(hi) << 64) | lo
			return fmt.Sprintf("%d", int64(result))
		}
	case xdr.ScValTypeScvBytes:
		return val.MustBytes()
	case xdr.ScValTypeScvString:
		return val.MustStr()
	case xdr.ScValTypeScvSymbol:
		return val.MustSym()
	case xdr.ScValTypeScvVec:
		vec := *val.MustVec()
		result := make([]interface{}, len(vec))
		for i, v := range vec {
			result[i] = ScValToInterface(v)
		}
		return result
	case xdr.ScValTypeScvMap:
		m := *val.MustMap()
		// Convert to array of key-value pairs (like soroban-rpc-indexer)
		result := make([]interface{}, 0, len(m))
		for _, entry := range m {
			key := ScValToInterface(entry.Key)
			val := ScValToInterface(entry.Val)
			result = append(result, map[string]interface{}{
				"key":   key,
				"value": val,
			})
		}
		return result
	case xdr.ScValTypeScvAddress:
		addr := val.MustAddress()
		addrStr, _ := addr.String()
		return addrStr
	case xdr.ScValTypeScvLedgerKeyContractInstance:
		// Return a proper JSON object for contract instance key
		return map[string]interface{}{
			"type": "LedgerKeyContractInstance",
		}
	case xdr.ScValTypeScvContractInstance:
		// Parse contract instance to match expected format
		instance := val.MustInstance()
		data := make(map[string]interface{})

		// Add executable info
		executable := make(map[string]interface{})
		executable["type"] = instance.Executable.Type.String()
		if instance.Executable.WasmHash != nil {
			executable["wasmHash"] = hex.EncodeToString((*instance.Executable.WasmHash)[:])
		}
		data["executable"] = executable

		// Add storage (convert ScMap to array of key-value pairs)
		if instance.Storage != nil {
			storage := make([]interface{}, 0, len(*instance.Storage))
			for _, entry := range *instance.Storage {
				key := ScValToInterface(entry.Key)
				val := ScValToInterface(entry.Val)
				storage = append(storage, map[string]interface{}{
					"key":   key,
					"value": val,
				})
			}
			data["storage"] = storage
		}

		return data
	default:
		// For unknown types, return nil instead of .String() to avoid invalid JSON
		// val.String() can return Go syntax like "<nil>" which breaks JSON encoding
		return nil
	}
}

// decodeEnvelope decodes base64 XDR transaction envelope
func decodeEnvelope(xdrStr string) (*xdr.TransactionEnvelope, error) {
	data, err := base64.StdEncoding.DecodeString(xdrStr)
	if err != nil {
		return nil, err
	}

	var envelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshal(data, &envelope); err != nil {
		return nil, err
	}

	return &envelope, nil
}

// decodeResult decodes base64 XDR transaction result
func decodeResult(xdrStr string) (*xdr.TransactionResult, error) {
	data, err := base64.StdEncoding.DecodeString(xdrStr)
	if err != nil {
		return nil, err
	}

	var result xdr.TransactionResult
	if err := xdr.SafeUnmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// parseMemo converts XDR memo to TypeItem
func parseMemo(memo xdr.Memo) *util.TypeItem {
	switch memo.Type {
	case xdr.MemoTypeMemoNone:
		return nil // No memo
	case xdr.MemoTypeMemoText:
		text := memo.MustText()
		return &util.TypeItem{
			Type:      "text",
			ItemValue: text,
		}
	case xdr.MemoTypeMemoId:
		id := memo.MustId()
		return &util.TypeItem{
			Type:      "id",
			ItemValue: fmt.Sprintf("%d", id),
		}
	case xdr.MemoTypeMemoHash:
		hash := memo.MustHash()
		return &util.TypeItem{
			Type:      "hash",
			ItemValue: hex.EncodeToString(hash[:]),
		}
	case xdr.MemoTypeMemoReturn:
		retHash := memo.MustRetHash()
		return &util.TypeItem{
			Type:      "return",
			ItemValue: hex.EncodeToString(retHash[:]),
		}
	default:
		return nil
	}
}

// parsePreconditionsV1 converts XDR preconditions (v1) to Preconditions
func parsePreconditionsV1(cond xdr.Preconditions) *util.Preconditions {
	switch cond.Type {
	case xdr.PreconditionTypePrecondNone:
		return nil // No preconditions
	case xdr.PreconditionTypePrecondTime:
		// V1 with only time bounds
		tb := cond.MustTimeBounds()
		return &util.Preconditions{
			TimeBounds: &util.Bonds{
				Min: int64(tb.MinTime),
				Max: int64(tb.MaxTime),
			},
		}
	case xdr.PreconditionTypePrecondV2:
		// V2 with full preconditions
		v2 := cond.MustV2()
		result := &util.Preconditions{}

		// Time bounds
		if v2.TimeBounds != nil {
			result.TimeBounds = &util.Bonds{
				Min: int64(v2.TimeBounds.MinTime),
				Max: int64(v2.TimeBounds.MaxTime),
			}
		}

		// Ledger bounds
		if v2.LedgerBounds != nil {
			result.LedgerBounds = &util.Bonds{
				Min: int64(v2.LedgerBounds.MinLedger),
				Max: int64(v2.LedgerBounds.MaxLedger),
			}
		}

		// Min sequence number
		if v2.MinSeqNum != nil {
			minSeq := int64(*v2.MinSeqNum)
			result.MinSeqNum = &minSeq
		}

		// Min sequence age
		minAge := int64(v2.MinSeqAge)
		result.MinSeqAge = &minAge

		// Min sequence ledger gap
		minGap := int32(v2.MinSeqLedgerGap)
		result.MinSeqLedgerGap = &minGap

		// Extra signers
		if len(v2.ExtraSigners) > 0 {
			signers := make([]util.SignerKey, 0, len(v2.ExtraSigners))
			for _, signer := range v2.ExtraSigners {
				signerKey := util.SignerKey{Type: signer.Type.String()}

				// Populate the appropriate field based on signer type
				switch signer.Type {
				case xdr.SignerKeyTypeSignerKeyTypeEd25519:
					signerKey.Ed25519 = base64.StdEncoding.EncodeToString(signer.Ed25519[:])
				case xdr.SignerKeyTypeSignerKeyTypePreAuthTx:
					signerKey.PreAuthTx = base64.StdEncoding.EncodeToString(signer.PreAuthTx[:])
				case xdr.SignerKeyTypeSignerKeyTypeHashX:
					signerKey.HashX = base64.StdEncoding.EncodeToString(signer.HashX[:])
				case xdr.SignerKeyTypeSignerKeyTypeEd25519SignedPayload:
					if payload, ok := signer.GetEd25519SignedPayload(); ok {
						signerKey.Ed25519SignedPayload = base64.StdEncoding.EncodeToString(payload.Ed25519[:])
					}
				}

				signers = append(signers, signerKey)
			}
			result.ExtraSigners = &signers
		}

		return result
	default:
		return nil
	}
}

// parsePreconditionsV0 converts V0 time bounds to Preconditions
func parsePreconditionsV0(tb *xdr.TimeBounds) *util.Preconditions {
	// V0 TimeBounds can be nil
	if tb == nil {
		return nil
	}
	return &util.Preconditions{
		TimeBounds: &util.Bonds{
			Min: int64(tb.MinTime),
			Max: int64(tb.MaxTime),
		},
	}
}

// parseSignatures converts XDR signatures to Signatures slice
func parseSignatures(sigs []xdr.DecoratedSignature) *util.Signatures {
	if len(sigs) == 0 {
		return nil
	}

	signatures := make(util.Signatures, 0, len(sigs))
	for _, sig := range sigs {
		signatures = append(signatures, util.Signature{
			Hint:      hex.EncodeToString(sig.Hint[:]),
			Signature: base64.StdEncoding.EncodeToString(sig.Signature),
		})
	}

	return &signatures
}
