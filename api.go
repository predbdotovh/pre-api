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
	RowCount int         `json:"rowCount"`
	Rows     []sphinxRow `json:"rows"`
	Offset   int         `json:"offset"`
	ReqCount int         `json:"reqCount"`
	Total    int         `json:"total"`
	Time     float64     `json:"time"`
}

type apiStatsData struct {
	Total int       `json:"total"`
	Date  time.Time `json:"date"`
	Time  float64   `json:"time"`
}

const indent = "    "

func apiSuccess(w http.ResponseWriter, o interface{}) error {
	return apiSend(w, apiResponse{"success", "", o})
}

func apiSuccessStr(w http.ResponseWriter, o interface{}, msg string) error {
	return apiSend(w, apiResponse{"success", msg, o})
}

func apiFail(w http.ResponseWriter, o interface{}) error {
	return apiSend(w, apiResponse{"error", "", o})
}

func apiErr(w http.ResponseWriter, msg string) error {
	return apiSend(w, apiResponse{"error", msg, nil})
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
