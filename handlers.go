package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

const defaultOffset = 0
const defaultCount = 20
const maxCount = 100

func pagination(query url.Values) (int, int, error) {
	var err error
	count := defaultCount
	offset := defaultOffset

	countStr := query.Get("count")
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			return 0, 0, fmt.Errorf("incorrect parameters (count)")
		}

	}
	if count > maxCount {
		count = maxCount
	}
	if count < 0 {
		count = defaultCount
	}

	pageStr := query.Get("page")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, fmt.Errorf("incorrect parameters (page)")
		}

		offset = (page - 1) * count
	}

	offsetStr := query.Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, fmt.Errorf("incorrect parameters (offset)")
		}
	}
	if offset < 0 {
		offset = defaultOffset
	}

	return offset, count, nil
}

func handleQuery(r *http.Request) (*apiRowData, error) {
	t := time.Now()

	query := r.URL.Query()
	idStr := query.Get("id")
	q := query.Get("q")

	offset, count, err := pagination(query)
	if err != nil {
		return nil, err
	}

	tx, err := sphinx.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	rows := make([]preRow, 0)
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		rows, err = getPresById(tx, id, true)
	} else if q == "" {
		rows, err = latestPres(tx, offset, count, true)
	} else {
		rows, err = searchPres(tx, q, offset, count, true)
	}
	if err != nil {
		return nil, err
	}

	meta, err := sphinxMeta(tx)
	if err != nil {
		return nil, err
	}

	total, err := strconv.Atoi(meta["total_found"])
	if err != nil {
		return nil, err
	}

	data := &apiRowData{
		RowCount: len(rows),
		Rows:     rows,
		Offset:   offset,
		ReqCount: count,
		Total:    total,
		Time:     time.Since(t).Seconds(),
	}

	return data, nil
}

func rootHandlerV1(w http.ResponseWriter, r *http.Request) {
	data, err := handleQuery(r)
	if err != nil {
		err = apiErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=60")
	err = apiSuccess(w, data)
	if err != nil {
		log.Println(err)
	}
}

func teamsHandlerV1(w http.ResponseWriter, _ *http.Request) {
	t := time.Now()

	tx, err := sphinx.Begin()
	if err != nil {
		err = apiErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer tx.Commit()

	rows, err := listTeams(tx)
	if err != nil {
		err = apiErr(w, err)
		if err != nil {
			log.Println(err)
		}

		return
	}

	data := &apiTeamsData{
		RowCount: len(rows),
		Rows:     rows,
		Offset:   0,
		ReqCount: maxListedTeams,
		Total:    len(rows),
		Time:     time.Since(t).Seconds(),
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	err = apiSuccess(w, data)
	if err != nil {
		log.Println(err)
	}
}

func liveHandlerV1(w http.ResponseWriter, r *http.Request) {
	data, err := handleQuery(r)
	if err != nil {
		err = apiErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	err = apiSuccess(w, data)
	if err != nil {
		log.Println(err)
	}
}

func websocketHandlerV1(w http.ResponseWriter, r *http.Request) {
	websocketUpgrader(w, r)
}

func preTriggerHandlerV1(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-forwarded-for") != "" {
		err := apiErr(w, errors.New("not authorized"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	var p preRow
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		err = apiFail(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer r.Body.Close()

	p.proc()

	backendUpdates <- triggerAction{Action: mux.Vars(r)["action"], Row: &p}
	err = apiSuccess(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func nukeTriggerHandlerV1(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-forwarded-for") != "" {
		err := apiErr(w, errors.New("not authorized"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	var n nuke
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		err = apiFail(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer r.Body.Close()

	p, err := getPre(sphinx, n.PreID, false)
	if err != nil {
		err = apiFail(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	p.setNuke(&n)

	backendUpdates <- triggerAction{Action: n.Type, Row: p}
	err = apiSuccess(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func statsHandlerV1(w http.ResponseWriter, _ *http.Request) {
	t := time.Now()

	tx, err := sphinx.Begin()
	if err != nil {
		err = apiErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer tx.Commit()

	_, err = latestPres(tx, 0, 0, false)
	if err != nil {
		err = apiErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	meta, err := sphinxMeta(tx)
	if err != nil {
		err = apiErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	total, err := strconv.Atoi(meta["total_found"])
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Cache-Control", "public, max-age=60")
	data := apiStatsData{
		Total: total,
		Date:  time.Now(),
		Time:  time.Since(t).Seconds(),
	}
	err = apiSuccess(w, data)
	if err != nil {
		log.Println(err)
	}
}

func rssHandlerV1(w http.ResponseWriter, r *http.Request) {
	data, err := handleQuery(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	feed := &feeds.Feed{
		Title:   "PreDB",
		Link:    &feeds.Link{Href: fmt.Sprintf("https://%s/", hostname)},
		Created: time.Now(),
	}

	for _, row := range data.Rows {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       row.Name,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://%s/?id=%d", hostname, row.ID)},
			Description: fmt.Sprintf("Cat:%s | Genre:%s | Size:%0.fMB | Files:%d | ID:%d", row.Cat, row.Genre, row.Size, row.Files, row.ID),
			Created:     time.Unix(int64(row.PreAt), 0),
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(rss))
	if err != nil {
		log.Println(err)
	}
}
