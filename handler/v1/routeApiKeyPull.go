package v1 // import jdel.org/acdc/v1

import (
	"fmt"
	"net/http"

	"jdel.org/acdc/api"
)

// RouteAPIKeyPull pulls the git remote associated to the key
func RouteAPIKeyPull(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)

		headShortHash, err := key.Pull()
		if err != nil {
			logRoute.WithField("key", apiKey).Error(err)
			jsonOutput(w, http.StatusInternalServerError,
				outputKey("Could not pull remote key", key.Unique))
			return
		}

		jsonOutput(w, http.StatusOK,
			outputKey(fmt.Sprintf("Pulled remote key\nCommit: %s", headShortHash), key.Unique))
	}
}
