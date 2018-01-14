package rtr

import (
	"net/http"

	"github.com/jdel/acdc/handler"
	"github.com/jdel/acdc/handler/v1"
)

type route struct {
	Name        string
	Method      string
	PathPrefix  string
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
		Pattern:     "/dockerhub",
		HandlerFunc: handler.RouteDockerHub,
	},
	route{
		Name:        "Github",
		Method:      "POST",
		Pattern:     "/github",
		HandlerFunc: handler.RouteGithub,
	},
	route{
		Name:        "APIKeys",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/keys",
		HandlerFunc: v1.RouteAPIKeyList,
	},
	route{
		Name:        "APIKeyCreate",
		Method:      "POST",
		PathPrefix:  "/v1",
		Pattern:     "/keys",
		HandlerFunc: v1.RouteAPIKeyCreate,
	},
	route{
		Name:        "APIKeyDelete",
		Method:      "DELETE",
		PathPrefix:  "/v1",
		Pattern:     "/keys/{apiKey}",
		HandlerFunc: v1.RouteAPIKeyDelete,
	},
	route{
		Name:        "APIKeyPull",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/pull",
		HandlerFunc: v1.RouteAPIKeyPull,
	},
	route{
		Name:        "HookUp",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}/up",
		HandlerFunc: v1.RouteHookUp,
	},
	route{
		Name:        "HookDown",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}/down",
		HandlerFunc: v1.RouteHookDown,
	},
	route{
		Name:        "HookStart",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}/start",
		HandlerFunc: v1.RouteHookStart,
	},
	route{
		Name:        "HookStop",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}/stop",
		HandlerFunc: v1.RouteHookStop,
	},
	route{
		Name:        "HookRestart",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}/restart",
		HandlerFunc: v1.RouteHookRestart,
	},
	route{
		Name:        "HookLogs",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}/logs",
		HandlerFunc: v1.RouteHookLogs,
	},
	route{
		Name:        "HookPull",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}/pull",
		HandlerFunc: v1.RouteHookPull,
	},
	route{
		Name:        "HookPs",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}",
		HandlerFunc: v1.RouteHookPs,
	},
	route{
		Name:        "HookCreate",
		Method:      "POST",
		PathPrefix:  "/v1",
		Pattern:     "/hooks",
		HandlerFunc: v1.RouteHookCreate,
	},
	route{
		Name:        "HookDelete",
		Method:      "DELETE",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}",
		HandlerFunc: v1.RouteHookDelete,
	},
	route{
		Name:        "HookAll",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks",
		HandlerFunc: v1.RouteHookGet,
	},
}
