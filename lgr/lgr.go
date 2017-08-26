package lgr

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Wrapper hijacks a http.Handler to log requests
func Wrapper(handler http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.WithFields(log.Fields{
			"module":   "handlers", // It's not really handlers but heh...
			"protocol": r.Proto,
			"method":   r.Method,
			"URI":      r.RequestURI,
			"name":     name,
			"time":     time.Since(start),
		}).Info("Request")

		handler.ServeHTTP(w, r)
	})
}
