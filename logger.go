package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var stdOutLog = log.New(os.Stdout, "", log.LstdFlags)

func httpLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		from := r.Header.Get("x-forwarded-for")
		if from == "" {
			from = r.RemoteAddr
		}

		stdOutLog.Printf(
			"%s %s %s %s %s",
			from,
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
