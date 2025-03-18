package strutil

import "strings"

// RemoveLeadingSlash removes a leading slash from a string.
//
// If the string does not start with a slash, it is returned unchanged.
func RemoveLeadingSlash(str string) string {
	return strings.TrimPrefix(str, "/")
}

// RemoveLeadingSlashes removes leading slashes from a string.
//
// If the string does not start with a slash, it is returned unchanged.
func RemoveLeadingSlashes(str string) string {
	for {
		if strings.HasPrefix(str, "/") {
			str = strings.TrimPrefix(str, "/")
		} else {
			break
		}
	}

	return str
}
