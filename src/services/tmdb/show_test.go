package tmdb_test

import (
	"testing"

	"github.com/ml-tv/tv-api/src/services/tmdb"
	"github.com/stretchr/testify/assert"
)

func TestValidShowByID(t *testing.T) {
	if tmdb.APIKey == "" {
		t.Skip("TMDb API key not defined")
	}

	show, err := tmdb.ShowByID(4607)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 4607, show.ID)
	assert.Equal(t, "Lost", show.Name)
	assert.Equal(t, "Lost", show.OriginalName)
	assert.Equal(t, "/pRPJj3oczSIxo8L80qNweXaWnrp.jpg", show.BackdropPath)
	assert.Equal(t, "/jyGspygDXJMydTOJj7iWNx9Elyd.jpg", show.PosterPath)
	assert.NotEmpty(t, show.Overview)
}

func TestInvalidShowByID(t *testing.T) {
	if tmdb.APIKey == "" {
		t.Skip("TMDb API key not defined")
	}

	testCases := []struct {
		description string
		id          int
		err         error
	}{
		{"Un-existing show", 4200000, tmdb.ErrNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			show, err := tmdb.ShowByID(tc.id)

			assert.Equal(t, tc.err, err)
			assert.Nil(t, show)
		})
	}
}
