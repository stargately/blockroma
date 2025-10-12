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

// QueryEventsByTopicContains queries events where topic JSONB contains the specified value
// Uses PostgreSQL's @> (contains) operator with GIN index support
// Example: QueryEventsByTopicContains(db, `["transfer"]`) finds events with "transfer" in topic array
func QueryEventsByTopicContains(db *gorm.DB, topicJSON string) ([]Event, error) {
	var events []Event
	err := db.Where("topic @> ?", topicJSON).Find(&events).Error
	return events, err
}

// QueryEventsByTopicElement queries events where a specific topic array element matches a value
// Uses PostgreSQL's -> operator to extract array elements and ->> for text comparison
// Example: QueryEventsByTopicElement(db, 0, "transfer") finds events where topic[0] = "transfer"
func QueryEventsByTopicElement(db *gorm.DB, index int, value string) ([]Event, error) {
	var events []Event
	// Use ->> for text extraction and comparison
	err := db.Where("topic ->> ? = ?", index, value).Find(&events).Error
	return events, err
}

// QueryEventsByContractAndTopic queries events by contract ID and topic pattern
// Combines regular index (contract_id) with JSONB GIN index (topic) for efficient filtering
func QueryEventsByContractAndTopic(db *gorm.DB, contractID, topicJSON string) ([]Event, error) {
	var events []Event
	err := db.Where("contract_id = ? AND topic @> ?", contractID, topicJSON).Find(&events).Error
	return events, err
}

// QueryEventsByValuePath queries events where a specific path in the value JSONB matches
// Uses PostgreSQL's #> operator for path extraction and comparison
// Example: QueryEventsByValuePath(db, "{amount}", `"1000"`) finds events where value.amount = 1000
func QueryEventsByValuePath(db *gorm.DB, path string, value string) ([]Event, error) {
	var events []Event
	err := db.Where("value #> ? = ?", path, value).Find(&events).Error
	return events, err
}

// QueryTokenTransferEvents queries transfer events for a specific token contract
// This is a convenience function for the common pattern of finding token transfer events
func QueryTokenTransferEvents(db *gorm.DB, contractID string, limit int, offset int) ([]Event, error) {
	var events []Event
	err := db.Where("contract_id = ? AND topic @> ?", contractID, `["transfer"]`).
		Order("ledger DESC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

// QueryEventsByLedgerRange queries events within a ledger range
// Uses regular B-tree index on ledger column for efficient range scans
func QueryEventsByLedgerRange(db *gorm.DB, startLedger, endLedger int32) ([]Event, error) {
	var events []Event
	err := db.Where("ledger >= ? AND ledger <= ?", startLedger, endLedger).
		Order("ledger ASC, tx_index ASC").
		Find(&events).Error
	return events, err
}
