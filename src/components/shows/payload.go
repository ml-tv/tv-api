package shows

import (
	"github.com/ml-tv/tv-api/src/core/storage/db"
)

// Payload represents a TV Show returnable to the clients
type Payload struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	OriginalName  string   `json:"original_name"`
	Synopsis      string   `json:"synopsis"`
	Poster        string   `json:"poster"`
	Backdrop      string   `json:"backdrop"`
	TMDbID        int      `json:"tmdb_id"`
	Status        int      `json:"status"`
	ReturningDate *db.Time `json:"returning_date,omitempty"`
	Website       string   `json:"website,omitempty"`
	Wikipedia     string   `json:"wikipedia"`
	ExtraLink     string   `json:"extra_link,omitempty"`
	OnNetflix     bool     `json:"on_netflix"`
}

// NewPayload turn a Show model into a payload
func NewPayload(m *Show) *Payload {
	return &Payload{
		ID:            m.ID,
		Name:          m.Name,
		OriginalName:  m.OriginalName,
		Synopsis:      m.Synopsis,
		Poster:        m.PosterURL(),
		Backdrop:      m.BackdropURL(),
		TMDbID:        m.TMDbID,
		Status:        m.Status,
		ReturningDate: m.ReturningDate,
		Website:       m.Website,
		Wikipedia:     m.Wikipedia,
		ExtraLink:     m.ExtraLink,
		OnNetflix:     m.OnNetflix,
	}
}
