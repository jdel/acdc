package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookRestart executes docker-compose start on the specified hook
func RouteHookRestart(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])

		output, err := hook.Restart().CombinedOutput()
		if err != nil {
			logRoute.WithField("route", "RouteHookRestart").Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not restart hook", hook.Name))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputHook(string(output), hook.Name))
	}
}
