package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/stretchr/testify/assert"

	"github.com/ml-tv/tv-api/src/core/security/auth/testdata"
)

// TestEndpointExecution test that an endpoint gets executed with no issue.
// It only tests that the main middleware does what it is supposed to do, and
// therefore does not tests details (like checking the params are parsed correctly)
func TestEndpointExecution(t *testing.T) {
	// Handler used for our request. We just need to know if it is called or not
	hdlr := func(req *router.Request) error {
		req.NoContent()
		return nil
	}

	u, s := testdata.NewAuth(t)
	defer lifecycle.PurgeModels(t)

	tests := []struct {
		description string
		endpoint    *router.Endpoint
		params      interface{}
		auth        *httptests.RequestAuth
		code        int
	}{
		{
			"Basic public GET",
			&router.Endpoint{Verb: "GET", Path: "/items", Handler: hdlr},
			struct{}{},
			nil,
			http.StatusNoContent,
		},
		{
			"Private GET as anonymous",
			&router.Endpoint{Verb: "GET", Path: "/items/{id}", Handler: hdlr, Auth: router.LoggedUserAccess},
			struct {
				ID string `json:"id"`
			}{ID: "item-id"},
			nil,
			http.StatusUnauthorized,
		},
		{
			"Private GET as logged user",
			&router.Endpoint{Verb: "GET", Path: "/items/{id}", Handler: hdlr, Auth: router.LoggedUserAccess},
			struct {
				ID string `json:"id"`
			}{ID: "item-id"},
			httptests.NewRequestAuth(s.ID, u.ID),
			http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := execHandler(t, tc.endpoint, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)
		})
	}
}

func execHandler(t *testing.T, e *router.Endpoint, params interface{}, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	r := mux.NewRouter()
	r.Methods(e.Verb).Path(e.Path).Handler(router.Handler(e))

	ri := &httptests.RequestInfo{
		Endpoint: e,
		Params:   params,
		Router:   r,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
