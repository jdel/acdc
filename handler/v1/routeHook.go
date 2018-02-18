package v1

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHook handles docker-compose actions
func RouteHook(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)
	actions := strings.Split(mux.Vars(r)["actions"], " ")

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])

		err := hook.Pull().Run()
		if err != nil {
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not pull images for hook", hook.Name))
			return
		}

		output, err := hook.ExecuteSequentially(actions...)
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not bring hook up", hook.Name))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputHook(string(output), hook.Name))
	}
}
