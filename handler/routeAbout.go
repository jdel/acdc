package handler // import jdel.org/acdc/handler

import (
	"time"
	"net/http"

	"jdel.org/acdc/cfg"
)

// RouteAbout displays some information about
// the running instance of acdc
func RouteAbout(w http.ResponseWriter, r *http.Request) {
	about := struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Maintainer string `json:"maintainer"`
		License    string `json:"license"`
		Year       int    `json:"year"`
	}{
		"acdc",
		cfg.Version,
		"jdel",
		"GNU GPL v3",
		time.Now().Year(),
	}
	jsonOutput(w, http.StatusOK, about)
}
