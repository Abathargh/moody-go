package main

import (
	"fmt"
	"gateway/communication"
	"gateway/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// [PUT] /neural_state {mode: DataMode, dataset: str}
// sets the neural engine to the passed mode, using the passed dataset
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
	log.Printf("New neural state: %+v\n", communication.NeuralState)
	w.WriteHeader(http.StatusOK)
}

// [GET] /actuator_mode
// returns the current actuator mode
func getActuatorMode(w http.ResponseWriter, r *http.Request) {
	resp := models.ActuatorState{Mode: communication.ActuatorMode}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, resp)
}

// [POST] /actuator_mode {mode: bool}
// sets the actuator mode
//		true: Actuator mode, the actuators act as MQTT nodes receiving data to actuate
//		false: Server Mode, the actuators act as HTTP Servers exposing their mappings API
func setActuatorMode(w http.ResponseWriter, r *http.Request) {
	newState := models.ActuatorState{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &newState); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}

	communication.ActuatorMode = newState.Mode
	if communication.ActuatorMode {
		// Clear list on server mode trigger in order not to carry old data
		for _, ip := range communication.ActuatorIPs.AsSlice() {
			end := fmt.Sprintf("http://%s/end", ip)
			log.Println(end)
			go http.PostForm(end, url.Values{})
		}
		communication.CommStopTickers()
		communication.ActuatorIPs.Clear()
	} else {
		// Send a msg on the actuator mode topic to switch every actutor to server mode
		communication.CommSwitchToActuatorServer()
	}

	resp := models.ActuatorState{Mode: communication.ActuatorMode}
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, resp)
	log.Printf("New actuator mode: %v", communication.ActuatorMode)
}

// [GET] /actuators
// returns a list of the connected actuators in server mode, along with their mappings
func getActuators(w http.ResponseWriter, r *http.Request) {
	actuators := []models.Actuator{}

	if !communication.ActuatorIPs.Empty() {
		actuatorServerList := communication.ActuatorIPs.AsSlice()
		for _, ip := range actuatorServerList {
			mappings, err := models.GetActuatorMapping(ip)
			if err != nil {
				models.RespondWithError(w, http.StatusInternalServerError, "server error")
				return
			}
			actuator := models.Actuator{IP: ip, Mappings: mappings.DataTable}
			actuators = append(actuators, actuator)
		}
	}

	resp := models.ConnectedServerActuatorList{
		Actuators:       actuators,
		ActuatorsNumber: len(actuators),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, resp)
}

// Helper function for the mapping handlers that checks that the specified actuator is
// connected to the gateway
func isActuatorConnected(ip string) bool {
	return communication.ActuatorIPs.Contains(ip)
}

// [POST] /actuators {ip: string, situation: int, action: int}
// adds a mapping to the mapping list of the target actuator, identified by the passed ip
func addMapping(w http.ResponseWriter, r *http.Request) {
	newMapping := models.IPMapping{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &newMapping); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}

	if ok := isActuatorConnected(newMapping.IP); !ok {
		models.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}

	data := url.Values{}
	data.Set("situation", strconv.Itoa(newMapping.Situation))
	data.Add("action", strconv.Itoa(newMapping.Action))

	actuatorUrl := fmt.Sprintf("http://%s/mapping", newMapping.IP)
	resp, err := http.PostForm(actuatorUrl, data)

	if err != nil || resp.StatusCode == http.StatusConflict {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "resource already exists")
		return
	}

	okMapping := &models.Mapping{}
	if err := models.ReadAndDecode(resp.Body, okMapping); err != nil {
		models.RespondWithError(w, http.StatusBadRequest, "server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, okMapping)
}

// [DELETE] /actuators {ip: string}
// delete the mappings of the actuator identified by the passed ip
func deleteMappings(w http.ResponseWriter, r *http.Request) {
	ipToDelete := models.IpToDelete{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &ipToDelete); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}

	if ok := isActuatorConnected(ipToDelete.IP); !ok {
		models.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}

	client := &http.Client{}
	urlDelete := fmt.Sprintf("http://%s/mapping", ipToDelete.IP)
	req, _ := http.NewRequest("DELETE", urlDelete, nil)

	_, err := client.Do(req)
	if err != nil {
		fullMsg := fmt.Sprintf("server error %s", err)
		models.RespondWithError(w, http.StatusBadRequest, fullMsg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// [POST] /sensor_service {serviceId: int}
// activates the service identified by serviceId
func activateService(w http.ResponseWriter, r *http.Request) {
	servReq := models.ServiceRequest{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &servReq); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}
	resp, err := http.Get(fmt.Sprintf("%sservice/%d", communication.ApiGatewayAddress, servReq.Service))
	if err != nil || resp.StatusCode == http.StatusInternalServerError {
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

// [GET] /sensor_service
// returns a list of the currently active services
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

// [DELETE] /sensor_service {name: string}
// deactivates the service identified by the passed name
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

// [PUT] /situation {id: int}
// changes the value of the current situation; this is used in the collecting mode of the neural engine
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

// [GET] /situation
// gets the value of the current situation
func getSituation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, communication.Situation)
}
