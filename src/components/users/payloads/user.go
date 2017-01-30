package payloads

import "github.com/ml-tv/tv-api/src/core/security/auth"

// User represents a user payload with non public field
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

// NewFullUser turns a user into an object that is safe to be
// returned by the API
func NewFullUser(u *auth.User) *User {
	pld := NewUser(u)
	pld.Email = u.Email

	return pld
}

// NewUser turns a user into an object that is safe to be
// returned by the API
func NewUser(u *auth.User) *User {
	return &User{
		ID:   u.ID,
		Name: u.Name,
	}
}
