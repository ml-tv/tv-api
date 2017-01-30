package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ml-tv/tv-api/src/components/users/handlers"
	"github.com/ml-tv/tv-api/src/components/users/payloads"
	"github.com/ml-tv/tv-api/src/components/users/routes"
	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/core/security/auth/testdata"
	"github.com/stretchr/testify/assert"
)

func callGetUser(t *testing.T, params *handlers.GetUserParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: routes.UserEndpoints[routes.EndpointGetUser],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}

func TestGetUser(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	u1, s1 := testdata.NewAuth(t)
	u2, s2 := testdata.NewAuth(t)

	tests := []struct {
		description string
		code        int
		params      *handlers.GetUserParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusOK,
			&handlers.GetUserParams{ID: u1.ID},
			nil,
		},
		{
			"Getting an other user",
			http.StatusOK,
			&handlers.GetUserParams{ID: u1.ID},
			httptests.NewRequestAuth(s2.ID, u2.ID),
		},
		{
			"Getting own data",
			http.StatusOK,
			&handlers.GetUserParams{ID: u1.ID},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Getting un-existing user with valid ID",
			http.StatusNotFound,
			&handlers.GetUserParams{ID: "f76700e7-988c-4ae9-9f02-ac3f9d7cd88e"},
			nil,
		},
		{
			"Getting un-existing user with invalid ID",
			http.StatusBadRequest,
			&handlers.GetUserParams{ID: "invalidID"},
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callGetUser(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if httptests.Is2XX(rec.Code) {
				var u payloads.User
				if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
					t.Fatal(err)
				}

				if assert.Equal(t, tc.params.ID, u.ID, "Not the same user") {
					// User access their own data
					if tc.auth != nil && u.ID == tc.auth.UserID {
						assert.NotEmpty(t, u.Email, "Same user needs their private data")
					} else { // user access an other user data
						assert.Empty(t, u.Email, "Should not return private data")
					}
				}
			}
		})
	}
}
