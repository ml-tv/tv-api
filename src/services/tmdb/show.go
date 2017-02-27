package tmdb

import "fmt"

const rootShowURL = rootURL + "/tv"

var (
	StatusReturning = "Returning Series"
	StatusEnded     = "Ended"
	StatusCanceled  = "Canceled"
)

// Show represents a TMDb payload representing a TV Show
type Show struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	BackdropPath string `json:"backdrop_path"`
	PosterPath   string `json:"poster_path"`
	Overview     string `json:"overview"`
	InProduction bool   `json:"in_production"`
	Status       string `json:"status"`
	Website      string `json:"homepage"`
	LastAirDate  string `json:"last_air_date"`
	FirstAirDate string `json:"first_air_date"`
}

// ShowByID returns a show using an ID
func ShowByID(id int) (*Show, error) {
	endpoint := fmt.Sprintf("%s/%d", rootShowURL, id)

	var dest Show
	err := Get(endpoint, &dest)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}
