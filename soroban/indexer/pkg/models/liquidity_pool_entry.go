package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LiquidityPoolEntry struct {
	LiquidityPoolID          []byte    `gorm:"column:liquidity_pool_id;type:bytea;primaryKey"`
	Type                     int32     `gorm:"column:type;type:int"`
	AssetAType               int32     `gorm:"column:asset_a_type;type:int"`
	AssetACode               []byte    `gorm:"column:asset_a_code;type:bytea"`
	AssetAIssuer             string    `gorm:"column:asset_a_issuer;type:varchar(64);index;not null"`
	AssetBType               int32     `gorm:"column:asset_b_type;type:int"`
	AssetBCode               []byte    `gorm:"column:asset_b_code;type:bytea"`
	AssetBIssuer             string    `gorm:"column:asset_b_issuer;type:varchar(64);index;not null"`
	Fee                      int32     `gorm:"column:fee;type:int"`
	ReserveA                 int64     `gorm:"column:reserve_a;type:bigint"`
	ReserveB                 int64     `gorm:"column:reserve_b;type:bigint"`
	TotalPoolShares          int64     `gorm:"column:total_pool_shares;type:bigint"`
	PoolSharesTrustLineCount int64     `gorm:"column:pool_shares_trust_line_count;type:bigint"`
	LastModifiedLedgerSeq    uint32    `gorm:"column:last_modified_ledger_seq;type:int;not null"`
	SponsoringID             string    `gorm:"column:sponsoring_id;type:varchar(64);index"`
	CreatedAt                time.Time `gorm:"column:created_at"`
	UpdatedAt                time.Time `gorm:"column:updated_at"`
}

func (LiquidityPoolEntry) TableName() string {
	return "liquidity_pool_entries"
}

func UpsertLiquidityPoolEntry(db *gorm.DB, entry *LiquidityPoolEntry) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "liquidity_pool_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"type", "asset_a_type", "asset_a_code", "asset_a_issuer",
			"asset_b_type", "asset_b_code", "asset_b_issuer", "fee",
			"reserve_a", "reserve_b", "total_pool_shares",
			"pool_shares_trust_line_count", "last_modified_ledger_seq",
			"sponsoring_id", "updated_at",
		}),
	}).Create(entry).Error
}
