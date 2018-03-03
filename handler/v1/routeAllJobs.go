package v1

import (
	"fmt"
	"net/http"

	"github.com/jdel/acdc/api"
	"github.com/jdel/acdc/util"
)

// RouteAllJobs lists all unfinished jobs
func RouteAllJobs(w http.ResponseWriter, r *http.Request) {
	_, authOK := api.BasicAuth(w, r)

	if authOK == true {
		jsonOutput(w, http.StatusOK,
			outputHook(fmt.Sprintf("There are %d jobs in line", len(util.InputQueue)), ""))
	}
}
