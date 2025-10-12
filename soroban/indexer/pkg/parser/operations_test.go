package parser

import (
	"encoding/base64"
	"testing"

	"github.com/stellar/go/xdr"
)

func TestParseOperations(t *testing.T) {
	// Create a simple transaction with operations
	sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	destination := xdr.MustAddress("GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWHF")

	// Create a payment operation
	paymentOp := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypePayment,
			PaymentOp: &xdr.PaymentOp{
				Destination: destination.ToMuxedAccount(),
				Asset:       xdr.MustNewNativeAsset(),
				Amount:      1000000000, // 100 XLM
			},
		},
	}

	// Create transaction
	tx := xdr.Transaction{
		SourceAccount: sourceAccount.ToMuxedAccount(),
		Fee:           100,
		SeqNum:        123456,
		Operations:    []xdr.Operation{paymentOp},
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

	// Parse operations
	txHash := "test-tx-hash"
	operations, err := ParseOperations(txHash, xdrString)
	if err != nil {
		t.Fatalf("ParseOperations() error = %v", err)
	}

	if len(operations) != 1 {
		t.Fatalf("ParseOperations() returned %d operations, want 1", len(operations))
	}

	op := operations[0]
	if op.TxHash != txHash {
		t.Errorf("Operation TxHash = %v, want %v", op.TxHash, txHash)
	}

	if op.OperationIndex != 0 {
		t.Errorf("Operation Index = %v, want 0", op.OperationIndex)
	}

	if op.OperationType != "OperationTypePayment" {
		t.Errorf("Operation Type = %v, want OperationTypePayment", op.OperationType)
	}

	if op.SourceAccount == "" {
		t.Error("Operation SourceAccount should not be empty")
	}
}

func TestParseOperations_MultipleOperations(t *testing.T) {
	sourceAccount := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	destination := xdr.MustAddress("GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWHF")

	// Create multiple operations
	paymentOp := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypePayment,
			PaymentOp: &xdr.PaymentOp{
				Destination: destination.ToMuxedAccount(),
				Asset:       xdr.MustNewNativeAsset(),
				Amount:      1000000000,
			},
		},
	}

	bumpOp := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeBumpSequence,
			BumpSequenceOp: &xdr.BumpSequenceOp{
				BumpTo: xdr.SequenceNumber(200000),
			},
		},
	}

	tx := xdr.Transaction{
		SourceAccount: sourceAccount.ToMuxedAccount(),
		Fee:           200,
		SeqNum:        123456,
		Operations:    []xdr.Operation{paymentOp, bumpOp},
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

	// Parse operations
	txHash := "test-tx-hash"
	operations, err := ParseOperations(txHash, xdrString)
	if err != nil {
		t.Fatalf("ParseOperations() error = %v", err)
	}

	if len(operations) != 2 {
		t.Fatalf("ParseOperations() returned %d operations, want 2", len(operations))
	}

	// Check first operation
	if operations[0].OperationIndex != 0 {
		t.Errorf("First operation index = %v, want 0", operations[0].OperationIndex)
	}

	// Check second operation
	if operations[1].OperationIndex != 1 {
		t.Errorf("Second operation index = %v, want 1", operations[1].OperationIndex)
	}

	// Check operation IDs
	expectedID0 := "test-tx-hash-0"
	if operations[0].ID != expectedID0 {
		t.Errorf("First operation ID = %v, want %v", operations[0].ID, expectedID0)
	}

	expectedID1 := "test-tx-hash-1"
	if operations[1].ID != expectedID1 {
		t.Errorf("Second operation ID = %v, want %v", operations[1].ID, expectedID1)
	}
}

func TestParseOperations_InvalidXDR(t *testing.T) {
	txHash := "test-tx-hash"
	_, err := ParseOperations(txHash, "invalid-base64-xdr")
	if err == nil {
		t.Error("ParseOperations() expected error for invalid XDR")
	}
}

func TestComputeClaimableBalanceID(t *testing.T) {
	tests := []struct {
		name              string
		sourceAccount     string
		seqNum            int64
		opIndex           int
		networkPassphrase string
		wantErr           bool
	}{
		{
			name:              "valid testnet account",
			sourceAccount:     "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			seqNum:            123456,
			opIndex:           0,
			networkPassphrase: "Test SDF Network ; September 2015",
			wantErr:           false,
		},
		{
			name:              "valid pubnet account",
			sourceAccount:     "GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWHF",
			seqNum:            999999,
			opIndex:           5,
			networkPassphrase: "Public Global Stellar Network ; September 2015",
			wantErr:           false,
		},
		{
			name:              "invalid account",
			sourceAccount:     "INVALID_ADDRESS",
			seqNum:            123456,
			opIndex:           0,
			networkPassphrase: "Test SDF Network ; September 2015",
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balanceID, err := ComputeClaimableBalanceID(tt.sourceAccount, tt.seqNum, tt.opIndex, tt.networkPassphrase)

			if tt.wantErr {
				if err == nil {
					t.Error("ComputeClaimableBalanceID() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("ComputeClaimableBalanceID() error = %v", err)
			}

			if balanceID == "" {
				t.Error("ComputeClaimableBalanceID() returned empty balance ID")
			}

			// Balance ID should be a 64-character hex string (32 bytes)
			if len(balanceID) != 64 {
				t.Errorf("ComputeClaimableBalanceID() balance ID length = %d, want 64", len(balanceID))
			}
		})
	}
}

func TestComputeClaimableBalanceID_Deterministic(t *testing.T) {
	// Same inputs should produce same output
	sourceAccount := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"
	seqNum := int64(123456)
	opIndex := 0
	networkPassphrase := "Test SDF Network ; September 2015"

	id1, err1 := ComputeClaimableBalanceID(sourceAccount, seqNum, opIndex, networkPassphrase)
	if err1 != nil {
		t.Fatalf("First call error = %v", err1)
	}

	id2, err2 := ComputeClaimableBalanceID(sourceAccount, seqNum, opIndex, networkPassphrase)
	if err2 != nil {
		t.Fatalf("Second call error = %v", err2)
	}

	if id1 != id2 {
		t.Errorf("ComputeClaimableBalanceID() not deterministic: %v != %v", id1, id2)
	}
}

func TestParseOperationDetails_CreateAccount(t *testing.T) {
	destination := xdr.MustAddress("GAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWHF")

	body := xdr.OperationBody{
		Type: xdr.OperationTypeCreateAccount,
		CreateAccountOp: &xdr.CreateAccountOp{
			Destination:     destination,
			StartingBalance: 10000000000, // 1000 XLM
		},
	}

	details, err := parseOperationDetails(body)
	if err != nil {
		t.Fatalf("parseOperationDetails() error = %v", err)
	}

	if details["destination"] != destination.Address() {
		t.Errorf("destination = %v, want %v", details["destination"], destination.Address())
	}

	if details["starting_balance"] != int64(10000000000) {
		t.Errorf("starting_balance = %v, want 10000000000", details["starting_balance"])
	}
}

func TestParseOperationDetails_ManageData(t *testing.T) {
	dataName := "test_key"
	dataValue := []byte("test_value")
	dataValueXDR := xdr.DataValue(dataValue)

	body := xdr.OperationBody{
		Type: xdr.OperationTypeManageData,
		ManageDataOp: &xdr.ManageDataOp{
			DataName:  xdr.String64(dataName),
			DataValue: &dataValueXDR,
		},
	}

	details, err := parseOperationDetails(body)
	if err != nil {
		t.Fatalf("parseOperationDetails() error = %v", err)
	}

	if details["data_name"] != dataName {
		t.Errorf("data_name = %v, want %v", details["data_name"], dataName)
	}

	if string(details["data_value"].([]byte)) != string(dataValue) {
		t.Errorf("data_value = %v, want %v", details["data_value"], dataValue)
	}
}

func TestAssetToMap(t *testing.T) {
	tests := []struct {
		name      string
		asset     xdr.Asset
		wantType  string
		wantCode  string
		wantEmpty bool
	}{
		{
			name:      "native asset",
			asset:     xdr.MustNewNativeAsset(),
			wantType:  "AssetTypeAssetTypeNative",
			wantCode:  "XLM",
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := assetToMap(tt.asset)

			if result["type"] != tt.wantType {
				t.Errorf("assetToMap() type = %v, want %v", result["type"], tt.wantType)
			}

			if tt.wantCode != "" && result["code"] != tt.wantCode {
				t.Errorf("assetToMap() code = %v, want %v", result["code"], tt.wantCode)
			}

			if tt.wantEmpty && len(result) > 1 {
				t.Errorf("assetToMap() returned non-empty details for empty asset")
			}
		})
	}
}

func TestPriceToMap(t *testing.T) {
	price := xdr.Price{N: 3, D: 2}
	result := priceToMap(price)

	if result["n"] != int32(3) {
		t.Errorf("priceToMap() n = %v, want 3", result["n"])
	}

	if result["d"] != int32(2) {
		t.Errorf("priceToMap() d = %v, want 2", result["d"])
	}
}
