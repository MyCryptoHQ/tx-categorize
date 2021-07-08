package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mycryptohq/tx-categorize/cmd/build"
	"github.com/mycryptohq/tx-categorize/cmd/categorize"
	"github.com/mycryptohq/tx-categorize/common/etherclient"
)

func Test_Schemas(t *testing.T) { // @todo: replace with replacement test name
	rpcUrl := os.Getenv("rpcUrl")
	if rpcUrl == "" {
		fmt.Printf("[Test_Schemas]: env variable 'rpcUrl' is required but not found\n")
		t.Errorf("[Test_Schemas]: env variable 'rpcUrl' is required but not found\n")
	}
	testListByte, err := ioutil.ReadFile("./testList.json")
	if err != nil {
		t.Errorf("[Test_Schemas]: Failed to fetch testList file")
	}
	type TestList map[string]string
	var testList TestList

	err = json.Unmarshal(testListByte, &testList)
	if err != nil {
		t.Errorf("[Test_Schemas]: Failed to unmarshal testlist")
	}
	client := etherclient.MakeETHClient(rpcUrl)
	schemas, _ := categorize.FetchAndWalkSchema("../schema/")
	if len(testList) != len(schemas) {
		fmt.Printf("[Test_Schemas]: There are %v untested schemas\n", len(schemas)-len(testList))
	}
	for schemaId, txHash := range testList {
		txHashes := []string{txHash}
		txConstructions := build.FetchTxReceipts(txHashes, *client)
		if len(txConstructions) == 0 {
			t.Errorf("[Key_Derivation_Test_%s]: Failed to fetch tx from node %s", schemaId, txHash)
		}
		txConstruction := categorize.FormatNormalTxsToStandard(txConstructions)[0]

		derivedTx, err := categorize.DetermineTxType(txConstruction, schemas)
		if err != nil {
			t.Errorf("[Key_Derivation_Test_%s]: Failed to format tx %s", schemaId, txConstruction.Hash)
		}

		if derivedTx.TxType != schemaId {
			t.Errorf("[Key_Derivation_Test_%s]: Test Success Status: %t\tInput: %v\tOutput: %v\tExpectedOutput: %v\n", schemaId, derivedTx.TxType == schemaId, txHash, derivedTx.TxType, schemaId)
		}
	}
	fmt.Printf("\n\n\n")
}
