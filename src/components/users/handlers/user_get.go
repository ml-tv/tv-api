package handlers

import (
	"github.com/ml-tv/tv-api/src/components/users/payloads"
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/security/auth"
)

// GetUserParams represents the params accepted by the GetUser handler
type GetUserParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// GetUser represents an API handler to get a user
func GetUser(req *router.Request) error {
	params := req.Params.(*GetUserParams)

	user, err := auth.GetUser(params.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return httperr.NewNotFound()
	}

	// if a user asks for their own data, we return as much as possible
	if req.User != nil && req.User.ID == user.ID {
		req.Ok(payloads.NewFullUser(user))
		return nil
	}
	req.Ok(payloads.NewUser(user))
	return nil
}
