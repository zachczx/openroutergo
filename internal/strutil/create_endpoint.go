package strutil

// CreateEndpoint creates an endpoint from a base URL and a path.
//
// It removes all the trailing slashes from the base URL and all the leading slashes
// from the path, then joins them together with a single slash in between.
func CreateEndpoint(baseURL, path string) string {
	return RemoveTrailingSlashes(baseURL) + "/" + RemoveLeadingSlashes(path)
}
