package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func HttpListenAndServer(port string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/", ServicesMux)
	router.HandleFunc("/{id}", ServiceMux)
	server := &http.Server{Addr: port, Handler: router}
	return server
}

// Get all the Services or add one
func ServicesMux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getServices(w, r)
	case http.MethodPost:
		postService(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Get/delete a situation or patch it to activate it
func ServiceMux(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getService(w, r)
	case http.MethodDelete:
		deleteService(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
