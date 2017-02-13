package medias

import (
	"github.com/gorilla/mux"
	"github.com/ml-tv/tv-api/src/components/medias/routes"
)

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	routes.ShowEndpoints.Activate(r)
}
