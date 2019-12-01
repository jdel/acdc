package api // import jdel.org/acdc/api

import (
	"crypto/subtle"
	"net/http"

	"jdel.org/acdc/cfg"
)

// BasicAuth provides basic auth
func BasicAuth(w http.ResponseWriter, r *http.Request) (apiKey string, ok bool) {
	expectedUsername := "api-key"
	w.Header().Set("WWW-Authenticate", `Basic realm="Authentication Required"`)
	user, pass, ok := r.BasicAuth()

	allKeys, err := AllAPIKeys()
	if err != nil {
		return "", false
	}

	_, keyAllowed := allKeys[pass]

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(expectedUsername)) != 1 || !keyAllowed {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorised"))
		return "", false
	}

	return pass, true
}

// BasicAuthMaster provides basic auth for master key
func BasicAuthMaster(w http.ResponseWriter, r *http.Request) bool {
	expectedUsername := "api-key"
	w.Header().Set("WWW-Authenticate", `Basic realm="Authentication Required"`)
	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(expectedUsername)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(cfg.GetMasterKey())) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorised"))
		return false
	}

	return true
}
