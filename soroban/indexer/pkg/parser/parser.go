package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
	"github.com/blockroma/soroban-indexer/pkg/client"
	"github.com/blockroma/soroban-indexer/pkg/models"
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

	// Extract details from envelope
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		if v0, ok := envelope.GetV0(); ok {
			// V0 uses Ed25519 public key, convert to address
			sourceAccount, _ = strkey.Encode(strkey.VersionByteAccountID, v0.Tx.SourceAccountEd25519[:])
			fee = int32(v0.Tx.Fee)
			sequence = int64(v0.Tx.SeqNum)
		}
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		if v1, ok := envelope.GetV1(); ok {
			sourceAccount = v1.Tx.SourceAccount.ToAccountId().Address()
			fee = int32(v1.Tx.Fee)
			sequence = int64(v1.Tx.SeqNum)
		}
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		if fb, ok := envelope.GetFeeBump(); ok {
			sourceAccount = fb.Tx.FeeSource.ToAccountId().Address()
			fee = int32(fb.Tx.Fee)
			// Fee bump doesn't have sequence, get from inner tx
			if innerV1, ok := fb.Tx.InnerTx.GetV1(); ok {
				sequence = int64(innerV1.Tx.SeqNum)
			}
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

	return scValToInterface(scVal), nil
}

// scValToInterface converts XDR ScVal to Go interface
func scValToInterface(val xdr.ScVal) interface{} {
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
		return map[string]interface{}{
			"lo": uint64(parts.Lo),
			"hi": uint64(parts.Hi),
		}
	case xdr.ScValTypeScvI128:
		parts := val.MustI128()
		return map[string]interface{}{
			"lo": uint64(parts.Lo),
			"hi": int64(parts.Hi),
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
			result[i] = scValToInterface(v)
		}
		return result
	case xdr.ScValTypeScvMap:
		m := *val.MustMap()
		result := make(map[string]interface{})
		for _, entry := range m {
			key := scValToInterface(entry.Key)
			val := scValToInterface(entry.Val)
			if keyStr, ok := key.(string); ok {
				result[keyStr] = val
			}
		}
		return result
	case xdr.ScValTypeScvAddress:
		addr := val.MustAddress()
		addrStr, _ := addr.String()
		return addrStr
	default:
		return val.String()
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
