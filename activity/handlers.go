package main

import (
	"activity/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type situationsResponse struct {
	Situations     []*models.Situation `json:"situations"`
	SituationCount int64               `json:"count"`
}

type servicesResponse struct {
	Services     []*models.Service `json:"services"`
	ServiceCount int64             `json:"count"`
}

// Service API handlers

func getServices(w http.ResponseWriter, _ *http.Request) {
	services, count, err := models.GetAllServices()
	if err != nil {
		services = []*models.Service{}
	}
	resp := servicesResponse{
		Services:     services,
		ServiceCount: count,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	MustEncode(w, resp)
}

func postService(w http.ResponseWriter, r *http.Request) {
	newService := &models.Service{}
	w.Header().Set("Content-Type", "application/json")
	ok := MustValidate(r, newService)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}

	if err := models.AddService(newService); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"Record with pk already exists"})
		return
	}
	w.WriteHeader(http.StatusOK)
	MustEncode(w, newService)
	if err := r.Body.Close(); err != nil {
		log.Println(err)
	}
}

// Get handler for /situation/{name}
func getService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}
	service, err := models.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		MustEncode(w, ErrorResponse{Error: "not found"})
		return
	}
	// here the response is the situation itself
	w.WriteHeader(http.StatusOK)
	MustEncode(w, service)
}

// Delete handler for /situation/{name}
func deleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}

	service, err := models.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		MustEncode(w, ErrorResponse{"The situation does not exist"})
		return
	}

	if err := models.DeleteService(service); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"An error occurred while trying to delete the situation"})
		return
	}
	w.WriteHeader(http.StatusOK)
	MustEncode(w, service)
}

// Situation API handlers

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
