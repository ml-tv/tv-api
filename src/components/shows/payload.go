package shows

// Payload represents a TV Show returnable to the clients
type Payload struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	OriginalName  string `json:"original_name"`
	Synopsis      string `json:"synopsis"`
	Poster        string `json:"poster"`
	Backdrop      string `json:"backdrop"`
	TMDbID        int    `json:"tmdb_id"`
	Status        int    `json:"status"`
	DayOfWeek     int    `json:"day_of_week"`
	ReturningDate string `json:"returning_date,omitempty"`
	Website       string `json:"website,omitempty"`
	Wikipedia     string `json:"wikipedia"`
	ExtraLink     string `json:"extra_link,omitempty"`
	OnNetflix     bool   `json:"on_netflix"`
}

// PayloadList represents a list of TV Show that can be returned to the clients
type PayloadList struct {
	Results []*Payload
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
		DayOfWeek:     int(m.DayOfWeek),
		Status:        m.Status,
		ReturningDate: m.ReturningDate,
		Website:       m.Website,
		Wikipedia:     m.Wikipedia,
		ExtraLink:     m.ExtraLink,
		OnNetflix:     m.OnNetflix,
	}
}

// NewPayloadList turns a list of Shows into a payload
func NewPayloadList(list []*Show) *PayloadList {
	pld := &PayloadList{}
	pld.Results = make([]*Payload, len(list))
	for i, show := range list {
		pld.Results[i] = NewPayload(show)
	}
	return pld
}
