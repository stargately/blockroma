package util

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Signature represents a transaction signature
type Signature struct {
	Hint      string `json:"hint"`
	Signature string `json:"signature"`
}

// Signatures is a slice of Signature with database serialization
type Signatures []Signature

// Value implements driver.Valuer for Signatures
func (s Signatures) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements sql.Scanner for Signatures
func (s *Signatures) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// TypeItem represents a typed value (used for Memo)
type TypeItem struct {
	Type      string `json:"type"`
	ItemValue string `json:"value"`
}

// Value implements driver.Valuer for TypeItem
func (t TypeItem) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan implements sql.Scanner for TypeItem
func (t *TypeItem) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, t)
}

// Bonds represents min/max bounds for time or ledger
type Bonds struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

// SignerKey represents a signer key with different types
type SignerKey struct {
	Type                 string `json:"type"`
	Ed25519              string `json:"ed25519,omitempty"`
	PreAuthTx            string `json:"pre_auth_tx,omitempty"`
	HashX                string `json:"hash_x,omitempty"`
	Ed25519SignedPayload string `json:"ed25519_signed_payload,omitempty"`
}

// Preconditions represents transaction preconditions
type Preconditions struct {
	TimeBounds      *Bonds       `json:"time_bounds,omitempty"`
	LedgerBounds    *Bonds       `json:"ledger_bounds,omitempty"`
	MinSeqNum       *int64       `json:"min_seq_num,omitempty"`
	MinSeqAge       *int64       `json:"min_seq_age,omitempty"`
	MinSeqLedgerGap *int32       `json:"min_seq_ledger_gap,omitempty"`
	ExtraSigners    *[]SignerKey `json:"extra_signers,omitempty"`
}

// Value implements driver.Valuer for Preconditions
func (p Preconditions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements sql.Scanner for Preconditions
func (p *Preconditions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, p)
}

// FeeBumpInfo represents fee bump transaction information
type FeeBumpInfo struct {
	Fee            int32   `json:"fee"`
	SourceAccount  *string `json:"source_account,omitempty"`
	MuxedAccountId *int64  `json:"muxed_account_id,omitempty"`
}

// Value implements driver.Valuer for FeeBumpInfo
func (f FeeBumpInfo) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan implements sql.Scanner for FeeBumpInfo
func (f *FeeBumpInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, f)
}
