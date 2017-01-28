package main

import "database/sql"

const nukeTable = "nuke"

type nuke struct {
	ID     int    `json:"id"`
	TypeID int    `json:"typeId"`
	Type   string `json:"type"`
	PreID  int    `json:"preId"`
	Reason string `json:"reason"`
	Net    string `json:"net"`
	source string
	NukeAt int64 `json:"nukeAt"`
}

var nukeTypes = map[int]string{
	1: "nuke",
	2: "unnuke",
	3: "modnuke",
	4: "delpre",
	5: "undelpre",
}

func (n *nuke) setType() {
	if n.Type == "" {
		n.Type = nukeTypes[n.TypeID]
	}
}

func getNuke(db *sql.DB, preID int) (*nuke, error) {
	sqlQuery := "SELECT id, pre_id, nuke_id, reason, net, source, UNIX_TIMESTAMP(set_at) FROM " + nukeTable +
		" WHERE pre_id = ? ORDER BY id DESC LIMIT 1"

	var r nuke

	err := db.QueryRow(sqlQuery, preID).Scan(&r.ID, &r.PreID, &r.TypeID, &r.Reason, &r.Net, &r.source, &r.NukeAt)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
