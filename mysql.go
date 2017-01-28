package main

import (
	"database/sql"
	"log"
)

var mysql *sql.DB

func newMysql(dbString string) {
	var err error
	mysql, err = sql.Open("mysql", dbString)
	if err != nil {
		log.Fatal(err)
	}
}
