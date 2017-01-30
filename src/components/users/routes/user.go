package routes

import (
	"github.com/ml-tv/tv-api/src/components/users/handlers"
	"github.com/ml-tv/tv-api/src/core/router"
)

// Contains the index of all Endpoints
const (
	EndpointAddUser = iota
	EndpointUpdateUser
	EndpointDeleteUser
	EndpointGetUser
)

// UserEndpoints is a list of endpoints for this components
var UserEndpoints = router.Endpoints{
	EndpointAddUser: {
		Verb:    "POST",
		Path:    "/users",
		Auth:    nil,
		Handler: handlers.AddUser,
		Params:  &handlers.AddUserParams{},
	},
	EndpointUpdateUser: {
		Verb:    "PATCH",
		Path:    "/users/{id}",
		Auth:    router.LoggedUser,
		Handler: handlers.UpdateUser,
		Params:  &handlers.UpdateUserParams{},
	},
	EndpointDeleteUser: {
		Verb:    "DELETE",
		Path:    "/users/{id}",
		Auth:    router.LoggedUser,
		Handler: handlers.DeleteUser,
		Params:  &handlers.DeleteUserParams{},
	},
	EndpointGetUser: {
		Verb:    "GET",
		Path:    "/users/{id}",
		Auth:    nil,
		Handler: handlers.GetUser,
		Params:  &handlers.GetUserParams{},
	},
}
