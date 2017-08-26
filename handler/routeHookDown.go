package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookDown executes docker-compose down on the specified hook
func RouteHookDown(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])

		output, err := hook.Down().CombinedOutput()
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not bring hook down", hook.Name))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputHook(string(output), hook.Name))
	}
}
