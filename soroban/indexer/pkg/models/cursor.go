package models

import (
	"time"

	"gorm.io/gorm"
)

type Cursor struct {
	ID         int       `gorm:"column:id;primaryKey"`
	LastLedger uint32    `gorm:"column:last_ledger;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Cursor) TableName() string {
	return "indexer_cursor"
}

func GetCursor(db *gorm.DB) (uint32, error) {
	var cursor Cursor
	err := db.First(&cursor, 1).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return cursor.LastLedger, err
}

func UpdateCursor(db *gorm.DB, ledger uint32) error {
	return db.Save(&Cursor{
		ID:         1,
		LastLedger: ledger,
		UpdatedAt:  time.Now(),
	}).Error
}
