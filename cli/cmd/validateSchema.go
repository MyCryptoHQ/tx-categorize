package cmd

import (
	"log"

	"github.com/mycryptohq/tx-categorize/cmd/categorize"
	"github.com/spf13/cobra"
	"gopkg.in/go-playground/validator.v9"
)

// validateSchemaCmd represents the validate-schema command
var validateSchemaCmd = &cobra.Command{
	Use:   "validate-schema",
	Short: "Fetches and validates schema files",
	Run: func(cmd *cobra.Command, args []string) {
		validate := validator.New()
		schemas, _ := categorize.FetchAndWalkSchema("./../schema/")
		var errOut error
		for _, schema := range schemas {
			err := validate.Struct(schema)
			if err != nil {
				errOut = err
				log.Printf("Schema %s has invalid type: %s\n", schema.Meta.Name, err.Error())
			}
		}
		if errOut != nil {
			log.Fatal("Error validating schemas.")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateSchemaCmd)
}
