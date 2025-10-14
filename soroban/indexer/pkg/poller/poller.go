package poller

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/xdr"
	"gorm.io/gorm"

	"github.com/blockroma/soroban-indexer/pkg/client"
	"github.com/blockroma/soroban-indexer/pkg/models"
	"github.com/blockroma/soroban-indexer/pkg/parser"
)

type Poller struct {
	rpcClient      *client.Client
	db             *gorm.DB
	logger         *logrus.Logger
	batchSize      uint
	maxConcurrency int // Maximum number of concurrent RPC requests

	// Network passphrase for transaction hashing
	networkPassphrase string

	// Statistics for empty hash responses (deprecated - now computed from envelope)
	emptyHashCount   int
	lastEmptyHashLog time.Time
}

// PollerConfig defines configuration for the poller
type PollerConfig struct {
	BatchSize      uint // Events per request (default: 1000)
	MaxConcurrency int  // Max concurrent RPC requests (default: 10)
}

func New(rpcClient *client.Client, db *gorm.DB, logger *logrus.Logger) *Poller {
	return NewWithConfig(rpcClient, db, logger, PollerConfig{
		BatchSize:      1000,
		MaxConcurrency: 10,
	})
}

func NewWithConfig(rpcClient *client.Client, db *gorm.DB, logger *logrus.Logger, config PollerConfig) *Poller {
	if config.BatchSize == 0 {
		config.BatchSize = 1000
	}
	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 10
	}

	return &Poller{
		rpcClient:      rpcClient,
		db:             db,
		logger:         logger,
		batchSize:      config.BatchSize,
		maxConcurrency: config.MaxConcurrency,
	}
}

// Start begins the polling loop (1 second interval)
func (p *Poller) Start(ctx context.Context) error {
	p.logger.Info("Starting poller with 1-second interval")

	// Check RPC health first
	if err := p.rpcClient.Health(ctx); err != nil {
		return fmt.Errorf("rpc health check failed: %w", err)
	}

	// Get network passphrase for transaction hashing
	networkInfo, err := p.rpcClient.GetNetwork(ctx)
	if err != nil {
		return fmt.Errorf("get network info: %w", err)
	}
	p.networkPassphrase = networkInfo.Passphrase
	p.logger.WithField("networkPassphrase", p.networkPassphrase).Info("Network passphrase configured")

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
					return fmt.Errorf("upsert token operation: %w", err)
				}
				tokenOpCount++
			}

			// Track unique tx hashes for transaction fetching
			if event.TxHash != "" {
				txHashes[event.TxHash] = true
			}
		}

		// Fetch and process transactions
		txCount := 0
		operationCount := 0
		contractCodeCount := 0
		contractDataCount := 0
		txWithMetaCount := 0

		for txHash := range txHashes {
			rpcTx, err := p.rpcClient.GetTransaction(ctx, txHash)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to fetch transaction")
				continue
			}

			// Determine the correct transaction hash
			actualTxHash := txHash // Start with the hash from events

			// If RPC returned an empty hash, compute it from the envelope
			if rpcTx.Hash == "" {
				computedHash, err := parser.ComputeTransactionHash(rpcTx.EnvelopeXdr, p.networkPassphrase)
				if err != nil {
					p.logger.WithError(err).WithField("txHash", txHash).Debug("Failed to compute transaction hash from envelope, using event hash")
				} else {
					actualTxHash = computedHash
					// Verify it matches the event hash
					if computedHash != txHash {
						p.logger.WithFields(logrus.Fields{
							"eventHash":    txHash,
							"computedHash": computedHash,
						}).Warn("Computed hash from envelope differs from event hash")
					}
				}
			} else if rpcTx.Hash != txHash {
				// RPC returned a hash but it doesn't match the event hash
				p.logger.WithFields(logrus.Fields{
					"eventHash": txHash,
					"rpcHash":   rpcTx.Hash,
				}).Warn("RPC returned different hash than event")
				// Use the RPC hash if available
				actualTxHash = rpcTx.Hash
			} else {
				// RPC hash matches event hash - use it
				actualTxHash = rpcTx.Hash
			}

			// Parse transaction with the determined hash
			dbTx, err := parser.ParseTransactionWithHash(*rpcTx, actualTxHash)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to parse transaction")
				continue
			}

			if err := models.UpsertTransaction(tx, dbTx); err != nil {
				return fmt.Errorf("upsert transaction: %w", err)
			}

			txCount++

			// Extract contract data from transaction metadata (passive indexing)
			// This captures contract storage changes from successful transactions
			if rpcTx.ResultMetaXdr != "" {
				txWithMetaCount++
				contractDataEntries, err := parser.ExtractContractDataFromMeta(txHash, rpcTx.ResultMetaXdr)
				if err != nil {
					p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to extract contract data from meta")
				} else if len(contractDataEntries) > 0 {
					p.logger.WithFields(logrus.Fields{
						"txHash":      txHash,
						"entryCount":  len(contractDataEntries),
					}).Info("Extracted contract data from transaction metadata")

					for _, entry := range contractDataEntries {
						if err := models.UpsertContractDataEntry(tx, entry); err != nil {
							p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to upsert contract data entry from meta")
						} else {
							contractDataCount++

							// Unmarshal JSONB bytes back to interface{} for parsing
							var keyInterface interface{}
							var valInterface interface{}
							if err := json.Unmarshal(entry.Key, &keyInterface); err != nil {
								p.logger.WithError(err).WithField("txHash", txHash).Debug("Failed to unmarshal key JSON")
								continue
							}
							if err := json.Unmarshal(entry.Val, &valInterface); err != nil {
								p.logger.WithError(err).WithField("txHash", txHash).Debug("Failed to unmarshal val JSON")
								continue
							}

							// Try to parse token metadata from this contract data entry
							if metadata := parser.ParseTokenMetadata(entry.ContractID, keyInterface, valInterface); metadata != nil {
								if err := models.UpsertTokenMetadata(tx, metadata); err != nil {
									p.logger.WithError(err).WithField("contractID", entry.ContractID).Warn("Failed to upsert token metadata from meta")
								}
							}

							// Try to parse token balance from this contract data entry
							if balance := parser.ParseTokenBalance(entry.ContractID, keyInterface, valInterface); balance != nil {
								if err := models.UpsertTokenBalance(tx, balance); err != nil {
									p.logger.WithError(err).WithField("contractID", entry.ContractID).Warn("Failed to upsert token balance from meta")
								}
							}
						}
					}
				}
			}

			// Parse and collect operations from this transaction
			operations, err := parser.ParseOperations(txHash, rpcTx.EnvelopeXdr)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to parse operations")
			} else {
				// Batch upsert all operations for this transaction
				if len(operations) > 0 {
					if err := models.BatchUpsertOperations(tx, operations); err != nil {
						p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to batch upsert operations")
					} else {
						operationCount += len(operations)
					}
				}
			}

			// Extract and store contract code (WASM) from UploadContractWasm operations
			contractCodes, err := parser.ExtractContractCodeFromEnvelope(txHash, rpcTx.Ledger, rpcTx.LedgerCloseTime, rpcTx.EnvelopeXdr)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Debug("Failed to extract contract code")
			} else {
				for _, code := range contractCodes {
					if err := models.UpsertContractCode(tx, code); err != nil {
						p.logger.WithError(err).WithField("hash", code.Hash).Warn("Failed to upsert contract code")
					} else {
						contractCodeCount++
					}
				}
			}
		}

		// Process contract data for discovered contracts
		// This extracts token metadata and balances from contract storage
		if err := p.processContractData(ctx, tx, contractIDs); err != nil {
			p.logger.WithError(err).Warn("Failed to process contract data")
			// Don't fail the whole batch if contract data processing fails
		}

		// Note: Account/trustline/offer/claimable balance processing removed - Soroban RPC returns corrupted XDR
		// See SOROBAN_RPC_LIMITATIONS.md for details - these tables cannot be populated via Soroban RPC

		// Update cursor to latest ledger
		if err := models.UpdateCursor(tx, latestLedger); err != nil {
			return fmt.Errorf("update cursor: %w", err)
		}

		p.logger.WithFields(logrus.Fields{
			"events":        eventCount,
			"transactions":  txCount,
			"txWithMeta":    txWithMetaCount,
			"operations":    operationCount,
			"contractCode":  contractCodeCount,
			"contractData":  contractDataCount,
			"tokenOps":      tokenOpCount,
			"contracts":     len(contractIDs),
			"ledger":        latestLedger,
			"duration":      time.Since(start),
		}).Info("Batch processed successfully")

		return nil
	})
}

// processContractData proactively fetches contract storage data for metadata and balances
// This is called after processing events to update contract state
func (p *Poller) processContractData(ctx context.Context, tx *gorm.DB, contractIDs map[string]bool) error {
	if len(contractIDs) == 0 {
		return nil
	}

	metadataCount := 0
	balanceCount := 0
	contractDataCount := 0

	// For each contract, proactively fetch its metadata and balance data
	for contractID := range contractIDs {
		// Step 1: Fetch contract metadata (token name, symbol, decimals)
		p.logger.WithField("contractID", contractID).Info("Attempting to fetch contract metadata")
		if err := p.fetchContractMetadata(ctx, tx, contractID); err != nil {
			p.logger.WithError(err).WithField("contractID", contractID).Warn("Failed to fetch contract metadata")
		} else {
			metadataCount++
			p.logger.WithField("contractID", contractID).Info("Successfully fetched contract metadata")
		}

		// Step 2: Query existing contract data entries for this contract (passive indexing)
		var entries []models.ContractDataEntry
		if err := p.db.Where("contract_id = ?", contractID).Find(&entries).Error; err != nil {
			p.logger.WithError(err).WithField("contractID", contractID).Warn("Failed to fetch contract data")
			continue
		}

		for _, entry := range entries {
			contractDataCount++

			// Unmarshal JSONB bytes back to interface{} for parsing
			var keyInterface interface{}
			var valInterface interface{}
			if err := json.Unmarshal(entry.Key, &keyInterface); err != nil {
				p.logger.WithError(err).WithField("contractID", contractID).Debug("Failed to unmarshal key JSON")
				continue
			}
			if err := json.Unmarshal(entry.Val, &valInterface); err != nil {
				p.logger.WithError(err).WithField("contractID", contractID).Debug("Failed to unmarshal val JSON")
				continue
			}

			// Try to parse token metadata from existing data
			if metadata := parser.ParseTokenMetadata(contractID, keyInterface, valInterface); metadata != nil {
				if err := models.UpsertTokenMetadata(tx, metadata); err != nil {
					p.logger.WithError(err).WithField("contractID", contractID).Warn("Failed to upsert token metadata")
				} else {
					metadataCount++
				}
			}

			// Try to parse token balance from existing data
			if balance := parser.ParseTokenBalance(contractID, keyInterface, valInterface); balance != nil {
				if err := models.UpsertTokenBalance(tx, balance); err != nil {
					p.logger.WithError(err).WithField("contractID", contractID).Warn("Failed to upsert token balance")
				} else {
					balanceCount++
				}
			}
		}
	}

	if metadataCount > 0 || balanceCount > 0 || contractDataCount > 0 {
		p.logger.WithFields(logrus.Fields{
			"metadata":     metadataCount,
			"balances":     balanceCount,
			"contractData": contractDataCount,
		}).Debug("Processed contract data")
	}

	return nil
}

// fetchContractMetadata proactively fetches metadata for a contract using getLedgerEntries RPC
func (p *Poller) fetchContractMetadata(ctx context.Context, tx *gorm.DB, contractID string) error {
	// Build metadata key (ScvLedgerKeyContractInstance)
	metadataKeyScVal := parser.BuildMetadataKey()

	// Build the proper LedgerKey for contract data
	ledgerKey, err := parser.BuildContractDataKey(contractID, metadataKeyScVal, xdr.ContractDataDurabilityPersistent)
	if err != nil {
		return fmt.Errorf("build contract data key: %w", err)
	}

	// Fetch contract data from RPC (persistent storage)
	resp, err := p.rpcClient.GetContractData(ctx, contractID, ledgerKey, "persistent")
	if err != nil {
		return fmt.Errorf("get contract data: %w", err)
	}

	// Parse the XDR response
	data, err := base64.StdEncoding.DecodeString(resp.XDR)
	if err != nil {
		return fmt.Errorf("decode xdr: %w", err)
	}

	var ledgerEntry xdr.LedgerEntry
	if err := xdr.SafeUnmarshal(data, &ledgerEntry); err != nil {
		// Known issue: Soroban RPC sometimes returns account entries with corrupted XDR
		// when querying for contract data. Treat this as "contract data not found".
		// See SOROBAN_RPC_LIMITATIONS.md for details.
		return fmt.Errorf("contract data not found")
	}

	// Parse and store contract data entry
	if ledgerEntry.Data.Type == xdr.LedgerEntryTypeContractData {
		contractData := ledgerEntry.Data.ContractData

		// Convert key and value to interface{} first
		keyInterface := parser.ScValToInterface(contractData.Key)
		valInterface := parser.ScValToInterface(contractData.Val)

		// Marshal to JSON bytes for JSONB storage
		keyBytes, err := json.Marshal(keyInterface)
		if err != nil {
			return fmt.Errorf("marshal key to JSON: %w", err)
		}
		valBytes, err := json.Marshal(valInterface)
		if err != nil {
			return fmt.Errorf("marshal val to JSON: %w", err)
		}

		// Create ledger key hash
		key, _ := ledgerEntry.LedgerKey()
		bin, _ := key.MarshalBinary()
		keyHash := sha256.Sum256(bin)
		hexKey := hex.EncodeToString(keyHash[:])

		keyXDR, _ := xdr.MarshalBase64(contractData.Key)
		valXDR, _ := xdr.MarshalBase64(contractData.Val)

		durability := "persistent"
		if contractData.Durability == xdr.ContractDataDurabilityTemporary {
			durability = "temporary"
		}

		// Store contract data entry
		contractDataEntry := &models.ContractDataEntry{
			KeyHash:    hexKey,
			ContractID: contractID,
			Key:        models.JSONB(keyBytes),
			KeyXdr:     keyXDR,
			Val:        models.JSONB(valBytes),
			ValXdr:     valXDR,
			Durability: durability,
		}

		if err := models.UpsertContractDataEntry(tx, contractDataEntry); err != nil {
			return fmt.Errorf("upsert contract data entry: %w", err)
		}

		// Try to parse token metadata (pass interface{} values, not JSONB bytes)
		if metadata := parser.ParseTokenMetadata(contractID, keyInterface, valInterface); metadata != nil {
			if err := models.UpsertTokenMetadata(tx, metadata); err != nil {
				return fmt.Errorf("upsert token metadata: %w", err)
			}
		}
	}

	return nil
}

// processLedgerEntries fetches and processes ledger entries for discovered accounts
func (p *Poller) processLedgerEntries(ctx context.Context, tx *gorm.DB, accountAddresses map[string]bool) error {
	if len(accountAddresses) == 0 {
		return nil
	}

	p.logger.WithField("accountCount", len(accountAddresses)).Debug("Building ledger keys for accounts")

	// Build ledger keys for all accounts
	var keys []string
	for address := range accountAddresses {
		key, err := parser.BuildAccountLedgerKey(address)
		if err != nil {
			p.logger.WithError(err).WithField("address", address).Warn("Failed to build ledger key")
			continue
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		p.logger.Warn("No valid ledger keys built from account addresses")
		return nil
	}

	p.logger.WithField("keyCount", len(keys)).Debug("Fetching ledger entries from RPC")

	// Process ledger entries
	accountCount := 0
	trustlineCount := 0
	offerCount := 0
	dataCount := 0
	claimableBalanceCount := 0
	liquidityPoolCount := 0

	// Batch ledger entry requests (RPC max is 200 keys per request)
	const maxKeysPerRequest = 200
	for i := 0; i < len(keys); i += maxKeysPerRequest {
		end := i + maxKeysPerRequest
		if end > len(keys) {
			end = len(keys)
		}
		batch := keys[i:end]

		p.logger.WithFields(logrus.Fields{
			"batchStart": i,
			"batchEnd":   end,
			"batchSize":  len(batch),
		}).Debug("Fetching ledger entry batch")

		// Fetch this batch of ledger entries from RPC
		resp, err := p.rpcClient.GetLedgerEntries(ctx, batch)
		if err != nil {
			p.logger.WithError(err).WithField("batchSize", len(batch)).Warn("Failed to fetch ledger entry batch")
			continue // Skip this batch but continue with others
		}

		p.logger.WithField("entriesReceived", len(resp.Entries)).Info("Received ledger entries from RPC")

		// Process each ledger entry in this batch
		for idx, entry := range resp.Entries {
		// Parse the ledger entry
		parsedModels, err := parser.ParseLedgerEntry(entry.XDR)
		if err != nil {
			// Known issue: Soroban RPC returns corrupted XDR for account entries
			// See SOROBAN_RPC_LIMITATIONS.md for details
			// Only log at Debug level to avoid flooding logs
			p.logger.WithError(err).Debug("Failed to parse ledger entry (expected for account entries from Soroban RPC)")
			continue
		}

		// DEBUG: Log what we parsed
		p.logger.WithFields(logrus.Fields{
			"entryIndex": idx,
			"modelCount": len(parsedModels),
		}).Debug("Parsed ledger entry")

		// Upsert each model to database
		for _, model := range parsedModels {
			switch m := model.(type) {
			case *models.AccountEntry:
				if err := models.UpsertAccountEntry(tx, m); err != nil {
					p.logger.WithError(err).WithField("accountID", m.AccountID).Warn("Failed to upsert account entry")
				} else {
					accountCount++
				}
			case *models.TrustLineEntry:
				if err := models.UpsertTrustLineEntry(tx, m); err != nil {
					p.logger.WithError(err).Warn("Failed to upsert trustline entry")
				} else {
					trustlineCount++
				}
			case *models.OfferEntry:
				if err := models.UpsertOfferEntry(tx, m); err != nil {
					p.logger.WithError(err).Warn("Failed to upsert offer entry")
				} else {
					offerCount++
				}
			case *models.DataEntry:
				if err := models.UpsertDataEntry(tx, m); err != nil {
					p.logger.WithError(err).Warn("Failed to upsert data entry")
				} else {
					dataCount++
				}
			case *models.ClaimableBalanceEntry:
				if err := models.UpsertClaimableBalanceEntry(tx, m); err != nil {
					p.logger.WithError(err).Warn("Failed to upsert claimable balance entry")
				} else {
					claimableBalanceCount++
				}
			case *models.LiquidityPoolEntry:
				if err := models.UpsertLiquidityPoolEntry(tx, m); err != nil {
					p.logger.WithError(err).Warn("Failed to upsert liquidity pool entry")
				} else {
					liquidityPoolCount++
				}
			}
		}
		}
	}

	p.logger.WithFields(logrus.Fields{
		"accounts":         accountCount,
		"trustlines":       trustlineCount,
		"offers":           offerCount,
		"data":             dataCount,
		"claimableBalance": claimableBalanceCount,
		"liquidityPools":   liquidityPoolCount,
	}).Info("Processed ledger entries")

	return nil
}

// processClaimableBalances fetches and processes claimable balance ledger entries
func (p *Poller) processClaimableBalances(ctx context.Context, tx *gorm.DB, balanceIDs map[string]bool) error {
	if len(balanceIDs) == 0 {
		return nil
	}

	// Build ledger keys for all claimable balances
	var keys []string
	for balanceID := range balanceIDs {
		key, err := parser.BuildClaimableBalanceLedgerKey(balanceID)
		if err != nil {
			p.logger.WithError(err).WithField("balanceID", balanceID).Debug("Failed to build claimable balance ledger key")
			continue
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return nil
	}

	// Fetch ledger entries from RPC
	resp, err := p.rpcClient.GetLedgerEntries(ctx, keys)
	if err != nil {
		return fmt.Errorf("get claimable balance ledger entries: %w", err)
	}

	// Process each ledger entry
	claimableBalanceCount := 0

	for _, entry := range resp.Entries {
		// Parse the ledger entry
		parsedModels, err := parser.ParseLedgerEntry(entry.XDR)
		if err != nil {
			p.logger.WithError(err).Debug("Failed to parse claimable balance ledger entry")
			continue
		}

		// Upsert each model to database
		for _, model := range parsedModels {
			if m, ok := model.(*models.ClaimableBalanceEntry); ok {
				if err := models.UpsertClaimableBalanceEntry(tx, m); err != nil {
					p.logger.WithError(err).WithField("balanceID", m.BalanceID).Warn("Failed to upsert claimable balance entry")
				} else {
					claimableBalanceCount++
				}
			}
		}
	}

	if claimableBalanceCount > 0 {
		p.logger.WithFields(logrus.Fields{
			"claimableBalances": claimableBalanceCount,
		}).Debug("Processed claimable balances")
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

// BackfillConfig defines configuration for backfill mode
type BackfillConfig struct {
	StartLedger uint32 // First ledger to process
	EndLedger   uint32 // Last ledger to process (0 = current ledger)
	BatchSize   uint32 // Number of ledgers to process per batch
	RateLimit   uint32 // Max requests per second
}

// Backfill processes historical ledgers sequentially
func (p *Poller) Backfill(ctx context.Context, config BackfillConfig) error {
	start := time.Now()

	// Validate configuration
	if config.StartLedger == 0 {
		return fmt.Errorf("start ledger must be > 0")
	}
	if config.BatchSize == 0 {
		config.BatchSize = 100 // Default batch size
	}
	if config.RateLimit == 0 {
		config.RateLimit = 10 // Default 10 requests/sec
	}

	// Get current ledger if end ledger not specified
	if config.EndLedger == 0 {
		latestLedger, err := p.rpcClient.GetLatestLedger(ctx)
		if err != nil {
			return fmt.Errorf("get latest ledger: %w", err)
		}
		config.EndLedger = latestLedger
	}

	// Validate range
	if config.StartLedger > config.EndLedger {
		return fmt.Errorf("start ledger (%d) must be <= end ledger (%d)", config.StartLedger, config.EndLedger)
	}

	totalLedgers := config.EndLedger - config.StartLedger + 1
	p.logger.WithFields(logrus.Fields{
		"startLedger":  config.StartLedger,
		"endLedger":    config.EndLedger,
		"totalLedgers": totalLedgers,
		"batchSize":    config.BatchSize,
		"rateLimit":    config.RateLimit,
	}).Info("Starting backfill")

	// Rate limiter: requests per second
	rateLimiter := time.NewTicker(time.Second / time.Duration(config.RateLimit))
	defer rateLimiter.Stop()

	// Progress tracking
	processedLedgers := uint32(0)
	totalEvents := 0
	totalTransactions := 0
	totalOperations := 0
	lastProgressLog := time.Now()

	// Process ledgers in sequential order
	currentLedger := config.StartLedger

	for currentLedger <= config.EndLedger {
		select {
		case <-ctx.Done():
			p.logger.WithFields(logrus.Fields{
				"processedLedgers": processedLedgers,
				"currentLedger":    currentLedger,
				"duration":         time.Since(start),
			}).Info("Backfill cancelled")
			return ctx.Err()
		case <-rateLimiter.C:
			// Rate limit: process one batch per tick
		}

		// Determine batch range
		endBatchLedger := currentLedger
		if currentLedger+config.BatchSize-1 < config.EndLedger {
			endBatchLedger = currentLedger + config.BatchSize - 1
		} else {
			endBatchLedger = config.EndLedger
		}

		// Process this batch of ledgers
		batchStart := time.Now()
		batchEvents, batchTxs, batchOps, err := p.processLedgerBatch(ctx, currentLedger, endBatchLedger)
		if err != nil {
			p.logger.WithError(err).WithFields(logrus.Fields{
				"startLedger": currentLedger,
				"endLedger":   endBatchLedger,
			}).Error("Failed to process ledger batch")

			// Try to resume from the failed ledger
			// For now, we skip and continue (can be made configurable)
			p.logger.Warn("Continuing to next batch after error")
		}

		processedLedgers += (endBatchLedger - currentLedger + 1)
		totalEvents += batchEvents
		totalTransactions += batchTxs
		totalOperations += batchOps

		// Log progress every 10 seconds
		if time.Since(lastProgressLog) >= 10*time.Second {
			progress := float64(processedLedgers) / float64(totalLedgers) * 100
			estimated := time.Duration(float64(time.Since(start)) / float64(processedLedgers) * float64(totalLedgers-processedLedgers))

			p.logger.WithFields(logrus.Fields{
				"progress":          fmt.Sprintf("%.2f%%", progress),
				"processedLedgers":  processedLedgers,
				"totalLedgers":      totalLedgers,
				"currentLedger":     endBatchLedger,
				"events":            totalEvents,
				"transactions":      totalTransactions,
				"operations":        totalOperations,
				"duration":          time.Since(start).Round(time.Second),
				"estimatedRemaining": estimated.Round(time.Second),
				"batchDuration":     batchStart.Sub(lastProgressLog).Round(time.Millisecond),
			}).Info("Backfill progress")
			lastProgressLog = time.Now()
		}

		// Move to next batch
		currentLedger = endBatchLedger + 1
	}

	// Final summary
	duration := time.Since(start)
	rate := float64(processedLedgers) / duration.Seconds()

	p.logger.WithFields(logrus.Fields{
		"totalLedgers":     processedLedgers,
		"events":           totalEvents,
		"transactions":     totalTransactions,
		"operations":       totalOperations,
		"duration":         duration.Round(time.Second),
		"ledgersPerSecond": fmt.Sprintf("%.2f", rate),
	}).Info("Backfill completed successfully")

	return nil
}

// processLedgerBatch processes a range of ledgers and returns counts
func (p *Poller) processLedgerBatch(ctx context.Context, startLedger, endLedger uint32) (events, txs, ops int, err error) {
	// Fetch events for this ledger range
	req := client.GetEventsRequest{
		StartLedger: startLedger,
		Pagination: &client.EventPaginationParams{
			Limit: p.batchSize,
		},
	}

	// Keep fetching until we've processed all events in this range
	for {
		resp, err := p.rpcClient.GetEvents(ctx, req)
		if err != nil {
			return events, txs, ops, fmt.Errorf("get events: %w", err)
		}

		if len(resp.Events) == 0 {
			break
		}

		// Process this page of events
		pageEvents, pageTxs, pageOps, err := p.processEventBatch(ctx, resp.Events)
		if err != nil {
			return events, txs, ops, err
		}

		events += pageEvents
		txs += pageTxs
		ops += pageOps

		// Check if we've gone past the end ledger
		if len(resp.Events) > 0 && resp.Events[len(resp.Events)-1].Ledger > endLedger {
			break
		}

		// If there's a cursor, use it for pagination
		if resp.LatestLedger > 0 && resp.LatestLedger >= endLedger {
			break
		}

		// No more pages
		if len(resp.Events) < int(p.batchSize) {
			break
		}

		// Update request to fetch next page (using last event's ledger)
		if len(resp.Events) > 0 {
			lastLedger := resp.Events[len(resp.Events)-1].Ledger
			req.StartLedger = lastLedger
		}
	}

	// Update cursor to end of this batch
	if err := models.UpdateCursor(p.db, endLedger); err != nil {
		return events, txs, ops, fmt.Errorf("update cursor: %w", err)
	}

	return events, txs, ops, nil
}

// processEventBatch processes a batch of events (similar to poll() but without cursor updates)
func (p *Poller) processEventBatch(ctx context.Context, eventList []client.Event) (events, txs, ops int, err error) {
	// Process events and transactions in a single transaction
	err = p.db.Transaction(func(tx *gorm.DB) error {
		eventCount := 0
		txHashes := make(map[string]bool)
		contractIDs := make(map[string]bool)

		for _, event := range eventList {
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
					return fmt.Errorf("upsert token operation: %w", err)
				}
			}

			// Track unique tx hashes for transaction fetching
			if event.TxHash != "" {
				txHashes[event.TxHash] = true
			}
		}

		// Fetch and process transactions
		txCount := 0
		operationCount := 0
		contractCodeCount := 0
		claimableBalanceIDs := make(map[string]bool)
		accountAddresses := make(map[string]bool)

		for txHash := range txHashes {
			rpcTx, err := p.rpcClient.GetTransaction(ctx, txHash)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to fetch transaction")
				continue
			}

			// Use the hash we requested (from the event) as the transaction ID
			dbTx, err := parser.ParseTransactionWithHash(*rpcTx, txHash)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to parse transaction")
				continue
			}

			if err := models.UpsertTransaction(tx, dbTx); err != nil {
				return fmt.Errorf("upsert transaction: %w", err)
			}

			txCount++

			// Extract source account for ledger entry processing
			if dbTx.SourceAccount != nil && *dbTx.SourceAccount != "" {
				accountAddresses[*dbTx.SourceAccount] = true
			}

			// Parse and collect operations from this transaction
			operations, err := parser.ParseOperations(txHash, rpcTx.EnvelopeXdr)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to parse operations")
			} else {
				// Batch upsert all operations for this transaction
				if len(operations) > 0 {
					if err := models.BatchUpsertOperations(tx, operations); err != nil {
						p.logger.WithError(err).WithField("txHash", txHash).Warn("Failed to batch upsert operations")
					} else {
						operationCount += len(operations)
					}
				}
			}

			// Extract and store contract code (WASM) from UploadContractWasm operations
			contractCodes, err := parser.ExtractContractCodeFromEnvelope(txHash, rpcTx.Ledger, rpcTx.LedgerCloseTime, rpcTx.EnvelopeXdr)
			if err != nil {
				p.logger.WithError(err).WithField("txHash", txHash).Debug("Failed to extract contract code")
			} else {
				for _, code := range contractCodes {
					if err := models.UpsertContractCode(tx, code); err != nil {
						p.logger.WithError(err).WithField("hash", code.Hash).Warn("Failed to upsert contract code")
					} else {
						contractCodeCount++
					}
				}
			}

			// Extract claimable balance IDs from transaction operations
			if balanceIDs, err := parser.ExtractClaimableBalanceIDs(rpcTx.EnvelopeXdr); err == nil {
				for _, balanceID := range balanceIDs {
					claimableBalanceIDs[balanceID] = true
				}
			}
		}

		// Process contract data for discovered contracts
		if err := p.processContractData(ctx, tx, contractIDs); err != nil {
			p.logger.WithError(err).Warn("Failed to process contract data")
		}

		// Process ledger entries for discovered accounts
		// Account addresses were collected during transaction processing above
		if len(accountAddresses) > 0 {
			p.logger.WithField("accountCount", len(accountAddresses)).Info("Processing account ledger entries")
			if err := p.processLedgerEntries(ctx, tx, accountAddresses); err != nil {
				p.logger.WithError(err).Warn("Failed to process ledger entries")
			}
		} else {
			p.logger.Debug("No account addresses found in transactions")
		}

		// Process claimable balance ledger entries
		if len(claimableBalanceIDs) > 0 {
			if err := p.processClaimableBalances(ctx, tx, claimableBalanceIDs); err != nil {
				p.logger.WithError(err).Warn("Failed to process claimable balances")
			}
		}

		events = eventCount
		txs = txCount
		ops = operationCount

		return nil
	})

	return events, txs, ops, err
}
