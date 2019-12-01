package v1 // import jdel.org/acdc/v1

import (
	"net/http"

	"jdel.org/acdc/api"
)

// RouteAPIKeyList lists API keys
func RouteAPIKeyList(w http.ResponseWriter, r *http.Request) {
	authOK := api.BasicAuthMaster(w, r)

	if authOK == true {
		keys, err := api.AllAPIKeys()
		if err != nil {
			logRoute.WithField("route", "RouteAPIKeyList").Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputKey("Could find keys", ""))
			return
		}

		jsonOutput(w, http.StatusOK,
			keys)
	}
}
