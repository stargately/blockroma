package parser

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
	"github.com/blockroma/soroban-indexer/pkg/models"
)

// ParseOperations extracts all operations from a transaction envelope XDR
func ParseOperations(txHash string, envelopeXdr string) ([]*models.Operation, error) {
	data, err := base64.StdEncoding.DecodeString(envelopeXdr)
	if err != nil {
		return nil, fmt.Errorf("decode envelope xdr: %w", err)
	}

	var envelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("unmarshal envelope: %w", err)
	}

	var operations []xdr.Operation
	var sourceAccount xdr.MuxedAccount

	// Extract operations and source account based on envelope type
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		if v0, ok := envelope.GetV0(); ok {
			operations = v0.Tx.Operations
			// TransactionV0 has SourceAccountEd25519 (Uint256), convert to MuxedAccount
			var ed25519Account xdr.Uint256 = v0.Tx.SourceAccountEd25519
			sourceAccount = xdr.MuxedAccount{
				Type: xdr.CryptoKeyTypeKeyTypeEd25519,
				Ed25519: &ed25519Account,
			}
		}
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		if v1, ok := envelope.GetV1(); ok {
			operations = v1.Tx.Operations
			sourceAccount = v1.Tx.SourceAccount
		}
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		if fb, ok := envelope.GetFeeBump(); ok {
			if innerV1, ok := fb.Tx.InnerTx.GetV1(); ok {
				operations = innerV1.Tx.Operations
				sourceAccount = innerV1.Tx.SourceAccount
			}
		}
	default:
		return nil, fmt.Errorf("unsupported envelope type: %v", envelope.Type)
	}

	// Parse each operation
	result := make([]*models.Operation, 0, len(operations))
	for i, op := range operations {
		parsedOp, err := parseOperation(txHash, i, op, sourceAccount)
		if err != nil {
			// Log error but continue processing other operations
			continue
		}
		result = append(result, parsedOp)
	}

	return result, nil
}

// parseOperation converts a single XDR operation to a database model
func parseOperation(txHash string, index int, op xdr.Operation, txSourceAccount xdr.MuxedAccount) (*models.Operation, error) {
	// Generate operation ID
	opID := fmt.Sprintf("%s-%d", txHash, index)

	// Get operation source account (use transaction source if not specified)
	var sourceAccount string
	if op.SourceAccount != nil {
		sourceAccount = op.SourceAccount.Address()
	} else {
		sourceAccount = txSourceAccount.Address()
	}

	// Get operation type
	opType := op.Body.Type.String()

	// Parse operation-specific details
	details, err := parseOperationDetails(op.Body)
	if err != nil {
		return nil, fmt.Errorf("parse operation details: %w", err)
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("marshal operation details: %w", err)
	}

	return &models.Operation{
		ID:               opID,
		TxHash:           txHash,
		OperationIndex:   int32(index),
		SourceAccount:    sourceAccount,
		OperationType:    opType,
		OperationDetails: detailsJSON,
	}, nil
}

// parseOperationDetails extracts type-specific details from an operation
func parseOperationDetails(body xdr.OperationBody) (map[string]interface{}, error) {
	details := make(map[string]interface{})

	switch body.Type {
	case xdr.OperationTypeCreateAccount:
		if op, ok := body.GetCreateAccountOp(); ok {
			details["destination"] = op.Destination.Address()
			details["starting_balance"] = int64(op.StartingBalance)
		}

	case xdr.OperationTypePayment:
		if op, ok := body.GetPaymentOp(); ok {
			details["destination"] = op.Destination.Address()
			details["asset"] = assetToMap(op.Asset)
			details["amount"] = int64(op.Amount)
		}

	case xdr.OperationTypePathPaymentStrictReceive:
		if op, ok := body.GetPathPaymentStrictReceiveOp(); ok {
			details["send_asset"] = assetToMap(op.SendAsset)
			details["send_max"] = int64(op.SendMax)
			details["destination"] = op.Destination.Address()
			details["dest_asset"] = assetToMap(op.DestAsset)
			details["dest_amount"] = int64(op.DestAmount)
			details["path"] = assetsToArray(op.Path)
		}

	case xdr.OperationTypePathPaymentStrictSend:
		if op, ok := body.GetPathPaymentStrictSendOp(); ok {
			details["send_asset"] = assetToMap(op.SendAsset)
			details["send_amount"] = int64(op.SendAmount)
			details["destination"] = op.Destination.Address()
			details["dest_asset"] = assetToMap(op.DestAsset)
			details["dest_min"] = int64(op.DestMin)
			details["path"] = assetsToArray(op.Path)
		}

	case xdr.OperationTypeManageSellOffer:
		if op, ok := body.GetManageSellOfferOp(); ok {
			details["selling"] = assetToMap(op.Selling)
			details["buying"] = assetToMap(op.Buying)
			details["amount"] = int64(op.Amount)
			details["price"] = priceToMap(op.Price)
			details["offer_id"] = int64(op.OfferId)
		}

	case xdr.OperationTypeManageBuyOffer:
		if op, ok := body.GetManageBuyOfferOp(); ok {
			details["selling"] = assetToMap(op.Selling)
			details["buying"] = assetToMap(op.Buying)
			details["buy_amount"] = int64(op.BuyAmount)
			details["price"] = priceToMap(op.Price)
			details["offer_id"] = int64(op.OfferId)
		}

	case xdr.OperationTypeCreatePassiveSellOffer:
		if op, ok := body.GetCreatePassiveSellOfferOp(); ok {
			details["selling"] = assetToMap(op.Selling)
			details["buying"] = assetToMap(op.Buying)
			details["amount"] = int64(op.Amount)
			details["price"] = priceToMap(op.Price)
		}

	case xdr.OperationTypeSetOptions:
		if op, ok := body.GetSetOptionsOp(); ok {
			if op.InflationDest != nil {
				details["inflation_dest"] = op.InflationDest.Address()
			}
			if op.ClearFlags != nil {
				details["clear_flags"] = uint32(*op.ClearFlags)
			}
			if op.SetFlags != nil {
				details["set_flags"] = uint32(*op.SetFlags)
			}
			if op.MasterWeight != nil {
				details["master_weight"] = uint32(*op.MasterWeight)
			}
			if op.LowThreshold != nil {
				details["low_threshold"] = uint32(*op.LowThreshold)
			}
			if op.MedThreshold != nil {
				details["med_threshold"] = uint32(*op.MedThreshold)
			}
			if op.HighThreshold != nil {
				details["high_threshold"] = uint32(*op.HighThreshold)
			}
			if op.HomeDomain != nil {
				details["home_domain"] = string(*op.HomeDomain)
			}
			if op.Signer != nil {
				details["signer"] = map[string]interface{}{
					"key":    op.Signer.Key.Address(),
					"weight": uint32(op.Signer.Weight),
				}
			}
		}

	case xdr.OperationTypeChangeTrust:
		if op, ok := body.GetChangeTrustOp(); ok {
			details["line"] = changeTrustAssetToMap(op.Line)
			details["limit"] = int64(op.Limit)
		}

	case xdr.OperationTypeAllowTrust:
		if op, ok := body.GetAllowTrustOp(); ok {
			details["trustor"] = op.Trustor.Address()
			// Extract asset code from AssetCode union
			var assetCode string
			switch op.Asset.Type {
			case xdr.AssetTypeAssetTypeCreditAlphanum4:
				if op.Asset.AssetCode4 != nil {
					assetCode = string((*op.Asset.AssetCode4)[:])
				}
			case xdr.AssetTypeAssetTypeCreditAlphanum12:
				if op.Asset.AssetCode12 != nil {
					assetCode = string((*op.Asset.AssetCode12)[:])
				}
			}
			details["asset"] = assetCode
			details["authorize"] = uint32(op.Authorize)
		}

	case xdr.OperationTypeAccountMerge:
		if destination, ok := body.GetDestination(); ok {
			details["destination"] = destination.Address()
		}

	case xdr.OperationTypeInflation:
		// No additional details for inflation

	case xdr.OperationTypeManageData:
		if op, ok := body.GetManageDataOp(); ok {
			details["data_name"] = string(op.DataName)
			if op.DataValue != nil {
				details["data_value"] = []byte(*op.DataValue)
			}
		}

	case xdr.OperationTypeBumpSequence:
		if op, ok := body.GetBumpSequenceOp(); ok {
			details["bump_to"] = int64(op.BumpTo)
		}

	case xdr.OperationTypeCreateClaimableBalance:
		if op, ok := body.GetCreateClaimableBalanceOp(); ok {
			details["asset"] = assetToMap(op.Asset)
			details["amount"] = int64(op.Amount)
			claimants := make([]map[string]interface{}, len(op.Claimants))
			for i, c := range op.Claimants {
				claimants[i] = map[string]interface{}{
					"destination": c.MustV0().Destination.Address(),
					"predicate":   "complex", // Simplified for now
				}
			}
			details["claimants"] = claimants
		}

	case xdr.OperationTypeClaimClaimableBalance:
		if op, ok := body.GetClaimClaimableBalanceOp(); ok {
			details["balance_id"] = op.BalanceId.V0.HexString()
		}

	case xdr.OperationTypeBeginSponsoringFutureReserves:
		if op, ok := body.GetBeginSponsoringFutureReservesOp(); ok {
			details["sponsored_id"] = op.SponsoredId.Address()
		}

	case xdr.OperationTypeEndSponsoringFutureReserves:
		// No additional details

	case xdr.OperationTypeRevokeSponsorship:
		if op, ok := body.GetRevokeSponsorshipOp(); ok {
			switch op.Type {
			case xdr.RevokeSponsorshipTypeRevokeSponsorshipLedgerEntry:
				details["type"] = "ledger_entry"
			case xdr.RevokeSponsorshipTypeRevokeSponsorshipSigner:
				details["type"] = "signer"
			}
		}

	case xdr.OperationTypeClawback:
		if op, ok := body.GetClawbackOp(); ok {
			details["asset"] = assetToMap(op.Asset)
			details["from"] = op.From.Address()
			details["amount"] = int64(op.Amount)
		}

	case xdr.OperationTypeClawbackClaimableBalance:
		if op, ok := body.GetClawbackClaimableBalanceOp(); ok {
			details["balance_id"] = op.BalanceId.V0.HexString()
		}

	case xdr.OperationTypeSetTrustLineFlags:
		if op, ok := body.GetSetTrustLineFlagsOp(); ok {
			details["trustor"] = op.Trustor.Address()
			details["asset"] = assetToMap(op.Asset)
			details["clear_flags"] = uint32(op.ClearFlags)
			details["set_flags"] = uint32(op.SetFlags)
		}

	case xdr.OperationTypeLiquidityPoolDeposit:
		if op, ok := body.GetLiquidityPoolDepositOp(); ok {
			details["liquidity_pool_id"] = fmt.Sprintf("%x", op.LiquidityPoolId)
			details["max_amount_a"] = int64(op.MaxAmountA)
			details["max_amount_b"] = int64(op.MaxAmountB)
			details["min_price"] = priceToMap(op.MinPrice)
			details["max_price"] = priceToMap(op.MaxPrice)
		}

	case xdr.OperationTypeLiquidityPoolWithdraw:
		if op, ok := body.GetLiquidityPoolWithdrawOp(); ok {
			details["liquidity_pool_id"] = fmt.Sprintf("%x", op.LiquidityPoolId)
			details["amount"] = int64(op.Amount)
			details["min_amount_a"] = int64(op.MinAmountA)
			details["min_amount_b"] = int64(op.MinAmountB)
		}

	case xdr.OperationTypeInvokeHostFunction:
		if op, ok := body.GetInvokeHostFunctionOp(); ok {
			details["host_function"] = op.HostFunction.Type.String()
			// Don't include full function details - too complex
		}

	case xdr.OperationTypeExtendFootprintTtl:
		if op, ok := body.GetExtendFootprintTtlOp(); ok {
			details["extend_to"] = uint32(op.ExtendTo)
		}

	case xdr.OperationTypeRestoreFootprint:
		// No additional details

	default:
		details["error"] = "unknown operation type"
	}

	return details, nil
}

// ComputeClaimableBalanceID computes the claimable balance ID for a CreateClaimableBalance operation
// The ID is computed as: sha256(sourceAccount || seqNum || opIndex || network_id)
func ComputeClaimableBalanceID(sourceAccount string, seqNum int64, opIndex int, networkPassphrase string) (string, error) {
	// Decode source account
	decoded, err := strkey.Decode(strkey.VersionByteAccountID, sourceAccount)
	if err != nil {
		return "", fmt.Errorf("decode source account: %w", err)
	}

	// Create the hash input: account_id + sequence_number + operation_index + CLAIMABLE_BALANCE_ID_TYPE + network_id
	h := sha256.New()

	// Write network ID hash
	networkHash := sha256.Sum256([]byte(networkPassphrase))
	h.Write(networkHash[:])

	// Write account ID
	h.Write(decoded)

	// Write sequence number (8 bytes, big endian)
	seqBytes := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		seqBytes[i] = byte(seqNum)
		seqNum >>= 8
	}
	h.Write(seqBytes)

	// Write operation index (4 bytes, big endian)
	opBytes := make([]byte, 4)
	for i := 3; i >= 0; i-- {
		opBytes[i] = byte(opIndex)
		opIndex >>= 8
	}
	h.Write(opBytes)

	// Write claimable balance ID type (0)
	h.Write([]byte{0, 0, 0, 0})

	result := h.Sum(nil)
	return fmt.Sprintf("%x", result), nil
}

// Helper functions

func assetToMap(asset xdr.Asset) map[string]interface{} {
	result := make(map[string]interface{})
	result["type"] = asset.Type.String()

	switch asset.Type {
	case xdr.AssetTypeAssetTypeNative:
		result["code"] = "XLM"
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		result["code"] = string(asset.AlphaNum4.AssetCode[:])
		result["issuer"] = asset.AlphaNum4.Issuer.Address()
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		result["code"] = string(asset.AlphaNum12.AssetCode[:])
		result["issuer"] = asset.AlphaNum12.Issuer.Address()
	}

	return result
}

func changeTrustAssetToMap(asset xdr.ChangeTrustAsset) map[string]interface{} {
	result := make(map[string]interface{})
	result["type"] = asset.Type.String()

	switch asset.Type {
	case xdr.AssetTypeAssetTypeNative:
		result["code"] = "XLM"
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		if asset.AlphaNum4 != nil {
			result["code"] = string(asset.AlphaNum4.AssetCode[:])
			result["issuer"] = asset.AlphaNum4.Issuer.Address()
		}
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		if asset.AlphaNum12 != nil {
			result["code"] = string(asset.AlphaNum12.AssetCode[:])
			result["issuer"] = asset.AlphaNum12.Issuer.Address()
		}
	case xdr.AssetTypeAssetTypePoolShare:
		if asset.LiquidityPool != nil {
			// Compute the liquidity pool ID from the parameters
			// Pool ID = sha256(LIQUIDITY_POOL_FEE_V18 || assetA || assetB || fee)
			result["liquidity_pool_params"] = "pool_share"
		}
	}

	return result
}

func assetsToArray(assets []xdr.Asset) []map[string]interface{} {
	result := make([]map[string]interface{}, len(assets))
	for i, asset := range assets {
		result[i] = assetToMap(asset)
	}
	return result
}

func priceToMap(price xdr.Price) map[string]interface{} {
	return map[string]interface{}{
		"n": int32(price.N),
		"d": int32(price.D),
	}
}
