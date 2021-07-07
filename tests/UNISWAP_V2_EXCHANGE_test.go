package schema

import (
	"fmt"
	"os"
	"testing"

	"github.com/mycryptohq/tx-categorize/cmd/build"
	"github.com/mycryptohq/tx-categorize/cmd/categorize"
	"github.com/mycryptohq/tx-categorize/common/etherclient"
)

func Test_UNISWAP_V2_EXCHANGE(t *testing.T) { // @todo: replace with replacement test name
	var expectedKey = "UNISWAP_V2_EXCHANGE"
	rpcUrl := os.Getenv("rpcUrl")
	if rpcUrl == "" {
		fmt.Printf("[Key_Derivation_Test_%s]: env variable 'rpcUrl' is required but not found\n", expectedKey)
		t.Errorf("[Key_Derivation_Test_%s]: env variable 'rpcUrl' is required but not found\n", expectedKey)
	}

	txHashes := []string{"0xfce920c5fbd9a51f0988448ac9c5224956f57a40f46a53edd8a6a4f0697a7ac0"} // @todo: replace with replacement txhashes

	for _, txHash := range txHashes {
		client := etherclient.MakeETHClient(rpcUrl)
		txHashes := []string{txHash}
		txConstructions := build.FetchTxReceipts(txHashes, *client)
		if len(txConstructions) == 0 {
			t.Errorf("[Key_Derivation_Test_%s]: Failed to fetch tx from node %s", expectedKey, txHash)
		}
		txConstruction := categorize.FormatNormalTxsToStandard(txConstructions)[0]
		schemas, _ := categorize.FetchAndWalkSchema("../schema/")
		derivedTx, err := categorize.DetermineTxType(txConstruction, schemas)
		if err != nil {
			t.Errorf("[Key_Derivation_Test_%s]: Failed to format tx %s", expectedKey, txConstruction.Hash)
		}

		if derivedTx.TxType != expectedKey {
			t.Errorf("[Key_Derivation_Test_%s]: Test Success Status: %t\tInput: %v\tOutput: %v\tExpectedOutput: %v\n", expectedKey, derivedTx.TxType == expectedKey, txHash, derivedTx.TxType, expectedKey)
		}
	}
	fmt.Printf("\n\n\n")
}
