package models

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BatchConfig defines configuration for batch operations
type BatchConfig struct {
	BatchSize int // Number of records per batch (default: 100)
}

// DefaultBatchConfig returns the default batch configuration
func DefaultBatchConfig() BatchConfig {
	return BatchConfig{
		BatchSize: 100,
	}
}

// BatchUpsertEvents upserts multiple events in batches
func BatchUpsertEvents(db *gorm.DB, events []*Event, config ...BatchConfig) error {
	if len(events) == 0 {
		return nil
	}

	batchSize := 100
	if len(config) > 0 && config[0].BatchSize > 0 {
		batchSize = config[0].BatchSize
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(events); i += batchSize {
			end := i + batchSize
			if end > len(events) {
				end = len(events)
			}

			batch := events[i:end]
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"tx_index", "type", "ledger", "ledger_closed_at", "contract_id",
					"paging_token", "topic", "value", "in_successful_contract_call",
					"last_modified_ledger_seq", "updated_at",
				}),
			}).Create(batch).Error; err != nil {
				return fmt.Errorf("batch upsert events: %w", err)
			}
		}
		return nil
	})
}

// BatchUpsertTransactions upserts multiple transactions in batches
func BatchUpsertTransactions(db *gorm.DB, transactions []*Transaction, config ...BatchConfig) error {
	if len(transactions) == 0 {
		return nil
	}

	batchSize := 100
	if len(config) > 0 && config[0].BatchSize > 0 {
		batchSize = config[0].BatchSize
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(transactions); i += batchSize {
			end := i + batchSize
			if end > len(transactions) {
				end = len(transactions)
			}

			batch := transactions[i:end]
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"status", "ledger", "ledger_created_at", "application_order", "fee_bump",
					"fee_bump_info", "fee", "fee_charged", "sequence", "source_account",
					"muxed_account_id", "memo", "preconditions", "signatures", "updated_at",
				}),
			}).Create(batch).Error; err != nil {
				return fmt.Errorf("batch upsert transactions: %w", err)
			}
		}
		return nil
	})
}

// BatchUpsertOperations upserts multiple operations in batches
func BatchUpsertOperations(db *gorm.DB, operations []*Operation, config ...BatchConfig) error {
	if len(operations) == 0 {
		return nil
	}

	batchSize := 100
	if len(config) > 0 && config[0].BatchSize > 0 {
		batchSize = config[0].BatchSize
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(operations); i += batchSize {
			end := i + batchSize
			if end > len(operations) {
				end = len(operations)
			}

			batch := operations[i:end]
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"tx_hash", "operation_index", "operation_type",
					"source_account", "operation_details", "updated_at",
				}),
			}).Create(batch).Error; err != nil {
				return fmt.Errorf("batch upsert operations: %w", err)
			}
		}
		return nil
	})
}

// BatchUpsertTokenOperations upserts multiple token operations in batches
func BatchUpsertTokenOperations(db *gorm.DB, tokenOps []*TokenOperation, config ...BatchConfig) error {
	if len(tokenOps) == 0 {
		return nil
	}

	batchSize := 100
	if len(config) > 0 && config[0].BatchSize > 0 {
		batchSize = config[0].BatchSize
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(tokenOps); i += batchSize {
			end := i + batchSize
			if end > len(tokenOps) {
				end = len(tokenOps)
			}

			batch := tokenOps[i:end]
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"type", "tx_index", "ledger", "ledger_closed_at", "contract_id",
					"from", "to", "amount", "authorized", "updated_at",
				}),
			}).Create(batch).Error; err != nil {
				return fmt.Errorf("batch upsert token operations: %w", err)
			}
		}
		return nil
	})
}

// BatchUpsertContractCode upserts multiple contract code entries in batches
func BatchUpsertContractCode(db *gorm.DB, codes []*ContractCode, config ...BatchConfig) error {
	if len(codes) == 0 {
		return nil
	}

	batchSize := 100
	if len(config) > 0 && config[0].BatchSize > 0 {
		batchSize = config[0].BatchSize
	}

	// Contract code uses special logic - only insert if not exists
	return db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(codes); i += batchSize {
			end := i + batchSize
			if end > len(codes) {
				end = len(codes)
			}

			batch := codes[i:end]
			// For contract code, we only insert, never update (immutable)
			for _, code := range batch {
				var existing ContractCode
				err := tx.Where("hash = ?", code.Hash).First(&existing).Error
				if err == gorm.ErrRecordNotFound {
					if err := tx.Create(code).Error; err != nil {
						return fmt.Errorf("batch insert contract code: %w", err)
					}
				} else if err != nil {
					return fmt.Errorf("check contract code: %w", err)
				}
				// Code already exists, skip
			}
		}
		return nil
	})
}

// BatchUpsertAccountEntries upserts multiple account entries in batches
func BatchUpsertAccountEntries(db *gorm.DB, accounts []*AccountEntry, config ...BatchConfig) error {
	if len(accounts) == 0 {
		return nil
	}

	batchSize := 100
	if len(config) > 0 && config[0].BatchSize > 0 {
		batchSize = config[0].BatchSize
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(accounts); i += batchSize {
			end := i + batchSize
			if end > len(accounts) {
				end = len(accounts)
			}

			batch := accounts[i:end]
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "account_id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"balance", "seq_num", "num_sub_entries", "flags",
					"home_domain", "signers", "ext", "inflation_dest", "thresholds", "updated_at",
				}),
			}).Create(batch).Error; err != nil {
				return fmt.Errorf("batch upsert account entries: %w", err)
			}
		}
		return nil
	})
}

// BatchResult contains results from batch operations
type BatchResult struct {
	TotalProcessed int
	SuccessCount   int
	ErrorCount     int
	Errors         []error
}

// BatchUpsertWithResults performs batch upsert and returns detailed results
func BatchUpsertWithResults(db *gorm.DB, items interface{}, upsertFunc func(*gorm.DB, interface{}) error) (*BatchResult, error) {
	result := &BatchResult{
		Errors: make([]error, 0),
	}

	err := upsertFunc(db, items)
	if err != nil {
		result.ErrorCount++
		result.Errors = append(result.Errors, err)
		return result, err
	}

	result.SuccessCount++
	return result, nil
}
