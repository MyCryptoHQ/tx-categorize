package etherclient

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetBalance(client ethclient.Client, address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}

func GetBlockByNum(client ethclient.Client, blockNum int64) (types.Block, error) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNum))
	if err != nil {
		return types.Block{}, err
	}
	return *block, nil
}

func GetBlockNumber(client ethclient.Client) (int64, error) {
	blockNumHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return int64(0), err
	}
	return blockNumHeader.Number.Int64(), nil
}

func Call(client ethclient.Client, callTx ethereum.CallMsg) ([]byte, error) {
	ethCallResult, err := client.CallContract(context.Background(), callTx, nil)
	if err != nil {
		return []byte(""), err
	}
	return ethCallResult, nil
}

func GetTransactionReceipt(client ethclient.Client, txHash string) (*types.Receipt, error) {
	txReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	return txReceipt, nil
}

func GetTransactionByHash(client ethclient.Client, txHash string) (*types.Transaction, bool, error) {
	txObj, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, isPending, err
	}
	return txObj, isPending, nil
}

func GetLogs(client ethclient.Client, logsFilter ethereum.FilterQuery) ([]types.Log, error) {
	logsArr, err := client.FilterLogs(context.Background(), logsFilter)
	if err != nil {
		return []types.Log{}, err
	}
	return logsArr, nil
}
