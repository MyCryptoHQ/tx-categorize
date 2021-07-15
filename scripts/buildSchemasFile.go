package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/mycryptohq/tx-categorize/cmd/categorize"
	"github.com/mycryptohq/tx-categorize/types"
)

var (
	tempSchemasFilePath = "/tmp/schemas.json"
)

type SchemaMeta struct {
	Protocol       types.Protocol       `json:"protocol" validate:"required"`
	ProtocolAction types.ProtocolAction `json:"type" validate:"required"`
}

type SchemaMetaMap map[string]SchemaMeta

func main() {
	schemas, err := categorize.FetchAndWalkSchema("./schema/")
	if err != nil {
		fmt.Println("[buildSchemasFile]: error fetching schema ", err)
		return
	}
	schemaMetaMap := make(map[string]SchemaMeta)
	for _, fullSchema := range schemas {
		schemaMetaMap[fullSchema.Meta.Name] = SchemaMeta{
			Protocol: fullSchema.Meta.Protocol,
			ProtocolAction: fullSchema.Meta.ProtocolAction,
		}
	}
	schemasFile,_ := json.MarshalIndent(schemaMetaMap, " ", "   ")
	err = ioutil.WriteFile(tempSchemasFilePath, schemasFile, 0644)
	if err != nil {
		fmt.Println("[buildSchemasFile]: can't write schema file ", err)
	}
}