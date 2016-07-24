package main

import (
	"fmt"
	"github.com/jkrecek/recbuilder/template"
	"log"
	"os"
	"regexp"

	_ "github.com/jkrecek/mysql"
)

const (
	ARG_REGEX_STR            = "^--([a-z-\\_]*)(?:=(.*))?$"
	UNDERSCORE_SEP_REGEX_STR = "_([a-z])"
	SQL_VALUE_TYPE_REGEX_STR = "^[a-z]+"
	HELP                     = `Tool for fast creating of repository+entity+collection structure.

Availible parameters:
	--help				Shows this table.
	--root=VALUE		Specifies in under which namespace should rec stack be created.
						Required.
						eg.: --root=App\Module\Repository
	--table=VALUE		Specifies for which table should structure be created.
						Required, unless specified --entire_database.
						eg.: --table=page_category
	--entire_database	Create structure for entire database.
						Required, unless specified --table=VALUE.
	--dsn				Database source name for connecting to SQL database. Creates Entity Structure when specified.
						eg.: --dsn=root:root@tcp(localhost:3306)/page_database

Examples:
	recbuilder --root=App\Model\Repository --dsn=root:usbw@tcp(localhost:3306)/zpmvcr --entire_database
						Creates structure for entire database.
	recbuilder --root=App\Model\Repository --table=page_category
						Creates structure only for table 'page_category'
	`
)

var (
	argumentRegex        = regexp.MustCompile(ARG_REGEX_STR)
	underscoreSepRegex   = regexp.MustCompile(UNDERSCORE_SEP_REGEX_STR)
	sqlValueTypeRegexStr = regexp.MustCompile(SQL_VALUE_TYPE_REGEX_STR)
	sourceFileMap        = map[string]string{
		"{{REPOSITORY_NAME}}.php":       template.REPOSITORY_CODE,
		"{{COLLECTION_NAME}}.php":       template.COLLECTION_CODE,
		"{{ENTITY_NAME}}.php":           template.ENTITY_CODE,
		"{{ENTITY_STRUCTURE_NAME}}.php": template.ENTITY_STRUCTURE_CODE,
	}
)

func main() {
	isHelp, _ := getExecuteArgumentValue("help")
	if isHelp {
		printHelp()
		return
	}

	rootFound, root := getExecuteArgumentValue("root")
	if !rootFound || len(root) == 0 {
		onError("You must specify --root parameter.")
	}

	_, tableName := getExecuteArgumentValue("table")

	_, dsn := getExecuteArgumentValue("dsn")
	entireDatabase, _ := getExecuteArgumentValue("entire_database")

	if entireDatabase && len(dsn) == 0 {
		onError("If you run with --entire_database dsn parameter is required.")
	}

	if len(tableName) == 0 && (!entireDatabase || len(dsn) == 0) {
		onError("You must either specify --table parameter or run with --entire_database.")
	}

	var db *database
	var tableNames []string
	if len(tableName) == 0 {
		db = &database{dsn: dsn}
		err := db.connect()
		if err != nil {
			onError(err)
		}

		tableNames = db.loadTables()
	} else {
		tableNames = []string{tableName}
	}

	for _, tblName := range tableNames {
		tbl := table{
			tableName:     tblName,
			rootNameSpace: root,
			db:            db,
		}
		tbl.process()
	}
}

func onError(v interface{}) {
	log.Println(v)
	os.Exit(1)
}

func getExecuteArgumentValue(argument string) (bool, string) {
	for _, osArg := range os.Args[1:] {
		subMatch := argumentRegex.FindStringSubmatch(osArg)
		if len(subMatch) == 0 {
			continue
		}

		if argument == subMatch[1] {
			return true, subMatch[2]
		}
	}

	return false, ""
}

func printHelp() {
	fmt.Print(HELP)
}
