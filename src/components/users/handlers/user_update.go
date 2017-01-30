package handlers

import (
	"github.com/ml-tv/tv-api/src/components/users/payloads"
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/security/auth"
)

// UpdateUserParams represents the params the UpdateUser handler accepts
type UpdateUserParams struct {
	ID              string `from:"url" json:"id"  params:"uuid"`
	Name            string `from:"form" json:"name" params:"trim"`
	Email           string `from:"form" json:"email" params:"trim"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
	NewPassword     string `from:"form" json:"new_password" params:"trim"`
}

// UpdateUser represents an API handler to update a user
func UpdateUser(req *router.Request) error {
	params := req.Params.(*UpdateUserParams)
	user := req.User

	if params.ID != user.ID {
		return httperr.NewForbidden()
	}

	// To change the email or the password we require the current password
	if params.NewPassword != "" || params.Email != "" {
		if !auth.IsPasswordValid(user.Password, params.CurrentPassword) {
			return httperr.NewUnauthorized()
		}
	}

	if params.Name != "" {
		user.Name = params.Name
	}

	if params.Email != "" {
		user.Email = params.Email
	}

	if params.NewPassword != "" {
		hashedPassword, err := auth.CryptPassword(params.NewPassword)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	if err := user.Save(); err != nil {
		return err
	}

	req.Ok(payloads.NewFullUser(user))
	return nil
}
