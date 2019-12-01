package v1 // import jdel.org/acdc/v1

import (
	"net/http"

	"jdel.org/acdc/api"
)

// RouteAPIKeyCreate creates a new API key
func RouteAPIKeyCreate(w http.ResponseWriter, r *http.Request) {
	authOK := api.BasicAuthMaster(w, r)

	if authOK == true {
		key, err := api.NewKey(r.FormValue("unique"), r.FormValue("remote"))
		if err != nil {
			logRoute.WithField("route", "RouteAPIKeyCreate").Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputKey("Could not create key", key.Unique))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputKey("Created key", key.Unique))
	}
}
