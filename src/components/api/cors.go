package api

import "github.com/gorilla/handlers"

// AllowedOrigins is a list containing all origins allowed to hit the API
var AllowedOrigins = handlers.AllowedOrigins([]string{
	"http://www.melvin.la",
	"http://orchid.melvin.la",
	"http://orchid.melvin.la:3000",
})

// AllowedMethods is a list containing all HTTP verb accepted by the API
var AllowedMethods = handlers.AllowedMethods([]string{
	"GET", "POST", "PATCH", "DELETE",
})

// AllowedHeaders is a list custom headers accepted by the API
var AllowedHeaders = handlers.AllowedHeaders([]string{
	"Content-Type", "Authorization",
})
