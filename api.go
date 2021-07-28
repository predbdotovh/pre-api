package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type apiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type apiRowData struct {
	RowCount int      `json:"rowCount"`
	Rows     []preRow `json:"rows"`
	Offset   int      `json:"offset"`
	ReqCount int      `json:"reqCount"`
	Total    int      `json:"total"`
	Time     float64  `json:"time"`
}

type apiStatsData struct {
	Total int       `json:"total"`
	Date  time.Time `json:"date"`
	Time  float64   `json:"time"`
}

type apiTeamsData struct {
	RowCount int       `json:"rowCount"`
	Rows     []teamRow `json:"rows"`
	Time     float64   `json:"time"`
}

const indent = "    "

func apiSuccess(w http.ResponseWriter, o interface{}) error {
	return apiSend(w, apiResponse{"success", "", o})
}

func apiFail(w http.ResponseWriter, o interface{}) error {
	return apiSend(w, apiResponse{"error", "", o})
}

func apiErr(w http.ResponseWriter, err error) error {
	return apiSend(w, apiResponse{"error", err.Error(), nil})
}

func apiSend(w http.ResponseWriter, o interface{}) error {
	headers := w.Header()
	headers.Set("Content-Type", "application/json; charset=utf-8")
	headers.Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetIndent("", indent)

	return enc.Encode(o)
}
