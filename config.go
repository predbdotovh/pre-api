package main

import (
	"encoding/json"
	"log"
	"os"
)

type configuration struct {
	SphinxDatabase string
	MysqlDatabase  string
	Addr           string
}

func loadConfig() *configuration {
	file, err := os.Open("conf/api.conf.json")
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	conf := configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}

	return &conf
}
