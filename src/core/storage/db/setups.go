package db

import "github.com/jmoiron/sqlx"

// Setup setup the database connection and init the Writer
func Setup(uri string) error {
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return err
	}

	Writer = db
	return nil
}
