package main

import (
	"database/sql"
)

const sphinxTable = "pre_plain, pre_rt"

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
