package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContractDataEntry struct {
	KeyHash             string      `gorm:"column:key_hash;type:text;primaryKey;not null"`
	ContractID          string      `gorm:"column:contract_id;type:varchar(64)"`
	KeyXdr              string      `gorm:"column:key_xdr;type:text"`
	ExpirationLedgerSeq uint32      `gorm:"column:expiration_ledger_seq;type:int"`
	Key                 interface{} `gorm:"column:key;type:jsonb"`
	Durability          string      `gorm:"column:durability;type:varchar(64)"`
	Flags               uint32      `gorm:"column:flags;type:int"`
	ValXdr              string      `gorm:"column:val_xdr;type:text"`
	Val                 interface{} `gorm:"column:val;type:jsonb"`
	CreatedAt           time.Time   `gorm:"column:created_at"`
	UpdatedAt           time.Time   `gorm:"column:updated_at"`
}

func (ContractDataEntry) TableName() string {
	return "contract_data_entries"
}

func UpsertContractDataEntry(db *gorm.DB, entry *ContractDataEntry) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key_hash"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"contract_id", "key_xdr", "expiration_ledger_seq",
			"key", "durability", "flags", "val_xdr", "val", "updated_at",
		}),
	}).Create(entry).Error
}
