package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/ml-tv/tv-api/src/components/medias/handlers"
	"github.com/ml-tv/tv-api/src/components/medias/routes"
	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/stretchr/testify/assert"

	"net/http"

	"github.com/ml-tv/tv-api/src/components/medias/testdata"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	authdata "github.com/ml-tv/tv-api/src/core/security/auth/testdata"
)

const theLyingGame = 38207

func callAddShow(t *testing.T, params *handlers.AddShowParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: routes.ShowEndpoints[routes.EndpointAddShow],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}

func TestAddShowDuplicate(t *testing.T) {
	defer lifecycle.PurgeModels(t)
	admin, adminSession := authdata.NewAuth(t)
	admin.IsAdmin = true
	admin.Save()

	// Create and save a show
	s := testdata.NewShow(t, theLyingGame)

	// Try to add the same show again
	params := &handlers.AddShowParams{TMDbID: s.TMDbID}
	auth := httptests.NewRequestAuth(adminSession.ID, admin.ID)
	rec := callAddShow(t, params, auth)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAddShowInvalid(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	user, userSession := authdata.NewAuth(t)
	admin, adminSession := authdata.NewAuth(t)
	admin.IsAdmin = true
	admin.Save()

	t.Run("parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			params      *handlers.AddShowParams
			auth        *httptests.RequestAuth
		}{
			{
				"User not logged",
				http.StatusUnauthorized,
				&handlers.AddShowParams{TMDbID: 0},
				nil,
			},
			{
				"User not admin",
				http.StatusUnauthorized,
				&handlers.AddShowParams{TMDbID: 0},
				httptests.NewRequestAuth(userSession.ID, user.ID),
			},
			{
				"Show does not exists",
				http.StatusNotFound,
				&handlers.AddShowParams{TMDbID: 0},
				httptests.NewRequestAuth(adminSession.ID, admin.ID),
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callAddShow(t, tc.params, tc.auth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}
