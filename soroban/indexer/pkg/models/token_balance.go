package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenBalance struct {
	ContractID string    `gorm:"column:contract_id;primaryKey;not null"`
	Address    string    `gorm:"column:address;primaryKey;not null"`
	Balance    string    `gorm:"column:balance"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (TokenBalance) TableName() string {
	return "token_balances"
}

func UpsertTokenBalance(db *gorm.DB, tokenBalance *TokenBalance) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "contract_id"}, {Name: "address"}},
		DoUpdates: clause.AssignmentColumns([]string{"balance", "updated_at"}),
	}).Create(tokenBalance).Error
}
