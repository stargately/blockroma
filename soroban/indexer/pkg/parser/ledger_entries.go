package parser

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/blockroma/soroban-indexer/pkg/models"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

// ParseLedgerEntry parses an XDR ledger entry and returns appropriate model(s)
func ParseLedgerEntry(xdrString string) ([]interface{}, error) {
	var entry xdr.LedgerEntry
	if err := xdr.SafeUnmarshalBase64(xdrString, &entry); err != nil {
		return nil, fmt.Errorf("unmarshal ledger entry: %w", err)
	}

	var results []interface{}

	// Get ledger key for hashing
	key, err := entry.LedgerKey()
	if err != nil {
		return nil, fmt.Errorf("get ledger key: %w", err)
	}

	bin, err := key.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal ledger key: %w", err)
	}
	keyHash := sha256.Sum256(bin)
	hexKey := hex.EncodeToString(keyHash[:])

	// Parse based on entry.Data.Type (not key type)
	// The key tells us what we requested, but entry.Data tells us what we got
	switch entry.Data.Type {
	case xdr.LedgerEntryTypeContractData:
		model := ParseContractDataEntry(entry, hexKey)
		if model != nil {
			results = append(results, model)
		}

	case xdr.LedgerEntryTypeAccount:
		model := ParseAccountEntry(entry)
		if model != nil {
			results = append(results, model)
		}

	case xdr.LedgerEntryTypeTrustline:
		model := ParseTrustLineEntry(entry)
		if model != nil {
			results = append(results, model)
		}

	case xdr.LedgerEntryTypeOffer:
		model := ParseOfferEntry(entry)
		if model != nil {
			results = append(results, model)
		}

	case xdr.LedgerEntryTypeData:
		model := ParseDataEntry(entry)
		if model != nil {
			results = append(results, model)
		}

	case xdr.LedgerEntryTypeClaimableBalance:
		model := ParseClaimableBalanceEntry(entry)
		if model != nil {
			results = append(results, model)
		}

	case xdr.LedgerEntryTypeLiquidityPool:
		model := ParseLiquidityPoolEntry(entry)
		if model != nil {
			results = append(results, model)
		}
	}

	return results, nil
}

// ParseContractDataEntry parses contract data entry
func ParseContractDataEntry(entry xdr.LedgerEntry, keyHash string) *models.ContractDataEntry {
	contractData := entry.Data.ContractData

	contractID := ""
	if contractData.Contract.ContractId != nil {
		contractID = strkey.MustEncode(strkey.VersionByteContract, contractData.Contract.ContractId[:])
	} else {
		return nil
	}

	keyXDR, _ := xdr.MarshalBase64(contractData.Key)
	valXDR, _ := xdr.MarshalBase64(contractData.Val)

	// Convert key and value to interface{} first
	keyInterface := ScValToInterface(contractData.Key)
	valInterface := ScValToInterface(contractData.Val)

	// Marshal to JSON bytes for JSONB storage
	keyBytes, err := json.Marshal(keyInterface)
	if err != nil {
		return nil // Return nil on error
	}
	valBytes, err := json.Marshal(valInterface)
	if err != nil {
		return nil // Return nil on error
	}

	durability := "persistent"
	if contractData.Durability == xdr.ContractDataDurabilityTemporary {
		durability = "temporary"
	}

	return &models.ContractDataEntry{
		KeyHash:    keyHash,
		ContractID: contractID,
		Key:        models.JSONB(keyBytes),
		KeyXdr:     keyXDR,
		Val:        models.JSONB(valBytes),
		ValXdr:     valXDR,
		Durability: durability,
	}
}

// ParseAccountEntry parses account entry
func ParseAccountEntry(entry xdr.LedgerEntry) *models.AccountEntry {
	account := entry.Data.Account
	inflationDest, _ := account.InflationDest.GetAddress()

	signersJSON := []byte("[]")
	if len(account.Signers) > 0 {
		signers := make([]map[string]interface{}, 0, len(account.Signers))
		for _, s := range account.Signers {
			signers = append(signers, map[string]interface{}{
				"key":    s.Key.Address(),
				"weight": uint32(s.Weight),
			})
		}
		signersJSON, _ = json.Marshal(signers)
	}

	extJSON, _ := json.Marshal(account.Ext)

	return &models.AccountEntry{
		AccountID:             account.AccountId.Address(),
		Balance:               int64(account.Balance),
		SeqNum:                int64(account.SeqNum),
		NumSubEntries:         uint32(account.NumSubEntries),
		InflationDest:         inflationDest,
		Flags:                 uint32(account.Flags),
		HomeDomain:            string(account.HomeDomain),
		Thresholds:            account.Thresholds[:],
		Signers:               signersJSON,
		Ext:                   extJSON,
		LastModifiedLedgerSeq: uint32(entry.LastModifiedLedgerSeq),
		SponsoringID:          getSponsoringID(&entry),
	}
}

// ParseTrustLineEntry parses trust line entry
func ParseTrustLineEntry(entry xdr.LedgerEntry) *models.TrustLineEntry {
	trustLine := entry.Data.TrustLine

	var poolID []byte
	if trustLine.Asset.LiquidityPoolId != nil {
		poolID = trustLine.Asset.LiquidityPoolId[:]
	}

	extJSON, _ := json.Marshal(trustLine.Ext)

	model := &models.TrustLineEntry{
		AccountID:             trustLine.AccountId.Address(),
		AssetType:             int32(trustLine.Asset.Type),
		LiquidityPoolID:       poolID,
		Balance:               int64(trustLine.Balance),
		Limit:                 int64(trustLine.Limit),
		Flags:                 uint32(trustLine.Flags),
		Ext:                   extJSON,
		LastModifiedLedgerSeq: uint32(entry.LastModifiedLedgerSeq),
		SponsoringID:          getSponsoringID(&entry),
	}

	switch trustLine.Asset.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		model.AssetCode = trustLine.Asset.AlphaNum4.AssetCode[:]
		model.AssetIssuer = trustLine.Asset.AlphaNum4.Issuer.Address()
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		model.AssetCode = trustLine.Asset.AlphaNum12.AssetCode[:]
		model.AssetIssuer = trustLine.Asset.AlphaNum12.Issuer.Address()
	case xdr.AssetTypeAssetTypePoolShare:
		model.AssetCode = trustLine.Asset.LiquidityPoolId[:]
	}

	return model
}

// ParseOfferEntry parses offer entry
func ParseOfferEntry(entry xdr.LedgerEntry) *models.OfferEntry {
	offer := entry.Data.Offer

	extJSON, _ := json.Marshal(offer.Ext)

	model := &models.OfferEntry{
		OfferID:               int64(offer.OfferId),
		SellerID:              offer.SellerId.Address(),
		SellingAssetType:      int32(offer.Selling.Type),
		BuyingAssetType:       int32(offer.Buying.Type),
		Amount:                int64(offer.Amount),
		Price:                 priceToString(offer.Price),
		Flags:                 uint32(offer.Flags),
		Ext:                   extJSON,
		LastModifiedLedgerSeq: uint32(entry.LastModifiedLedgerSeq),
		SponsoringID:          getSponsoringID(&entry),
	}

	switch offer.Selling.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		model.SellingAssetCode = offer.Selling.AlphaNum4.AssetCode[:]
		model.SellingAssetIssuer = offer.Selling.AlphaNum4.Issuer.Address()
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		model.SellingAssetCode = offer.Selling.AlphaNum12.AssetCode[:]
		model.SellingAssetIssuer = offer.Selling.AlphaNum12.Issuer.Address()
	}

	switch offer.Buying.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		model.BuyingAssetCode = offer.Buying.AlphaNum4.AssetCode[:]
		model.BuyingAssetIssuer = offer.Buying.AlphaNum4.Issuer.Address()
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		model.BuyingAssetCode = offer.Buying.AlphaNum12.AssetCode[:]
		model.BuyingAssetIssuer = offer.Buying.AlphaNum12.Issuer.Address()
	}

	return model
}

// ParseDataEntry parses account data entry
func ParseDataEntry(entry xdr.LedgerEntry) *models.DataEntry {
	data := entry.Data.Data

	extJSON, _ := json.Marshal(data.Ext)
	dataValueJSON, _ := json.Marshal(data.DataValue)

	return &models.DataEntry{
		AccountID:             data.AccountId.Address(),
		DataName:              string(data.DataName),
		DataValue:             dataValueJSON,
		Ext:                   extJSON,
		SponsoringID:          getSponsoringID(&entry),
		LastModifiedLedgerSeq: uint32(entry.LastModifiedLedgerSeq),
	}
}

// ParseClaimableBalanceEntry parses claimable balance entry
func ParseClaimableBalanceEntry(entry xdr.LedgerEntry) *models.ClaimableBalanceEntry {
	claimableBalance := entry.Data.ClaimableBalance

	claimantsJSON, _ := json.Marshal(claimableBalance.Claimants)

	model := &models.ClaimableBalanceEntry{
		BalanceID:             claimableBalance.BalanceId.V0.HexString(),
		Claimants:             claimantsJSON,
		AssetType:             int32(claimableBalance.Asset.Type),
		Amount:                int64(claimableBalance.Amount),
		LastModifiedLedgerSeq: uint32(entry.LastModifiedLedgerSeq),
		SponsoringID:          getSponsoringID(&entry),
	}

	switch claimableBalance.Asset.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		model.AssetCode = claimableBalance.Asset.AlphaNum4.AssetCode[:]
		model.AssetIssuer = claimableBalance.Asset.AlphaNum4.Issuer.Address()
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		model.AssetCode = claimableBalance.Asset.AlphaNum12.AssetCode[:]
		model.AssetIssuer = claimableBalance.Asset.AlphaNum12.Issuer.Address()
	}

	return model
}

// ParseLiquidityPoolEntry parses liquidity pool entry
func ParseLiquidityPoolEntry(entry xdr.LedgerEntry) *models.LiquidityPoolEntry {
	lp := entry.Data.LiquidityPool

	model := &models.LiquidityPoolEntry{
		LiquidityPoolID:          lp.LiquidityPoolId[:],
		Type:                     int32(lp.Body.Type),
		Fee:                      int32(lp.Body.ConstantProduct.Params.Fee),
		ReserveA:                 int64(lp.Body.ConstantProduct.ReserveA),
		ReserveB:                 int64(lp.Body.ConstantProduct.ReserveB),
		TotalPoolShares:          int64(lp.Body.ConstantProduct.TotalPoolShares),
		PoolSharesTrustLineCount: int64(lp.Body.ConstantProduct.PoolSharesTrustLineCount),
		LastModifiedLedgerSeq:    uint32(entry.LastModifiedLedgerSeq),
		SponsoringID:             getSponsoringID(&entry),
	}

	assetA := lp.Body.ConstantProduct.Params.AssetA
	model.AssetAType = int32(assetA.Type)
	switch assetA.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		model.AssetACode = assetA.AlphaNum4.AssetCode[:]
		model.AssetAIssuer = assetA.AlphaNum4.Issuer.Address()
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		model.AssetACode = assetA.AlphaNum12.AssetCode[:]
		model.AssetAIssuer = assetA.AlphaNum12.Issuer.Address()
	}

	assetB := lp.Body.ConstantProduct.Params.AssetB
	model.AssetBType = int32(assetB.Type)
	switch assetB.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		model.AssetBCode = assetB.AlphaNum4.AssetCode[:]
		model.AssetBIssuer = assetB.AlphaNum4.Issuer.Address()
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		model.AssetBCode = assetB.AlphaNum12.AssetCode[:]
		model.AssetBIssuer = assetB.AlphaNum12.Issuer.Address()
	}

	return model
}

// Helper functions

func getSponsoringID(entry *xdr.LedgerEntry) string {
	if sponsor := entry.SponsoringID(); sponsor != nil {
		return sponsor.Address()
	}
	return ""
}

func priceToString(p xdr.Price) string {
	if p.D == 0 {
		return "0"
	}
	// Calculate decimal representation: N/D
	return fmt.Sprintf("%d", p.N) + "/" + fmt.Sprintf("%d", p.D)
}

// BuildAccountLedgerKey builds a base64-encoded ledger key for an account address
func BuildAccountLedgerKey(accountAddress string) (string, error) {
	// Decode the Stellar address to get the public key
	decoded, err := strkey.Decode(strkey.VersionByteAccountID, accountAddress)
	if err != nil {
		return "", fmt.Errorf("decode account address: %w", err)
	}

	// Create AccountId from decoded bytes
	var accountID xdr.AccountId
	var uint256 xdr.Uint256
	copy(uint256[:], decoded)
	accountID.Type = xdr.PublicKeyTypePublicKeyTypeEd25519
	accountID.Ed25519 = &uint256

	// Create ledger key for account
	ledgerKey := xdr.LedgerKey{
		Type: xdr.LedgerEntryTypeAccount,
		Account: &xdr.LedgerKeyAccount{
			AccountId: accountID,
		},
	}

	// Marshal to base64
	xdrBytes, err := ledgerKey.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("marshal ledger key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(xdrBytes), nil
}

// BuildClaimableBalanceLedgerKey builds a base64-encoded ledger key for a claimable balance ID
func BuildClaimableBalanceLedgerKey(balanceIDHex string) (string, error) {
	// Decode hex string to bytes
	balanceIDBytes, err := hex.DecodeString(balanceIDHex)
	if err != nil {
		return "", fmt.Errorf("decode balance ID hex: %w", err)
	}

	// Create ClaimableBalanceID from bytes
	var hash xdr.Hash
	if len(balanceIDBytes) != 32 {
		return "", fmt.Errorf("invalid balance ID length: expected 32 bytes, got %d", len(balanceIDBytes))
	}
	copy(hash[:], balanceIDBytes)

	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &hash,
	}

	// Create ledger key for claimable balance
	ledgerKey := xdr.LedgerKey{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		ClaimableBalance: &xdr.LedgerKeyClaimableBalance{
			BalanceId: balanceID,
		},
	}

	// Marshal to base64
	xdrBytes, err := ledgerKey.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("marshal ledger key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(xdrBytes), nil
}

// ExtractClaimableBalanceIDs extracts claimable balance IDs from a transaction envelope XDR
func ExtractClaimableBalanceIDs(envelopeXdr string) ([]string, error) {
	data, err := base64.StdEncoding.DecodeString(envelopeXdr)
	if err != nil {
		return nil, fmt.Errorf("decode envelope xdr: %w", err)
	}

	var envelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("unmarshal envelope: %w", err)
	}

	var balanceIDs []string
	var operations []xdr.Operation

	// Extract operations based on envelope type
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		if v0, ok := envelope.GetV0(); ok {
			operations = v0.Tx.Operations
		}
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		if v1, ok := envelope.GetV1(); ok {
			operations = v1.Tx.Operations
		}
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		if fb, ok := envelope.GetFeeBump(); ok {
			if innerV1, ok := fb.Tx.InnerTx.GetV1(); ok {
				operations = innerV1.Tx.Operations
			}
		}
	}

	// Look for claimable balance operations
	for _, op := range operations {
		switch op.Body.Type {
		case xdr.OperationTypeCreateClaimableBalance:
			// For CreateClaimableBalance, we need to compute the balance ID
			// The balance ID is derived from the operation source and sequence number
			// This is complex, so we'll skip it for now and rely on claim operations
			continue
		case xdr.OperationTypeClaimClaimableBalance:
			if claimOp, ok := op.Body.GetClaimClaimableBalanceOp(); ok {
				balanceIDs = append(balanceIDs, claimOp.BalanceId.V0.HexString())
			}
		}
	}

	return balanceIDs, nil
}

// BuildContractDataKey builds a base64-encoded ledger key for contract data
func BuildContractDataKey(contractID string, key xdr.ScVal, durability xdr.ContractDataDurability) (string, error) {
	// Decode the contract ID
	decoded, err := strkey.Decode(strkey.VersionByteContract, contractID)
	if err != nil {
		return "", fmt.Errorf("decode contract ID: %w", err)
	}

	// Create contract address
	var hash xdr.Hash
	copy(hash[:], decoded)

	// Convert Hash to ContractId (they are the same type)
	contractIDXDR := xdr.ContractId(hash)

	contractAddress := xdr.ScAddress{
		Type:       xdr.ScAddressTypeScAddressTypeContract,
		ContractId: &contractIDXDR,
	}

	// Create ledger key for contract data
	ledgerKey := xdr.LedgerKey{
		Type: xdr.LedgerEntryTypeContractData,
		ContractData: &xdr.LedgerKeyContractData{
			Contract:   contractAddress,
			Key:        key,
			Durability: durability,
		},
	}

	// Marshal to base64
	xdrBytes, err := ledgerKey.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("marshal ledger key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(xdrBytes), nil
}

// BuildMetadataKey creates an ScVal for the METADATA key
func BuildMetadataKey() xdr.ScVal {
	return xdr.ScVal{
		Type: xdr.ScValTypeScvLedgerKeyContractInstance,
	}
}

// BuildBalanceKey creates an ScVal for a Balance(address) key
func BuildBalanceKey(address string) (xdr.ScVal, error) {
	// Decode the address to get the raw bytes
	decoded, err := strkey.Decode(strkey.VersionByteAccountID, address)
	if err != nil {
		// Try as contract address
		decoded, err = strkey.Decode(strkey.VersionByteContract, address)
		if err != nil {
			return xdr.ScVal{}, fmt.Errorf("decode address: %w", err)
		}
		// Contract address - use ContractId type
		var hash xdr.Hash
		copy(hash[:], decoded)
		contractID := xdr.ContractId(hash)

		scAddress := xdr.ScAddress{
			Type:       xdr.ScAddressTypeScAddressTypeContract,
			ContractId: &contractID,
		}

		balanceSymbol := xdr.ScSymbol("Balance")
		vec := xdr.ScVec{
			// ["Balance", Address]
			xdr.ScVal{
				Type: xdr.ScValTypeScvSymbol,
				Sym:  &balanceSymbol,
			},
			xdr.ScVal{
				Type:    xdr.ScValTypeScvAddress,
				Address: &scAddress,
			},
		}
		vecPtr := &vec

		return xdr.ScVal{
			Type: xdr.ScValTypeScvVec,
			Vec:  &vecPtr,
		}, nil
	}

	// Account address
	var uint256 xdr.Uint256
	copy(uint256[:], decoded)

	accountID := xdr.AccountId{
		Type:    xdr.PublicKeyTypePublicKeyTypeEd25519,
		Ed25519: &uint256,
	}

	scAddress := xdr.ScAddress{
		Type:      xdr.ScAddressTypeScAddressTypeAccount,
		AccountId: &accountID,
	}

	balanceSymbol := xdr.ScSymbol("Balance")
	vec := xdr.ScVec{
		// ["Balance", Address]
		xdr.ScVal{
			Type: xdr.ScValTypeScvSymbol,
			Sym:  &balanceSymbol,
		},
		xdr.ScVal{
			Type:    xdr.ScValTypeScvAddress,
			Address: &scAddress,
		},
	}
	vecPtr := &vec

	return xdr.ScVal{
		Type: xdr.ScValTypeScvVec,
		Vec:  &vecPtr,
	}, nil
}
