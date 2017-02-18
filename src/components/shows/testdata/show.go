package testdata

import (
	"testing"

	"github.com/ml-tv/tv-api/src/components/shows"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/services/tmdb"
)

// NewShow returns a show that can be use for testing
func NewShow(t *testing.T, id int) *shows.Show {
	show, err := tmdb.ShowByID(id)
	if err != nil {
		t.Fatal(err)
	}

	// Save the show to the database
	s := shows.NewFromTMDb(show)
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, s)
	return s
}
