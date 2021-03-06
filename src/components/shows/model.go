package shows

import (
	"fmt"

	"time"

	"github.com/ml-tv/tv-api/src/core/storage/db"
	"github.com/ml-tv/tv-api/src/services/tmdb"
)

const (
	ShowStatusShowing = iota
	ShowStatusPaused
	ShowStatusFinished
	ShowStatusCanceled

	// ShowStatusEndOfList represent the total number of possible status
	ShowStatusEndOfList
)

const ShowDateFormat = "January 02 2006"

// Show represents a TV show from the database
//go:generate tv-api-cli generate model Show -t shows
type Show struct {
	ID            string       `db:"id"`
	CreatedAt     *db.Time     `db:"created_at"`
	UpdatedAt     *db.Time     `db:"updated_at"`
	DeletedAt     *db.Time     `db:"deleted_at"`
	Name          string       `db:"name"`
	OriginalName  string       `db:"original_name"`
	YearReleased  int          `db:"year_released"`
	Synopsis      string       `db:"synopsis"`
	PosterPath    string       `db:"poster_path"`
	BackdropPath  string       `db:"backdrop_path"`
	TMDbID        int          `db:"tmdb_id"`
	Status        int          `db:"status"`
	DayOfWeek     time.Weekday `db:"day_of_week"`
	ReturningDate string       `db:"returning_date"`
	Website       string       `db:"website"`
	Wikipedia     string       `db:"wikipedia"`
	ExtraLink     string       `db:"extra_link"`
	IsOnNetflix   bool         `db:"is_on_netflix"`
}

// BackdropURL returns a URL to the backdrop
func (m *Show) BackdropURL() string {
	return fmt.Sprintf("https://image.tmdb.org/t/p/original/%s", m.BackdropPath)
}

// PosterURL returns a URL to the Poster
func (m *Show) PosterURL() string {
	return fmt.Sprintf("https://image.tmdb.org/t/p/original/%s", m.PosterPath)
}

// NewFromTMDb turns a TMDb.Show into a Show
func NewFromTMDb(show *tmdb.Show) (*Show, error) {
	s := &Show{
		TMDbID:       show.ID,
		Name:         show.Name,
		OriginalName: show.OriginalName,
		Synopsis:     show.Overview,
		PosterPath:   show.PosterPath,
		BackdropPath: show.BackdropPath,
		Website:      show.Website,
	}

	firstEpisode, err := time.Parse(tmdb.DateFormat, show.FirstAirDate)
	if err != nil {
		return nil, err
	}
	s.YearReleased = firstEpisode.Year()
	s.DayOfWeek = firstEpisode.Weekday()

	// Set the status
	switch show.Status {
	case tmdb.StatusCanceled:
		s.Status = ShowStatusCanceled
	case tmdb.StatusEnded:
		s.Status = ShowStatusFinished
	default:
		s.Status = ShowStatusShowing
	}

	// Set the returning date, improve the day of the week and the status
	if show.LastAirDate != "" {
		lastEpisode, err := time.Parse(tmdb.DateFormat, show.LastAirDate)
		if err != nil {
			return nil, err
		}
		if !lastEpisode.IsZero() {
			// it's obviously more accurate to use the last air day for the DayOfWeek
			s.DayOfWeek = lastEpisode.Weekday()

			// no need to bother checking the ReturningDate if the show has ben stopped
			if s.Status == ShowStatusShowing {
				// check if lastEpisode is the future
				if lastEpisode.After(time.Now()) {
					s.ReturningDate = lastEpisode.Format(ShowDateFormat)
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
	}

	return s, nil
}

// NewListFromTMDb turns a list of TMDb.Show into a list of Show
func NewListFromTMDb(shows []*tmdb.Show) ([]*Show, error) {
	var output = make([]*Show, len(shows))

	for i, show := range shows {
		s, err := NewFromTMDb(show)
		if err != nil {
			return nil, err
		}
		output[i] = s
	}
	return output, nil
}
