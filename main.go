package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/ml-tv/tv-api/src/components/api"
)

func main() {
	args := api.Setup()

	r := api.GetRouter()
	port := ":" + args.Port

	handler := handlers.CORS(api.AllowedOrigins, api.AllowedMethods, api.AllowedHeaders)
	http.ListenAndServe(port, handler(r))
}
