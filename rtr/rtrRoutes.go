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
		Name:        "Gogs",
		Method:      "POST",
		Pattern:     "/gogs",
		HandlerFunc: handler.RouteGogs,
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
		Name:        "HookActions",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/hooks/{hookName}/{actions}",
		HandlerFunc: v1.RouteHookActions,
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
	route{
		Name:        "Job",
		Method:      "GET",
		PathPrefix:  "/v1",
		Pattern:     "/jobs/{ticket}",
		HandlerFunc: v1.RouteJobs,
	},
}
