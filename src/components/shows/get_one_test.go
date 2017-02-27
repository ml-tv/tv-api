package shows_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"net/http"

	"github.com/ml-tv/tv-api/src/components/shows"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	authdata "github.com/ml-tv/tv-api/src/core/security/auth/testdata"
)

const breakingBadID = 1396

func callGetOne(t *testing.T, params *shows.GetOneParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: shows.Endpoints[shows.EndpointGetOne],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}

func TestGetOne(t *testing.T) {
	defer lifecycle.PurgeModels(t)
	regUser, regUserSession := authdata.NewAuth(t)

	// Create and save a show
	s := NewShow(t, &shows.Show{
		TMDbID:       breakingBadID,
		Name:         "Breaking Bad",
		OriginalName: "Breaking Bad",
		Synopsis:     "Breaking Bad is the story of Walter White, a struggling high school chemistry teacher who is diagnosed with inoperable lung cancer at the beginning of the series. He turns to a life of crime, producing and selling methamphetamine, in order to secure his family's financial future before he dies, teaming with his former student, Jesse Pinkman.",
		Status:       shows.ShowStatusFinished,
		DayOfWeek:    time.Sunday,
		Website:      "http://www.amctv.com/shows/breaking-bad",
		BackdropPath: "/eSzpy96DwBujGFj0xMbXBcGcfxX.jpg",
		PosterPath:   "/1yeVJox3rjo2jBKrrihIMj7uoS9.jpg",
	})

	t.Run("Parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			params      *shows.GetOneParams
			auth        *httptests.RequestAuth
		}{
			{
				"Invalid UUID",
				http.StatusBadRequest,
				&shows.GetOneParams{ID: "Invalid"},
				httptests.NewRequestAuth(regUserSession.ID, regUser.ID),
			},
			{
				"Unexisting UUID",
				http.StatusNotFound,
				&shows.GetOneParams{ID: "8cd5a72e-c285-48bc-b03e-d842c453ec4b"},
				httptests.NewRequestAuth(regUserSession.ID, regUser.ID),
			},
			{
				"Get Breaking Bad",
				http.StatusOK,
				&shows.GetOneParams{ID: s.ID},
				httptests.NewRequestAuth(regUserSession.ID, regUser.ID),
			},
		}
		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callGetOne(t, tc.params, tc.auth)
				require.Equal(t, tc.code, rec.Code)

				if tc.code == http.StatusOK {
					var pld shows.Payload
					if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, s.Name, pld.Name)
					assert.Equal(t, s.OriginalName, pld.Name)
					assert.Equal(t, s.Synopsis, pld.Synopsis)
					assert.Equal(t, s.Status, pld.Status)
					assert.Equal(t, int(s.DayOfWeek), pld.DayOfWeek)
					assert.Equal(t, s.Website, pld.Website)
				}
			})
		}
	})
}
