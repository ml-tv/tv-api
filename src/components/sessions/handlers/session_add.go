package handlers

import (
	"github.com/ml-tv/tv-api/src/components/sessions/payloads"
	"github.com/ml-tv/tv-api/src/core/network/http/httperr"
	"github.com/ml-tv/tv-api/src/core/router"
	"github.com/ml-tv/tv-api/src/core/security/auth"
	"github.com/ml-tv/tv-api/src/core/storage/db"
)

// AddSessionParams represent the request params accepted by HandlerAdd
type AddSessionParams struct {
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

// AddSession represents an API handler to create a new user session
func AddSession(req *router.Request) error {
	params := req.Params.(*AddSessionParams)

	var user auth.User
	stmt := "SELECT * FROM users WHERE email=$1 LIMIT 1"
	err := db.Get(&user, stmt, params.Email)
	if err != nil {
		return err
	}

	if user.ID == "" || !auth.IsPasswordValid(user.Password, params.Password) {
		return httperr.NewBadRequest("Bad email/password")
	}

	s := &auth.Session{
		UserID: user.ID,
	}
	if err := s.Save(); err != nil {
		return err
	}

	req.Created(payloads.NewSession(s))
	return nil
}
