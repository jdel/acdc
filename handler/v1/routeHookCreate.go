package v1

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jdel/acdc/api"
	"github.com/jdel/acdc/cfg"
)

// RouteHookCreate uploads compose file to hookName.yml
func RouteHookCreate(w http.ResponseWriter, r *http.Request) {
	apiKey, authOK := api.BasicAuth(w, r)

	if authOK == true {
		key := api.FindKey(apiKey)

		if key.IsRemote() {
			logRoute.WithField("key", key.Unique).Error("Cannot create hook on a remote key")
			jsonOutput(w, http.StatusUnprocessableEntity,
				outputKey("Cannot create hook on a remote key", key.Remote))
			return
		}

		alphanum := regexp.MustCompile("[^a-zA-Z0-9]+")
		paramHookName := alphanum.ReplaceAllString(strings.TrimSpace(r.FormValue("name")), "")

		if paramHookName == "" {
			logRoute.WithField("route", "RouteHookCreate").Error("Missing hook name")
			jsonOutput(w, http.StatusBadRequest,
				outputHook("Missing hook name", paramHookName))
			return
		}

		inFile, handler, err := r.FormFile("compose")
		if err != nil {
			logRoute.WithField("route", "RouteHookCreate").Error(err)
			jsonOutput(w, http.StatusBadRequest,
				outputHook("Corrupted file", paramHookName))
			return
		}
		defer inFile.Close()

		logRoute.Debugf("%v", handler.Header)

		outFile, err := os.OpenFile(fmt.Sprintf("%s.yml", filepath.Join(cfg.GetComposeDir(), apiKey, paramHookName)), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			logRoute.WithField("route", "RouteHookCreate").Error(err)
			jsonOutput(w, http.StatusUnprocessableEntity,
				outputHook("Could not create hook", paramHookName))
			return
		}
		defer outFile.Close()
		io.Copy(outFile, inFile)

		jsonOutput(w, http.StatusOK,
			outputHook("Created hook", paramHookName))
	}
}
