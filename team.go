package main

import (
	"database/sql"
	"fmt"
	"log"
)

type teamRow struct {
	Team      string `json:"team"`
	FirstPre  int    `json:"firstPre"`
	LatestPre int    `json:"latestPre"`
	Count     int    `json:"count"`
}

const maxListedTeams = 1000

func listTeams(tx *sql.Tx) ([]teamRow, error) {
	sqlQuery := fmt.Sprintf("SELECT team, MIN(pre_at), MAX(pre_at), COUNT(*) AS count FROM %s GROUP BY team ORDER BY count DESC LIMIT %d", sphinxTable, maxListedTeams)

	rows, err := tx.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]teamRow, 0)
	for rows.Next() {
		var r teamRow
		err := rows.Scan(&r.Team, &r.FirstPre, &r.LatestPre, &r.Count)
		if err != nil {
			log.Print(err)
			continue
		}

		res = append(res, r)
	}

	return res, nil
}
