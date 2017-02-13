package handlers

import (
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/security/auth"
	"github.com/ml-tv/tv-api/src/core/storage/db"
)

// DeleteSessionParams represent the request params accepted by HandlerDelete
type DeleteSessionParams struct {
	Token           string `from:"url" json:"token" params:"uuid"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// DeleteSession represent an API handler to remove a session
func DeleteSession(req *router.Request) error {
	params := req.Params.(*DeleteSessionParams)

	if !auth.IsPasswordValid(req.User.Password, params.CurrentPassword) {
		return httperr.NewUnauthorized()
	}

	var session auth.Session
	stmt := "SELECT * FROM sessions WHERE id=$1 AND deleted_at IS NULL LIMIT 1"
	err := db.Get(&session, stmt, params.Token)
	if err != nil {
		return err
	}

	// We always return a 404 in case of a user error to avoid brute-force
	if session.ID == "" || session.UserID != req.User.ID {
		return httperr.NewNotFound()
	}

	if err := session.Delete(); err != nil {
		return err
	}

	req.NoContent()
	return nil
}
