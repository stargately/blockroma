package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8000")
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.endpoint != "http://localhost:8000" {
		t.Errorf("endpoint = %v, want http://localhost:8000", client.endpoint)
	}
	if client.httpClient == nil {
		t.Error("httpClient is nil")
	}
}

func TestClient_GetLatestLedger(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify request body
		var req jsonRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Method != "getLatestLedger" {
			t.Errorf("Method = %v, want getLatestLedger", req.Method)
		}

		// Send response
		resp := jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"sequence": 12345}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ledger, err := client.GetLatestLedger(context.Background())
	if err != nil {
		t.Fatalf("GetLatestLedger() error = %v", err)
	}

	if ledger != 12345 {
		t.Errorf("GetLatestLedger() = %v, want 12345", ledger)
	}
}

func TestClient_GetEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Method != "getEvents" {
			t.Errorf("Method = %v, want getEvents", req.Method)
		}

		// Send mock events response
		mockResponse := GetEventsResponse{
			Events: []Event{
				{
					ID:                       "0000012345-0000000001",
					Type:                     "contract",
					Ledger:                   12345,
					LedgerClosedAt:           "2024-01-01T00:00:00Z",
					ContractID:               "CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM",
					PagingToken:              "12345-1",
					Topic:                    []string{"transfer"},
					Value:                    "1000000",
					InSuccessfulContractCall: true,
					TxHash:                   "abc123",
				},
			},
			LatestLedger: 12345,
		}

		respJSON, _ := json.Marshal(mockResponse)
		resp := jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(respJSON),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	eventsReq := GetEventsRequest{
		StartLedger: 12345,
	}

	eventsResp, err := client.GetEvents(context.Background(), eventsReq)
	if err != nil {
		t.Fatalf("GetEvents() error = %v", err)
	}

	if len(eventsResp.Events) != 1 {
		t.Errorf("Events count = %v, want 1", len(eventsResp.Events))
	}

	if eventsResp.Events[0].ID != "0000012345-0000000001" {
		t.Errorf("Event ID = %v, want 0000012345-0000000001", eventsResp.Events[0].ID)
	}

	if eventsResp.LatestLedger != 12345 {
		t.Errorf("LatestLedger = %v, want 12345", eventsResp.LatestLedger)
	}
}

func TestClient_GetTransaction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Method != "getTransaction" {
			t.Errorf("Method = %v, want getTransaction", req.Method)
		}

		// Send mock transaction response
		mockTx := Transaction{
			Hash:             "test-hash-123",
			Status:           "SUCCESS",
			Ledger:           12345,
			ApplicationOrder: 1,
			LedgerCloseTime:  time.Now().Unix(),
			EnvelopeXdr:      "envelope-xdr-data",
			ResultXdr:        "result-xdr-data",
			ResultMetaXdr:    "meta-xdr-data",
		}

		respJSON, _ := json.Marshal(mockTx)
		resp := jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(respJSON),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	tx, err := client.GetTransaction(context.Background(), "test-hash-123")
	if err != nil {
		t.Fatalf("GetTransaction() error = %v", err)
	}

	if tx.Hash != "test-hash-123" {
		t.Errorf("Hash = %v, want test-hash-123", tx.Hash)
	}
	if tx.Status != "SUCCESS" {
		t.Errorf("Status = %v, want SUCCESS", tx.Status)
	}
}

func TestClient_GetNetwork(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockNetwork := NetworkInfo{
			FriendbotURL:    "https://friendbot.stellar.org",
			Passphrase:      "Public Global Stellar Network ; September 2015",
			ProtocolVersion: 20,
		}

		respJSON, _ := json.Marshal(mockNetwork)
		resp := jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      1,
			Result:  json.RawMessage(respJSON),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	network, err := client.GetNetwork(context.Background())
	if err != nil {
		t.Fatalf("GetNetwork() error = %v", err)
	}

	if network.Passphrase != "Public Global Stellar Network ; September 2015" {
		t.Errorf("Passphrase = %v", network.Passphrase)
	}
	if network.ProtocolVersion != 20 {
		t.Errorf("ProtocolVersion = %v, want 20", network.ProtocolVersion)
	}
}

func TestClient_Health(t *testing.T) {
	tests := []struct {
		name       string
		response   jsonRPCResponse
		wantErr    bool
		errMessage string
	}{
		{
			name: "healthy status",
			response: jsonRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Result:  json.RawMessage(`{"status": "healthy"}`),
			},
			wantErr: false,
		},
		{
			name: "unhealthy status",
			response: jsonRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Result:  json.RawMessage(`{"status": "unhealthy"}`),
			},
			wantErr:    true,
			errMessage: "unhealthy",
		},
		{
			name: "rpc error",
			response: jsonRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Error: &rpcError{
					Code:    -32603,
					Message: "data stores are not initialized",
				},
			},
			wantErr:    true,
			errMessage: "data stores are not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClient(server.URL)
			err := client.Health(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.errMessage != "" && err != nil {
				if err.Error() != tt.errMessage && !contains(err.Error(), tt.errMessage) {
					t.Errorf("Health() error message = %v, want to contain %v", err.Error(), tt.errMessage)
				}
			}
		})
	}
}

func TestClient_GetLedgerEntries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Method != "getLedgerEntries" {
			t.Errorf("Method = %v, want getLedgerEntries", req.Method)
		}

		mockResponse := GetLedgerEntriesResponse{
			Entries: []LedgerEntryResult{
				{
					Key:                "key-123",
					XDR:                "xdr-data",
					LastModifiedLedger: 12345,
					LiveUntilLedgerSeq: 12445,
				},
			},
			LatestLedger: 12345,
		}

		respJSON, _ := json.Marshal(mockResponse)
		resp := jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(respJSON),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	entries, err := client.GetLedgerEntries(context.Background(), []string{"key-123"})
	if err != nil {
		t.Fatalf("GetLedgerEntries() error = %v", err)
	}

	if len(entries.Entries) != 1 {
		t.Errorf("Entries count = %v, want 1", len(entries.Entries))
	}

	if entries.Entries[0].Key != "key-123" {
		t.Errorf("Entry key = %v, want key-123", entries.Entries[0].Key)
	}
}

func TestClient_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetLatestLedger(context.Background())
	if err == nil {
		t.Error("Expected error for HTTP 500, got nil")
	}
}

func TestClient_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetLatestLedger(context.Background())
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestClient_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	client.httpClient.Timeout = 100 * time.Millisecond

	ctx := context.Background()
	_, err := client.GetLatestLedger(ctx)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestRPCError_Error(t *testing.T) {
	err := &rpcError{
		Code:    -32603,
		Message: "internal error",
	}

	expected := "rpc error -32603: internal error"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr))
}
