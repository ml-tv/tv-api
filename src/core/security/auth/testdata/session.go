package testdata

import (
	"testing"

	"github.com/ml-tv/tv-api/src/core/security/auth"
	"github.com/ml-tv/tv-api/src/core/tests/testhelpers"
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

	testhelpers.SaveModels(t, session)
	return user, session
}
