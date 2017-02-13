package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ml-tv/tv-api/src/components/sessions/handlers"
	"github.com/ml-tv/tv-api/src/components/sessions/payloads"
	"github.com/ml-tv/tv-api/src/components/sessions/routes"
	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/core/security/auth"
	"github.com/ml-tv/tv-api/src/core/security/auth/testdata"
	"github.com/stretchr/testify/assert"
)

func TestAddSession(t *testing.T) {
	defer lifecycle.PurgeModels(t)
	u1 := testdata.NewUser(t, nil)

	t.Run("parallel", func(t *testing.T) {
		tests := []struct {
			description string
			code        int
			params      *handlers.AddSessionParams
		}{
			{
				"Invalid email",
				http.StatusBadRequest,
				&handlers.AddSessionParams{Email: "invalid@fake.com", Password: "fake"},
			},
			{
				"Invalid password",
				http.StatusBadRequest,
				&handlers.AddSessionParams{Email: u1.Email, Password: "invalid"},
			},
			{
				"Valid Request",
				http.StatusCreated,
				&handlers.AddSessionParams{Email: u1.Email, Password: "fake"},
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callAddUser(t, tc.params)
				assert.Equal(t, tc.code, rec.Code)

				if httptests.Is2XX(rec.Code) {
					var session payloads.Session
					if err := json.NewDecoder(rec.Body).Decode(&session); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, session.Token)
					assert.Equal(t, u1.ID, session.UserID)

					// clean the test
					(&auth.Session{ID: session.Token}).FullyDelete()
				}
			})
		}
	})
}

func callAddUser(t *testing.T, params *handlers.AddSessionParams) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: routes.SessionEndpoints[routes.EndpointAddSession],
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}
