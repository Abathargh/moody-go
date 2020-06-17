package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"service/model"
	"strconv"
)

type servicesResponse struct {
	Services     []*model.Service `json:"services"`
	ServiceCount int64            `json:"count"`
}

func getServices(w http.ResponseWriter, _ *http.Request) {
	services, count, err := model.GetAllServices()
	if err != nil {
		services = []*model.Service{}
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
	newService := &model.Service{}
	w.Header().Set("Content-Type", "application/json")
	ok := MustValidate(r, newService)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		MustEncode(w, ErrorResponse{"Bad syntax"})
		return
	}

	if err := model.AddService(newService); err != nil {
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
	service, err := model.GetService(id)
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

	service, err := model.GetService(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		MustEncode(w, ErrorResponse{"The situation does not exist"})
		return
	}

	if err := model.DeleteService(service); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		MustEncode(w, ErrorResponse{"An error occurred while trying to delete the situation"})
		return
	}
	w.WriteHeader(http.StatusOK)
	MustEncode(w, service)
}
