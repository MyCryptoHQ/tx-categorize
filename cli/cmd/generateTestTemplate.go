package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mycryptohq/tx-categorize/cmd/build"
	"github.com/mycryptohq/tx-categorize/cmd/categorize"
	"github.com/mycryptohq/tx-categorize/common/etherclient"
	"github.com/mycryptohq/tx-categorize/types"
	"github.com/spf13/cobra"
)

// generateTemplateCmd represents the generate-template command
var generateTestTemplateCmd = &cobra.Command{
	Use:   "generate-test-template",
	Short: "Generates a test template based on an input tx hash",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpc")
		if rpcUrl == "" {
			log.Fatal(fmt.Errorf("rpcUrl is required but not found"))
		}
		schemaId, _ := cmd.Flags().GetString("schemaId")
		if schemaId == "" {
			log.Fatal(fmt.Errorf("schemaId is required but not found"))
		}
		txHash, _ := cmd.Flags().GetString("txHash")
		if txHash == "" {
			log.Fatal(fmt.Errorf("txHash is required but not found"))
		}
		client := etherclient.MakeETHClient(rpcUrl)
		txHashes := []string{txHash}
		txConstructions := build.FetchTxReceipts(txHashes, *client)
		formattedNormalTxs := categorize.FormatNormalTxsToStandard(txConstructions)
		schemas, _ := categorize.FetchAndWalkSchema("./../schema/")
		var out []types.StandardTx
		for _, tx := range formattedNormalTxs {
			finalTx, err := categorize.DetermineTxType(tx, schemas)
			if err != nil {
				fmt.Println("[generate-test-template]: error - ", err)
			} else {
				out = append(out, finalTx)
			}
		}
		file, _ := json.MarshalIndent(out, "", "  ")
		fileLoc := fmt.Sprintf("../tests/test_txs/%s.json", schemaId)
		_ = ioutil.WriteFile(fileLoc, file, 0644)
		fmt.Println("A empty schema template file has been generated at", fileLoc)
	},
}

func init() {
	rootCmd.AddCommand(generateTestTemplateCmd)
	generateTestTemplateCmd.Flags().StringP("rpc", "r", "", "Ethereum node rpc url")
	generateTestTemplateCmd.Flags().StringP("schemaId", "k", "", "The unique schema key that this test is for")
	generateTestTemplateCmd.Flags().StringP("txHash", "t", "", "The txHash to generate the test with")
}
