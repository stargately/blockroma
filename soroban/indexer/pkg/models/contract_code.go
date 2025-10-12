package models

import (
	"time"

	"gorm.io/gorm"
)

// ContractCode represents deployed contract WASM code
type ContractCode struct {
	Hash        string    `gorm:"column:hash;primaryKey"`      // SHA256 hash of WASM bytecode
	Wasm        []byte    `gorm:"column:wasm;type:bytea"`      // WASM bytecode
	DeployedAt  time.Time `gorm:"column:deployed_at"`          // First time this code was seen
	Ledger      uint32    `gorm:"column:ledger"`               // Ledger where first deployed
	TxHash      string    `gorm:"column:tx_hash"`              // Transaction hash that deployed it
	SizeBytes   int       `gorm:"column:size_bytes"`           // Size of WASM in bytes
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

// TableName returns the table name for ContractCode
func (ContractCode) TableName() string {
	return "contract_code"
}

// UpsertContractCode inserts or updates contract code in the database
func UpsertContractCode(db *gorm.DB, code *ContractCode) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Check if code already exists
		var existing ContractCode
		err := tx.Where("hash = ?", code.Hash).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// Insert new code
			return tx.Create(code).Error
		} else if err != nil {
			return err
		}

		// Code already exists, no update needed (immutable)
		return nil
	})
}

// GetContractCodeByHash retrieves contract code by its hash
func GetContractCodeByHash(db *gorm.DB, hash string) (*ContractCode, error) {
	var code ContractCode
	err := db.Where("hash = ?", hash).First(&code).Error
	if err != nil {
		return nil, err
	}
	return &code, nil
}

// GetAllContractCodes retrieves all contract codes
func GetAllContractCodes(db *gorm.DB, limit, offset int) ([]ContractCode, error) {
	var codes []ContractCode
	err := db.Order("deployed_at DESC").Limit(limit).Offset(offset).Find(&codes).Error
	return codes, err
}
