package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"situation/models"
	"strconv"
)

type situationsResponse struct {
	Situations     []*models.Situation `json:"situations"`
	SituationCount int64               `json:"count"`
}

// Get handler for /situations
func getSituations(w http.ResponseWriter, _ *http.Request) {
	situations, count, err := models.GetAllSituations()
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
	newSituation := &models.Situation{}
	w.Header().Set("Content-Type", "application/json")
	ok := MustValidate(r, newSituation)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}
	if err := models.AddSituation(newSituation); err != nil {
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
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}
	situation, err := models.GetSituation(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		MustEncode(w, ErrorResponse{Error: "not found"})
		return
	}
	// here the response is the situation itself
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	MustEncode(w, situation)
}

// Delete handler for /situation/{name}
func deleteSituation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}
	situation, err := models.GetSituation(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		MustEncode(w, ErrorResponse{"The situation does not exist"})
	}

	if err := models.DeleteSituation(situation); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"An error occurred while trying to delete the situation"})
	}
	w.WriteHeader(http.StatusOK)
	MustEncode(w, situation)
}
