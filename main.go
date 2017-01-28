package main

import (
	"log"
	"net/http"
)

const httpListen = "127.0.0.1:8081"

func main() {

	conf := loadConfig()
	newSphinx(conf.SphinxDatabase)
	newMysql(conf.MysqlDatabase)

	router := newRouter()

	log.Fatal(http.ListenAndServe(httpListen, router))
}
