package etherclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ClientTypes "github.com/ethereum/go-ethereum/core/types"
)

func FetchBlockMinGas(txs ClientTypes.Transactions) int64 {
	var minGasPrice int64

	for _, tx := range txs {
		if tx.GasPrice().Int64() < minGasPrice || minGasPrice == 0 {
			minGasPrice = tx.GasPrice().Int64()
		}
	}
	if minGasPrice == 0 {
		minGasPrice = 100000000
	}

	return minGasPrice
}

func ParseTransferLog(logEvent types.Log, decimal int, symbol string, rate float64) (TokenTransferProcessedObj, error) {
	val := new(big.Int)
	val.SetBytes(logEvent.Data).Uint64()
	userReadableAmount := ConvertFromBase(*val, decimal)
	return TokenTransferProcessedObj{
		ContractAddress: common.BytesToAddress(logEvent.Address.Bytes()),
		From:            common.BytesToAddress(logEvent.Topics[1].Bytes()),
		To:              common.BytesToAddress(logEvent.Topics[2].Bytes()),
		Amount:          userReadableAmount,
		Symbol:          symbol,
		FiatAmount:      rate * userReadableAmount,
	}, nil
}
