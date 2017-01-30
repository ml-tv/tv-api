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
	"github.com/ml-tv/tv-api/src/core/security/auth"
	"github.com/ml-tv/tv-api/src/core/security/auth/testdata"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	u1, s1 := testdata.NewAuth(t)
	u2, s2 := testdata.NewAuth(t)

	tests := []struct {
		description string
		code        int
		params      *handlers.UpdateUserParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&handlers.UpdateUserParams{ID: u1.ID},
			nil,
		},
		{
			"Updating an other user",
			http.StatusForbidden,
			&handlers.UpdateUserParams{ID: u1.ID},
			httptests.NewRequestAuth(s2.ID, u2.ID),
		},
		{
			"Updating email without providing password",
			http.StatusUnauthorized,
			&handlers.UpdateUserParams{ID: u1.ID, Email: "melvin@fake.io"},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Updating password without providing current Password",
			http.StatusUnauthorized,
			&handlers.UpdateUserParams{ID: u1.ID, NewPassword: "TestUpdateUser"},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Updating regular field",
			http.StatusOK,
			&handlers.UpdateUserParams{ID: u1.ID, Name: "Melvin"},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		{
			"Updating email to a used one",
			http.StatusConflict,
			&handlers.UpdateUserParams{ID: u1.ID, CurrentPassword: "fake", Email: u2.Email},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
		// Keep this one last for u1 as it changes the password
		{
			"Updating password",
			http.StatusOK,
			&handlers.UpdateUserParams{ID: u1.ID, CurrentPassword: "fake", NewPassword: "TestUpdateUser"},
			httptests.NewRequestAuth(s1.ID, u1.ID),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerUpdate(t, tc.params, tc.auth)
			assert.Equal(t, tc.code, rec.Code)

			if httptests.Is2XX(rec.Code) {
				var u payloads.User
				if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
					t.Fatal(err)
				}

				if tc.params.Name != "" {
					assert.NotEmpty(t, tc.params.Name, u.Name)
				}

				if tc.params.Email != "" {
					assert.NotEmpty(t, tc.params.Email, u.Email)
				}

				if tc.params.NewPassword != "" {
					// To check the password has been updated with need to get the
					// encrypted version, and compare it to the raw one
					updatedUser, err := auth.GetUser(tc.params.ID)
					if err != nil {
						t.Fatal(err)
					}

					hash := updatedUser.Password
					assert.True(t, auth.IsPasswordValid(hash, tc.params.NewPassword))
				}
			}
		})
	}
}

func callHandlerUpdate(t *testing.T, params *handlers.UpdateUserParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: routes.UserEndpoints[routes.EndpointUpdateUser],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}
