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
			result := scValToInterface(tt.scVal)
			if !tt.checkFn(result) {
				t.Errorf("scValToInterface() check failed for %v (%T)", result, result)
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
