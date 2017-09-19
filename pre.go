package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

type sphinxRow struct {
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

func (p *sphinxRow) proc() {
	p.Size /= 1000
}

func scanPresRows(rows *sql.Rows, appendNukes bool) []sphinxRow {
	res := make([]sphinxRow, 0)
	for rows.Next() {
		var r sphinxRow
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

func getPre(db *sql.DB, preID int, withNuke bool) (*sphinxRow, error) {
	sqlQuery := "SELECT id, name, team, cat, genre, url, size, files, pre_at FROM " + sphinxTable +
		" WHERE id = ? OPTION reverse_scan = 1"

	var r sphinxRow

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

func searchPres(tx *sql.Tx, q string, offsetInt, countInt int, withNukes bool) ([]sphinxRow, error) {
	sqlQuery := "SELECT id, name, team, cat, genre, url, size, files, pre_at FROM " + sphinxTable +
		" WHERE MATCH(?) ORDER BY id DESC" +
		" LIMIT " + strconv.Itoa(offsetInt) + "," + strconv.Itoa(countInt) +
		" OPTION reverse_scan = 1"

	sqlRows, err := tx.Query(sqlQuery, replacer.Replace(q))
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	return scanPresRows(sqlRows, withNukes), nil
}

func latestPres(tx *sql.Tx, offsetInt, countInt int, withNukes bool) ([]sphinxRow, error) {
	sqlQuery := "SELECT id, name, team, cat, genre, url, size, files, pre_at FROM " + sphinxTable +
		" ORDER BY id DESC" +
		" LIMIT " + strconv.Itoa(offsetInt) + "," + strconv.Itoa(countInt) +
		" OPTION reverse_scan = 1"

	sqlRows, err := tx.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	return scanPresRows(sqlRows, withNukes), nil
}

func (p *sphinxRow) fetchNuke(db *sql.DB) {
	n, err := getNuke(db, p.ID)
	if err != nil {
		return
	}
	p.setNuke(n)
}

func (p *sphinxRow) setNuke(n *nuke) {
	if n != nil {
		n.setType()
		p.Nuke = n
	}
}
