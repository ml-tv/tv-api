package tmdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// APIKey is the TMDb API key that will be used to query the TMDb API
var APIKey string

const rootURL = "https://api.themoviedb.org/3"

var (
	// ErrNotFound is an error returned when TMDb cannot find a resource
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidKey is an error returned when TMDb cannot validate the API key
	ErrInvalidKey = errors.New("invalid API key")
)

type apiError struct {
	Code          int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// Get is used to trigger a GET request against the TMDb API
func Get(endpoint string, dest interface{}) error {
	// Setup the URl
	e, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	qs := e.Query()
	qs.Add("language", "en-US")
	qs.Add("api_key", APIKey)
	e.RawQuery = qs.Encode()

	// Make the request
	reqBody := bytes.NewBufferString("")
	req, err := http.NewRequest("GET", e.String(), reqBody)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Parse the response
	switch res.StatusCode {
	case 200:
		return json.NewDecoder(res.Body).Decode(dest)
	case 401:
		return ErrInvalidKey
	case 404:
		return ErrNotFound
	default:
		var fullError apiError
		err := json.NewDecoder(res.Body).Decode(&fullError)
		if err != nil {
			return fmt.Errorf("unknown error with code %d: %s", res.StatusCode, err)
		}
		return fmt.Errorf("unknown error with code %d: %s", res.StatusCode, fullError.StatusMessage)
	}
}
