package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mycryptohq/tx-categorize/cmd/categorize"
	"github.com/spf13/cobra"
)

// listUniqueTypesCmd represents the list-unique-types command
var listUniqueTypesCmd = &cobra.Command{
	Use:   "list-unique-types",
	Short: "Lists unique types and protocols",
	Run: func(cmd *cobra.Command, args []string) {
		schemas, _ := categorize.FetchAndWalkSchema("./../schema/")
		typesMap := make(map[string]int)
		protocolsMap := make(map[string]int)
		for _, schema := range schemas {
			if string(schema.Meta.ProtocolAction) != "" {
				typesMap[string(schema.Meta.ProtocolAction)] += 1
			} else {
				fmt.Println("Missing type in:", schema.Meta.Name)
			}

			if string(schema.Meta.Protocol) != "" {
				protocolsMap[string(schema.Meta.Protocol)] += 1
			} else {
				fmt.Println("Missing protocol in:", schema.Meta.Name)
			}
		}
		protocolsJson, _ := json.MarshalIndent(protocolsMap, "", "  ")
		fmt.Println("Protocols & instances:", string(protocolsJson))

		typesJson, _ := json.MarshalIndent(typesMap, "", "  ")
		fmt.Println("Types & instances:", string(typesJson))
	},
}

func init() {
	rootCmd.AddCommand(listUniqueTypesCmd)
}
