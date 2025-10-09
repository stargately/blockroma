package models

import (
	"encoding/json"
	"time"

	"github.com/blockroma/soroban-indexer/pkg/models/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenOperation struct {
	ID               string       `gorm:"column:id;primaryKey"`
	Type             string       `gorm:"column:type"`
	TxIndex          int32        `gorm:"column:tx_index"`
	Ledger           int32        `gorm:"column:ledger"`
	LedgerClosedAt   string       `gorm:"column:ledger_closed_at"`
	ContractID       string       `gorm:"column:contract_id"`
	From             string       `gorm:"column:from"`
	To               *string      `gorm:"column:to"`
	Amount           *util.Int128 `gorm:"column:amount"`
	Authorized       *bool        `gorm:"column:authorized"`
	ExpirationLedger *int32       `gorm:"column:expiration_ledger"`
	CreatedAt        time.Time    `gorm:"column:created_at"`
	UpdatedAt        time.Time    `gorm:"column:updated_at"`
}

func (TokenOperation) TableName() string {
	return "token_operations"
}

func UpsertTokenOperation(db *gorm.DB, tokenOp *TokenOperation) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"type", "tx_index", "ledger", "ledger_closed_at", "contract_id",
			"from", "to", "amount", "authorized", "expiration_ledger", "updated_at",
		}),
	}).Create(tokenOp).Error
}

func NewTokenOperation(data []byte) (*TokenOperation, error) {
	var to TokenOperation
	err := json.Unmarshal(data, &to)
	return &to, err
}
