package main

import (
	"log"
	"net/http"
	"time"
)

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		from := r.Header.Get("x-forwarded-for")
		if from == "" {
			from = r.RemoteAddr
		}

		log.Printf(
			"%s %s %s %s %s",
			from,
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
