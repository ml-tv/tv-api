package shows

// Code auto-generated; DO NOT EDIT

import (
	"errors"
	"fmt"

	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/storage/db"
	uuid "github.com/satori/go.uuid"
)

// JoinShowSQL returns a string ready to be embed in a JOIN query
func JoinShowSQL(prefix string) string {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "name", "original_name", "synopsis", "poster_path", "backdrop_path", "tmdb_id", "status", "day_of_week", "returning_date", "website", "wikipedia", "extra_link", "is_on_netflix" }
	output := ""

	for i, field := range fields {
		if i != 0 {
			output += ", "
		}

		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\"", fullName, fullName)
	}

	return output
}

// GetShow finds and returns an active show by ID
func GetShow(id string) (*Show, error) {
	s := &Show{}
	stmt := "SELECT * from shows WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := db.Get(s, stmt, id)
	// We want to return nil if a show is not found
	if s.ID == "" {
		return nil, err
	}
	return s, err
}

// ShowExists checks if a show exists for a specific ID
func ShowExists(id string) (bool, error) {
	exists := false
	stmt := "SELECT exists(SELECT 1 FROM shows WHERE id=$1 and deleted_at IS NULL)"
	err := db.Writer.Get(&exists, stmt, id)
	return exists, err
}

// Save creates or updates the show depending on the value of the id
func (s *Show) Save() error {
	return s.SaveQ(db.Writer)
}

// SaveQ creates or updates the article depending on the value of the id using
// a transaction
func (s *Show) SaveQ(q db.Queryable) error {
	if s == nil {
		return httperr.NewServerError("show is not instanced")
	}

	if s.ID == "" {
		return s.CreateQ(q)
	}

	return s.UpdateQ(q)
}

// Create persists a show in the database
func (s *Show) Create() error {
	return s.CreateQ(db.Writer)
}

// Create persists a show in the database
func (s *Show) CreateQ(q db.Queryable) error {
	if s == nil {
		return httperr.NewServerError("show is not instanced")
	}

	if s.ID != "" {
		return httperr.NewServerError("cannot persist a show that already has an ID")
	}

	return s.doCreate(q)
}

// doCreate persists a show in the database using a Node
func (s *Show) doCreate(q db.Queryable) error {
	if s == nil {
		return errors.New("show not instanced")
	}

	s.ID = uuid.NewV4().String()
	s.CreatedAt = db.Now()
	s.UpdatedAt = db.Now()

	stmt := "INSERT INTO shows (id, created_at, updated_at, deleted_at, name, original_name, synopsis, poster_path, backdrop_path, tmdb_id, status, day_of_week, returning_date, website, wikipedia, extra_link, is_on_netflix) VALUES (:id, :created_at, :updated_at, :deleted_at, :name, :original_name, :synopsis, :poster_path, :backdrop_path, :tmdb_id, :status, :day_of_week, :returning_date, :website, :wikipedia, :extra_link, :is_on_netflix)"
	_, err := q.NamedExec(stmt, s)

  return err
}

// Update updates most of the fields of a persisted show.
// Excluded fields are id, created_at, deleted_at, etc.
func (s *Show) Update() error {
	return s.UpdateQ(db.Writer)
}

// Update updates most of the fields of a persisted show using a transaction
// Excluded fields are id, created_at, deleted_at, etc.
func (s *Show) UpdateQ(q db.Queryable) error {
	if s == nil {
		return httperr.NewServerError("show is not instanced")
	}

	if s.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted show")
	}

	return s.doUpdate(q)
}

// doUpdate updates a show in the database using an optional transaction
func (s *Show) doUpdate(q db.Queryable) error {
	if s == nil {
		return httperr.NewServerError("show is not instanced")
	}

	if s.ID == "" {
		return httperr.NewServerError("cannot update a non-persisted show")
	}

	s.UpdatedAt = db.Now()

	stmt := "UPDATE shows SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, name=:name, original_name=:original_name, synopsis=:synopsis, poster_path=:poster_path, backdrop_path=:backdrop_path, tmdb_id=:tmdb_id, status=:status, day_of_week=:day_of_week, returning_date=:returning_date, website=:website, wikipedia=:wikipedia, extra_link=:extra_link, is_on_netflix=:is_on_netflix WHERE id=:id"
	_, err := q.NamedExec(stmt, s)

	return err
}

// FullyDelete removes a show from the database
func (s *Show) FullyDelete() error {
	return s.FullyDeleteQ(db.Writer)
}

// FullyDeleteQ removes a show from the database using a transaction
func (s *Show) FullyDeleteQ(q db.Queryable) error {
	if s == nil {
		return errors.New("show not instanced")
	}

	if s.ID == "" {
		return errors.New("show has not been saved")
	}

	stmt := "DELETE FROM shows WHERE id=$1"
	_, err := q.Exec(stmt, s.ID)

	return err
}

// Delete soft delete a show.
func (s *Show) Delete() error {
	return s.DeleteQ(db.Writer)
}

// DeleteQ soft delete a show using a transaction
func (s *Show) DeleteQ(q db.Queryable) error {
	return s.doDelete(q)
}

// doDelete performs a soft delete operation on a show using an optional transaction
func (s *Show) doDelete(q db.Queryable) error {
	if s == nil {
		return httperr.NewServerError("show is not instanced")
	}

	if s.ID == "" {
		return httperr.NewServerError("cannot delete a non-persisted show")
	}

	s.DeletedAt = db.Now()

	stmt := "UPDATE shows SET deleted_at = $2 WHERE id=$1"
	_, err := q.Exec(stmt, s.ID, s.DeletedAt)
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (s *Show) IsZero() bool {
	return s == nil || s.ID == ""
}