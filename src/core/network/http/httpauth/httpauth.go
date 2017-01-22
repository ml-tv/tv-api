package httpauth

import (
	"fmt"
	"net/http"
)

// SetWWWAuthenticate set an auth error to the response
func SetWWWAuthenticate(res http.ResponseWriter, typ, realm string) {
	realmStr := ""
	if realm != "" {
		realmStr = fmt.Sprintf(`realm="%s"`, realm)
	}

	authStr := fmt.Sprintf("%s %s", typ, realmStr)
	res.Header().Set("WWW-Authenticate", authStr)
}
