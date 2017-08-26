package rtr

import (
	"strings"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/lgr"
	log "github.com/sirupsen/logrus"
)

var logRouter = log.WithFields(log.Fields{
	"module": "router",
})

// NewRouter is the mux router that handle main routes
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// Create parh routes from appRoutes in rtr/routes.go
	for _, route := range appRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(lgr.Wrapper(route.HandlerFunc, route.Name))
	}
	// Create path prefix routes from appPrefixes in rtr/routes.go
	for _, prefix := range appPrefixes {
		router.
			Methods(prefix.Method).
			PathPrefix(prefix.Pattern()).
			Name(prefix.Name).
			Handler(lgr.Wrapper(prefix.HandlerFunc, prefix.Name))
	}
	// Prints routes in debug log
	router.Walk(printRoutes)
	return router
}

// Prints all available routes for debugging purpose
func printRoutes(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	template, err := route.GetPathTemplate()
	if err != nil {
		return err
	}
	methods, err := route.GetMethods()
	if err != nil {
		return err
	}
	logRouter.WithFields(log.Fields{
		"name":     route.GetName(),
		"methods":  strings.Join(methods, ","),
		"template": template,
	}).Debug("Route details")
	return nil
}
