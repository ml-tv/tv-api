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
	EndpointDelete
	EndpointGet
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
	// EndpointUpdate: {
	// 	Verb:    "PATCH",
	// 	Path:    "/medias/shows/{id}",
	// 	Auth:    router.LoggedUser,
	// 	Handler: handlers.Update,
	// 	Params:  &handlers.UpdateParams{},
	// },
	// EndpointDelete: {
	// 	Verb:    "DELETE",
	// 	Path:    "/medias/shows/{id}",
	// 	Auth:    router.LoggedUser,
	// 	Handler: handlers.Delete,
	// 	Params:  &handlers.DeleteParams{},
	// },
	// EndpointGet: {
	// 	Verb:    "GET",
	// 	Path:    "/medias/shows/{id}",
	// 	Auth:    nil,
	// 	Handler: handlers.Get,
	// 	Params:  &handlers.GetParams{},
	// },
}

// SetRoutes is used to set all the routes of the article
func SetRoutes(r *mux.Router) {
	Endpoints.Activate(r)
}
