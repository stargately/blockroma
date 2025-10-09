package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Event struct {
	ID                       string      `gorm:"column:id;primaryKey"`
	TxIndex                  int32       `gorm:"column:tx_index"`
	EventType                string      `gorm:"column:type"`
	Ledger                   int32       `gorm:"column:ledger"`
	LedgerClosedAt           string      `gorm:"column:ledger_closed_at"`
	ContractID               string      `gorm:"column:contract_id"`
	PagingToken              string      `gorm:"column:paging_token"`
	Topic                    interface{} `gorm:"column:topic;type:jsonb"`
	Value                    interface{} `gorm:"column:value;type:jsonb"`
	InSuccessfulContractCall bool        `gorm:"column:in_successful_contract_call"`
	LastModifiedLedgerSeq    uint32      `gorm:"column:last_modified_ledger_seq;type:int;not null"`
	CreatedAt                time.Time   `gorm:"column:created_at"`
	UpdatedAt                time.Time   `gorm:"column:updated_at"`
}

func (Event) TableName() string {
	return "events"
}

func UpsertEvent(db *gorm.DB, event *Event) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"tx_index", "type", "ledger", "ledger_closed_at", "contract_id",
			"paging_token", "topic", "value", "in_successful_contract_call",
			"last_modified_ledger_seq", "updated_at",
		}),
	}).Create(event).Error
}

func NewEvent(data []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	return &event, err
}
