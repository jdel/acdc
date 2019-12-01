package handler // import jdel.org/acdc/handler

import (
	"encoding/gob"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var logRoute = log.WithFields(log.Fields{
	"module": "handlers",
})

func getScheme(r *http.Request) string {
	scheme := "http"
	if r.URL.IsAbs() {
		scheme = r.URL.Scheme
	}
	return scheme
}

func jsonOutput(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		logRoute.Fatal(err)
	}
	return err
}

func plainOutput(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(status)
	err := gob.NewEncoder(w).Encode(v)
	if err != nil {
		logRoute.Fatal(err)
	}
	return err
}
