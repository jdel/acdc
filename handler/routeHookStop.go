package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookStop executes docker-compose start on the specified hook
func RouteHookStop(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])

		output, err := hook.Stop().CombinedOutput()
		if err != nil {
			logRoute.Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not stop hook", hook.Name))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputHook(string(output), hook.Name))
	}
}
