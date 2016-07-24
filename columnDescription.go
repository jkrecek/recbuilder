package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type columnDescription struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

func (cd *columnDescription) getColumnName() string {
	return cd.Field
}

func (cd *columnDescription) getVariableName() string {
	return underscoreSepRegex.ReplaceAllStringFunc(cd.Field, func(matches string) string {
		return strings.ToUpper(string(matches[1]))
	})
}

func (cd *columnDescription) getPhpType() string {
	sqlBaseType := sqlValueTypeRegexStr.FindString(cd.Type)
	if strings.Contains(sqlBaseType, "int") {
		return "int"
	} else if strings.Contains(sqlBaseType, "char") || sqlBaseType == "text" {
		return "string"
	} else if sqlBaseType == "timestamp" || sqlBaseType == "datetime" {
		return "\\DateTime"
	} else {
		return "string"
	}
}

func (cd *columnDescription) isPrimary() bool {
	return cd.Key == "PRI"
}

func (cd *columnDescription) toOutputString() string {
	var phpDocParams []string
	phpDocParams = append(phpDocParams, fmt.Sprintf("@var %s", cd.getPhpType()))
	phpDocParams = append(phpDocParams, fmt.Sprintf("@Column(\"%s\")", cd.getColumnName()))
	if cd.isPrimary() {
		phpDocParams = append(phpDocParams, "@Primary")
	}

	var buffer bytes.Buffer
	indent := "\n    "
	buffer.WriteString(indent + "/**")

	for _, phpDocLine := range phpDocParams {
		buffer.WriteString(indent + " * " + phpDocLine)
	}

	buffer.WriteString(indent + " */")
	buffer.WriteString(indent + fmt.Sprintf("private $%s;\n", cd.getVariableName()))

	return buffer.String()
}
