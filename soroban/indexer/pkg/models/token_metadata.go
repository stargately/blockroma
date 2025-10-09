package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenMetadata struct {
	ContractID   string    `gorm:"column:contract_id;primaryKey"`
	AdminAddress string    `gorm:"column:admin_address"`
	Decimal      uint32    `gorm:"column:decimal"`
	Name         string    `gorm:"column:name"`
	Symbol       string    `gorm:"column:symbol"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (TokenMetadata) TableName() string {
	return "token_metadata"
}

func UpsertTokenMetadata(db *gorm.DB, metadata *TokenMetadata) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "contract_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"admin_address", "decimal", "name", "symbol", "updated_at"}),
	}).Create(metadata).Error
}

func NewTokenMetadata(data []byte) (*TokenMetadata, error) {
	var tm TokenMetadata
	err := json.Unmarshal(data, &tm)
	return &tm, err
}
