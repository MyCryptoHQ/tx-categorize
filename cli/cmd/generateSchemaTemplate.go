package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mycryptohq/tx-categorize/types"
	"github.com/spf13/cobra"
)

// generateTemplateCmd represents the generate-template command
var generateSchemaTemplateCmd = &cobra.Command{
	Use:   "generate-template",
	Short: "Generates a log topics-based schema template",
	Run: func(cmd *cobra.Command, args []string) {
		emptySchema := types.FullTxLabelSchema{}
		file, _ := json.MarshalIndent(emptySchema, "", "  ")
		_ = ioutil.WriteFile("schema_template.json", file, 0644)
		fmt.Println("A empty schema template file has been generated at ./template_schema.json")
	},
}

func init() {
	rootCmd.AddCommand(generateSchemaTemplateCmd)
}
