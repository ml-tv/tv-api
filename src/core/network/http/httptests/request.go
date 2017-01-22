package httptests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/gorilla/mux"
)

// RequestAuth represents the auth data for a request
type RequestAuth struct {
	SessionID string
	UserID    string
}

// NewRequestAuth creates new request auth
func NewRequestAuth(sessionID string, userID string) *RequestAuth {
	return &RequestAuth{
		SessionID: sessionID,
		UserID:    userID,
	}
}

// RequestInfo represents the params accepted by NewRequest
type RequestInfo struct {
	Endpoint *router.Endpoint
	URI      string       // Optional
	Params   interface{}  // Optional
	Auth     *RequestAuth // Optional

	// Router is used to parse Mux Variables. Default on the api router
	Router *mux.Router
}

// NewRequest simulates a new http request executed against the api
func NewRequest(t *testing.T, info *RequestInfo) *httptest.ResponseRecorder {
	params := bytes.NewBufferString("")

	if info.Params != nil {
		jsonDump, err := json.Marshal(info.Params)
		if err != nil {
			t.Fatalf("could not create request %s", err)
		}

		params = bytes.NewBuffer(jsonDump)
	}

	req, err := http.NewRequest(info.Endpoint.Verb, info.URI, params)
	if err != nil {
		t.Fatalf("could not execute request %s", err)
	}

	if info.Auth != nil {
		req.Header.Add("X-Session-Token", info.Auth.SessionID)
		req.Header.Add("X-User-Id", info.Auth.UserID)
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
