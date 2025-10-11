package poller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/blockroma/soroban-indexer/pkg/client"
	"github.com/blockroma/soroban-indexer/pkg/models"
	"github.com/blockroma/soroban-indexer/pkg/parser"
)

type Poller struct {
	rpcClient *client.Client
	db        *gorm.DB
	logger    *logrus.Logger
	batchSize uint
}

func New(rpcClient *client.Client, db *gorm.DB, logger *logrus.Logger) *Poller {
	return &Poller{
		rpcClient: rpcClient,
		db:        db,
		logger:    logger,
		batchSize: 1000, // events per request
	}
}

// Start begins the polling loop (1 second interval)
func (p *Poller) Start(ctx context.Context) error {
	p.logger.Info("Starting poller with 1-second interval")

	// Check RPC health first
	if err := p.rpcClient.Health(ctx); err != nil {
		return fmt.Errorf("rpc health check failed: %w", err)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Poller stopped")
			return nil

		case <-ticker.C:
			if err := p.poll(ctx); err != nil {
				p.logger.WithError(err).Error("Poll failed")
			}
		}
	}
}

// poll fetches and processes new data
func (p *Poller) poll(ctx context.Context) error {
	start := time.Now()

	// Get current cursor (last processed ledger)
	cursor, err := models.GetCursor(p.db)
	if err != nil {
		return fmt.Errorf("get cursor: %w", err)
	}

	// Get latest ledger from RPC
	latestLedger, err := p.rpcClient.GetLatestLedger(ctx)
	if err != nil {
		return fmt.Errorf("get latest ledger: %w", err)
	}

	// No new ledgers
	if cursor >= latestLedger {
		return nil
	}

	// Fetch events since cursor
	startLedger := cursor
	if startLedger == 0 {
		// First run - start from current ledger
		startLedger = latestLedger
	}

	req := client.GetEventsRequest{
		StartLedger: startLedger,
		Pagination: &client.EventPaginationParams{
			Limit: p.batchSize,
		},
	}

	resp, err := p.rpcClient.GetEvents(ctx, req)
	if err != nil {
		return fmt.Errorf("get events: %w", err)
	}

	if len(resp.Events) == 0 {
		// Update cursor even if no events
		if err := models.UpdateCursor(p.db, latestLedger); err != nil {
			return fmt.Errorf("update cursor: %w", err)
		}
		return nil
	}

	p.logger.WithFields(logrus.Fields{
		"events":       len(resp.Events),
		"startLedger":  startLedger,
		"latestLedger": latestLedger,
		"duration":     time.Since(start),
	}).Info("Processing batch")

	// Process events and transactions in a single transaction
	return p.db.Transaction(func(tx *gorm.DB) error {
		// Process events
		eventCount := 0
		tokenOpCount := 0
		txHashes := make(map[string]bool)
		contractIDs := make(map[string]bool)

		for _, event := range resp.Events {
			dbEvent, err := parser.ParseEvent(event)
			if err != nil {
				p.logger.WithError(err).WithField("eventID", event.ID).Warn("Failed to parse event")
				continue
			}

			if err := models.UpsertEvent(tx, dbEvent); err != nil {
				return fmt.Errorf("upsert event: %w", err)
			}

			eventCount++

			// Track contract IDs for later processing
			if event.ContractID != "" {
				contractIDs[event.ContractID] = true
			}

			// Try to parse token operation from this event
			ledgerClosedAt, _ := time.Parse(time.RFC3339, dbEvent.LedgerClosedAt)
			topics := []interface{}{}
			json.Unmarshal([]byte(dbEvent.Topic.(string)), &topics)
			var value interface{}
			json.Unmarshal([]byte(dbEvent.Value.(string)), &value)

			if tokenOp := parser.ParseTokenOperation(
				event.ID,
				event.ContractID,
				event.Ledger,
				ledgerClosedAt,
				dbEvent.TxIndex,
				topics,
				value,
			); tokenOp != nil {
				if err := models.UpsertTokenOperation(tx, tokenOp); err != nil {
					p.logger.WithError(err).WithField("eventID", event.ID).Warn("Failed to upsert token operation")
				} else {
					tokenOpCount++
				}
			}

			// Track unique tx hashes for transaction fetching
			if event.TxHash != "" {
				txHashes[event.TxHash] = true
			}
		}

		// Fetch and process transactions
		txCount := 0
		for txHash := range txHashes {
			rpcTx, err := p.rpcClient.GetTransaction(ctx, txHash)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to fetch transaction")
				continue
			}

			// Validate that RPC returned a hash
			if rpcTx.Hash == "" {
				p.logger.WithField("requestedHash", txHash).Warn("RPC returned transaction with empty hash, using requested hash")
			} else if rpcTx.Hash != txHash {
				p.logger.WithFields(logrus.Fields{
					"requestedHash": txHash,
					"receivedHash":  rpcTx.Hash,
				}).Warn("RPC returned different hash than requested")
			}

			// Use the hash we requested (from the event) as the transaction ID
			// This ensures we don't have duplicate IDs if RPC returns wrong/empty hashes
			dbTx, err := parser.ParseTransactionWithHash(*rpcTx, txHash)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to parse transaction")
				continue
			}

			if err := models.UpsertTransaction(tx, dbTx); err != nil {
				return fmt.Errorf("upsert transaction: %w", err)
			}

			txCount++
		}

		// Process contract data for discovered contracts
		// This extracts token metadata and balances from contract storage
		if err := p.processContractData(ctx, tx, contractIDs); err != nil {
			p.logger.WithError(err).Warn("Failed to process contract data")
			// Don't fail the whole batch if contract data processing fails
		}

		// Update cursor to latest ledger
		if err := models.UpdateCursor(tx, latestLedger); err != nil {
			return fmt.Errorf("update cursor: %w", err)
		}

		p.logger.WithFields(logrus.Fields{
			"events":        eventCount,
			"transactions":  txCount,
			"tokenOps":      tokenOpCount,
			"contracts":     len(contractIDs),
			"ledger":        latestLedger,
			"duration":      time.Since(start),
		}).Info("Batch processed successfully")

		return nil
	})
}

// processContractData processes contract data entries to extract token metadata and balances
// This is called after processing events to update contract state
func (p *Poller) processContractData(ctx context.Context, tx *gorm.DB, contractIDs map[string]bool) error {
	if len(contractIDs) == 0 {
		return nil
	}

	metadataCount := 0
	balanceCount := 0

	// For each contract, process its stored data
	for contractID := range contractIDs {
		// Query contract data entries for this contract
		var entries []models.ContractDataEntry
		if err := p.db.Where("contract_id = ?", contractID).Find(&entries).Error; err != nil {
			p.logger.WithError(err).WithField("contractID", contractID).Warn("Failed to fetch contract data")
			continue
		}

		for _, entry := range entries {
			// Try to parse token metadata
			if metadata := parser.ParseTokenMetadata(contractID, entry.Key, entry.Val); metadata != nil {
				if err := models.UpsertTokenMetadata(tx, metadata); err != nil {
					p.logger.WithError(err).WithField("contractID", contractID).Warn("Failed to upsert token metadata")
				} else {
					metadataCount++
				}
			}

			// Try to parse token balance
			if balance := parser.ParseTokenBalance(contractID, entry.Key, entry.Val); balance != nil {
				if err := models.UpsertTokenBalance(tx, balance); err != nil {
					p.logger.WithError(err).WithField("contractID", contractID).Warn("Failed to upsert token balance")
				} else {
					balanceCount++
				}
			}
		}
	}

	if metadataCount > 0 || balanceCount > 0 {
		p.logger.WithFields(logrus.Fields{
			"metadata": metadataCount,
			"balances": balanceCount,
		}).Debug("Processed contract data")
	}

	return nil
}

// GetStats returns poller statistics
func (p *Poller) GetStats() (map[string]interface{}, error) {
	cursor, err := models.GetCursor(p.db)
	if err != nil {
		return nil, err
	}

	var eventCount int64
	if err := p.db.Model(&models.Event{}).Count(&eventCount).Error; err != nil {
		return nil, err
	}

	var txCount int64
	if err := p.db.Model(&models.Transaction{}).Count(&txCount).Error; err != nil {
		return nil, err
	}

	var tokenOpCount int64
	if err := p.db.Model(&models.TokenOperation{}).Count(&tokenOpCount).Error; err != nil {
		return nil, err
	}

	var contractDataCount int64
	if err := p.db.Model(&models.ContractDataEntry{}).Count(&contractDataCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"lastLedger":         cursor,
		"totalEvents":        eventCount,
		"totalTransactions":  txCount,
		"totalTokenOps":      tokenOpCount,
		"totalContractData":  contractDataCount,
	}, nil
}
