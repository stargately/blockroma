package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DataEntry struct {
	AccountID             string      `gorm:"column:account_id;type:varchar(64);primaryKey;not null"`
	DataName              string      `gorm:"column:data_name;type:varchar(64);primaryKey;not null"`
	DataValue             interface{} `gorm:"column:data_value;type:jsonb;not null"`
	Ext                   interface{} `gorm:"column:ext;type:jsonb;not null"`
	SponsoringID          string      `gorm:"column:sponsoring_id;type:varchar(64);index"`
	LastModifiedLedgerSeq uint32      `gorm:"column:last_modified_ledger_seq;type:int;not null"`
	CreatedAt             time.Time   `gorm:"column:created_at"`
	UpdatedAt             time.Time   `gorm:"column:updated_at"`
}

func (DataEntry) TableName() string {
	return "data_entries"
}

func UpsertDataEntry(db *gorm.DB, entry *DataEntry) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "account_id"}, {Name: "data_name"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"data_value", "ext", "sponsoring_id", "last_modified_ledger_seq", "updated_at",
		}),
	}).Create(entry).Error
}
