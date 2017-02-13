package routes

import (
	"github.com/ml-tv/tv-api/src/components/medias/handlers"
	"github.com/ml-tv/tv-api/src/core/router"
)

// Contains the index of all Endpoints
const (
	EndpointAddShow = iota
	EndpointSearchShow
	EndpointUpdateShow
	EndpointDeleteShow
	EndpointGetShow
)

// ShowEndpoints is a list of endpoints for this components
var ShowEndpoints = router.Endpoints{
	EndpointAddShow: {
		Verb:    "POST",
		Path:    "/medias/shows",
		Auth:    router.AdminAccess,
		Handler: handlers.AddShow,
		Params:  &handlers.AddShowParams{},
	},
	// EndpointSearchShow: {
	// 	Verb:    "GET",
	// 	Path:    "/medias/shows",
	// 	Auth:    nil,
	// 	Handler: handlers.SearchShow,
	// 	Params:  &handlers.SearchShowParams{},
	// },
	// EndpointUpdateShow: {
	// 	Verb:    "PATCH",
	// 	Path:    "/medias/shows/{id}",
	// 	Auth:    router.LoggedUser,
	// 	Handler: handlers.UpdateShow,
	// 	Params:  &handlers.UpdateShowParams{},
	// },
	// EndpointDeleteShow: {
	// 	Verb:    "DELETE",
	// 	Path:    "/medias/shows/{id}",
	// 	Auth:    router.LoggedUser,
	// 	Handler: handlers.DeleteShow,
	// 	Params:  &handlers.DeleteShowParams{},
	// },
	// EndpointGetShow: {
	// 	Verb:    "GET",
	// 	Path:    "/medias/shows/{id}",
	// 	Auth:    nil,
	// 	Handler: handlers.GetShow,
	// 	Params:  &handlers.GetShowParams{},
	// },
}
