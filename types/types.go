package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type NormalTx struct {
	BlockNumber       int       `json:"blockNumber,string"`
	TimeStamp         time.Time `json:"timeStamp"`
	Hash              string    `json:"hash"`
	Nonce             int       `json:"nonce,string"`
	BlockHash         string    `json:"blockHash"`
	TransactionIndex  int       `json:"transactionIndex,string"`
	From              string    `json:"from"`
	To                string    `json:"to"`
	Value             *big.Int  `json:"value"`
	Gas               int       `json:"gas,string"`
	GasPrice          *big.Int  `json:"gasPrice"`
	IsError           int       `json:"isError,string"`
	TxReceiptStatus   string    `json:"txreceipt_status"`
	Input             string    `json:"input"`
	ContractAddress   string    `json:"contractAddress"`
	CumulativeGasUsed int       `json:"cumulativeGasUsed,string"`
	GasUsed           int       `json:"gasUsed,string"`
	Confirmations     int       `json:"confirmations,string"`
}

type ParsedStandardTx struct {
	GasUsed uint64       `json:"gasUsed"`
	Tx      NormalTx     `json:"tx"`
	Logs    []*types.Log `json:"logs"`
}

type PreDeterminedStandardTx struct {
	To               string       `json:"to"`
	From             string       `json:"from"`
	Value            *big.Int     `json:"value"`
	BlockNumber      int          `json:"blockNumber"`
	TimeStamp        int64        `json:"timestamp"`
	ContractAddress  string       `json:"contractAddress"`
	GasUsed          int          `json:"gasUsed"`
	GasLimit         int          `json:"gasLimit"`
	GasPrice         *big.Int     `json:"gasPrice"`
	Status           TxStatus     `json:"status"`
	Nonce            int          `json:"nonce"`
	Hash             string       `json:"hash"`
	Logs             []*types.Log `json:"logs"`
	RecipientAddress string       `json:"recipientAddress"`
	Data             string       `json:"data"`
}

type StandardTx struct {
	To               string          `json:"to"`
	From             string          `json:"from"`
	Value            string          `json:"value"`
	BlockNumber      string          `json:"blockNumber"`
	TimeStamp        int64           `json:"timestamp"`
	GasLimit         string          `json:"gasLimit"`
	GasUsed          string          `json:"gasUsed"`
	GasPrice         string          `json:"gasPrice"`
	Status           TxStatus        `json:"status"`
	Nonce            string          `json:"nonce"`
	ERC20Transfers   []ERC20Transfer `json:"erc20Transfers"`
	RecipientAddress string          `json:"recipientAddress"`
	Hash             string          `json:"hash"`
	TxType           string          `json:"txType"`
	Data             string          `json:"data"`
}

type ERC20Transfer struct {
	From            string         `json:"from"`
	To              string         `json:"to"`
	ContractAddress common.Address `json:"contractAddress"`
	Amount          string         `json:"amount"`
}

type TxLabelMeta struct {
	Name           string         `json:"name"`
	Platform       Platform       `json:"platform"`
	PlatformAction PlatformAction `json:"platformAction"`
	Priority       int            `json:"priority"`
}

type TxLabelSchema struct {
	Key       string   `json:"key"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
	Topics    []string `json:"topics"`
	And       []string `json:"and"`
	Statuses  []string `json:"statuses"`
}

type FullTxLabelSchema struct {
	Schema  TxLabelSchema `json:"schema"`
	Meta    TxLabelMeta   `json:"meta"`
	Version int           `json:"version"`
}
