package shows

import (
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/storage/db"
	"github.com/ml-tv/tv-api/src/services/tmdb"
)

// AddParams represents the params needed by the AddShow handler
type AddParams struct {
	TMDbID int `from:"form" json:"tmdb_id" params:"required"`
}

// Add is an API handler to
func Add(req *router.Request) error {
	params := req.Params.(*AddParams)

	// Get the show from TMDb
	show, err := tmdb.ShowByID(params.TMDbID)
	if err != nil {
		switch err {
		case tmdb.ErrNotFound:
			return httperr.NewNotFound()
		default:
			return err
		}
	}

	// Save the show to the database
	s, err := NewFromTMDb(show)
	if err != nil {
		return err
	}
	if err := s.Save(); err != nil {
		if db.IsDup(err) {
			return httperr.NewConflict("show already added")
		}
		return err
	}
	return nil
}
