package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookUp executes docker-compose up -d on the specified hook
func RouteHookUp(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])

		err := hook.Pull().Run()
		if err != nil {
			jsonOutput(w, http.StatusInternalServerError,
				outputHook("Could not pull images for hook", hook.Name))
			return
		}

		output, err := hook.Up().CombinedOutput()
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
