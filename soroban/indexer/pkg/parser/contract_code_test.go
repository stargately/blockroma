package parser

import (
	"encoding/base64"
	"testing"

	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExtractContractCode_NotInvokeHostFunction tests that non-InvokeHostFunction operations return nil
func TestExtractContractCode_NotInvokeHostFunction(t *testing.T) {
	tests := []struct {
		name string
		body xdr.OperationBody
	}{
		{
			name: "payment operation",
			body: xdr.OperationBody{Type: xdr.OperationTypePayment},
		},
		{
			name: "create account operation",
			body: xdr.OperationBody{Type: xdr.OperationTypeCreateAccount},
		},
		{
			name: "manage data operation",
			body: xdr.OperationBody{Type: xdr.OperationTypeManageData},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := ExtractContractCode("test-tx-hash", 12345, 1234567890, tt.body)
			assert.NoError(t, err)
			assert.Nil(t, code)
		})
	}
}

// TestExtractContractCode_NotUploadContractWasm tests that non-UploadContractWasm host functions return nil
func TestExtractContractCode_NotUploadContractWasm(t *testing.T) {
	// Create InvokeHostFunction operation with CreateContract
	hostFn := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeCreateContract,
	}

	body := xdr.OperationBody{
		Type: xdr.OperationTypeInvokeHostFunction,
		InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
			HostFunction: hostFn,
		},
	}

	code, err := ExtractContractCode("test-tx-hash", 12345, 1234567890, body)
	assert.NoError(t, err)
	assert.Nil(t, code)
}

// TestExtractContractCode_Success tests successful WASM extraction
func TestExtractContractCode_Success(t *testing.T) {
	// Create sample WASM bytecode
	wasm := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00} // WASM magic number

	// Create UploadContractWasm host function
	hostFn := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm,
		Wasm: &wasm,
	}

	body := xdr.OperationBody{
		Type: xdr.OperationTypeInvokeHostFunction,
		InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
			HostFunction: hostFn,
		},
	}

	txHash := "test-tx-hash-123"
	ledger := uint32(12345)
	ledgerCloseTime := int64(1234567890)

	code, err := ExtractContractCode(txHash, ledger, ledgerCloseTime, body)
	require.NoError(t, err)
	require.NotNil(t, code)

	// Verify contract code fields
	assert.Equal(t, txHash, code.TxHash)
	assert.Equal(t, ledger, code.Ledger)
	assert.Equal(t, wasm, code.Wasm)
	assert.Equal(t, len(wasm), code.SizeBytes)
	assert.NotEmpty(t, code.Hash)
	// Hash should be 64 characters (SHA256 hex)
	assert.Equal(t, 64, len(code.Hash))
}

// TestExtractContractCode_HashDeterministic tests that the same WASM produces the same hash
func TestExtractContractCode_HashDeterministic(t *testing.T) {
	wasm := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}

	hostFn := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm,
		Wasm: &wasm,
	}

	body := xdr.OperationBody{
		Type: xdr.OperationTypeInvokeHostFunction,
		InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
			HostFunction: hostFn,
		},
	}

	// Extract twice
	code1, err1 := ExtractContractCode("tx1", 100, 1000, body)
	require.NoError(t, err1)

	code2, err2 := ExtractContractCode("tx2", 200, 2000, body)
	require.NoError(t, err2)

	// Hashes should match
	assert.Equal(t, code1.Hash, code2.Hash)
}

// TestExtractContractCode_DifferentWasmDifferentHash tests that different WASM produces different hashes
func TestExtractContractCode_DifferentWasmDifferentHash(t *testing.T) {
	wasm1 := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	wasm2 := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x01} // Different last byte

	hostFn1 := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm,
		Wasm: &wasm1,
	}

	hostFn2 := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm,
		Wasm: &wasm2,
	}

	body1 := xdr.OperationBody{
		Type: xdr.OperationTypeInvokeHostFunction,
		InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
			HostFunction: hostFn1,
		},
	}

	body2 := xdr.OperationBody{
		Type: xdr.OperationTypeInvokeHostFunction,
		InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
			HostFunction: hostFn2,
		},
	}

	code1, err1 := ExtractContractCode("tx1", 100, 1000, body1)
	require.NoError(t, err1)

	code2, err2 := ExtractContractCode("tx2", 100, 1000, body2)
	require.NoError(t, err2)

	// Hashes should be different
	assert.NotEqual(t, code1.Hash, code2.Hash)
}

// TestExtractContractCodeFromEnvelope_NoOperations tests empty envelope
func TestExtractContractCodeFromEnvelope_NoOperations(t *testing.T) {
	// Create envelope with no operations
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

	// Encode to XDR
	data, err := envelope.MarshalBinary()
	require.NoError(t, err)
	envelopeXdr := base64.StdEncoding.EncodeToString(data)

	codes, err := ExtractContractCodeFromEnvelope("tx-hash", 100, 1000, envelopeXdr)
	assert.NoError(t, err)
	assert.Empty(t, codes)
}

// TestExtractContractCodeFromEnvelope_InvalidBase64 tests invalid base64 input
func TestExtractContractCodeFromEnvelope_InvalidBase64(t *testing.T) {
	codes, err := ExtractContractCodeFromEnvelope("tx-hash", 100, 1000, "invalid-base64!!!")
	assert.Error(t, err)
	assert.Nil(t, codes)
	assert.Contains(t, err.Error(), "decode envelope xdr")
}

// TestExtractContractCodeFromEnvelope_InvalidXDR tests invalid XDR data
func TestExtractContractCodeFromEnvelope_InvalidXDR(t *testing.T) {
	invalidData := []byte{0x00, 0x01, 0x02, 0x03}
	envelopeXdr := base64.StdEncoding.EncodeToString(invalidData)

	codes, err := ExtractContractCodeFromEnvelope("tx-hash", 100, 1000, envelopeXdr)
	assert.Error(t, err)
	assert.Nil(t, codes)
	assert.Contains(t, err.Error(), "unmarshal envelope")
}

// TestExtractContractCodeFromEnvelope_WithUploadWasm tests envelope with UploadContractWasm
func TestExtractContractCodeFromEnvelope_WithUploadWasm(t *testing.T) {
	wasm := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}

	hostFn := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm,
		Wasm: &wasm,
	}

	op := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeInvokeHostFunction,
			InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
				HostFunction: hostFn,
			},
		},
	}

	sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

	tx := xdr.Transaction{
		SourceAccount: sourceAccount.ToMuxedAccount(),
		Fee:           100,
		SeqNum:        123456,
		Operations:    []xdr.Operation{op},
	}

	envelope := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: tx,
		},
	}

	data, err := envelope.MarshalBinary()
	require.NoError(t, err)
	envelopeXdr := base64.StdEncoding.EncodeToString(data)

	txHash := "test-tx-hash"
	ledger := uint32(12345)
	ledgerCloseTime := int64(1234567890)

	codes, err := ExtractContractCodeFromEnvelope(txHash, ledger, ledgerCloseTime, envelopeXdr)
	require.NoError(t, err)
	require.Len(t, codes, 1)

	code := codes[0]
	assert.Equal(t, txHash, code.TxHash)
	assert.Equal(t, ledger, code.Ledger)
	assert.Equal(t, wasm, code.Wasm)
	assert.Equal(t, len(wasm), code.SizeBytes)
	assert.NotEmpty(t, code.Hash)
}

// TestExtractContractCodeFromEnvelope_MultipleOperations tests envelope with multiple operations
func TestExtractContractCodeFromEnvelope_MultipleOperations(t *testing.T) {
	wasm1 := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	wasm2 := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x01}

	hostFn1 := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm,
		Wasm: &wasm1,
	}

	hostFn2 := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm,
		Wasm: &wasm2,
	}

	// First operation: UploadContractWasm
	op1 := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeInvokeHostFunction,
			InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
				HostFunction: hostFn1,
			},
		},
	}

	// Second operation: UploadContractWasm
	op2 := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeInvokeHostFunction,
			InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
				HostFunction: hostFn2,
			},
		},
	}

	sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

	tx := xdr.Transaction{
		SourceAccount: sourceAccount.ToMuxedAccount(),
		Fee:           100,
		SeqNum:        123456,
		Operations:    []xdr.Operation{op1, op2},
	}

	envelope := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: tx,
		},
	}

	data, err := envelope.MarshalBinary()
	require.NoError(t, err)
	envelopeXdr := base64.StdEncoding.EncodeToString(data)

	codes, err := ExtractContractCodeFromEnvelope("tx-hash", 100, 1000, envelopeXdr)
	require.NoError(t, err)
	require.Len(t, codes, 2)

	// Verify both contract codes were extracted
	assert.Equal(t, wasm1, codes[0].Wasm)
	assert.Equal(t, wasm2, codes[1].Wasm)
	assert.NotEqual(t, codes[0].Hash, codes[1].Hash)
}

// TestExtractContractCodeFromEnvelope_V0Envelope tests V0 envelope type
func TestExtractContractCodeFromEnvelope_V0Envelope(t *testing.T) {
	wasm := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}

	hostFn := xdr.HostFunction{
		Type: xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm,
		Wasm: &wasm,
	}

	op := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeInvokeHostFunction,
			InvokeHostFunctionOp: &xdr.InvokeHostFunctionOp{
				HostFunction: hostFn,
			},
		},
	}

	// Create account ID for source account
	var accountID xdr.Uint256
	copy(accountID[:], []byte("test-account-id-for-v0-envelope-"))

	tx := xdr.TransactionV0{
		Operations: []xdr.Operation{op},
		SourceAccountEd25519: accountID,
	}

	envelope := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTxV0,
		V0: &xdr.TransactionV0Envelope{
			Tx: tx,
		},
	}

	data, err := envelope.MarshalBinary()
	require.NoError(t, err)
	envelopeXdr := base64.StdEncoding.EncodeToString(data)

	codes, err := ExtractContractCodeFromEnvelope("tx-hash", 100, 1000, envelopeXdr)
	require.NoError(t, err)
	require.Len(t, codes, 1)

	assert.Equal(t, wasm, codes[0].Wasm)
}
