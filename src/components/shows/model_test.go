package shows_test

import (
	"testing"

	"time"

	"github.com/ml-tv/tv-api/src/components/shows"
	"github.com/ml-tv/tv-api/src/services/tmdb"
	"github.com/stretchr/testify/assert"
)

func TestNewFromTMDbWithFinishedShow(t *testing.T) {
	s := &tmdb.Show{
		ID:           1396,
		Name:         "Breaking Bad",
		OriginalName: "Breaking Bad",
		Overview:     "Breaking Bad is the story of Walter White, a struggling high school chemistry teacher who is diagnosed with inoperable lung cancer at the beginning of the series. He turns to a life of crime, producing and selling methamphetamine, in order to secure his family's financial future before he dies, teaming with his former student, Jesse Pinkman",
		BackdropPath: "/eSzpy96DwBujGFj0xMbXBcGcfxX.jpg",
		PosterPath:   "/1yeVJox3rjo2jBKrrihIMj7uoS9.jpg",
		InProduction: false,
		Status:       "Ended",
		Website:      "http://www.amctv.com/shows/breaking-bad",
		LastAirDate:  "2013-09-29",
		FirstAirDate: "2008-01-19",
	}

	show, err := shows.NewFromTMDb(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2008, show.YearReleased)
	assert.Equal(t, shows.ShowStatusFinished, show.Status)
	assert.Equal(t, time.Sunday, show.DayOfWeek)
	assert.Equal(t, "", show.ReturningDate)

	assert.Equal(t, s.Name, show.Name)
	assert.Equal(t, s.OriginalName, show.OriginalName)
	assert.Equal(t, s.ID, show.TMDbID)
	assert.Equal(t, s.Overview, show.Synopsis)
	assert.Equal(t, s.BackdropPath, show.BackdropPath)
	assert.Equal(t, s.PosterPath, show.PosterPath)
	assert.Equal(t, s.Website, show.Website)
}

func TestNewFromTMDbWithPausedShow(t *testing.T) {
	in8Days := time.Now().Add(24 * 9 * time.Hour)

	s := &tmdb.Show{
		ID:           1396,
		Name:         "The Walking Dead",
		OriginalName: "The Walking Dead",
		Overview:     "Sheriff's deputy Rick Grimes awakens from a coma to find a post-apocalyptic world dominated by flesh-eating zombies. He sets out to find his family and encounters many other survivors along the way.",
		BackdropPath: "/eSzpy96DwBujGFj0xMbXBcGcfxX.jpg",
		PosterPath:   "/vxuoMW6YBt6UsxvMfRNwRl9LtWS.jpg",
		InProduction: true,
		Status:       "Returning Series",
		Website:      "http://www.amctv.com/shows/the-walking-dead/",
		LastAirDate:  in8Days.Format(tmdb.DateFormat),
		FirstAirDate: "2010-10-31",
	}

	show, err := shows.NewFromTMDb(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2010, show.YearReleased)
	assert.Equal(t, shows.ShowStatusPaused, show.Status)
	assert.Equal(t, in8Days.Weekday(), show.DayOfWeek)

	assert.Equal(t, s.Name, show.Name)
	assert.Equal(t, s.OriginalName, show.OriginalName)
	assert.Equal(t, s.ID, show.TMDbID)
	assert.Equal(t, s.Overview, show.Synopsis)
	assert.Equal(t, s.BackdropPath, show.BackdropPath)
	assert.Equal(t, s.PosterPath, show.PosterPath)
	assert.Equal(t, s.Website, show.Website)
}
