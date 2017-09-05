package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookLogs executes docker-compose logs on the specified hook
func RouteHookLogs(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])

		output, err := hook.Logs().CombinedOutput()
		if err != nil {

			logRoute.WithField("route", "RouteHookLogs").Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not get logs for hook", hook.Name))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputHook(string(output), hook.Name))
	}
}
