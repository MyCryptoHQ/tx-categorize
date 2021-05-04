package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mycryptohq/tx-categorize/cmd/categorize"
	"github.com/spf13/cobra"
)

// generateTemplateCmd represents the generate-template command
var listUniqueTypesCmd = &cobra.Command{
	Use:   "list-unique-types",
	Short: "Lists unique types and protocols",
	Run: func(cmd *cobra.Command, args []string) {
		schemas, _ := categorize.FetchAndWalkSchema("./../schema/")
		typesMap := make(map[string]int)
		protocolsMap := make(map[string]int)
		//var types []string
		//var protocols []string
		for _, schema := range schemas {
			//if string(schema.Meta.ProtocolAction) != "" {
			typesMap[string(schema.Meta.ProtocolAction)] += 1
			//}

			//if string(schema.Meta.Protocol) != "" {
			protocolsMap[string(schema.Meta.Protocol)] += 1
			//}
		}
		// for protocol, _ := range protocolsMap {
		// 	protocols = append(protocols, protocol)
		// }

		// for t, _ := range typesMap {
		// 	types = append(types, t)
		// }
		//sort.Strings(protocolsMap)
		protocolsJson, _ := json.MarshalIndent(protocolsMap, "", "  ")
		fmt.Println("Protocols & instances:", string(protocolsJson))

		//sort.Strings(types)
		typesJson, _ := json.MarshalIndent(typesMap, "", "  ")
		fmt.Println("Types & instances:", string(typesJson))
	},
}

func init() {
	rootCmd.AddCommand(listUniqueTypesCmd)
}
