package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ml-tv/tv-api/src/components/users/handlers"
	"github.com/ml-tv/tv-api/src/components/users/routes"
	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/core/security/auth"
	"github.com/ml-tv/tv-api/src/core/security/auth/testdata"
	"github.com/ml-tv/tv-api/src/core/storage/db"
	"github.com/stretchr/testify/assert"
)

func callDeleteUser(t *testing.T, params *handlers.DeleteUserParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: routes.UserEndpoints[routes.EndpointDeleteUser],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}

func TestDeleteUser(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	u1, s1 := testdata.NewAuth(t)
	u2, s2 := testdata.NewAuth(t)

	tests := []struct {
		description string
		code        int
		params      *handlers.DeleteUserParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&handlers.DeleteUserParams{ID: u1.ID},
			nil,
		},
		{
			"Deleting an other user",
			http.StatusForbidden,
			&handlers.DeleteUserParams{ID: u1.ID},
			httptests.NewRequestAuth(s2.ID, u2.ID),
		},
		{
			"Deleting without providing password",
			http.StatusUnauthorized,
			&handlers.DeleteUserParams{ID: u1.ID},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		// Keep this one last for u1 as it deletes the user
		{
			"Deleting user",
			http.StatusNoContent,
			&handlers.DeleteUserParams{ID: u1.ID, CurrentPassword: "fake"},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callDeleteUser(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if httptests.Is2XX(rec.Code) {
				// We check that the user is still in DB but is flagged for deletion
				var user auth.User
				stmt := "SELECT * FROM users WHERE id=$1 LIMIT 1"
				err := db.Get(&user, stmt, tc.params.ID)
				if err != nil {
					t.Fatal(err)
				}

				if assert.NotEmpty(t, user.ID, "User fully deleted") {
					assert.NotNil(t, user.DeletedAt, "User not marked for deletion")
				}
			}
		})
	}
}
