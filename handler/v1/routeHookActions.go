package v1

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookActions handles docker-compose actions
func RouteHookActions(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)
	actions := strings.Split(mux.Vars(r)["actions"], "+")

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])
		if hook == nil {
			jsonOutput(w, http.StatusNotFound,
				outputHook("Could not find hook", mux.Vars(r)["hookName"]))
			return
		}

		tickets, err := hook.ExecuteSequentially(actions...)
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not execute actions for", hook.Name))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputHook(tickets[:len(tickets)-1], hook.Name))
	}
}
