package shows_test

import (
	"net/http/httptest"
	"testing"

	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/stretchr/testify/assert"

	"net/http"

	"github.com/ml-tv/tv-api/src/components/shows"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	authdata "github.com/ml-tv/tv-api/src/core/security/auth/testdata"
)

const theLyingGameID = 38207

func callAdd(t *testing.T, params *shows.AddParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: shows.Endpoints[shows.EndpointAdd],
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
	s := NewShowFromTMDb(t, theLyingGameID)

	// Try to add the same show again
	params := &shows.AddParams{TMDbID: s.TMDbID}
	auth := httptests.NewRequestAuth(adminSession.ID, admin.ID)
	rec := callAdd(t, params, auth)
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
			params      *shows.AddParams
			auth        *httptests.RequestAuth
		}{
			{
				"User not logged",
				http.StatusUnauthorized,
				&shows.AddParams{TMDbID: 0},
				nil,
			},
			{
				"User not admin",
				http.StatusUnauthorized,
				&shows.AddParams{TMDbID: 0},
				httptests.NewRequestAuth(userSession.ID, user.ID),
			},
			{
				"Show does not exists",
				http.StatusNotFound,
				&shows.AddParams{TMDbID: 0},
				httptests.NewRequestAuth(adminSession.ID, admin.ID),
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callAdd(t, tc.params, tc.auth)
				assert.Equal(t, tc.code, rec.Code)
			})
		}
	})
}
