package handler

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
		apiKey := mux.Vars(r)["apiKey"]
		key := api.FindKey(apiKey)
		err := key.Delete()
		if err != nil {
			logRoute.Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputKey("Could not delete key", key.Unique))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputKey("Deleted key", key.Unique))
	}
}
