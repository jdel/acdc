package v1 // import jdel.org/acdc/v1

import (
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

var logRoute = log.WithFields(log.Fields{
	"module": "v1",
})

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

func outputHook(message string, details string) interface{} {
	output := struct {
		Message  []string `json:"message"`
		HookName string   `json:"hook-name"`
	}{
		Message:  strings.Split(message, "\n"),
		HookName: details,
	}
	return output
}

func outputKey(message string, details string) interface{} {
	output := struct {
		Message  []string `json:"message"`
		HookName string   `json:"key-unique"`
	}{
		Message:  strings.Split(message, "\n"),
		HookName: details,
	}
	return output
}
