package shows

import (
	"github.com/gorilla/mux"
	"github.com/ml-tv/tv-api/src/core/router"
)

// Contains the index of all Endpoints
const (
	EndpointAdd = iota
	EndpointSearch
	EndpointUpdate
	EndpointGetOne
	EndpointDelete
)

// Endpoints is a list of endpoints for this components
var Endpoints = router.Endpoints{
	EndpointAdd: {
		Verb:    "POST",
		Path:    "/shows",
		Auth:    router.AdminAccess,
		Handler: Add,
		Params:  &AddParams{},
	},
	EndpointSearch: {
		Verb:    "GET",
		Path:    "/shows",
		Auth:    nil,
		Handler: Search,
		Params:  &SearchParams{},
	},
	EndpointUpdate: {
		Verb:    "PATCH",
		Path:    "/shows/{id}",
		Auth:    router.AdminAccess,
		Handler: Update,
		Params:  &UpdateParams{},
	},
	EndpointGetOne: {
		Verb:    "GET",
		Path:    "/shows/{id}",
		Auth:    nil,
		Handler: GetOne,
		Params:  &GetOneParams{},
	},
	EndpointDelete: {
		Verb:    "DELETE",
		Path:    "/shows/{id}",
		Auth:    router.AdminAccess,
		Handler: Delete,
		Params:  &DeleteParams{},
	},
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
