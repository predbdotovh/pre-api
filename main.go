package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var sphinx *sql.DB
var mysql *sql.DB

func main() {
	sphinxDatabase := getEnv("SPHINX_DATABASE", "tcp(127.0.0.1:9306)/?interpolateParams=true")
	nukesDatabase := getEnv("NUKES_DATABASE", "")
	listenAddr := getEnv("LISTEN_ADDRESS", "127.0.0.1:8088")

	sphinx = newMysql(sphinxDatabase)
	mysql = newMysql(nukesDatabase)

	router := newRouter()

	log.Fatal(http.ListenAndServe(listenAddr, router))
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == "" {
			log.Fatal("Missing mandatory env variable : " + key)
		}

		return defaultValue
	}

	return value
}
