package communication

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

const (
	httpServerPort = ":6666"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func MustEncode(w http.ResponseWriter, response interface{}) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

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

func HttpListenAndServer() *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/situation", SituationsMux)
	router.HandleFunc("/situation/{id}", SituationMux)
	router.HandleFunc("/service", ServicesMux)
	router.HandleFunc("/service/{id}", ServiceMux)
	server := &http.Server{Addr: httpServerPort, Handler: router}
	//http.Handle("/", router)
	return server
}

// Get all the Situations or add one
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

// Get or delete a situation
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
	case http.MethodPatch:
		patchService(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
