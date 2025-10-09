package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OfferEntry struct {
	OfferID  int64  `gorm:"column:offer_id;type:bigint;primaryKey;not null"`
	SellerID string `gorm:"column:seller_id;type:varchar(64);primaryKey;not null"`

	SellingAssetType   int32  `gorm:"column:selling_asset_type;type:int"`
	SellingAssetCode   []byte `gorm:"column:selling_asset_code;type:bytea"`
	SellingAssetIssuer string `gorm:"column:selling_asset_issuer;type:varchar(64);index;not null"`

	BuyingAssetType   int32  `gorm:"column:buying_asset_type;type:int"`
	BuyingAssetCode   []byte `gorm:"column:buying_asset_code;type:bytea"`
	BuyingAssetIssuer string `gorm:"column:buying_asset_issuer;type:varchar(64);index;not null"`

	Amount                int64       `gorm:"column:amount;type:bigint;not null"`
	Price                 string      `gorm:"column:price;type:numeric"`
	Flags                 uint32      `gorm:"column:flags;type:int;not null"`
	Ext                   interface{} `gorm:"column:ext;type:jsonb;not null"`
	SponsoringID          string      `gorm:"column:sponsoring_id;type:varchar(64);index"`
	LastModifiedLedgerSeq uint32      `gorm:"column:last_modified_ledger_seq;type:int;not null"`
	CreatedAt             time.Time   `gorm:"column:created_at"`
	UpdatedAt             time.Time   `gorm:"column:updated_at"`
}

func (OfferEntry) TableName() string {
	return "offer_entries"
}

func UpsertOfferEntry(db *gorm.DB, offer *OfferEntry) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "offer_id"}, {Name: "seller_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"selling_asset_type", "selling_asset_code", "selling_asset_issuer",
			"buying_asset_type", "buying_asset_code", "buying_asset_issuer",
			"amount", "price", "flags", "ext", "sponsoring_id",
			"last_modified_ledger_seq", "updated_at",
		}),
	}).Create(offer).Error
}
