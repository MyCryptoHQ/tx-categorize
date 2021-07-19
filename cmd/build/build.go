package build

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mycryptohq/tx-categorize/common/etherclient"
	"github.com/mycryptohq/tx-categorize/types"
)

func FetchTxReceipts(txHashes []string, ethclient ethclient.Client) []types.ParsedStandardTx {
	// parallelize the tx receipt api reqs
	parsedStandardTxs := make([]types.ParsedStandardTx, 0, len(txHashes))
	ch := make(chan types.ParsedStandardTx, len(txHashes))
	wg := sync.WaitGroup{}
	for _, txHash := range txHashes {
		wg.Add(1)
		go fetchTxReceipt(txHash, ethclient, ch, &wg)
	}

	// now we wait for everyone to finish - again, not a must.
	// you can just receive from the channel N times, and use a timeout or something for safety
	wg.Wait()
	// we need to close the channel or the following loop will get stuck
	close(ch)
	for parsedStandardTx := range ch {
		parsedStandardTxs = append(parsedStandardTxs, parsedStandardTx)
	}
	return parsedStandardTxs
}

func fetchTxReceipt(txHash string, ethclient ethclient.Client, ch chan types.ParsedStandardTx, wg *sync.WaitGroup) {
	blockNumber, err := ethclient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("[fetchTxReceipt]: Couldn't fetch block number from node ", txHash, ". Err: ", err)
		wg.Done()
	}
	txReceipt, err := ethclient.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		log.Println("[fetchTxReceipt]: Couldn't fetch txReceipt from node ", txHash, ". Err: ", err)
		wg.Done()
	}
	tx, _, err := ethclient.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		log.Println("[fetchTxReceipt]: Couldn't fetch transaction from node ", txHash, ". Err: ", err)
		wg.Done()
	}
	block, err := ethclient.BlockByHash(context.Background(), txReceipt.BlockHash)
	if err != nil {
		log.Println("[fetchTxReceipt]: Couldn't fetch block from node ", txReceipt.BlockHash, ". Err: ", err)
		wg.Done()
	}
	sender, err := ethclient.TransactionSender(context.Background(), tx, txReceipt.BlockHash, txReceipt.TransactionIndex)
	if err != nil {
		log.Println("[fetchTxReceipt]: Couldn't fetch sender from node ", txReceipt.BlockHash, ". Err: ", err)
		wg.Done()
	}
	blockTime := block.Time()
	// push the standard tx object down the channel
	ch <- types.ParsedStandardTx{
		Logs:    txReceipt.Logs,
		GasUsed: txReceipt.GasUsed,
		Tx: types.NormalTx{
			BlockNumber:       int(txReceipt.BlockNumber.Int64()),
			TimeStamp:         time.Unix(int64(blockTime), 0),
			Hash:              txReceipt.TxHash.String(),
			Nonce:             int(tx.Nonce()),
			BlockHash:         txReceipt.BlockHash.String(),
			TransactionIndex:  int(txReceipt.TransactionIndex),
			From:              sender.String(),
			To:                tx.To().String(),
			Value:             tx.Value(),
			Gas:               int(tx.Gas()),
			GasPrice:          tx.GasPrice(),
			IsError:           int(txReceipt.Status),
			TxReceiptStatus:   fmt.Sprintf("%v", txReceipt.Status), //deriveTxStatus(isPending, txReceipt.Status),
			Input:             etherclient.ConvertFromHexByteToHexString(tx.Data()),
			ContractAddress:   txReceipt.ContractAddress.String(), // contract address created by contract construction tx
			CumulativeGasUsed: int(txReceipt.CumulativeGasUsed),
			GasUsed:           int(txReceipt.GasUsed),
			Confirmations:     int(blockNumber.Number.Int64() - txReceipt.BlockNumber.Int64()),
		},
	}
	// let the wait group know we finished
	wg.Done()
}
