package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ml-tv/tv-api/src/components/users"
	"github.com/ml-tv/tv-api/src/core/network/http/httpres"
)

func notFound(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := fmt.Sprintf(`{"error":"%s"}`, http.StatusText(http.StatusNotFound))
	httpres.ErrorJSON(w, err, http.StatusNotFound)
}

// GetRouter return the api router with all the routes
func GetRouter() *mux.Router {
	r := mux.NewRouter()
	users.SetRoutes(r)
	// blog.SetRoutes(baseURI, r)
	// sessions.SetRoutes(baseURI, r)
	r.NotFoundHandler = http.HandlerFunc(notFound)

	return r
}
