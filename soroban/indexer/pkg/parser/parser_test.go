package parser

import (
	"testing"

	"github.com/stellar/go/xdr"
	"github.com/blockroma/soroban-indexer/pkg/client"
)

func TestParseEvent(t *testing.T) {
	tests := []struct {
		name    string
		event   client.Event
		wantErr bool
	}{
		{
			name: "basic event with string topic",
			event: client.Event{
				ID:                       "0000012345-0000000001",
				Type:                     "contract",
				Ledger:                   12345,
				LedgerClosedAt:           "2024-01-01T00:00:00Z",
				ContractID:               "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM",
				PagingToken:              "12345-1",
				Topic:                    []string{createStringScVal("transfer")},
				Value:                    createU64ScVal(1000000),
				InSuccessfulContractCall: true,
				TxHash:                   "abc123",
			},
			wantErr: false,
		},
		{
			name: "event with multiple topics",
			event: client.Event{
				ID:                       "0000012345-0000000002",
				Type:                     "contract",
				Ledger:                   12345,
				LedgerClosedAt:           "2024-01-01T00:00:00Z",
				ContractID:               "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM",
				Topic:                    []string{createStringScVal("mint"), createAddressScVal(), createAddressScVal()},
				Value:                    createU128ScVal(1000000, 0),
				InSuccessfulContractCall: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := ParseEvent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if event.ID != tt.event.ID {
					t.Errorf("ParseEvent() ID = %v, want %v", event.ID, tt.event.ID)
				}
				if event.Ledger != int32(tt.event.Ledger) {
					t.Errorf("ParseEvent() Ledger = %v, want %v", event.Ledger, tt.event.Ledger)
				}
				if event.ContractID != tt.event.ContractID {
					t.Errorf("ParseEvent() ContractID = %v, want %v", event.ContractID, tt.event.ContractID)
				}
			}
		})
	}
}

func TestScValToInterface(t *testing.T) {
	tests := []struct {
		name     string
		scVal    xdr.ScVal
		checkFn  func(interface{}) bool
	}{
		{
			name:     "bool true",
			scVal:    xdr.ScVal{Type: xdr.ScValTypeScvBool, B: boolPtr(true)},
			checkFn:  func(v interface{}) bool { return v == true },
		},
		{
			name:     "void",
			scVal:    xdr.ScVal{Type: xdr.ScValTypeScvVoid},
			checkFn:  func(v interface{}) bool { return v == nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScValToInterface(tt.scVal)
			if !tt.checkFn(result) {
				t.Errorf("ScValToInterface() check failed for %v (%T)", result, result)
			}
		})
	}
}

func TestScValToInterface_Vec(t *testing.T) {
	// Skip complex XDR structure test - covered by integration tests
	t.Skip("XDR Vec structure complex, covered by integration tests")
}

func TestParseTransaction(t *testing.T) {
	// Skip complex XDR transaction test - covered by integration tests
	t.Skip("XDR transaction structure complex, covered by integration tests")
}

func TestComputeTransactionHash(t *testing.T) {
	// Create a simple transaction envelope for testing
	createTestEnvelope := func() string {
		// Create a minimal V1 transaction
		sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

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

	tests := []struct {
		name              string
		envelopeXDR       string
		networkPassphrase string
		wantErr           bool
		checkHash         bool
	}{
		{
			name:              "valid testnet transaction",
			envelopeXDR:       createTestEnvelope(),
			networkPassphrase: "Test SDF Network ; September 2015",
			wantErr:           false,
			checkHash:         true,
		},
		{
			name:              "valid pubnet transaction",
			envelopeXDR:       createTestEnvelope(),
			networkPassphrase: "Public Global Stellar Network ; September 2015",
			wantErr:           false,
			checkHash:         true,
		},
		{
			name:              "invalid base64 envelope",
			envelopeXDR:       "not-valid-base64!!!",
			networkPassphrase: "Test SDF Network ; September 2015",
			wantErr:           true,
			checkHash:         false,
		},
		{
			name:              "empty envelope",
			envelopeXDR:       "",
			networkPassphrase: "Test SDF Network ; September 2015",
			wantErr:           true,
			checkHash:         false,
		},
		{
			name:              "valid envelope with empty passphrase",
			envelopeXDR:       createTestEnvelope(),
			networkPassphrase: "",
			wantErr:           true, // network.HashTransactionInEnvelope requires passphrase
			checkHash:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := ComputeTransactionHash(tt.envelopeXDR, tt.networkPassphrase)

			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeTransactionHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkHash && err == nil {
				// Verify hash format (64 hex characters)
				if len(hash) != 64 {
					t.Errorf("ComputeTransactionHash() hash length = %d, want 64", len(hash))
				}

				// Verify hash is valid hex
				for _, c := range hash {
					if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
						t.Errorf("ComputeTransactionHash() hash contains invalid hex character: %c", c)
						break
					}
				}

				t.Logf("Computed hash: %s", hash)
			}
		})
	}
}

func TestComputeTransactionHash_Deterministic(t *testing.T) {
	// Create test envelope
	sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

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

	envelopeXDR, _ := xdr.MarshalBase64(envelope)
	networkPassphrase := "Test SDF Network ; September 2015"

	// Compute hash multiple times
	hash1, err1 := ComputeTransactionHash(envelopeXDR, networkPassphrase)
	if err1 != nil {
		t.Fatalf("First hash computation failed: %v", err1)
	}

	hash2, err2 := ComputeTransactionHash(envelopeXDR, networkPassphrase)
	if err2 != nil {
		t.Fatalf("Second hash computation failed: %v", err2)
	}

	// Verify hashes are identical (deterministic)
	if hash1 != hash2 {
		t.Errorf("ComputeTransactionHash() is not deterministic: hash1=%s, hash2=%s", hash1, hash2)
	}
}

func TestComputeTransactionHash_DifferentPassphrases(t *testing.T) {
	// Create test envelope
	sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

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

	envelopeXDR, _ := xdr.MarshalBase64(envelope)

	// Compute with testnet passphrase
	testnetHash, err1 := ComputeTransactionHash(envelopeXDR, "Test SDF Network ; September 2015")
	if err1 != nil {
		t.Fatalf("Testnet hash computation failed: %v", err1)
	}

	// Compute with pubnet passphrase
	pubnetHash, err2 := ComputeTransactionHash(envelopeXDR, "Public Global Stellar Network ; September 2015")
	if err2 != nil {
		t.Fatalf("Pubnet hash computation failed: %v", err2)
	}

	// Verify hashes are different (same tx on different networks has different hash)
	if testnetHash == pubnetHash {
		t.Errorf("ComputeTransactionHash() should produce different hashes for different network passphrases")
	}

	t.Logf("Testnet hash: %s", testnetHash)
	t.Logf("Pubnet hash: %s", pubnetHash)
}

// Helper functions to create test XDR values

func createStringScVal(s string) string {
	str := xdr.ScString(s)
	scVal := xdr.ScVal{Type: xdr.ScValTypeScvString, Str: &str}
	encoded, _ := xdr.MarshalBase64(scVal)
	return encoded
}

func createU64ScVal(val uint64) string {
	v := xdr.Uint64(val)
	scVal := xdr.ScVal{Type: xdr.ScValTypeScvU64, U64: &v}
	encoded, _ := xdr.MarshalBase64(scVal)
	return encoded
}

func createU128ScVal(lo, hi uint64) string {
	scVal := xdr.ScVal{
		Type: xdr.ScValTypeScvU128,
		U128: &xdr.UInt128Parts{Lo: xdr.Uint64(lo), Hi: xdr.Uint64(hi)},
	}
	encoded, _ := xdr.MarshalBase64(scVal)
	return encoded
}

func createAddressScVal() string {
	// Create a test account address
	var accountID xdr.AccountId
	accountID.SetAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	addr := xdr.ScAddress{Type: xdr.ScAddressTypeScAddressTypeAccount, AccountId: &accountID}
	scVal := xdr.ScVal{Type: xdr.ScValTypeScvAddress, Address: &addr}
	encoded, _ := xdr.MarshalBase64(scVal)
	return encoded
}

// Helper pointer functions
func boolPtr(b bool) *bool { return &b }

// TestParseMemo tests memo parsing from XDR
func TestParseMemo(t *testing.T) {
	tests := []struct {
		name     string
		memo     xdr.Memo
		wantType string
		wantNil  bool
	}{
		{
			name:    "none memo",
			memo:    xdr.Memo{Type: xdr.MemoTypeMemoNone},
			wantNil: true,
		},
		{
			name:     "text memo",
			memo:     xdr.Memo{Type: xdr.MemoTypeMemoText, Text: strPtr("Hello World")},
			wantType: "text",
		},
		{
			name:     "id memo",
			memo:     xdr.Memo{Type: xdr.MemoTypeMemoId, Id: uint64Ptr(12345)},
			wantType: "id",
		},
		{
			name:     "hash memo",
			memo:     createHashMemo(),
			wantType: "hash",
		},
		{
			name:     "return memo",
			memo:     createReturnMemo(),
			wantType: "return",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMemo(tt.memo)
			if tt.wantNil {
				if result != nil {
					t.Errorf("parseMemo() = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Errorf("parseMemo() = nil, want non-nil")
					return
				}
				if result.Type != tt.wantType {
					t.Errorf("parseMemo() Type = %v, want %v", result.Type, tt.wantType)
				}
			}
		})
	}
}

// TestParsePreconditionsV1 tests precondition parsing for V1 transactions
func TestParsePreconditionsV1(t *testing.T) {
	tests := []struct {
		name    string
		cond    xdr.Preconditions
		wantNil bool
	}{
		{
			name:    "no preconditions",
			cond:    xdr.Preconditions{Type: xdr.PreconditionTypePrecondNone},
			wantNil: true,
		},
		{
			name: "time bounds only",
			cond: xdr.Preconditions{
				Type: xdr.PreconditionTypePrecondTime,
				TimeBounds: &xdr.TimeBounds{
					MinTime: xdr.TimePoint(1000),
					MaxTime: xdr.TimePoint(2000),
				},
			},
			wantNil: false,
		},
		{
			name:    "v2 preconditions",
			cond:    createV2Preconditions(),
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePreconditionsV1(tt.cond)
			if tt.wantNil {
				if result != nil {
					t.Errorf("parsePreconditionsV1() = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Errorf("parsePreconditionsV1() = nil, want non-nil")
				}
			}
		})
	}
}

// TestParsePreconditionsV0 tests precondition parsing for V0 transactions
func TestParsePreconditionsV0(t *testing.T) {
	tests := []struct {
		name    string
		tb      *xdr.TimeBounds
		wantNil bool
	}{
		{
			name:    "nil time bounds",
			tb:      nil,
			wantNil: true,
		},
		{
			name: "valid time bounds",
			tb: &xdr.TimeBounds{
				MinTime: xdr.TimePoint(1000),
				MaxTime: xdr.TimePoint(2000),
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePreconditionsV0(tt.tb)
			if tt.wantNil {
				if result != nil {
					t.Errorf("parsePreconditionsV0() = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Errorf("parsePreconditionsV0() = nil, want non-nil")
					return
				}
				if result.TimeBounds == nil {
					t.Errorf("parsePreconditionsV0() TimeBounds = nil, want non-nil")
					return
				}
				if result.TimeBounds.Min != int64(tt.tb.MinTime) {
					t.Errorf("parsePreconditionsV0() Min = %v, want %v", result.TimeBounds.Min, tt.tb.MinTime)
				}
				if result.TimeBounds.Max != int64(tt.tb.MaxTime) {
					t.Errorf("parsePreconditionsV0() Max = %v, want %v", result.TimeBounds.Max, tt.tb.MaxTime)
				}
			}
		})
	}
}

// TestParseSignatures tests signature parsing
func TestParseSignatures(t *testing.T) {
	tests := []struct {
		name    string
		sigs    []xdr.DecoratedSignature
		wantNil bool
		wantLen int
	}{
		{
			name:    "empty signatures",
			sigs:    []xdr.DecoratedSignature{},
			wantNil: true,
		},
		{
			name:    "single signature",
			sigs:    []xdr.DecoratedSignature{createSignature()},
			wantNil: false,
			wantLen: 1,
		},
		{
			name:    "multiple signatures",
			sigs:    []xdr.DecoratedSignature{createSignature(), createSignature()},
			wantNil: false,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSignatures(tt.sigs)
			if tt.wantNil {
				if result != nil {
					t.Errorf("parseSignatures() = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Errorf("parseSignatures() = nil, want non-nil")
					return
				}
				if len(*result) != tt.wantLen {
					t.Errorf("parseSignatures() length = %v, want %v", len(*result), tt.wantLen)
				}
			}
		})
	}
}

// Helper functions for creating test data

func strPtr(s string) *string {
	return &s
}

func uint64Ptr(u uint64) *xdr.Uint64 {
	v := xdr.Uint64(u)
	return &v
}

func createHashMemo() xdr.Memo {
	var hash xdr.Hash
	copy(hash[:], []byte("0123456789abcdef0123456789abcdef"))
	return xdr.Memo{Type: xdr.MemoTypeMemoHash, Hash: &hash}
}

func createReturnMemo() xdr.Memo {
	var hash xdr.Hash
	copy(hash[:], []byte("fedcba9876543210fedcba9876543210"))
	return xdr.Memo{Type: xdr.MemoTypeMemoReturn, RetHash: &hash}
}

func createV2Preconditions() xdr.Preconditions {
	minSeq := xdr.SequenceNumber(100)
	return xdr.Preconditions{
		Type: xdr.PreconditionTypePrecondV2,
		V2: &xdr.PreconditionsV2{
			TimeBounds: &xdr.TimeBounds{
				MinTime: xdr.TimePoint(1000),
				MaxTime: xdr.TimePoint(2000),
			},
			LedgerBounds: &xdr.LedgerBounds{
				MinLedger: xdr.Uint32(100),
				MaxLedger: xdr.Uint32(200),
			},
			MinSeqNum:       &minSeq,
			MinSeqAge:       xdr.Duration(300),
			MinSeqLedgerGap: xdr.Uint32(10),
			ExtraSigners:    []xdr.SignerKey{},
		},
	}
}

func createSignature() xdr.DecoratedSignature {
	var hint xdr.SignatureHint
	copy(hint[:], []byte{0x01, 0x02, 0x03, 0x04})

	sig := make([]byte, 64)
	for i := range sig {
		sig[i] = byte(i)
	}

	return xdr.DecoratedSignature{
		Hint:      hint,
		Signature: xdr.Signature(sig),
	}
}
