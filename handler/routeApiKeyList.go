package handler

import (
	"net/http"

	"github.com/jdel/acdc/api"
)

// RouteAPIKeyList lists API keys
func RouteAPIKeyList(w http.ResponseWriter, r *http.Request) {
	authOK := api.BasicAuthMaster(w, r)

	if authOK == true {
		keys, err := api.AllAPIKeys()
		if err != nil {
			logRoute.Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputKey("Could find keys", ""))
			return
		}

		jsonOutput(w, http.StatusOK,
			keys)
	}
}
