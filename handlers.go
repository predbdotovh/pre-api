package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

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
			return 0, 0, fmt.Errorf("Incorrect parameters (count)")
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
			return 0, 0, fmt.Errorf("Incorrect parameters (page)")
		}

		offset = (page - 1) * count
	}

	offsetStr := query.Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, fmt.Errorf("Incorrect parameters (offset)")
		}
	}
	if offset < 0 {
		offset = defaultOffset
	}

	return offset, count, nil
}

func handleQuery(w http.ResponseWriter, r *http.Request) (*apiRowData, error) {
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

	rows := make([]sphinxRow, 0)
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
	data, err := handleQuery(w, r)
	if err != nil {
		apiErr(w, err.Error())
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=60")
	err = apiSuccess(w, data)
	if err != nil {
		log.Println(err)
	}
}

func liveHandlerV1(w http.ResponseWriter, r *http.Request) {
	data, err := handleQuery(w, r)
	if err != nil {
		apiErr(w, err.Error())
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
		apiErr(w, "Not authorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var p sphinxRow
	err := decoder.Decode(&p)
	if err != nil {
		apiFail(w, err)
		return
	}
	defer r.Body.Close()

	p.proc()

	backendUpdates <- triggerAction{Action: mux.Vars(r)["action"], Row: &p}
	apiSuccess(w, nil)
}

func nukeTriggerHandlerV1(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-forwarded-for") != "" {
		apiErr(w, "Not authorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var n nuke
	err := decoder.Decode(&n)
	if err != nil {
		apiFail(w, err)
		return
	}
	defer r.Body.Close()

	p, err := getPre(sphinx, n.PreID, false)
	if err != nil {
		apiFail(w, err)
		return
	}

	p.setNuke(&n)

	backendUpdates <- triggerAction{Action: n.Type, Row: p}
	apiSuccess(w, nil)
}

func statsHandlerV1(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	tx, err := sphinx.Begin()
	if err != nil {
		apiErr(w, err.Error())
		return
	}
	defer tx.Commit()

	_, err = latestPres(tx, 0, 0, false)
	if err != nil {
		apiErr(w, err.Error())
		return
	}

	meta, err := sphinxMeta(tx)
	if err != nil {
		apiErr(w, err.Error())
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
