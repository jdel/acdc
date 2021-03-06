package v1 // import jdel.org/acdc/v1

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jdel.org/acdc/util"

	"jdel.org/acdc/api"
)

// RouteJobs lists all hooks
func RouteJobs(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		ticket, _ := strconv.Atoi(mux.Vars(r)["ticket"])
		o, ok := util.OutputLog[apiKey][ticket]
		if !ok {
			jsonOutput(w, http.StatusNotFound,
				outputKey("Could not find ticket", apiKey))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputKey(o, apiKey))
	}
}
