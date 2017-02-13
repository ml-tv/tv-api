package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ml-tv/tv-api/src/components/sessions/handlers"
	"github.com/ml-tv/tv-api/src/components/sessions/routes"
	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/core/security/auth"
	"github.com/ml-tv/tv-api/src/core/security/auth/testdata"
	"github.com/ml-tv/tv-api/src/core/storage/db"
	"github.com/stretchr/testify/assert"
)

func TestHandlerDelete(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	u1, s1 := testdata.NewAuth(t)
	u2, s2 := testdata.NewAuth(t)

	t.Run("parallel", func(t *testing.T) {
		tests := []struct {
			description string
			code        int
			params      *handlers.DeleteSessionParams
			auth        *httptests.RequestAuth
		}{
			{
				"Not logged",
				http.StatusUnauthorized,
				&handlers.DeleteSessionParams{Token: s1.ID},
				nil,
			},
			{
				"Deleting an other user sessions",
				http.StatusNotFound,
				&handlers.DeleteSessionParams{Token: s1.ID, CurrentPassword: "fake"},
				httptests.NewRequestAuth(s2.ID, u2.ID),
			},
			{
				"Deleting an invalid ID",
				http.StatusBadRequest,
				&handlers.DeleteSessionParams{Token: "invalid", CurrentPassword: "fake"},
				httptests.NewRequestAuth(s1.ID, u1.ID),
			},
			{
				"Deleting without providing password",
				http.StatusUnauthorized,
				&handlers.DeleteSessionParams{Token: s1.ID},
				httptests.NewRequestAuth(s1.ID, u1.ID),
			},
			// Keep this one last for u1 as it deletes the session
			{
				"Deleting session",
				http.StatusNoContent,
				&handlers.DeleteSessionParams{Token: s1.ID, CurrentPassword: "fake"},
				httptests.NewRequestAuth(s1.ID, u1.ID),
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callHandlerDelete(t, tc.params, tc.auth)
				assert.Equal(t, tc.code, rec.Code)

				if httptests.Is2XX(rec.Code) {
					// We check that the user is still in DB but is flagged for deletion
					var session auth.Session
					stmt := "SELECT * FROM sessions WHERE id=$1 LIMIT 1"
					err := db.Get(&session, stmt, tc.params.Token)
					if err != nil {
						t.Fatal(err)
					}

					if assert.NotEmpty(t, session.ID, "session fully deleted") {
						assert.NotNil(t, session.DeletedAt, "User not marked for deletion")
					}
				}
			})
		}
	})
}

func callHandlerDelete(t *testing.T, params *handlers.DeleteSessionParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: routes.SessionEndpoints[routes.EndpointDeleteSession],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
