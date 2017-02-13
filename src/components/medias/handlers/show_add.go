package handlers

import (
	"github.com/ml-tv/tv-api/src/components/medias/models"
	"github.com/ml-tv/tv-api/src/components/medias/services/tmdb"
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/storage/db"
)

// AddShowParams represents the params needed by the AddShow handler
type AddShowParams struct {
	TMDbID int `from:"form" json:"tmdb_id" params:"required"`
}

// AddShow is an API handler to
func AddShow(req *router.Request) error {
	params := req.Params.(*AddShowParams)

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
	s := models.NewShowFromTMDb(show)
	if err := s.Save(); err != nil {
		if db.IsDup(err) {
			return httperr.NewConflict("show already added")
		}
		return err
	}
	return nil
}
