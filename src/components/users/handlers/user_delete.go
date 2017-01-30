package handlers

import (
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/security/auth"
)

// DeleteUserParams represents the params accepted by the DeleteUser handler
type DeleteUserParams struct {
	ID              string `from:"url" json:"id" params:"uuid"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// DeleteUser represents an API handler to remove a user
func DeleteUser(req *router.Request) error {
	params := req.Params.(*DeleteUserParams)
	user := req.User

	if params.ID != user.ID {
		return httperr.NewForbidden()
	}

	if !auth.IsPasswordValid(user.Password, params.CurrentPassword) {
		return httperr.NewUnauthorized()
	}

	if err := user.Delete(); err != nil {
		return err
	}

	req.NoContent()
	return nil
}
