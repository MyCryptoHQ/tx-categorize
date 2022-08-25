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
	uniqueSchemaActions = "./tmp/uniqueSchemaActions.json"
)

type SchemaMeta struct {
	Protocol       types.Protocol       `json:"protocol" validate:"required"`
	ProtocolAction types.ProtocolAction `json:"type" validate:"required"`
}

type UniqueSchema struct {
	Actions	map[string]int `json:"actions" validate:"required"`
	Protocols map[string]int `json:"protocols" validate:"required"`
}

type SchemaMetaMap map[string]SchemaMeta

func main() {
	schemas, err := categorize.FetchAndWalkSchema(false)
	if err != nil {
		fmt.Println("[buildSchemasFile]: error fetching schema", err)
		return
	}
	var fullSchemas []types.FullTxLabelSchema
	schemaMetaMap := make(map[string]SchemaMeta)
	actionsMap := make(map[string]int)
	protocolsMap := make(map[string]int)
	for _, fullSchema := range schemas {
		if fullSchema.Meta.ExcludeFromBuild {
			continue
		}
		if string(fullSchema.Meta.ProtocolAction) != "" {
			actionsMap[string(fullSchema.Meta.ProtocolAction)] += 1
		} else {
			fmt.Println("Missing type in:", fullSchema.Meta.Name)
		}

		if string(fullSchema.Meta.Protocol) != "" {
			protocolsMap[string(fullSchema.Meta.Protocol)] += 1
		} else {
			fmt.Println("Missing protocol in:", fullSchema.Meta.Name)
		}
		schemaMetaMap[fullSchema.Meta.Name] = SchemaMeta{
			Protocol:       fullSchema.Meta.Protocol,
			ProtocolAction: fullSchema.Meta.ProtocolAction,
		}
		fullSchemas = append(fullSchemas, fullSchema)
	}
	uniqueLists := UniqueSchema{
		Actions: actionsMap,
		Protocols: protocolsMap,
	}
	schemasFile, _ := json.MarshalIndent(schemaMetaMap, " ", "   ")
	err = ioutil.WriteFile(tempSchemasFilePath, schemasFile, 0644)
	if err != nil {
		fmt.Println("[buildSchemasFile]: can't write schema file", err)
	}
	fullSchemasFile, _ := json.MarshalIndent(schemas, " ", "   ")
	err = ioutil.WriteFile(tempFullSchemasFilePath, fullSchemasFile, 0644)
	if err != nil {
		fmt.Println("[buildSchemasFile]: can't write full schemas file", err)
	}
	schemasActionsFile, _ := json.MarshalIndent(uniqueLists, " ", "   ")
	err = ioutil.WriteFile(uniqueSchemaActions, schemasActionsFile, 0644)
	if err != nil {
		fmt.Println("[buildSchemasFile]: can't write full schemas file", err)
	}
}
