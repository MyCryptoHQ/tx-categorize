package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	testsDir = "../tests"
)

// generateTemplateCmd represents the generate-template command
var generateTestTemplateCmd = &cobra.Command{
	Use:   "generate-test-template",
	Short: "Generates a test template based on an input tx hash",
	Run: func(cmd *cobra.Command, args []string) {
		schemaId, _ := cmd.Flags().GetString("schemaId")
		if schemaId == "" {
			log.Fatal(fmt.Errorf("schemaId is required but not found"))
		}
		txHash, _ := cmd.Flags().GetString("txHash")
		if txHash == "" {
			log.Fatal(fmt.Errorf("txHash is required but not found"))
		}

		exampleTestFileLoc := testsDir + "/example_test.txt"
		exampleTestReader, err := ioutil.ReadFile(exampleTestFileLoc)
		if err != nil {
			log.Fatal(fmt.Errorf("could not read example_test.txt file at %s", exampleTestFileLoc))
		}

		finishedTestFile := setupTestFile(string(exampleTestReader), schemaId, txHash)

		fileLoc := fmt.Sprintf(testsDir+"/%s_test.go", schemaId)
		_ = ioutil.WriteFile(fileLoc, []byte(finishedTestFile), 0644)
		fmt.Println("A empty schema test file has been generated at", fileLoc)
	},
}

func init() {
	rootCmd.AddCommand(generateTestTemplateCmd)
	generateTestTemplateCmd.Flags().StringP("rpc", "r", "", "Ethereum node rpc url")
	generateTestTemplateCmd.Flags().StringP("schemaId", "k", "", "The unique schema key that this test is for")
	generateTestTemplateCmd.Flags().StringP("txHash", "t", "", "The txHash to generate the test with")
}

func setupTestFile(input string, key string, txHash string) string {
	replaceKey := "%expected"
	replaceTxHash := "%txHashes"
	var replaceKeyReg = regexp.MustCompile(replaceKey)
	var replaceTxHashReg = regexp.MustCompile(replaceTxHash)
	rep := replaceKeyReg.ReplaceAllString(input, key)
	out := replaceTxHashReg.ReplaceAllString(rep, txHash)
	return out
}
