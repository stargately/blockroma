package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/blockroma/soroban-indexer/pkg/models"
	"github.com/blockroma/soroban-indexer/pkg/parser"
)

func main() {
	// Get database connection
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Check total contract_data_entries
	var totalCount int64
	if err := db.Model(&models.ContractDataEntry{}).Count(&totalCount).Error; err != nil {
		log.Fatalf("Failed to count contract_data_entries: %v", err)
	}
	fmt.Printf("Total contract_data_entries: %d\n", totalCount)

	// Sample some keys to see their format
	var entries []models.ContractDataEntry
	if err := db.Limit(100).Find(&entries).Error; err != nil {
		log.Fatalf("Failed to fetch contract_data_entries: %v", err)
	}

	fmt.Printf("\nSampled %d entries:\n", len(entries))

	keyTypes := make(map[string]int)
	instanceKeyCount := 0
	balanceKeyCount := 0
	metadataExtractable := 0

	for _, entry := range entries {
		// Unmarshal key to see its structure
		var keyInterface interface{}
		if err := json.Unmarshal(entry.Key, &keyInterface); err != nil {
			fmt.Printf("Failed to unmarshal key: %v\n", err)
			continue
		}

		// Unmarshal value
		var valInterface interface{}
		if err := json.Unmarshal(entry.Val, &valInterface); err != nil {
			fmt.Printf("Failed to unmarshal val: %v\n", err)
			continue
		}

		// Determine key type
		switch k := keyInterface.(type) {
		case string:
			keyTypes[fmt.Sprintf("string:%s", k)]++

			// Check if it's an instance key
			if k == "ScvLedgerKeyContractInstance" || k == "\"ScvLedgerKeyContractInstance\"" {
				instanceKeyCount++
				fmt.Printf("\n✓ Found contract instance key (string format): %s\n", k)
				fmt.Printf("  Contract ID: %s\n", entry.ContractID)

				// Try to parse metadata
				if metadata := parser.ParseTokenMetadata(entry.ContractID, keyInterface, valInterface); metadata != nil {
					metadataExtractable++
					fmt.Printf("  ✓ Metadata extractable: name=%s, symbol=%s, decimal=%d\n",
						metadata.Name, metadata.Symbol, metadata.Decimal)
				} else {
					fmt.Printf("  ✗ Metadata NOT extractable from this entry\n")
				}
			}
		case map[string]interface{}:
			if typeVal, ok := k["type"]; ok {
				keyTypes[fmt.Sprintf("map:%v", typeVal)]++

				// Check if it's an instance key (old format)
				if typeStr, ok := typeVal.(string); ok && typeStr == "LedgerKeyContractInstance" {
					instanceKeyCount++
					fmt.Printf("\n✓ Found contract instance key (map format): %v\n", k)
					fmt.Printf("  Contract ID: %s\n", entry.ContractID)

					// Try to parse metadata
					if metadata := parser.ParseTokenMetadata(entry.ContractID, keyInterface, valInterface); metadata != nil {
						metadataExtractable++
						fmt.Printf("  ✓ Metadata extractable: name=%s, symbol=%s, decimal=%d\n",
							metadata.Name, metadata.Symbol, metadata.Decimal)
					} else {
						fmt.Printf("  ✗ Metadata NOT extractable from this entry\n")
					}
				}
			} else {
				keyTypes["map:other"]++
			}
		case []interface{}:
			if len(k) >= 2 {
				if k[0] == "Balance" {
					balanceKeyCount++
				}
				keyTypes[fmt.Sprintf("array:%v", k[0])]++
			} else {
				keyTypes["array:empty"]++
			}
		default:
			keyTypes[fmt.Sprintf("unknown:%T", keyInterface)]++
		}
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total entries sampled: %d\n", len(entries))
	fmt.Printf("Contract instance keys found: %d\n", instanceKeyCount)
	fmt.Printf("Balance keys found: %d\n", balanceKeyCount)
	fmt.Printf("Metadata extractable entries: %d\n", metadataExtractable)

	fmt.Printf("\nKey type distribution:\n")
	for keyType, count := range keyTypes {
		fmt.Printf("  %s: %d\n", keyType, count)
	}

	// Now check if token_metadata table has any entries
	var metadataCount int64
	if err := db.Model(&models.TokenMetadata{}).Count(&metadataCount).Error; err != nil {
		log.Fatalf("Failed to count token_metadata: %v", err)
	}
	fmt.Printf("\ntoken_metadata table count: %d\n", metadataCount)

	// Check token_balances
	var balanceCount int64
	if err := db.Model(&models.TokenBalance{}).Count(&balanceCount).Error; err != nil {
		log.Fatalf("Failed to count token_balances: %v", err)
	}
	fmt.Printf("token_balances table count: %d\n", balanceCount)
}
