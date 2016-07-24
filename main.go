package main

import (
	"fmt"
	"github.com/jkrecek/recbuilder/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	ARG_REGEX_STR        = "^--([a-z-]*)(?:=(.*))?$"
	FROM_TABLE_REGEX_STR = "(^|_)([a-z])"
)

var (
	argumentRegex  = regexp.MustCompile(ARG_REGEX_STR)
	fromTableRegex = regexp.MustCompile(FROM_TABLE_REGEX_STR)
	variableMap    = make(map[string]string)
	sourceFileMap  = map[string]string{
		"{{REPOSITORY_NAME}}.php":       template.REPOSITORY_CODE,
		"{{COLLECTION_NAME}}.php":       template.COLLECTION_CODE,
		"{{ENTITY_NAME}}.php":           template.ENTITY_CODE,
		"{{ENTITY_STRUCTURE_NAME}}.php": template.ENTITY_STRUCTURE_CODE,
	}
)

func main() {
	tableName := getExecuteArgumentValue("table")
	root := getExecuteArgumentValue("root")

	if len(tableName) == 0 || len(root) == 0 {
		log.Println("You must specify --table and --root parameter.")
		os.Exit(1)
	}

	setVariables(tableName, root)

	outputFileMap := make(map[string]string)
	for fileName, fileContent := range sourceFileMap {
		outputFileMap[customize(fileName)] = customize(fileContent)
	}

	os.Mkdir(tableName, 0755)
	for outFile, outContent := range outputFileMap {
		ioutil.WriteFile(tableName+"/"+outFile, []byte(outContent), 0644)
	}
}

func setVariables(tableName, root string) {
	entityBase := fromTableRegex.ReplaceAllStringFunc(tableName, func(matches string) string {
		lastChar := matches[len(matches)-1]
		return strings.ToUpper(string(lastChar))
	})

	readable := strings.Replace(tableName, "_", " ", -1)
	readable = strings.ToUpper(string(readable[0])) + readable[1:]

	variableMap["NAMESPACE"] = strings.Trim(root, "\\") + "\\" + entityBase
	variableMap["READABLE_NAME"] = readable

	variableMap["TABLE_NAME"] = tableName
	variableMap["REPOSITORY_NAME"] = entityBase + "Repository"
	variableMap["COLLECTION_NAME"] = entityBase + "Collection"
	variableMap["ENTITY_NAME"] = entityBase + "Entity"
	variableMap["ENTITY_STRUCTURE_NAME"] = entityBase + "EntityStructure"
}

func getExecuteArgumentValue(argument string) string {
	for _, osArg := range os.Args[1:] {
		subMatch := argumentRegex.FindStringSubmatch(osArg)
		if len(subMatch) == 0 {
			continue
		}

		if argument == subMatch[1] {
			return subMatch[2]
		}
	}

	return ""
}

func customize(str string) string {
	for old, new := range variableMap {
		str = strings.Replace(str, fmt.Sprintf("{{%s}}", old), new, -1)
	}

	return str
}
