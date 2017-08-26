package rtr

import (
	"net/http"

	"github.com/jdel/acdc/handler"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

// This variable is responsible for configuring all routes
var appRoutes = routes{
	route{
		Name:        "About",
		Method:      "GET",
		Pattern:     "/about",
		HandlerFunc: handler.RouteAbout,
	},
	route{
		Name:        "Slack",
		Method:      "POST",
		Pattern:     "/slack",
		HandlerFunc: handler.RouteSlack,
	},
	route{
		Name:        "DockerHub",
		Method:      "POST",
		Pattern:     "/dockerhub/{apiKey}/{tag}",
		HandlerFunc: handler.RouteDockerHub,
	},
	route{
		Name:        "Github",
		Method:      "POST",
		Pattern:     "/github",
		HandlerFunc: handler.RouteGithub,
	},
	route{
		Name:        "APIKeyCreate",
		Method:      "POST",
		Pattern:     "/key/new",
		HandlerFunc: handler.RouteAPIKeyCreate,
	},
	route{
		Name:        "APIKeyPull",
		Method:      "GET",
		Pattern:     "/key/pull",
		HandlerFunc: handler.RouteAPIKeyPull,
	},
	route{
		Name:        "APIKeyDelete",
		Method:      "DELETE",
		Pattern:     "/key/{apiKey}",
		HandlerFunc: handler.RouteAPIKeyDelete,
	},
	route{
		Name:        "APIKeyLise",
		Method:      "GET",
		Pattern:     "/key",
		HandlerFunc: handler.RouteAPIKeyList,
	},
	route{
		Name:        "HookUp",
		Method:      "GET",
		Pattern:     "/{hookName}/up",
		HandlerFunc: handler.RouteHookUp,
	},
	route{
		Name:        "HookDown",
		Method:      "GET",
		Pattern:     "/{hookName}/down",
		HandlerFunc: handler.RouteHookDown,
	},
	route{
		Name:        "HookStart",
		Method:      "GET",
		Pattern:     "/{hookName}/start",
		HandlerFunc: handler.RouteHookStart,
	},
	route{
		Name:        "HookStop",
		Method:      "GET",
		Pattern:     "/{hookName}/stop",
		HandlerFunc: handler.RouteHookStop,
	},
	route{
		Name:        "HookRestart",
		Method:      "GET",
		Pattern:     "/{hookName}/restart",
		HandlerFunc: handler.RouteHookRestart,
	},
	route{
		Name:        "HookLogs",
		Method:      "GET",
		Pattern:     "/{hookName}/logs",
		HandlerFunc: handler.RouteHookLogs,
	},
	route{
		Name:        "HookPull",
		Method:      "GET",
		Pattern:     "/{hookName}/pull",
		HandlerFunc: handler.RouteHookPull,
	},
	route{
		Name:        "HookPs",
		Method:      "GET",
		Pattern:     "/{hookName}",
		HandlerFunc: handler.RouteHookPs,
	},
	route{
		Name:        "HookCreate",
		Method:      "POST",
		Pattern:     "/{hookName}",
		HandlerFunc: handler.RouteHookCreate,
	},
	route{
		Name:        "HookDelete",
		Method:      "DELETE",
		Pattern:     "/{hookName}",
		HandlerFunc: handler.RouteHookDelete,
	},
	route{
		Name:        "HookAll",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: handler.RouteHookGet,
	},
}
