package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteHookDelete deletes the specific hook yml file
func RouteHookDelete(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)
		hook := key.GetHook(mux.Vars(r)["hookName"])
		if hook == nil {
			jsonOutput(w, http.StatusNotFound,
				outputHook("Could not find hook", mux.Vars(r)["hookName"]))
			return
		}

		if key.IsRemote() {
			logRoute.WithField("key", key.Unique).Error("Cannot create hook on a remote key")
			jsonOutput(w, http.StatusUnprocessableEntity,
				outputHook("Cannot delete hook on a remote key", key.Remote))
			return
		}

		if err := hook.Down().Run(); err != nil {
			logRoute.WithField("key", key.Unique).Error("Could not bring hook down")
			jsonOutput(w, http.StatusUnprocessableEntity,
				outputHook("Could not bring hook down", hook.Name))
			return
		}

		if err := hook.Delete(); err != nil {
			logRoute.WithField("key", key.Unique).Error(err)
			jsonOutput(w, http.StatusUnprocessableEntity,
				outputHook("Could not delete hook", hook.Name))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputHook("Deleted hook", hook.Name))
	}
}
