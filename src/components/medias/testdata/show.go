package testdata

import (
	"testing"

	"github.com/ml-tv/tv-api/src/components/medias/models"
	"github.com/ml-tv/tv-api/src/components/medias/services/tmdb"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
)

// NewShow returns a show that can be use for testing
func NewShow(t *testing.T, id int) *models.Show {
	show, err := tmdb.ShowByID(id)
	if err != nil {
		t.Fatal(err)
	}

	// Save the show to the database
	s := models.NewShowFromTMDb(show)
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, s)
	return s
}
