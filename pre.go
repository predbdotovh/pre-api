package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type preRow struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Team  string  `json:"team"`
	Cat   string  `json:"cat"`
	Genre string  `json:"genre"`
	URL   string  `json:"url"`
	Size  float64 `json:"size"`
	Files int     `json:"files"`
	PreAt int     `json:"preAt"`
	Nuke  *nuke   `json:"nuke"`
}

var replacer = strings.NewReplacer("(", "\\(", ")", "\\)")

const preColumns = "id, name, team, cat, genre, url, size, files, pre_at"
const defaultMaxMatches = 1000

func (p *preRow) proc() {
	p.Size /= 1000
}

func scanPresRows(rows *sql.Rows, appendNukes bool) []preRow {
	res := make([]preRow, 0)
	for rows.Next() {
		var r preRow
		err := rows.Scan(&r.ID, &r.Name, &r.Team, &r.Cat, &r.Genre, &r.URL, &r.Size, &r.Files, &r.PreAt)
		if err != nil {
			log.Print(err)
			continue
		}

		r.proc()
		if appendNukes {
			r.fetchNuke(mysql)
		}

		res = append(res, r)
	}

	return res
}

func getPre(db *sql.DB, preID int, withNuke bool) (*preRow, error) {
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", preColumns, sphinxTable)

	var r preRow

	err := db.QueryRow(sqlQuery, preID).Scan(&r.ID, &r.Name, &r.Team, &r.Cat, &r.Genre, &r.URL, &r.Size, &r.Files, &r.PreAt)
	if err != nil {
		return nil, err
	}

	r.proc()
	if withNuke {
		r.fetchNuke(mysql)
	}

	return &r, nil
}

func getPresById(tx *sql.Tx, preID int, withNuke bool) ([]preRow, error) {
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", preColumns, sphinxTable)

	var r preRow

	err := tx.QueryRow(sqlQuery, preID).Scan(&r.ID, &r.Name, &r.Team, &r.Cat, &r.Genre, &r.URL, &r.Size, &r.Files, &r.PreAt)
	if err != nil {
		return nil, err
	}

	r.proc()
	if withNuke {
		r.fetchNuke(mysql)
	}

	return []preRow{r}, nil
}

func searchPres(tx *sql.Tx, q string, offsetInt, countInt int, withNukes bool) ([]preRow, error) {
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s WHERE MATCH(?) ORDER BY id DESC LIMIT %d,%d", preColumns, sphinxTable, offsetInt, countInt)

	if offsetInt+countInt > defaultMaxMatches {
		sqlQuery += ", max_matches = " + strconv.Itoa(offsetInt+countInt)
	}

	sqlRows, err := tx.Query(sqlQuery, replacer.Replace(q))
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	return scanPresRows(sqlRows, withNukes), nil
}

func latestPres(tx *sql.Tx, offsetInt, countInt int, withNukes bool) ([]preRow, error) {
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s ORDER BY id DESC LIMIT %d,%d", preColumns, sphinxTable, offsetInt, countInt)

	if offsetInt+countInt > defaultMaxMatches {
		sqlQuery += ", max_matches = " + strconv.Itoa(offsetInt+countInt)
	}

	sqlRows, err := tx.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	return scanPresRows(sqlRows, withNukes), nil
}

func (p *preRow) fetchNuke(db *sql.DB) {
	n, err := getNuke(db, p.ID)
	if err != nil {
		return
	}
	p.setNuke(n)
}

func (p *preRow) setNuke(n *nuke) {
	if n != nil {
		n.setType()
		p.Nuke = n
	}
}
