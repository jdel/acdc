package util

import (
	"fmt"
	"regexp"
)

// Prefixes a string with / if the string
// doesn't start with /
func addSlashPrefix(value string) string {
	if match, _ := regexp.MatchString("^/", value); !match {
		value = fmt.Sprintf("/%s", value)
	}
	return value
}
