package httptests

// Is2XX returns a true if the HTTP code is a 2XX, false otherwise
func Is2XX(code int) bool {
	return code >= 200 && code < 300
}
