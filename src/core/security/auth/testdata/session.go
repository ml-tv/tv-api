package testdata

import (
	"testing"

	"github.com/ml-tv/tv-api/src/core/primitives/models/lifecycle"
	"github.com/ml-tv/tv-api/src/core/security/auth"
)

// NewAuth creates a new user and their session
func NewAuth(t *testing.T) (*auth.User, *auth.Session) {
	user := NewUser(t, nil)
	session := &auth.Session{
		UserID: user.ID,
	}

	if err := session.Create(); err != nil {
		t.Fatal(err)
	}

	lifecycle.SaveModels(t, session)
	return user, session
}
