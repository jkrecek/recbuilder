package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type table struct {
	tableName          string
	rootNameSpace      string
	db                 *database
	transformVariables map[string]string
}

func (t *table) process() {
	t.setVariables(t.tableName, t.rootNameSpace)

	outputFileMap := make(map[string]string)
	for fileName, fileContent := range sourceFileMap {
		outputFileMap[t.customize(fileName)] = t.customize(fileContent)
	}

	_, err := os.Stat(t.tableName)
	if err == nil {
		fmt.Printf("Entity structure for `%s` already exists. Skipping...\n", t.tableName)
		return
	}

	os.Mkdir(t.tableName, 0755)
	for outFile, outContent := range outputFileMap {
		err := ioutil.WriteFile(t.tableName+"/"+outFile, []byte(outContent), 0644)
		if err != nil {
			log.Println(err)
			return
		}
	}

	fmt.Printf("Entity structure for `%s` successfully created.\n", t.tableName)
}

func (t *table) setVariables(tableName, root string) {
	entityBase := underscoreSepRegex.ReplaceAllStringFunc(tableName, func(matches string) string {
		return strings.ToUpper(string(matches[1]))
	})
	entityBase = strings.ToUpper(string(entityBase[0])) + entityBase[1:]

	readable := strings.Replace(tableName, "_", " ", -1)
	readable = strings.ToUpper(string(readable[0])) + readable[1:]

	t.transformVariables = make(map[string]string)
	t.transformVariables["NAMESPACE"] = strings.Trim(root, "\\") + "\\" + entityBase
	t.transformVariables["READABLE_NAME"] = readable

	t.transformVariables["TABLE_NAME"] = tableName
	t.transformVariables["REPOSITORY_NAME"] = entityBase + "Repository"
	t.transformVariables["COLLECTION_NAME"] = entityBase + "Collection"
	t.transformVariables["ENTITY_NAME"] = entityBase + "Entity"
	t.transformVariables["ENTITY_STRUCTURE_NAME"] = entityBase + "EntityStructure"

	t.transformVariables["ENTITY_VALUES"] = t.getEntityValues()
}

func (t *table) getEntityValues() string {
	if t.db == nil {
		return ""
	}

	columns := t.db.loadTableStructure(t.tableName)

	var buffer bytes.Buffer
	for _, column := range columns {
		buffer.WriteString(column.toOutputString())
	}

	return buffer.String()
}

func (t *table) customize(str string) string {
	for old, new := range t.transformVariables {
		str = strings.Replace(str, fmt.Sprintf("{{%s}}", old), new, -1)
	}

	return str
}
