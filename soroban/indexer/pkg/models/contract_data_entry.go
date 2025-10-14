package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// JSONB is a custom type for PostgreSQL JSONB columns that properly handles JSON serialization
type JSONB []byte

// Value implements the driver.Valuer interface for database serialization
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	// Return as string for PostgreSQL JSONB
	return string(j), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = v
		return nil
	case string:
		*j = []byte(v)
		return nil
	default:
		return errors.New("failed to scan JSONB value: unexpected type")
	}
}

// GormDataType implements gorm.io/gorm/schema.GormDataTypeInterface to ensure GORM uses Value() method
func (JSONB) GormDataType() string {
	return "jsonb"
}

// GormValue implements gorm.io/gorm/schema.GormValuerInterface to ensure GORM uses Value() method
func (j JSONB) GormValue(ctx interface{}, db *gorm.DB) interface{} {
	if len(j) == 0 {
		return nil
	}
	return string(j)
}

type ContractDataEntry struct {
	KeyHash             string    `gorm:"column:key_hash;type:text;primaryKey;not null"`
	ContractID          string    `gorm:"column:contract_id;type:varchar(64)"`
	KeyXdr              string    `gorm:"column:key_xdr;type:text"`
	ExpirationLedgerSeq uint32    `gorm:"column:expiration_ledger_seq;type:int"`
	Key                 JSONB     `gorm:"column:key;type:jsonb"`
	Durability          string    `gorm:"column:durability;type:varchar(64)"`
	Flags               uint32    `gorm:"column:flags;type:int"`
	ValXdr              string    `gorm:"column:val_xdr;type:text"`
	Val                 JSONB     `gorm:"column:val;type:jsonb"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
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
