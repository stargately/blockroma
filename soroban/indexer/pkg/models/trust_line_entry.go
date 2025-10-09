package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TrustLineEntry struct {
	AccountID             string      `gorm:"column:account_id;type:varchar(64);primaryKey;not null"`
	Balance               int64       `gorm:"column:balance;type:bigint;not null"`
	Limit                 int64       `gorm:"column:limit;type:bigint;not null"`
	AssetType             int32       `gorm:"column:asset_type;type:int;primaryKey"`
	AssetCode             []byte      `gorm:"column:asset_code;type:bytea;primaryKey"`
	AssetIssuer           string      `gorm:"column:asset_issuer;type:varchar(64);primaryKey"`
	LiquidityPoolID       []byte      `gorm:"column:liquidity_pool_id;type:bytea;index"`
	Flags                 uint32      `gorm:"column:flags;type:int;not null"`
	Ext                   interface{} `gorm:"column:ext;type:jsonb;not null"`
	SponsoringID          string      `gorm:"column:sponsoring_id;type:varchar(64);index"`
	LastModifiedLedgerSeq uint32      `gorm:"column:last_modified_ledger_seq;type:int;not null"`
	CreatedAt             time.Time   `gorm:"column:created_at"`
	UpdatedAt             time.Time   `gorm:"column:updated_at"`
}

func (TrustLineEntry) TableName() string {
	return "trust_line_entries"
}

func UpsertTrustLineEntry(db *gorm.DB, entry *TrustLineEntry) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "account_id"},
			{Name: "asset_type"},
			{Name: "asset_code"},
			{Name: "asset_issuer"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"balance", "limit", "liquidity_pool_id", "flags", "ext",
			"sponsoring_id", "last_modified_ledger_seq", "updated_at",
		}),
	}).Create(entry).Error
}
