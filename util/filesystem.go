package util // import jdel.org/acdc/util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var logUtil = log.WithFields(log.Fields{
	"module": "util",
})

// CreateDir Creates a directory if it doesn't exist
func CreateDir(dir string) (string, error) {
	var err error
	if !FileExists(dir) {
		logUtil.Debug("Creating directory: ", dir)
		if err = os.Mkdir(dir, 0755); err != nil {
			logUtil.Error(err)
		}
	}
	return dir, err
}

// FileExists returns true if the dir exists on filesystem
func FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
