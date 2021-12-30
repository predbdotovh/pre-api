package main

import (
	"errors"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func newRouter(hostname string) *mux.Router {
	backendUpdates = make(chan triggerAction)
	go backendPump()

	router := mux.NewRouter()

	publicRouter := router.
		StrictSlash(true).
		Host(hostname).
		PathPrefix("/api").
		Subrouter()

	for ver, jRoutes := range jsonRoutes {
		versionRouter := publicRouter.PathPrefix("/" + ver).Subrouter()
		for _, r := range jRoutes {
			versionRouter.
				Methods(r.Method).
				Path(r.Pattern).
				Name(r.Name).
				Handler(httpLogger(r.Handler, r.Name))
		}
	}

	privateRouter := router.
		StrictSlash(true).
		Host("localhost").
		PathPrefix("/trigger").
		Subrouter()

	for _, r := range triggerRoutes {
		privateRouter.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(httpLogger(r.Handler, r.Name))
	}

	router.NotFoundHandler = http.HandlerFunc(notFound)

	return router
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	err := apiErr(w, errors.New("not Found"))
	if err != nil {
		log.Println(err)
	}
}
