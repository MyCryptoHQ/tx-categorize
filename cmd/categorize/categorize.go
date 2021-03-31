package categorize

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethertypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mycryptohq/tx-categorize/types"
)

func contains(a []string, x string) bool {
	for _, n := range a {
		if strings.EqualFold(x, n) {
			return true
		}
	}
	return false
}

func DetermineTxType(tx types.PreDeterminedStandardTx, schemaList []types.FullTxLabelSchema) (types.StandardTx, error) {
	sortedSchemas := sortSchemaListByPriority(schemaList)
	appliedSchemas := make(map[string]types.FullTxLabelSchema)
	for _, schemaItem := range sortedSchemas {
		schema := schemaItem.Schema
		var listToUse []string
		switch schema.Type {
		case "topics":
			listToUse = schema.Topics
		case "addresses":
			listToUse = schema.Addresses
		case "statuses":
			listToUse = schema.Statuses
		}
		if handleBaseKeys(tx, schema.Key, listToUse) {

			if len(schemaItem.Schema.And) != 0 {
				if otherRequirementsPresent(appliedSchemas, schemaItem.Schema.And) {
					appliedSchemas[schemaItem.Meta.Name] = schemaItem
				}
			} else {
				appliedSchemas[schemaItem.Meta.Name] = schemaItem
			}
		}
	}
	if len(appliedSchemas) == 0 {
		if tx.Data == "0x" {
			appliedSchemas["STANDARD"] = types.FullTxLabelSchema{}
		} else {
			appliedSchemas["GENERIC_CONTRACT_CALL"] = types.FullTxLabelSchema{}
		}
	}

	txTypeDerived := interpretAppliedSchemas(appliedSchemas)
	return types.StandardTx{
		To:               tx.To,
		From:             tx.From,
		Value:            fmt.Sprintf("%#x", tx.Value),
		BlockNumber:      fmt.Sprintf("%#x", tx.BlockNumber),
		TimeStamp:        tx.TimeStamp,
		GasUsed:          fmt.Sprintf("%#x", tx.GasUsed),
		GasLimit:         fmt.Sprintf("%#x", tx.GasLimit),
		GasPrice:         fmt.Sprintf("%#x", tx.GasPrice),
		Status:           tx.Status,
		Nonce:            fmt.Sprintf("%#x", tx.Nonce),
		ERC20Transfers:   ExtractERC20Transfers(tx.Logs),
		RecipientAddress: tx.RecipientAddress,
		Hash:             tx.Hash,
		TxType:           txTypeDerived,
		Data:             tx.Data,
	}, nil
}

func handleBaseKeys(tx types.PreDeterminedStandardTx, key string, referenceList []string) bool {
	switch key {
	case "to":
		return handleTo(tx, referenceList)
	case "from":
		return handleFrom(tx, referenceList)
	case "status":
		return handleStatus(tx, referenceList)
	case "topics":
		return handleLogTopics(tx, referenceList)
	case "logAddress":
		return handleLogAddress(tx, referenceList)
	}
	return false
}

func handleTo(tx types.PreDeterminedStandardTx, addresses []string) bool {
	return contains(addresses, tx.To)
}

func handleFrom(tx types.PreDeterminedStandardTx, addresses []string) bool {
	return contains(addresses, tx.From)
}

func handleLogTopics(tx types.PreDeterminedStandardTx, topics []string) bool {
	for _, log := range tx.Logs {
		for _, logTopic := range log.Topics {
			if contains(topics, logTopic.String()) {
				return true
			}
		}
	}
	return false
}

func handleStatus(tx types.PreDeterminedStandardTx, statuses []string) bool {
	return types.TxStatus(statuses[0]) == tx.Status
}

func handleLogAddress(tx types.PreDeterminedStandardTx, logAddresses []string) bool {
	for _, log := range tx.Logs {
		if contains(logAddresses, log.Address.String()) {
			return true
		}
	}
	return false
}

func otherRequirementsPresent(applied map[string]types.FullTxLabelSchema, requiredList []string) bool {
	for _, required := range requiredList {
		if _, ok := applied[required]; !ok {
			return false
		}
	}
	return true
}

func interpretAppliedSchemas(applied map[string]types.FullTxLabelSchema) string {
	var txType string
	priority := 0
	if _, ok := applied["STANDARD"]; ok {
		return "STANDARD"
	} else if _, ok := applied["GENERIC_CONTRACT_CALL"]; ok {
		return "GENERIC_CONTRACT_CALL"
	}

	for _, item := range applied {
		if item.Meta.Priority > priority {
			txType = item.Meta.Name
			priority = item.Meta.Priority
		}
	}
	return txType
}

func sortSchemaListByPriority(schemas []types.FullTxLabelSchema) []types.FullTxLabelSchema {
	sort.Slice(schemas[:], func(i, j int) bool {
		return schemas[i].Meta.Priority < schemas[j].Meta.Priority
	})
	return schemas
}

func ExtractERC20Transfers(logs []*ethertypes.Log) []types.ERC20Transfer {
	transfers := []types.ERC20Transfer{}
	for _, log := range logs {
		if log.Topics[0] == common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef") {
			transfers = append(transfers, types.ERC20Transfer{
				ContractAddress: log.Address,
				From:            convertHashToString(log.Topics[1]),
				To:              convertHashToString(log.Topics[2]),
				Amount:          convertFromHexByteToHexString(log.Data),
			})
		}
	}
	return transfers
}

func convertHashToString(hash common.Hash) string {
	trimmedHash := common.TrimLeftZeroes(hash.Bytes())
	return common.BytesToAddress(trimmedHash).String()
}

func convertFromHexByteToHexString(amount []byte) string {
	trimmedAmt := common.TrimLeftZeroes(amount)
	enc := make([]byte, len(trimmedAmt)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], trimmedAmt)
	return string(enc)
}
