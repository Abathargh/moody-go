package communication

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func MustEncode(w http.ResponseWriter, response interface{}) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func HttpListenAndServe() {
	router := mux.NewRouter()
	router.HandleFunc("/situation", SituationsMux)
	router.HandleFunc("/situation/{name}", SituationMux)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":6666", nil))
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
