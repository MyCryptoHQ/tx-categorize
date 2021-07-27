package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mycryptohq/tx-categorize/cmd/categorize"
	"github.com/mycryptohq/tx-categorize/types"
)

var (
	tempSchemasFilePath     = "./tmp/schemas.json"
	tempFullSchemasFilePath = "./tmp/fullschemas.json"
)

type SchemaMeta struct {
	Protocol       types.Protocol       `json:"protocol" validate:"required"`
	ProtocolAction types.ProtocolAction `json:"type" validate:"required"`
}

type SchemaMetaMap map[string]SchemaMeta

func main() {
	schemas, err := categorize.FetchAndWalkSchema(true)
	if err != nil {
		fmt.Println("[buildSchemasFile]: error fetching schema", err)
		return
	}
	var fullSchemas []types.FullTxLabelSchema
	schemaMetaMap := make(map[string]SchemaMeta)
	for _, fullSchema := range schemas {
		if fullSchema.Meta.ExcludeFromBuild {
			continue
		}
		schemaMetaMap[fullSchema.Meta.Name] = SchemaMeta{
			Protocol:       fullSchema.Meta.Protocol,
			ProtocolAction: fullSchema.Meta.ProtocolAction,
		}
		fullSchemas = append(fullSchemas, fullSchema)
	}
	schemasFile, _ := json.MarshalIndent(schemaMetaMap, " ", "   ")
	err = ioutil.WriteFile(tempSchemasFilePath, schemasFile, 0644)
	if err != nil {
		fmt.Println("[buildSchemasFile]: can't write schema file", err)
	}
	fullSchemasFile, _ := json.MarshalIndent(fullSchemas, " ", "   ")
	err = ioutil.WriteFile(tempFullSchemasFilePath, fullSchemasFile, 0644)
	if err != nil {
		fmt.Println("[buildSchemasFile]: can't write full schemas file", err)
	}
}
