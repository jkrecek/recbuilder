package main

import (
	"database/sql"
	"fmt"
)

type database struct {
	dsn        string
	connection *sql.DB
}

func (db *database) connect() (err error) {
	db.connection, err = sql.Open("mysql", db.dsn)
	if err != nil {
		return
	}

	err = db.connection.Ping()
	return
}

func (db *database) loadTables() []string {
	rows, err := db.connection.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			panic(err)
		}
		tables = append(tables, table)
	}

	return tables
}

func (db *database) loadTableStructure(tableName string) []columnDescription {
	rows, err := db.connection.Query(fmt.Sprintf("DESCRIBE %s", tableName))
	if err != nil {
		panic(err)
	}

	var columns []columnDescription
	for rows.Next() {
		var column columnDescription
		err := rows.Scan(&column.Field, &column.Type, &column.Null, &column.Key, &column.Default, &column.Extra)
		if err != nil {
			panic(err)
		}

		columns = append(columns, column)
	}

	return columns
}
