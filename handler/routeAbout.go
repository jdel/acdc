package handler

import (
	"net/http"

	"github.com/jdel/acdc/cfg"
)

// RouteAbout displays some information about
// the running instance of acdc
func RouteAbout(w http.ResponseWriter, r *http.Request) {
	about := struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Maintainer string `json:"maintainer"`
		License    string `json:"license"`
		Year       uint   `json:"year"`
	}{
		"acdc",
		cfg.Version,
		"jdel",
		"GNU GPL v3",
		2017,
	}
	jsonOutput(w, http.StatusOK, about)
}
