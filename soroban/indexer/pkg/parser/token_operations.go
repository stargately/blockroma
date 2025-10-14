package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blockroma/soroban-indexer/pkg/models"
	"github.com/blockroma/soroban-indexer/pkg/models/util"
)

// ParseTokenOperation extracts token operations from contract events
// Returns nil if the event is not a recognized token operation
func ParseTokenOperation(eventID string, contractID string, ledger uint32, ledgerClosedAt time.Time, txIndex int32, topics []interface{}, value interface{}) *models.TokenOperation {
	// Convert topics to strings
	topicStrs := make([]string, 0, len(topics))
	for _, t := range topics {
		switch v := t.(type) {
		case string:
			topicStrs = append(topicStrs, strings.Trim(v, "\""))
		default:
			topicStrs = append(topicStrs, fmt.Sprintf("%v", v))
		}
	}

	if len(topicStrs) == 0 {
		return nil
	}

	opType := topicStrs[0]

	// Initialize base token operation
	initOp := func(from string) models.TokenOperation {
		return models.TokenOperation{
			ID:             eventID,
			Type:           opType,
			TxIndex:        txIndex,
			Ledger:         int32(ledger),
			LedgerClosedAt: ledgerClosedAt.Format(time.RFC3339),
			ContractID:     contractID,
			From:           from,
		}
	}

	switch opType {
	case "transfer":
		if len(topicStrs) < 3 {
			return nil
		}
		amount := getInt128FromInterface(value)
		op := initOp(topicStrs[1])
		op.To = &topicStrs[2]
		op.Amount = &amount
		return &op

	case "mint":
		if len(topicStrs) < 3 {
			return nil
		}
		amount := getInt128FromInterface(value)
		op := initOp(topicStrs[1])
		op.To = &topicStrs[2]
		op.Amount = &amount
		return &op

	case "burn":
		if len(topicStrs) < 2 {
			return nil
		}
		amount := getInt128FromInterface(value)
		op := initOp(topicStrs[1])
		op.Amount = &amount
		return &op

	case "clawback":
		if len(topicStrs) < 3 {
			return nil
		}
		amount := getInt128FromInterface(value)
		op := initOp(topicStrs[2])
		op.To = &topicStrs[1]
		op.Amount = &amount
		return &op

	case "approve":
		if len(topicStrs) < 3 {
			return nil
		}
		// Value is an array: [amount, expiration_ledger]
		var data []interface{}
		valueBytes, _ := json.Marshal(value)
		json.Unmarshal(valueBytes, &data)

		if len(data) >= 2 {
			amount := getInt128FromInterface(data[0])
			expiration := int32(getIntFromInterface(data[1]))

			op := initOp(topicStrs[1])
			op.To = &topicStrs[2]
			op.Amount = &amount
			op.ExpirationLedger = &expiration
			return &op
		}
		return nil

	case "set_authorized":
		if len(topicStrs) < 3 {
			return nil
		}
		authorized := getBoolFromInterface(value)
		op := initOp(topicStrs[1])
		op.To = &topicStrs[2]
		op.Authorized = &authorized
		return &op

	case "set_admin":
		if len(topicStrs) < 2 {
			return nil
		}
		admin := getStringFromInterface(value)
		op := initOp(topicStrs[1])
		op.To = &admin
		return &op

	default:
		// Not a recognized token operation
		return nil
	}
}

// ParseTokenMetadata extracts token metadata from contract instance data
// Key should be an object with type "LedgerKeyContractInstance" or the legacy string "ScvLedgerKeyContractInstance"
// Value should be the contract instance structure
func ParseTokenMetadata(contractID string, key interface{}, value interface{}) *models.TokenMetadata {
	// Check if key is the contract instance key
	isContractInstanceKey := false

	switch k := key.(type) {
	case string:
		// Legacy format: plain string
		if k == "ScvLedgerKeyContractInstance" || k == "\"ScvLedgerKeyContractInstance\"" {
			isContractInstanceKey = true
		}
	case map[string]interface{}:
		// New format: object with type field
		if typeVal, ok := k["type"]; ok {
			if typeStr, ok := typeVal.(string); ok && typeStr == "LedgerKeyContractInstance" {
				isContractInstanceKey = true
			}
		}
	default:
		// Try to parse as JSON object
		keyBytes, _ := json.Marshal(key)
		var keyMap map[string]interface{}
		if json.Unmarshal(keyBytes, &keyMap) == nil {
			if typeVal, ok := keyMap["type"]; ok {
				if typeStr, ok := typeVal.(string); ok && typeStr == "LedgerKeyContractInstance" {
					isContractInstanceKey = true
				}
			}
		}
	}

	if !isContractInstanceKey {
		return nil
	}

	// Parse contract instance
	var instance struct {
		Storage []struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
		} `json:"storage"`
		Executable interface{} `json:"executable"`
	}

	// Value might already be a map or need to be unmarshaled from JSON
	switch v := value.(type) {
	case map[string]interface{}:
		// Already a map, marshal and unmarshal to fit the struct
		valueBytes, _ := json.Marshal(v)
		if err := json.Unmarshal(valueBytes, &instance); err != nil {
			return nil
		}
	case string:
		// It's a JSON string, unmarshal directly
		if err := json.Unmarshal([]byte(v), &instance); err != nil {
			return nil
		}
	default:
		// Try to marshal and unmarshal
		valueBytes, _ := json.Marshal(value)
		if err := json.Unmarshal(valueBytes, &instance); err != nil {
			return nil
		}
	}

	metadata := &models.TokenMetadata{
		ContractID: contractID,
	}

	for _, item := range instance.Storage {
		if item.Key == "METADATA" {
			// Parse metadata structure
			metaBytes, _ := json.Marshal(item.Value)
			var metaPairs []struct {
				Key   string      `json:"key"`
				Value interface{} `json:"value"`
			}
			json.Unmarshal(metaBytes, &metaPairs)

			for _, pair := range metaPairs {
				switch pair.Key {
				case "name":
					metadata.Name = getStringFromInterface(pair.Value)
				case "symbol":
					metadata.Symbol = getStringFromInterface(pair.Value)
				case "decimal":
					metadata.Decimal = uint32(getIntFromInterface(pair.Value))
				}
			}
		} else if item.Key == "Admin" {
			metadata.AdminAddress = getStringFromInterface(item.Value)
		}
	}

	// Only return if we found valid metadata
	if metadata.Name != "" || metadata.Symbol != "" {
		return metadata
	}

	return nil
}

// ParseTokenBalance extracts token balance from contract data
// Key format: ["Balance", "ADDRESS"]
// Value format: either a number or a struct with "amount" field
func ParseTokenBalance(contractID string, key interface{}, value interface{}) *models.TokenBalance {
	// Parse key
	var keyArray []string
	keyBytes, _ := json.Marshal(key)
	if err := json.Unmarshal(keyBytes, &keyArray); err != nil {
		return nil
	}

	if len(keyArray) != 2 || keyArray[0] != "Balance" {
		return nil
	}

	address := keyArray[1]

	// Parse value
	var balance string

	// Try as string/number first
	if str := getStringFromInterface(value); str != "" {
		balance = str
	} else {
		// Try as struct with "amount" field
		var pairs []struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
		}
		valueBytes, _ := json.Marshal(value)
		if err := json.Unmarshal(valueBytes, &pairs); err == nil {
			for _, pair := range pairs {
				if pair.Key == "amount" {
					balance = getStringFromInterface(pair.Value)
					break
				}
			}
		}
	}

	if balance == "" {
		return nil
	}

	return &models.TokenBalance{
		ContractID: contractID,
		Address:    address,
		Balance:    balance,
	}
}

// Helper functions

func getInt128FromInterface(v interface{}) util.Int128 {
	str := getStringFromInterface(v)
	var i128 util.Int128
	i128.SetString(str, 10)
	return i128
}

func getIntFromInterface(v interface{}) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	case string:
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	default:
		return 0
	}
}

func getBoolFromInterface(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true"
	default:
		return false
	}
}

func getStringFromInterface(v interface{}) string {
	switch val := v.(type) {
	case string:
		return strings.Trim(val, "\"")
	default:
		return fmt.Sprintf("%v", val)
	}
}
