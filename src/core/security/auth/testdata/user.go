package testdata

import (
	"fmt"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/core/security/auth"
)

// NewUser creates a new user with "fake" as password
func NewUser(t *testing.T, u *auth.User) *auth.User {
	if u == nil {
		u = &auth.User{}
	}

	if u.Email == "" {
		u.Email = fmt.Sprintf("fake+%s@melvin.la", uniuri.New())
	}

	if u.Name == "" {
		u.Name = "Fake Account"
	}

	if u.Password == "" {
		var err error
		u.Password, err = auth.CryptPassword("fake")
		if err != nil {
			t.Fatalf("failed to create password: %s", err)
		}
	}

	if err := u.Create(); err != nil {
		t.Fatalf("failed to create user: %s", err)
	}

	lifecycle.SaveModels(t, u)
	return u
}
