package models

import (
	"fmt"

	"time"

	"github.com/ml-tv/tv-api/src/components/medias/services/tmdb"
	"github.com/ml-tv/tv-api/src/core/storage/db"
)

const (
	ShowStatusShowing = iota
	ShowStatusPaused
	ShowStatusFinished
	ShowStatusCanceled
)

// Show represents a TV show from the database
//go:generate tv-api-cli generate model Show -t shows
type Show struct {
	ID            string       `db:"id"`
	CreatedAt     *db.Time     `db:"created_at"`
	UpdatedAt     *db.Time     `db:"updated_at"`
	DeletedAt     *db.Time     `db:"deleted_at"`
	Name          string       `db:"name"`
	OriginalName  string       `db:"original_name"`
	Synopsis      string       `db:"synopsis"`
	PosterPath    string       `db:"poster_path"`
	BackdropPath  string       `db:"backdrop_path"`
	TMDbID        int          `db:"tmdb_id"`
	Status        int          `db:"status"`
	DayOfWeek     time.Weekday `db:"day_of_week"`
	ReturningDate *db.Time     `db:"returning_date"`
	Website       string       `db:"website"`
	Wikipedia     string       `db:"wikipedia"`
	ExtraLink     string       `db:"extra_link"`
	OnNetflix     bool         `db:"on_netflix"`
}

// BackdropURL returns a URL to the backdrop
func (m *Show) BackdropURL() string {
	return fmt.Sprintf("https://image.tmdb.org/t/p/original/%s", m.BackdropPath)
}

// PosterURL returns a URL to the Poster
func (m *Show) PosterURL() string {
	return fmt.Sprintf("https://image.tmdb.org/t/p/original/%s", m.PosterPath)
}

// NewShowFromTMDb turns a tmdb.Show int a Show
func NewShowFromTMDb(show *tmdb.Show) *Show {
	s := &Show{
		TMDbID:       show.ID,
		Name:         show.Name,
		OriginalName: show.OriginalName,
		Synopsis:     show.Overview,
		PosterPath:   show.PosterPath,
		BackdropPath: show.BackdropPath,
		Website:      show.Website,
	}

	// Set the status, the returning date, and the day of the week
	switch show.Status {
	case tmdb.StatusCanceled:
		s.Status = ShowStatusCanceled
	case tmdb.StatusEnded:
		s.Status = ShowStatusFinished
	default:
		s.Status = ShowStatusShowing

		lastEpisode, err := time.Parse(tmdb.DateFormat, show.LastAirDate)
		if err != nil && !lastEpisode.IsZero() {
			s.DayOfWeek = lastEpisode.Weekday()

			// check if lastEpisode is the future
			if lastEpisode.After(time.Now()) {
				s.ReturningDate = &db.Time{Time: lastEpisode}
			}

			// If the next episode is in over a week or was more than a week ago,
			// the show is on break
			oneWeek := time.Duration(7 * 24 * time.Hour)
			nextWeek := time.Now().Add(oneWeek)
			lastWeek := time.Now().Add(-oneWeek)
			if lastEpisode.Before(lastWeek) || lastEpisode.After(nextWeek) {
				s.Status = ShowStatusPaused
			}
		}
	}

	return s
}
