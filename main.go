package main

import (
	"log"
	"net/http"
)

func main() {

	conf := loadConfig()
	newSphinx(conf.SphinxDatabase)
	newMysql(conf.MysqlDatabase)

	router := newRouter()

	log.Fatal(http.ListenAndServe(conf.Addr, router))
}
