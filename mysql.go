package main

import (
	"database/sql"
	"log"
)

func newMysql(dbString string) *sql.DB {
	mysql, err := sql.Open("mysql", dbString)
	if err != nil {
		log.Fatal(err)
	}

	return mysql
}
