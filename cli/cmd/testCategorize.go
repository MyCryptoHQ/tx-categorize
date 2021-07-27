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

// categorizeCmd represents the categorize command
var categorizeCmd = &cobra.Command{
	Use:   "categorize",
	Short: "Categorizes a tx using txhash as parameter",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpc")
		if rpcUrl == "" {
			log.Fatal(fmt.Errorf("rpcUrl is required but not found"))
		}
		fmt.Println("[rpcUrl]: ", rpcUrl)
		txHash, _ := cmd.Flags().GetString("txHash")
		if txHash == "" {
			log.Fatal(fmt.Errorf("txHash is required but not found"))
		}
		fmt.Println("[txHash]: ", txHash)
		client := etherclient.MakeETHClient(rpcUrl)
		txHashes := []string{txHash}
		txConstructions := build.FetchTxReceipts(txHashes, *client)
		formattedNormalTxs := categorize.FormatNormalTxsToStandard(txConstructions)
		schemas, err := categorize.FetchAndWalkSchema()
		if err != nil {
			fmt.Println("Error fetching schemas: ", err)
		}
		var txs []types.StandardTx
		for _, txConstruction := range formattedNormalTxs {
			derivedTx, err := categorize.DetermineTxType(txConstruction, schemas)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to format tx %s", txConstruction.Hash))
			}
			fmt.Printf("%s categorized as: %s\n", txConstruction.Hash, derivedTx.TxType)
			txs = append(txs, derivedTx)
		}

		file, _ := json.MarshalIndent(txs, "", "  ")
		//fileLoc := fmt.Sprintf(, string(file))
		_ = ioutil.WriteFile("./testfile.json", file, 0644)
		fmt.Println("A empty schema template file has been generated at testfile.json")
	},
}

func init() {
	rootCmd.AddCommand(categorizeCmd)
	categorizeCmd.Flags().StringP("rpc", "r", "", "Ethereum node rpc url")
	categorizeCmd.Flags().StringP("schemaId", "k", "", "The unique schema key that this test is for")
	categorizeCmd.Flags().StringP("txHash", "t", "", "The txHash to generate the test with")
}
