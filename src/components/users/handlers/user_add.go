package handlers

import (
	"github.com/ml-tv/tv-api/src/components/users/payloads"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/security/auth"
)

// AddUserParams represents the params needed by the AddUser handler
type AddUserParams struct {
	Name     string `from:"form" json:"name" params:"required,trim"`
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

// AddUser is an HTTP handler to create a new user
func AddUser(req *router.Request) error {
	params := req.Params.(*AddUserParams)

	encryptedPassword, err := auth.CryptPassword(params.Password)
	if err != nil {
		return err
	}

	user := &auth.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: encryptedPassword,
	}

	if err := user.Save(); err != nil {
		return err
	}

	req.Created(payloads.NewFullUser(user))
	return nil
}
