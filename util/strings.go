package util

// IsStringInSlice returns true if the strings slice sl
// contains an element equal to s
func IsStringInSlice(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}

	}
	return false
}
