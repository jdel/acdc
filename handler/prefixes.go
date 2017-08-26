package handler

import (
	"net/http"
	"strings"

	"github.com/jdel/acdc/cfg"
)

// PrefixStatic serves anything under /static/
func PrefixStatic(w http.ResponseWriter, r *http.Request) {
	// Translates prefix from URL to local FS
	localPath := strings.Replace(r.URL.Path, cfg.GetStaticPrefix(), cfg.GetStaticDir(), 1)
	http.ServeFile(w, r, localPath)
}
