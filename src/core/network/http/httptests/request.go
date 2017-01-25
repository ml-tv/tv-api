package httptests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"

	"github.com/gorilla/mux"
	"github.com/ml-tv/tv-api/src/core/router"
)

// RequestAuth represents the auth data for a request
type RequestAuth struct {
	SessionID string
	UserID    string
}

// ToBasicAuth returns the data using the basic auth format
func (ra *RequestAuth) ToBasicAuth() string {
	authValue := fmt.Sprintf("%s:%s", ra.UserID, ra.SessionID)

	encoded := base64.StdEncoding.EncodeToString([]byte(authValue))
	return "basic " + encoded
}

func NewRequestAuth(sessionID string, userID string) *RequestAuth {
	return &RequestAuth{
		SessionID: sessionID,
		UserID:    userID,
	}
}

// RequestInfo represents the params accepted by NewRequest
type RequestInfo struct {
	Endpoint *router.Endpoint
	Body     interface{} // Optional
	URL      map[string]string
	Auth     *RequestAuth // Optional

	// Router is used to parse Mux Variables. Default on the api router
	Router *mux.Router
}

// NewRequest simulates a new http request executed against the api
func NewRequest(t *testing.T, info *RequestInfo) *httptest.ResponseRecorder {
	body := bytes.NewBufferString("")

	// Parse the body as a JSON object
	if info.Body != nil {
		jsonDump, err := json.Marshal(info.Body)
		if err != nil {
			t.Fatalf("could not create request %s", err)
		}

		body = bytes.NewBuffer(jsonDump)
	}

	// Parse the URL
	url := info.Endpoint.Path
	for param, value := range info.URL {
		url = strings.Replace(url, "{"+param+"}", value, -1)
	}

	req, err := http.NewRequest(info.Endpoint.Verb, url, body)
	if err != nil {
		t.Fatalf("could not execute request %s", err)
	}

	if info.Auth != nil {
		req.Header.Add("Authorization", info.Auth.ToBasicAuth())
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// If no router is provided we assume that we want to execute a regular endpoint
	if info.Router == nil {
		info.Router = router.CurrentRouter
	}

	rec := httptest.NewRecorder()
	info.Router.ServeHTTP(rec, req)
	return rec
}
