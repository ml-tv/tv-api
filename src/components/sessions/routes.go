package sessions

import (
	"github.com/gorilla/mux"
	"github.com/ml-tv/tv-api/src/components/sessions/routes"
)

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	routes.SessionEndpoints.Activate(r)
}
