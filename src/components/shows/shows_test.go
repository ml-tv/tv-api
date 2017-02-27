package shows_test

import (
	"testing"

	"github.com/ml-tv/tv-api/src/components/api"
	"github.com/ml-tv/tv-api/src/components/shows"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/services/tmdb"
)

func init() {
	api.SetupIfNeeded()
}

// NewShowFromTMDb returns a show that can be use for testing
func NewShowFromTMDb(t *testing.T, id int) *shows.Show {
	show, err := tmdb.ShowByID(id)
	if err != nil {
		t.Fatal(err)
	}

	// Save the show to the database
	s, err := shows.NewFromTMDb(show)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, s)
	return s
}

// NewShow returns a show that can be use for testing
func NewShow(t *testing.T, s *shows.Show) *shows.Show {
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}
	lifecycle.SaveModels(t, s)
	return s
}
