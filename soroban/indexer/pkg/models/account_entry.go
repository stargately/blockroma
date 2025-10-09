package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountEntry struct {
	AccountID             string      `gorm:"column:account_id;type:varchar(64);primaryKey;not null"`
	Balance               int64       `gorm:"column:balance;type:bigint;not null"`
	SeqNum                int64       `gorm:"column:seq_num;type:bigint;not null"`
	NumSubEntries         uint32      `gorm:"column:num_sub_entries;type:int;not null"`
	Flags                 uint32      `gorm:"column:flags;type:int;not null"`
	HomeDomain            string      `gorm:"column:home_domain;type:varchar(32);index"`
	Signers               interface{} `gorm:"column:signers;type:jsonb"`
	Ext                   interface{} `gorm:"column:ext;type:jsonb;not null"`
	InflationDest         string      `gorm:"column:inflation_dest;type:varchar(64);not null"`
	Thresholds            []byte      `gorm:"column:thresholds;type:bytea"`
	SponsoringID          string      `gorm:"column:sponsoring_id;type:varchar(64);index"`
	LastModifiedLedgerSeq uint32      `gorm:"column:last_modified_ledger_seq;type:int;not null"`
	CreatedAt             time.Time   `gorm:"column:created_at"`
	UpdatedAt             time.Time   `gorm:"column:updated_at"`
}

func (AccountEntry) TableName() string {
	return "account_entries"
}

func UpsertAccountEntry(db *gorm.DB, entry *AccountEntry) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "account_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"balance", "seq_num", "num_sub_entries", "flags",
			"home_domain", "signers", "ext", "inflation_dest",
			"thresholds", "sponsoring_id", "last_modified_ledger_seq", "updated_at",
		}),
	}).Create(entry).Error
}
