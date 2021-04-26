package etherclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TokenObject struct {
	Contract common.Address
	Wallet   common.Address
	Name     string
	Symbol   string
	Balance  *big.Int
	ETH      *big.Int
	Decimals int64
	Block    int64
	ctx      context.Context
}

// TokenTransfer represents a Transfer event raised by the Token contract.
type TokenTransferProcessedObj struct {
	Symbol          string         `json:"symbol,omitempty"`
	From            common.Address `json:"from,omitempty"`
	To              common.Address `json:"to,omitempty"`
	ContractAddress common.Address `json:"contractAddress,omitempty"`
	Amount          float64        `json:"amount,omitempty"`
	FiatAmount      float64        `json:"fiatAmount,omitempty"`
}
