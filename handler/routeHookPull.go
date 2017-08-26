package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookPull executes docker-compose start on the specified hook
func RouteHookPull(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])

		output, err := hook.Pull().CombinedOutput()
		if err != nil {
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not pull images for hook", hook.Name))
		}

		jsonOutput(w, http.StatusOK,
			outputHook(string(output), hook.Name))
	}
}
