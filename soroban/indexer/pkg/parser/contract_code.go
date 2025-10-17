package parser

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/blockroma/soroban-indexer/pkg/models"
	"github.com/stellar/go/xdr"
)

// ExtractContractCode extracts WASM code from InvokeHostFunction operations
// Returns contract code if the operation is an UploadContractWasm, nil otherwise
func ExtractContractCode(txHash string, ledger uint32, ledgerCloseTime int64, body xdr.OperationBody) (*models.ContractCode, error) {
	// Check if this is an InvokeHostFunction operation
	if body.Type != xdr.OperationTypeInvokeHostFunction {
		return nil, nil // Not a contract code operation
	}

	op, ok := body.GetInvokeHostFunctionOp()
	if !ok {
		return nil, nil
	}

	// Check if this is an UploadContractWasm host function
	if op.HostFunction.Type != xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm {
		return nil, nil // Not an upload operation
	}

	// Extract WASM bytecode
	wasm, ok := op.HostFunction.GetWasm()
	if !ok {
		return nil, fmt.Errorf("failed to get WASM from UploadContractWasm")
	}

	// Compute SHA256 hash of WASM
	hash := sha256.Sum256(wasm)
	hashHex := fmt.Sprintf("%x", hash)

	return &models.ContractCode{
		Hash:       hashHex,
		Wasm:       wasm,
		DeployedAt: time.Unix(ledgerCloseTime, 0),
		Ledger:     ledger,
		TxHash:     txHash,
		SizeBytes:  len(wasm),
	}, nil
}

// ExtractContractCodeFromEnvelope extracts all contract code from a transaction envelope
func ExtractContractCodeFromEnvelope(txHash string, ledger uint32, ledgerCloseTime int64, envelopeXdr string) ([]*models.ContractCode, error) {
	// Decode envelope
	data, err := base64.StdEncoding.DecodeString(envelopeXdr)
	if err != nil {
		return nil, fmt.Errorf("decode envelope xdr: %w", err)
	}

	var envelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("unmarshal envelope: %w", err)
	}

	// Extract operations based on envelope type
	var operations []xdr.Operation
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		if v0, ok := envelope.GetV0(); ok {
			operations = v0.Tx.Operations
		}
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		if v1, ok := envelope.GetV1(); ok {
			operations = v1.Tx.Operations
		}
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		if fb, ok := envelope.GetFeeBump(); ok {
			if innerV1, ok := fb.Tx.InnerTx.GetV1(); ok {
				operations = innerV1.Tx.Operations
			}
		}
	default:
		return nil, fmt.Errorf("unsupported envelope type: %v", envelope.Type)
	}

	// Extract contract code from each operation
	var codes []*models.ContractCode
	for _, op := range operations {
		code, err := ExtractContractCode(txHash, ledger, ledgerCloseTime, op.Body)
		if err != nil {
			// Log error but continue processing
			continue
		}
		if code != nil {
			codes = append(codes, code)
		}
	}

	return codes, nil
}
