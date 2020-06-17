package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func HttpListenAndServer(port string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/", SituationsMux)
	router.HandleFunc("/{id}", SituationMux)
	server := &http.Server{Addr: port, Handler: router}
	return server
}

// Get all the Services or add one
func SituationsMux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSituations(w, r)
	case http.MethodPost:
		postSituation(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Get/delete a situation or patch it to activate it
func SituationMux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSituation(w, r)
	case http.MethodDelete:
		deleteSituation(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
