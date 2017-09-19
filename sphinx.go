package main

import (
	"database/sql"
	"log"
)

const sphinxTable = "pre_plain, pre_rt"

var sphinx *sql.DB

func newSphinx(dbString string) {
	var err error
	sphinx, err = sql.Open("mysql", dbString)
	if err != nil {
		log.Fatal(err)
	}
}

func sphinxMeta(tx *sql.Tx) (map[string]string, error) {
	sqlRows, err := tx.Query("SHOW META")
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	rows := make(map[string]string, 0)
	for sqlRows.Next() {
		var name, value string
		err := sqlRows.Scan(&name, &value)
		if err != nil {
			continue
		}
		rows[name] = value
	}

	return rows, nil
}
