package strutil

import "strings"

// RemoveTrailingSlash removes a trailing slash from a string.
//
// If the string does not end with a slash, it is returned unchanged.
func RemoveTrailingSlash(str string) string {
	return strings.TrimSuffix(str, "/")
}

// RemoveTrailingSlashes removes all trailing slashes from a string.
//
// If the string does not end with a slash, it is returned unchanged.
func RemoveTrailingSlashes(str string) string {
	for {
		if strings.HasSuffix(str, "/") {
			str = strings.TrimSuffix(str, "/")
		} else {
			break
		}
	}

	return str
}
