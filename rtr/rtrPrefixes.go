package rtr // import jdel.org/acdc/rtr

import (
	"net/http"

	"jdel.org/acdc/handler"
)

type pathPrefix struct {
	Name        string
	Method      string
	Pattern     func() string
	HandlerFunc http.HandlerFunc
}

type pathPrefixes []pathPrefix

// This variable is responsible for configuring all prefixes
// prefix funcs are in handlers/prefixes.go
var appPrefixes = pathPrefixes{
	pathPrefix{
		Name:        "Static",
		Method:      "GET",
		Pattern:     func() string { return "/static/" },
		HandlerFunc: handler.PrefixStatic,
	},
}
