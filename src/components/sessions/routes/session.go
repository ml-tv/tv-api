package routes

import (
	"github.com/ml-tv/tv-api/src/components/sessions/handlers"
	"github.com/ml-tv/tv-api/src/core/router"
)

// Contains the index of all Endpoints
const (
	EndpointAddSession = iota
	EndpointDeleteSession
)

// SessionEndpoints is a list of endpoints for this components
var SessionEndpoints = router.Endpoints{
	EndpointAddSession: {
		Verb:    "POST",
		Path:    "/sessions",
		Handler: handlers.AddSession,
		Params:  &handlers.AddSessionParams{},
	},
	EndpointDeleteSession: {
		Verb:    "DELETE",
		Path:    "/sessions/{token}",
		Handler: handlers.DeleteSession,
		Params:  &handlers.DeleteSessionParams{},
		Auth:    router.LoggedUserAccess,
	},
}
