package db

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/blockroma/soroban-indexer/pkg/models"
)

type DB struct {
	*gorm.DB
}

// Connect establishes connection to PostgreSQL
func Connect(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Run custom migrations before auto-migrate
	if err := runCustomMigrations(db); err != nil {
		return nil, fmt.Errorf("run custom migrations: %w", err)
	}

	// Auto-migrate tables
	if err := db.AutoMigrate(
		&models.Event{},
		&models.Transaction{},
		&models.Operation{},
		&models.Cursor{},
		&models.TokenMetadata{},
		&models.TokenOperation{},
		&models.TokenBalance{},
		&models.ContractDataEntry{},
		&models.ContractCode{},
		&models.AccountEntry{},
		&models.TrustLineEntry{},
		&models.OfferEntry{},
		&models.LiquidityPoolEntry{},
		&models.ClaimableBalanceEntry{},
		&models.DataEntry{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	// Create indexes for better query performance
	if err := createIndexes(db); err != nil {
		return nil, fmt.Errorf("create indexes: %w", err)
	}

	logrus.Info("Database connected and migrated")

	return &DB{DB: db}, nil
}

// Ping checks database connectivity
func (db *DB) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// WithTransaction runs function within a transaction
func (db *DB) WithTransaction(fn func(*gorm.DB) error) error {
	return db.Transaction(fn)
}

// runCustomMigrations handles custom database migrations that AutoMigrate can't handle
func runCustomMigrations(db *gorm.DB) error {
	// Check if events table exists
	if !db.Migrator().HasTable("events") {
		// Table doesn't exist yet, AutoMigrate will create it with all columns
		return nil
	}

	// Check if last_modified_ledger_seq column exists
	columnExists := db.Migrator().HasColumn(&models.Event{}, "last_modified_ledger_seq")
	if !columnExists {
		// Column doesn't exist, add it properly
		logrus.Info("Adding last_modified_ledger_seq column with default value for existing rows")

		// Add column as nullable first
		if err := db.Exec("ALTER TABLE events ADD COLUMN last_modified_ledger_seq integer").Error; err != nil {
			return fmt.Errorf("add last_modified_ledger_seq column: %w", err)
		}

		// Set default value (0) for existing NULL rows
		if err := db.Exec("UPDATE events SET last_modified_ledger_seq = 0 WHERE last_modified_ledger_seq IS NULL").Error; err != nil {
			return fmt.Errorf("set default value for last_modified_ledger_seq: %w", err)
		}

		// Now make the column NOT NULL (skip for SQLite as it doesn't support this syntax easily)
		if err := setColumnNotNull(db, "events", "last_modified_ledger_seq"); err != nil {
			return fmt.Errorf("set last_modified_ledger_seq NOT NULL: %w", err)
		}

		logrus.Info("Successfully migrated last_modified_ledger_seq column")
		return nil
	}

	// Column exists - check if there are any NULL values and fix them
	var nullCount int64
	if err := db.Raw("SELECT COUNT(*) FROM events WHERE last_modified_ledger_seq IS NULL").Scan(&nullCount).Error; err != nil {
		// Query failed, might be NOT NULL already - that's ok
		logrus.Debug("Could not check for NULL values in last_modified_ledger_seq, column may already be NOT NULL")
		return nil
	}

	if nullCount > 0 {
		// Column exists but has NULL values, fill them and add constraint
		logrus.Info("Updating existing last_modified_ledger_seq column with NULL values")
		if err := db.Exec("UPDATE events SET last_modified_ledger_seq = 0 WHERE last_modified_ledger_seq IS NULL").Error; err != nil {
			return fmt.Errorf("set default value for last_modified_ledger_seq: %w", err)
		}

		if err := setColumnNotNull(db, "events", "last_modified_ledger_seq"); err != nil {
			// This might fail if constraint already exists or on SQLite, that's ok
			logrus.Debug("Could not add NOT NULL constraint, may already exist")
		}

		logrus.Info("Successfully updated last_modified_ledger_seq column")
	}

	return nil
}

// setColumnNotNull sets a column to NOT NULL in a database-agnostic way
func setColumnNotNull(db *gorm.DB, tableName, columnName string) error {
	// Get the database dialect
	dialector := db.Dialector.Name()

	switch dialector {
	case "postgres":
		// PostgreSQL syntax
		sql := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL", tableName, columnName)
		return db.Exec(sql).Error
	case "sqlite":
		// SQLite doesn't support ALTER COLUMN easily - would need to recreate the table
		// For testing purposes, we skip this as SQLite doesn't enforce NOT NULL retroactively
		logrus.Debug("SQLite detected - skipping NOT NULL constraint (not supported)")
		return nil
	default:
		// Try PostgreSQL syntax as default
		sql := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL", tableName, columnName)
		return db.Exec(sql).Error
	}
}

// createIndexes creates database indexes for optimal query performance
func createIndexes(db *gorm.DB) error {
	dialector := db.Dialector.Name()

	// Only create specialized indexes for PostgreSQL
	if dialector != "postgres" {
		logrus.Debug("Skipping index creation for non-PostgreSQL database")
		return nil
	}

	indexes := []struct {
		name  string
		query string
	}{
		{
			name:  "idx_events_topic_gin",
			query: "CREATE INDEX IF NOT EXISTS idx_events_topic_gin ON events USING gin (topic)",
		},
		{
			name:  "idx_events_value_gin",
			query: "CREATE INDEX IF NOT EXISTS idx_events_value_gin ON events USING gin (value)",
		},
		{
			name:  "idx_events_contract_id",
			query: "CREATE INDEX IF NOT EXISTS idx_events_contract_id ON events (contract_id)",
		},
		{
			name:  "idx_events_ledger",
			query: "CREATE INDEX IF NOT EXISTS idx_events_ledger ON events (ledger)",
		},
		{
			name:  "idx_events_type",
			query: "CREATE INDEX IF NOT EXISTS idx_events_type ON events (type)",
		},
		{
			name:  "idx_transactions_ledger",
			query: "CREATE INDEX IF NOT EXISTS idx_transactions_ledger ON transactions (ledger)",
		},
		{
			name:  "idx_operations_tx_hash",
			query: "CREATE INDEX IF NOT EXISTS idx_operations_tx_hash ON operations (tx_hash)",
		},
		{
			name:  "idx_operations_type",
			query: "CREATE INDEX IF NOT EXISTS idx_operations_type ON operations (operation_type)",
		},
		{
			name:  "idx_token_operations_contract_id",
			query: "CREATE INDEX IF NOT EXISTS idx_token_operations_contract_id ON token_operations (contract_id)",
		},
		{
			name:  "idx_token_operations_from_address",
			query: "CREATE INDEX IF NOT EXISTS idx_token_operations_from_address ON token_operations (from_address)",
		},
		{
			name:  "idx_token_operations_to_address",
			query: "CREATE INDEX IF NOT EXISTS idx_token_operations_to_address ON token_operations (to_address)",
		},
	}

	for _, idx := range indexes {
		logrus.WithField("index", idx.name).Debug("Creating index")
		if err := db.Exec(idx.query).Error; err != nil {
			logrus.WithError(err).WithField("index", idx.name).Warn("Failed to create index (may already exist)")
			// Continue with other indexes even if one fails
		}
	}

	logrus.Info("Database indexes created/verified")
	return nil
}
