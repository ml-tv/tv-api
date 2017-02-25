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
	"github.com/ml-tv/tv-api/src/core/primitives/ptrs"
	authdata "github.com/ml-tv/tv-api/src/core/security/auth/testdata"
)

const scrubsID = 4556
const malcolmInTheMiddleID = 2004

func callUpdate(t *testing.T, params *shows.UpdateParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: shows.Endpoints[shows.EndpointUpdate],
		Params:   params,
		Auth:     auth,
	}
	return httptests.NewRequest(t, ri)
}

func TestValidUpdate(t *testing.T) {
	defer lifecycle.PurgeModels(t)
	admin, adminSession := authdata.NewAuth(t)
	admin.IsAdmin = true
	admin.Save()

	// Create and save a show
	s := NewShow(t, &shows.Show{
		TMDbID:       scrubsID,
		Name:         "Scrubs",
		OriginalName: "Scrubs",
		Synopsis:     "In the unreal world of Sacred Heart Hospital, John \"J.D\" Dorian learns the ways of medicine, friendship and life.",
		Status:       shows.ShowStatusCanceled,
		DayOfWeek:    time.Wednesday,
	})

	testCases := []struct {
		description string
		params      *shows.UpdateParams
		auth        *httptests.RequestAuth
	}{
		{
			"Update name",
			&shows.UpdateParams{ID: s.ID, Name: "New name"},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
		{
			"Update original name",
			&shows.UpdateParams{ID: s.ID, OriginalName: "New original name"},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
		{
			"Update Synopsis",
			&shows.UpdateParams{ID: s.ID, Synopsis: "New Synopsis"},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
		{
			"Update Status",
			&shows.UpdateParams{ID: s.ID, Status: ptrs.NewInt(shows.ShowStatusShowing)},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
		{
			"Update DayOfWeek",
			&shows.UpdateParams{ID: s.ID, DayOfWeek: ptrs.NewInt(int(time.Sunday))},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
		{
			"Update Website",
			&shows.UpdateParams{ID: s.ID, Website: ptrs.NewString("http://m.com")},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
		{
			"Update Wikipedia",
			&shows.UpdateParams{ID: s.ID, Wikipedia: ptrs.NewString("http://wiki.org")},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
		{
			"Update ExtraLink",
			&shows.UpdateParams{ID: s.ID, ExtraLink: ptrs.NewString("http://whatever.com")},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
		{
			"Update OnNetflix",
			&shows.UpdateParams{ID: s.ID, IsOnNetflix: ptrs.NewBool(true)},
			httptests.NewRequestAuth(adminSession.ID, admin.ID),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			rec := callUpdate(t, tc.params, tc.auth)
			require.Equal(t, http.StatusOK, rec.Code)

			var pld shows.Payload
			if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
				t.Fatal(err)
			}

			if tc.params.Name != "" {
				assert.Equal(t, tc.params.Name, pld.Name)
			}
			if tc.params.OriginalName != "" {
				assert.Equal(t, tc.params.OriginalName, pld.OriginalName)
			}
			if tc.params.Synopsis != "" {
				assert.Equal(t, tc.params.Synopsis, pld.Synopsis)
			}
			if tc.params.Status != nil {
				assert.Equal(t, *tc.params.Status, pld.Status)
			}
			if tc.params.DayOfWeek != nil {
				assert.Equal(t, *tc.params.DayOfWeek, pld.DayOfWeek)
			}
			if tc.params.Website != nil {
				assert.Equal(t, *tc.params.Website, pld.Website)
			}
			if tc.params.Wikipedia != nil {
				assert.Equal(t, *tc.params.Wikipedia, pld.Wikipedia)
			}
			if tc.params.ExtraLink != nil {
				assert.Equal(t, *tc.params.ExtraLink, pld.ExtraLink)
			}
			if tc.params.IsOnNetflix != nil {
				assert.Equal(t, *tc.params.IsOnNetflix, pld.IsOnNetflix)
			}
		})
	}
}
