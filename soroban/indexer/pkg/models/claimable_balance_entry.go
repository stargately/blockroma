package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClaimableBalanceEntry struct {
	BalanceID             string      `gorm:"column:balance_id;type:varchar(64);primaryKey;not null"`
	Claimants             interface{} `gorm:"column:claimants;type:jsonb"`
	AssetType             int32       `gorm:"column:asset_type;type:int"`
	AssetCode             []byte      `gorm:"column:asset_code;type:bytea"`
	AssetIssuer           string      `gorm:"column:asset_issuer;type:varchar(64);index;not null"`
	Amount                int64       `gorm:"column:amount;type:bigint;not null"`
	Ext                   interface{} `gorm:"column:ext;type:jsonb"`
	SponsoringID          string      `gorm:"column:sponsoring_id;type:varchar(64);index"`
	LastModifiedLedgerSeq uint32      `gorm:"column:last_modified_ledger_seq;type:int;not null"`
	CreatedAt             time.Time   `gorm:"column:created_at"`
	UpdatedAt             time.Time   `gorm:"column:updated_at"`
}

func (ClaimableBalanceEntry) TableName() string {
	return "claimable_balance_entries"
}

func UpsertClaimableBalanceEntry(db *gorm.DB, entry *ClaimableBalanceEntry) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "balance_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"claimants", "asset_type", "asset_code", "asset_issuer", "amount",
			"ext", "sponsoring_id", "last_modified_ledger_seq", "updated_at",
		}),
	}).Create(entry).Error
}
