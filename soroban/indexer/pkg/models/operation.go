package models

import (
	"gorm.io/gorm"
)

// Operation represents a Stellar operation within a transaction
type Operation struct {
	// Composite primary key: tx_hash + operation_index
	ID               string `gorm:"primaryKey;type:varchar(128)" json:"id"`
	TxHash           string `gorm:"index;type:varchar(64);not null" json:"tx_hash"`
	OperationIndex   int32  `gorm:"not null" json:"operation_index"`
	SourceAccount    string `gorm:"type:varchar(56)" json:"source_account"`
	OperationType    string `gorm:"type:varchar(64);not null" json:"operation_type"`
	OperationDetails []byte `gorm:"type:jsonb" json:"operation_details"` // JSONB for type-specific data

	// Timestamps
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Operation) TableName() string {
	return "operations"
}

// UpsertOperation inserts or updates an operation
func UpsertOperation(db *gorm.DB, operation *Operation) error {
	return db.Save(operation).Error
}

// GetOperationsByTxHash retrieves all operations for a transaction
func GetOperationsByTxHash(db *gorm.DB, txHash string) ([]Operation, error) {
	var operations []Operation
	err := db.Where("tx_hash = ?", txHash).Order("operation_index ASC").Find(&operations).Error
	return operations, err
}

// GetOperationByID retrieves a specific operation by ID
func GetOperationByID(db *gorm.DB, id string) (*Operation, error) {
	var operation Operation
	err := db.Where("id = ?", id).First(&operation).Error
	if err != nil {
		return nil, err
	}
	return &operation, nil
}
