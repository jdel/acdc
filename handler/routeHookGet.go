package handler

import (
	"net/http"
	"strings"

	"github.com/jdel/acdc/api"
)

// RouteHookGet lists all hooks
func RouteHookGet(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hooks := key.AllHooks()
		var hookNames []string

		for _, hook := range hooks {
			hookNames = append(hookNames, hook.Name)
		}

		jsonOutput(w, http.StatusOK,
			outputKey(strings.Join(hookNames, "\n"), key.Unique))
	}
}
