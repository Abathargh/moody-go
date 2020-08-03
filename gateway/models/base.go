package models

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/validator.v2"
)

var (
	NotFound = fmt.Errorf("record Not Found")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// Encode struct to json, if something goes wrong the application exits.
// Used internally to create the body for the http responses.
func MustEncode(w http.ResponseWriter, response interface{}) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}
}

// Decode and validate json data to struct, returning a bool
// containing the outcome of the operation.
func MustValidate(r *http.Request, dest interface{}) (outcome bool) {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		log.Println(err)
		return false
	}
	if err := validator.Validate(dest); err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Read decodes and closes from a io.ReadCloser instance into a reference
// of a given object.
func ReadAndDecode(closer io.ReadCloser, dest interface{}) error {
	defer closeBody(closer)
	body, err := ioutil.ReadAll(closer)
	if err != nil {
		return err
	}
	if jsonErr := json.Unmarshal(body, dest); jsonErr != nil {
		return err
	}
	return nil
}

// Passes the passed error to the response
func RespondWithError(w http.ResponseWriter, status int, error string) {
	w.WriteHeader(status)
	MustEncode(w, ErrorResponse{Error: error})
}

// Helper function to close an io.ReadCloser
func closeBody(rc io.ReadCloser) {
	if err := rc.Close(); err != nil {
		log.Printf("Error while closing the response body: %s", err)
	}
}
