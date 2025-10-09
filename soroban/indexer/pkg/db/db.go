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

	// Auto-migrate tables
	if err := db.AutoMigrate(
		&models.Event{},
		&models.Transaction{},
		&models.Cursor{},
		&models.TokenMetadata{},
		&models.TokenOperation{},
		&models.TokenBalance{},
		&models.ContractDataEntry{},
		&models.AccountEntry{},
		&models.TrustLineEntry{},
		&models.OfferEntry{},
		&models.LiquidityPoolEntry{},
		&models.ClaimableBalanceEntry{},
		&models.DataEntry{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
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
