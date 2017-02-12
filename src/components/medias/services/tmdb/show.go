package tmdb

import "fmt"

const rootShowURL = rootURL + "/tv"

// Show represents a TMDb payload representing a TV Show
type Show struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	BackdropPath string `json:"backdrop_path"`
	PosterPath   string `json:"poster_path"`
	Overview     string `json:"overview"`
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
