package models

import (
	"encoding/json"
	"time"

	"github.com/blockroma/soroban-indexer/pkg/models/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction struct {
	ID               string              `gorm:"column:id;primaryKey"`
	Status           string              `gorm:"column:status"`
	Ledger           *uint32             `gorm:"column:ledger"`
	LedgerCreatedAt  *int64              `gorm:"column:ledger_created_at"`
	ApplicationOrder *int32              `gorm:"column:application_order"`
	FeeBump          *bool               `gorm:"column:fee_bump"`
	FeeBumpInfo      *util.FeeBumpInfo   `gorm:"column:fee_bump_info;type:jsonb"`
	Fee              *int32              `gorm:"column:fee"`
	FeeCharged       *int32              `gorm:"column:fee_charged"`
	Sequence         *int64              `gorm:"column:sequence"`
	SourceAccount    *string             `gorm:"column:source_account"`
	MuxedAccountId   *int64              `gorm:"column:muxed_account_id"`
	Memo             *util.TypeItem      `gorm:"column:memo;type:jsonb"`
	Preconditions    *util.Preconditions `gorm:"column:preconditions;type:jsonb"`
	Signatures       *util.Signatures    `gorm:"column:signatures;type:jsonb"`
	CreatedAt        time.Time           `gorm:"column:created_at"`
	UpdatedAt        time.Time           `gorm:"column:updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func UpsertTransaction(db *gorm.DB, tx *Transaction) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"status", "ledger", "ledger_created_at", "application_order", "fee_bump",
			"fee_bump_info", "fee", "fee_charged", "sequence", "source_account",
			"muxed_account_id", "memo", "preconditions", "signatures", "updated_at",
		}),
	}).Create(tx).Error
}

func NewTransaction(data []byte) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal(data, &tx)
	return &tx, err
}
