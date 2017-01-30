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
	"github.com/stretchr/testify/assert"
)

func callAddUser(t *testing.T, params *handlers.AddUserParams) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: routes.UserEndpoints[routes.EndpointAddUser],
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}

func TestAddUser(t *testing.T) {
	globalT := t
	defer lifecycle.PurgeModels(t)

	tests := []struct {
		description string
		code        int
		params      *handlers.AddUserParams
	}{
		{
			"Empty User",
			http.StatusBadRequest,
			&handlers.AddUserParams{},
		},
		{
			"Valid User",
			http.StatusCreated,
			&handlers.AddUserParams{Name: "Name", Email: "email+TestHandlerAdd@fake.com", Password: "password"},
		},
		{
			"Duplicate Email",
			http.StatusConflict,
			&handlers.AddUserParams{Name: "Name", Email: "email+TestHandlerAdd@fake.com", Password: "password"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callAddUser(t, tc.params)
			assert.Equal(t, tc.code, rec.Code)

			if httptests.Is2XX(rec.Code) {
				var u payloads.User
				if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, u.ID)
				assert.Equal(t, tc.params.Email, u.Email)
				lifecycle.SaveModels(globalT, &auth.User{ID: u.ID})
			}
		})
	}
}
