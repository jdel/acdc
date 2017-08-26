package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookPs executes docker-compose ps hookName
func RouteHookPs(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])

		output, err := hook.Ps().CombinedOutput()
		if err != nil {
			logRoute.Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not get hook", hook.Name))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputHook(string(output), hook.Name))
	}
}
