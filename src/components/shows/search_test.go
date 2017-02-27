package shows_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ml-tv/tv-api/src/core/network/http/httptests"
	"github.com/ml-tv/tv-api/src/core/paginator"
	"github.com/ml-tv/tv-api/src/core/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"net/http"

	"time"

	"fmt"

	"strconv"

	"github.com/ml-tv/tv-api/src/components/shows"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/core/primitives/ptrs"
)

const (
	memoryLostID      = 68386
	lostID            = 4607
	lostGirlID        = 33852
	californicationID = 1215
	batesMotelID      = 46786
	theWalkingDeadID  = 1402
)

func callSearch(t *testing.T, params *shows.SearchParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: shows.Endpoints[shows.EndpointSearch],
		Params:   params,
		Auth:     auth,
	}

	return httptests.NewRequest(t, ri)
}

func TestSearch(t *testing.T) {
	defer lifecycle.PurgeModels(t)
	setupSearchData(t)

	t.Run("parallel", func(t *testing.T) {
		t.Run("No params", searchTestNoParams)
		t.Run("Full text", searchTestFullText)
		t.Run("Filter Day Of Week", searchTestFilterDayOfWeek)
		t.Run("Order By", searchTestOrderBy)
		t.Run("Filter Status", searchTestFilterStatus)
		t.Run("Pagination", searchTestPagination)
	})
}

func searchTestNoParams(t *testing.T) {
	t.Parallel()
	rec := callSearch(t, &shows.SearchParams{}, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	var pld shows.PayloadList
	if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 6, len(pld.Results))
}

func searchTestPagination(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description  string
		totalResults int
		page         *int
		perPage      *int
		orderedID    []int
	}{
		{
			"Default values",
			6,
			nil, nil,
			[]int{theWalkingDeadID, batesMotelID, californicationID, lostGirlID, lostID, memoryLostID},
		},
		{
			"Page 2, default results per page",
			0,
			ptrs.NewInt(2), nil,
			[]int{},
		},
		{
			"Page 1, 5 results per page",
			5,
			ptrs.NewInt(1), ptrs.NewInt(5),
			[]int{theWalkingDeadID, batesMotelID, californicationID, lostGirlID, lostID},
		},
		{
			"Page 2, 5 results per page",
			1,
			ptrs.NewInt(2), ptrs.NewInt(5),
			[]int{memoryLostID},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			pagination := paginator.HandlerParams{Page: tc.page, PerPage: tc.perPage}
			params := &shows.SearchParams{HandlerParams: pagination}

			rec := callSearch(t, params, nil)
			require.Equal(t, http.StatusOK, rec.Code)

			var pld shows.PayloadList
			if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tc.totalResults, len(pld.Results))
			for i, res := range pld.Results {
				assert.Equal(t, tc.orderedID[i], res.TMDbID)
			}
		})
	}
}

func searchTestFilterStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description  string
		totalResults int
		params       *shows.SearchParams
		orderedID    []int
	}{
		{
			"Filter Showing",
			1,
			&shows.SearchParams{Status: strconv.Itoa(shows.ShowStatusShowing)},
			[]int{theWalkingDeadID},
		},
		{
			"Filter Finished",
			4,
			&shows.SearchParams{Status: strconv.Itoa(shows.ShowStatusFinished)},
			[]int{californicationID, lostGirlID, lostID, memoryLostID},
		},
		{
			"Filter Cancelled",
			0,
			&shows.SearchParams{Status: strconv.Itoa(shows.ShowStatusCanceled)},
			[]int{},
		},
		{
			"Filter Showing|Paused",
			2,
			&shows.SearchParams{Status: fmt.Sprintf("%d|%d", shows.ShowStatusShowing, shows.ShowStatusPaused)},
			[]int{theWalkingDeadID, batesMotelID},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			rec := callSearch(t, tc.params, nil)
			require.Equal(t, http.StatusOK, rec.Code)

			var pld shows.PayloadList
			if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tc.totalResults, len(pld.Results))
			for i, res := range pld.Results {
				assert.Equal(t, tc.orderedID[i], res.TMDbID)
			}
		})
	}
}

func searchTestFilterDayOfWeek(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description  string
		totalResults int
		params       *shows.SearchParams
		orderedID    []int
	}{
		{
			"Filter Sunday",
			3,
			&shows.SearchParams{DayOfWeek: ptrs.NewInt(int(time.Sunday))},
			[]int{theWalkingDeadID, lostGirlID, lostID},
		},
		{
			"Filter Monday",
			2,
			&shows.SearchParams{DayOfWeek: ptrs.NewInt(int(time.Monday))},
			[]int{batesMotelID, memoryLostID},
		},
		{
			"Filter Thursday",
			1,
			&shows.SearchParams{DayOfWeek: ptrs.NewInt(int(time.Thursday))},
			[]int{californicationID},
		},
		{
			"Filter Friday",
			0,
			&shows.SearchParams{DayOfWeek: ptrs.NewInt(int(time.Friday))},
			[]int{},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			rec := callSearch(t, tc.params, nil)
			require.Equal(t, http.StatusOK, rec.Code)

			var pld shows.PayloadList
			if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tc.totalResults, len(pld.Results))
			for i, res := range pld.Results {
				assert.Equal(t, tc.orderedID[i], res.TMDbID)
			}
		})
	}
}

func searchTestOrderBy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description  string
		totalResults int
		params       *shows.SearchParams
		orderedID    []int
	}{
		{
			"Order all by name",
			6,
			&shows.SearchParams{OrderBy: "name"},
			[]int{batesMotelID, californicationID, lostID, lostGirlID, memoryLostID, theWalkingDeadID},
		},
		{
			"Order all by -name",
			6,
			&shows.SearchParams{OrderBy: "-name"},
			[]int{theWalkingDeadID, memoryLostID, lostGirlID, lostID, californicationID, batesMotelID},
		},
		{
			"Order all by -day_of_week and name",
			6,
			&shows.SearchParams{OrderBy: "-day_of_week|name"},
			[]int{californicationID, batesMotelID, memoryLostID, lostID, lostGirlID, theWalkingDeadID},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			rec := callSearch(t, tc.params, nil)
			require.Equal(t, http.StatusOK, rec.Code)

			var pld shows.PayloadList
			if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tc.totalResults, len(pld.Results))
			for i, res := range pld.Results {
				assert.Equal(t, tc.orderedID[i], res.TMDbID)
			}
		})
	}
}

func searchTestFullText(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description  string
		totalResults int
		params       *shows.SearchParams
		orderedID    []int
	}{
		{
			"Search for Lost",
			3,
			&shows.SearchParams{Name: "Lost"},
			[]int{lostID, memoryLostID, lostGirlID},
		},
		{
			"Search for lost",
			3,
			&shows.SearchParams{Name: "lost"},
			[]int{lostID, memoryLostID, lostGirlID},
		},
		{
			"Search for girl",
			1,
			&shows.SearchParams{Name: "girl"},
			[]int{lostGirlID},
		},
		{
			"Search for lost girl",
			1,
			&shows.SearchParams{Name: "lost girl"},
			[]int{lostGirlID},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			rec := callSearch(t, tc.params, nil)
			require.Equal(t, http.StatusOK, rec.Code)

			var pld shows.PayloadList
			if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tc.totalResults, len(pld.Results))
			for i, res := range pld.Results {
				assert.Equal(t, tc.orderedID[i], res.TMDbID)
			}
		})
	}
}

func setupSearchData(t *testing.T) {
	NewShow(t, &shows.Show{
		TMDbID:       memoryLostID,
		Name:         "Memory Lost",
		OriginalName: "美人为馅",
		Synopsis:     "A spoiled rich kid lost her identity and her memory after getting kidnapped by a criminal organization. Her boyfriend, who never gave up looking for her, became a police as result. Many years later, a string of criminal activities reunited them, though they did not immediately recognize each other. Although they were now strangers, their complementary skillsets made them formidable partners within a police task force called Black Shield. As they grew closer, her past came back to haunt them.",
		Status:       shows.ShowStatusFinished,
		DayOfWeek:    time.Monday,
		CreatedAt:    &db.Time{Time: time.Now()},
	})

	NewShow(t, &shows.Show{
		TMDbID:       lostID,
		Name:         "Lost",
		OriginalName: "Lost",
		Synopsis:     "Lost is a drama series containing elements of science fiction and the supernatural that follows the survivors of the crash of a commercial passenger jet flying between Sydney and Los Angeles, on a mysterious tropical island somewhere in the South Pacific Ocean.",
		Status:       shows.ShowStatusFinished,
		DayOfWeek:    time.Sunday,
		CreatedAt:    &db.Time{Time: time.Now().Add(10 * time.Second)},
	})

	NewShow(t, &shows.Show{
		TMDbID:       lostGirlID,
		Name:         "Lost Girl",
		OriginalName: "Lost Girl",
		Synopsis:     "Lost Girl focuses on the gorgeous and charismatic Bo, a supernatural being called a succubus who feeds on the energy of humans, sometimes with fatal results. Refusing to embrace her supernatural clan and its rigid hierarchy, Bo is a renegade who takes up the fight for the underdog while searching for the truth about her own mysterious origins.",
		Status:       shows.ShowStatusFinished,
		DayOfWeek:    time.Sunday,
		CreatedAt:    &db.Time{Time: time.Now().Add(20 * time.Second)},
	})

	NewShow(t, &shows.Show{
		TMDbID:       californicationID,
		Name:         "Californication",
		OriginalName: "Californication",
		Synopsis:     "A self-loathing, alcoholic writer attempts to repair his damaged relationships with his daughter and her mother while combating sex addiction, a budding drug problem, and the seeming inability to avoid making bad decisions.",
		Status:       shows.ShowStatusFinished,
		DayOfWeek:    time.Thursday,
		CreatedAt:    &db.Time{Time: time.Now().Add(30 * time.Second)},
	})

	NewShow(t, &shows.Show{
		TMDbID:       batesMotelID,
		Name:         "Bates Motel",
		OriginalName: "Bates Motel",
		Synopsis:     "A contemporary prequel to the 1960 film Psycho, depicting the life of Norman Bates and his mother Norma prior to the events portrayed in Hitchcock's film, albeit in a different fictional town and in a modern setting.",
		Status:       shows.ShowStatusPaused,
		DayOfWeek:    time.Monday,
		CreatedAt:    &db.Time{Time: time.Now().Add(40 * time.Second)},
	})

	NewShow(t, &shows.Show{
		TMDbID:       theWalkingDeadID,
		Name:         "The Walking Dead",
		OriginalName: "The Walking Dead",
		Synopsis:     "Sheriff's deputy Rick Grimes awakens from a coma to find a post-apocalyptic world dominated by flesh-eating zombies. He sets out to find his family and encounters many other survivors along the way.",
		Status:       shows.ShowStatusShowing,
		DayOfWeek:    time.Sunday,
		CreatedAt:    &db.Time{Time: time.Now().Add(50 * time.Second)},
	})
}
