package api

import log "github.com/sirupsen/logrus"

var logAPI = log.WithFields(log.Fields{
	"module": "auth",
})
