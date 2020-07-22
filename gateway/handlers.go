package main

import (
	"fmt"
	"gateway/communication"
	"gateway/models"
	"log"
	"net/http"
)

func setNeuralState(w http.ResponseWriter, r *http.Request) {
	newState := models.NeuralState{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &newState); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}
	_, err := newState.DatasetKeysIfExists()
	if err != nil {
		if err == models.NotFound {
			models.RespondWithError(w, http.StatusNotFound, "not found")
		} else {
			models.RespondWithError(w, http.StatusInternalServerError, "server error")
		}
		communication.NeuralState = models.NeuralState{
			Mode:    models.Stopped,
			Dataset: "",
		}
		return
	}
	communication.NeuralState = newState
}

func setActuatorMode(w http.ResponseWriter, r *http.Request) {
	newState := models.ActuatorState{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &newState); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}
	communication.ActuatorMode = newState.Mode
	// Clear list on server mode trigger so not to carry old data
	if !communication.ActuatorMode {
		communication.ActuatorIPs = nil
	}
}

func activateService(w http.ResponseWriter, r *http.Request) {
	servReq := models.ServiceRequest{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &servReq); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}
	resp, err := http.Get(fmt.Sprintf("%sservice/%d", communication.ApiGatewayAddress, servReq.Service))
	if err != nil || resp.StatusCode == http.StatusInternalServerError {
		if err != nil {
			log.Println(err)
		}
		models.RespondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		models.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}
	newService := &models.Service{}
	if err := models.ReadAndDecode(resp.Body, newService); err != nil {
		models.RespondWithError(w, http.StatusBadRequest, "server error")
		return
	}
	communication.ActiveServices.Add(newService.Name)
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, newService)
}

func getActivatedServices(w http.ResponseWriter, _ *http.Request) {
	activeServices := communication.ActiveServices.AsSlice()
	resp := models.ServicesResponse{
		Services:     activeServices,
		ServiceCount: len(activeServices),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, resp)
}

func deactivateService(w http.ResponseWriter, r *http.Request) {
	servReq := models.ServiceByName{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &servReq); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}
	if !communication.ActiveServices.Contains(servReq.ServiceName) {
		models.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}
	communication.ActiveServices.Remove(servReq.ServiceName)
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, servReq)
}

func setSituation(w http.ResponseWriter, r *http.Request) {
	situationRequest := models.SituationRequest{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &situationRequest); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}

	resp, err := http.Get(fmt.Sprintf("%ssituation/%d", communication.ApiGatewayAddress, situationRequest.SituationId))
	if err != nil || resp.StatusCode == http.StatusInternalServerError {
		models.RespondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		models.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}
	newSituation := &models.Situation{}
	if err := models.ReadAndDecode(r.Body, newSituation); err != nil {
		models.RespondWithError(w, http.StatusBadRequest, "server error")
		return
	}
	communication.Situation = newSituation
}

func getSituation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, communication.Situation)
}
