package payloads

import "github.com/ml-tv/tv-api/src/core/security/auth"

// Session represents a Session that can be safely returned by the API
type Session struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

// NewSession turns a Session into an object that is safe to be
// returned by the API
func NewSession(s *auth.Session) *Session {
	return &Session{
		Token:  s.ID,
		UserID: s.UserID,
	}
}
