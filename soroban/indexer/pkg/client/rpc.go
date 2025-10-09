package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	endpoint   string
	httpClient *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// JSON-RPC request/response types
type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message)
}

// call performs a JSON-RPC call
func (c *Client) call(ctx context.Context, method string, params interface{}, result interface{}) error {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status: %d", resp.StatusCode)
	}

	var rpcResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if rpcResp.Error != nil {
		return rpcResp.Error
	}

	if result != nil {
		if err := json.Unmarshal(rpcResp.Result, result); err != nil {
			return fmt.Errorf("unmarshal result: %w", err)
		}
	}

	return nil
}

// GetLatestLedger gets the latest ledger number
func (c *Client) GetLatestLedger(ctx context.Context) (uint32, error) {
	var result struct {
		Sequence uint32 `json:"sequence"`
	}

	if err := c.call(ctx, "getLatestLedger", nil, &result); err != nil {
		return 0, err
	}

	return result.Sequence, nil
}

// Event represents a contract event from RPC
type Event struct {
	ID                       string   `json:"id"`
	Type                     string   `json:"type"`
	Ledger                   uint32   `json:"ledger"`
	LedgerClosedAt           string   `json:"ledgerClosedAt"`
	ContractID               string   `json:"contractId"`
	PagingToken              string   `json:"pagingToken"`
	Topic                    []string `json:"topic"`
	Value                    string   `json:"value"`
	InSuccessfulContractCall bool     `json:"inSuccessfulContractCall"`
	TxHash                   string   `json:"txHash"`
}

// GetEventsRequest parameters for getEvents
type GetEventsRequest struct {
	StartLedger uint32                 `json:"startLedger"`
	Filters     []EventFilter          `json:"filters,omitempty"`
	Pagination  *EventPaginationParams `json:"pagination,omitempty"`
}

type EventFilter struct {
	Type         *string  `json:"type,omitempty"`
	ContractIDs  []string `json:"contractIds,omitempty"`
	Topics       []string `json:"topics,omitempty"`
}

type EventPaginationParams struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  uint   `json:"limit,omitempty"`
}

// GetEventsResponse response from getEvents
type GetEventsResponse struct {
	Events       []Event `json:"events"`
	LatestLedger uint32  `json:"latestLedger"`
}

// GetEvents fetches events from the RPC
func (c *Client) GetEvents(ctx context.Context, req GetEventsRequest) (*GetEventsResponse, error) {
	var result GetEventsResponse

	if err := c.call(ctx, "getEvents", req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Transaction represents a transaction from RPC
type Transaction struct {
	Hash             string `json:"hash"`
	Status           string `json:"status"`
	Ledger           uint32 `json:"ledger"`
	ApplicationOrder int32  `json:"applicationOrder"`
	LedgerCloseTime  int64  `json:"ledgerCloseTime"`
	EnvelopeXdr      string `json:"envelopeXdr"`
	ResultXdr        string `json:"resultXdr"`
	ResultMetaXdr    string `json:"resultMetaXdr"`
}

// GetTransaction fetches a transaction by hash
func (c *Client) GetTransaction(ctx context.Context, hash string) (*Transaction, error) {
	params := struct {
		Hash string `json:"hash"`
	}{
		Hash: hash,
	}

	var result Transaction
	if err := c.call(ctx, "getTransaction", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTransactions fetches transactions from events (helper method)
func (c *Client) GetTransactionsFromEvents(ctx context.Context, events []Event) (map[string]*Transaction, error) {
	txs := make(map[string]*Transaction)
	seen := make(map[string]bool)

	for _, event := range events {
		if event.TxHash == "" || seen[event.TxHash] {
			continue
		}
		seen[event.TxHash] = true

		tx, err := c.GetTransaction(ctx, event.TxHash)
		if err != nil {
			// Log error but continue
			continue
		}
		txs[event.TxHash] = tx
	}

	return txs, nil
}

// NetworkInfo contains network information
type NetworkInfo struct {
	FriendbotURL      string `json:"friendbotUrl"`
	Passphrase        string `json:"passphrase"`
	ProtocolVersion   int    `json:"protocolVersion"`
}

// GetNetwork fetches network information
func (c *Client) GetNetwork(ctx context.Context) (*NetworkInfo, error) {
	var result NetworkInfo
	if err := c.call(ctx, "getNetwork", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Health checks RPC health
func (c *Client) Health(ctx context.Context) error {
	var result struct {
		Status string `json:"status"`
	}

	if err := c.call(ctx, "getHealth", nil, &result); err != nil {
		return err
	}

	if result.Status != "healthy" {
		return fmt.Errorf("rpc status: %s", result.Status)
	}

	return nil
}

// LedgerEntryResult represents a ledger entry from RPC
type LedgerEntryResult struct {
	Key              string `json:"key"`
	XDR              string `json:"xdr"`
	LastModifiedLedger uint32 `json:"lastModifiedLedgerSeq,omitempty"`
	LiveUntilLedgerSeq uint32 `json:"liveUntilLedgerSeq,omitempty"`
}

// GetLedgerEntriesRequest parameters for getLedgerEntries
type GetLedgerEntriesRequest struct {
	Keys []string `json:"keys"`
}

// GetLedgerEntriesResponse response from getLedgerEntries
type GetLedgerEntriesResponse struct {
	Entries      []LedgerEntryResult `json:"entries"`
	LatestLedger uint32              `json:"latestLedger"`
}

// GetLedgerEntries fetches ledger entries by keys
func (c *Client) GetLedgerEntries(ctx context.Context, keys []string) (*GetLedgerEntriesResponse, error) {
	req := GetLedgerEntriesRequest{Keys: keys}
	var result GetLedgerEntriesResponse

	if err := c.call(ctx, "getLedgerEntries", req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
