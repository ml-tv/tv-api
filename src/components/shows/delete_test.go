package shows_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/stretchr/testify/require"

	"net/http"

	"github.com/ml-tv/tv-api/src/components/shows"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	authdata "github.com/ml-tv/tv-api/src/core/security/auth/testdata"
)

const betterCallSaulID = 60059
const sonsOfAnarchyID = 1409

func callDelete(t *testing.T, params *shows.DeleteParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: shows.Endpoints[shows.EndpointDelete],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}

func TestInvalidDelete(t *testing.T) {
	defer lifecycle.PurgeModels(t)
	regUser, regUserSession := authdata.NewAuth(t)
	admin, adminSession := authdata.NewAuth(t)
	admin.IsAdmin = true
	admin.Save()

	// Create and save a show
	s := NewShow(t, &shows.Show{
		TMDbID:       betterCallSaulID,
		Name:         "Better Call Saul",
		OriginalName: "Better Call Saul",
		Synopsis:     `Six years before Saul Goodman meets Walter White. We meet him when the man who will become Saul Goodman is known as Jimmy McGill, a small-time lawyer searching for his destiny, and, more immediately, hustling to make ends meet. Working alongside, and, often, against Jimmy, is “fixer” Mike Erhmantraut. The series will track Jimmy's transformation into Saul Goodman, the man who puts “criminal” in “criminal lawyer".`,
		Status:       shows.ShowStatusPaused,
		DayOfWeek:    time.Sunday,
		Website:      "http://www.amctv.com/shows/breaking-bad",
		BackdropPath: "/ljik3PqnobCL9fNYJRrDD8eTuFe.jpg",
		PosterPath:   "/b6tpspJMOauCQqud0KOcwNSq3F3.jpg",
	})

	t.Run("Parallel", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			params      *shows.DeleteParams
			auth        *httptests.RequestAuth
		}{
			{
				"Invalid UUID",
				http.StatusBadRequest,
				&shows.DeleteParams{ID: "Invalid"},
				httptests.NewRequestAuth(adminSession.ID, admin.ID),
			},
			{
				"Unexisting UUID",
				http.StatusNotFound,
				&shows.DeleteParams{ID: "8cd5a72e-c285-48bc-b03e-d842c453ec4b"},
				httptests.NewRequestAuth(adminSession.ID, admin.ID),
			},
			{
				"Anonymous user",
				http.StatusUnauthorized,
				&shows.DeleteParams{ID: s.ID},
				nil,
			},
			{
				"Regular user",
				http.StatusUnauthorized,
				&shows.DeleteParams{ID: s.ID},
				httptests.NewRequestAuth(regUserSession.ID, regUser.ID),
			},
		}
		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callDelete(t, tc.params, tc.auth)
				require.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func TestValidDelete(t *testing.T) {
	defer lifecycle.PurgeModels(t)
	admin, adminSession := authdata.NewAuth(t)
	admin.IsAdmin = true
	admin.Save()

	// Create and save a show
	s := NewShow(t, &shows.Show{
		TMDbID:       sonsOfAnarchyID,
		Name:         "Sons of Anarchy",
		OriginalName: "Sons of Anarchy",
		Synopsis:     `An adrenalized drama with darkly comedic undertones that explores a notorious outlaw motorcycle club’s (MC) desire to protect its livelihood while ensuring that their simple, sheltered town of Charming, California remains exactly that, charming. The MC must confront threats from drug dealers, corporate developers, and overzealous law officers. Behind the MC’s familial lifestyle and legally thriving automotive shop is a ruthless and illegal arms business driven by the seduction of money, power, and blood.`,
		Status:       shows.ShowStatusFinished,
		DayOfWeek:    time.Sunday,
		Website:      "http://www.fxnetworks.com/shows/sons-of-anarchy/about",
		BackdropPath: "/fZ8j6F8dxZPA8wE5sGS9oiKzXzM.jpg",
		PosterPath:   "/2qg0MOwPD1G0FcYpDPeu6AOjh8i.jpg",
	})

	testCases := []struct {
		description string
		params      *shows.DeleteParams
		auth        *httptests.RequestAuth
	}{
		{
			"Remove Sons of Anarchy",
			&shows.DeleteParams{ID: s.ID},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			rec := callDelete(t, tc.params, tc.auth)
			require.Equal(t, http.StatusNoContent, rec.Code)
		})
	}
}
