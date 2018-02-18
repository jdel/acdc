package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdel/acdc/api"
)

// RouteAPIKeyDelete deletes an API key
// TODO: When deleting a key, cycle through all hooks and invoke docker-compose down before deleting
func RouteAPIKeyDelete(w http.ResponseWriter, r *http.Request) {
	authOK := api.BasicAuthMaster(w, r)

	if authOK == true {
		key := api.FindKey(mux.Vars(r)["apiKey"])
		err := key.Delete()
		if err != nil {
			logRoute.WithField("route", "RouteAPIKeyDelete").Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputKey("Could not delete key", key.Unique))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputKey("Deleted key", key.Unique))
	}
}
