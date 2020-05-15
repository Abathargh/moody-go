package communication

import (
	"encoding/json"
	"github.com/Abathargh/moody-go/db"
	"github.com/Abathargh/moody-go/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type situationsResponse struct {
	Situations     []*models.Situation `json:"situations"`
	SituationCount int64               `json:"count"`
}

// Get handler for /situations
func getSituations(w http.ResponseWriter, _ *http.Request) {
	situations, count, err := db.GetAllSituations()
	if err != nil {
		situations = []*models.Situation{}
	}
	resp := situationsResponse{
		Situations:     situations,
		SituationCount: count,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	MustEncode(w, resp)
}

// Post handler for /situations
func postSituation(w http.ResponseWriter, r *http.Request) {
	var newSituation models.Situation
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&newSituation); err != nil || newSituation.SituationName == "" {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}

	if err := db.AddSituation(&newSituation); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"Record with pk already exists"})
		return
	}
	w.WriteHeader(http.StatusOK)
	MustEncode(w, newSituation)
	if err := r.Body.Close(); err != nil {
		log.Println(err)
	}
}

// Get handler for /situation/{name}
func getSituation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	situation, err := db.GetSituation(name)
	if err != nil {
		situation = &models.Situation{}
	}
	// here the response is the situation itself
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	MustEncode(w, situation)
}

// Delete handler for /situation/{name}
func deleteSituation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	w.Header().Set("Content-Type", "application/json")

	situation, err := db.GetSituation(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		MustEncode(w, ErrorResponse{"The situation does not exist"})
	}

	// TODO remove all situation settings and test

	if err := db.DeleteSituation(situation); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"An error occurred while trying to delete the situation"})

	}
}
