package v1 // import jdel.org/acdc/v1

import (
	"fmt"
	"net/http"

	"jdel.org/acdc/api"
	"jdel.org/acdc/util"
)

// RouteJobsList lists all unfinished jobs
func RouteJobsList(w http.ResponseWriter, r *http.Request) {
	_, authOK := api.BasicAuth(w, r)

	if authOK == true {
		jsonOutput(w, http.StatusOK,
			outputHook(fmt.Sprintf("There are %d jobs in line", len(util.InputQueue)), ""))
	}
}
