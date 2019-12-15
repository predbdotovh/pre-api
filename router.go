package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const hostname = "predb.ovh"

func newRouter() *mux.Router {
	backendUpdates = make(chan triggerAction)
	go backendPump()

	router := mux.NewRouter().
		StrictSlash(true).
		Host(hostname).
		PathPrefix("/api").
		Subrouter()

	for ver, jRoutes := range jsonRoutes {
		versionRouter := router.PathPrefix("/" + ver).Subrouter()
		for _, r := range jRoutes {
			versionRouter.
				Methods(r.Method).
				Path(r.Pattern).
				Name(r.Name).
				Handler(logger(r.Handler, r.Name))
		}
	}

	router.NotFoundHandler = http.HandlerFunc(notFound)

	return router
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	apiErr(w, "404 Not Found")
}
