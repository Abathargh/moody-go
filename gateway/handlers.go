package main

import (
	"bytes"
	"fmt"
	"gateway/communication"
	"gateway/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Forwards incoming requests that need to query an external service to the API GW
func forwardToApiGW(w http.ResponseWriter, r *http.Request) {
	url := communication.ApiGatewayAddress
	url.Path = r.URL.Path
	w.Header().Set("Content-Type", "application/json")

	newReq, err := http.NewRequest(r.Method, url.String(), r.Body)
	if err != nil {
		log.Println(err)
		models.RespondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	newReq.Header.Set("Host", r.Host)
	newReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
	for header, headerValues := range r.Header {
		for _, headerValue := range headerValues {
			newReq.Header.Add(header, headerValue)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		log.Println(err)
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(buf.Bytes())
}

// [GET] /neural_state
// returns a view of the dataset state of the app, containing information about
// the dataset engine state and the current dataset in use
func getNeuralState(w http.ResponseWriter, r *http.Request) {
	state := communication.NeuralState
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, state)
}

func DatasetKeysIfExists(state models.NeuralState) ([]string, error) {
	url := communication.ApiGatewayAddress
	url.Path = fmt.Sprintf("/dataset/%s", state.Dataset)
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, models.NotFound
	}

	var meta models.DatasetMeta
	if err := models.ReadAndDecode(resp.Body, &meta); err != nil {
		return nil, err
	}
	return meta.Keys, nil
}

// [PUT] /neural_state {mode: DataMode, dataset: str}
// sets the dataset engine to the passed mode, using the passed dataset
func setNeuralState(w http.ResponseWriter, r *http.Request) {
	newState := models.NeuralState{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &newState); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}

	if newState.Mode == models.Stopped {
		communication.NeuralState = models.NeuralState{
			Mode:    models.Stopped,
			Dataset: "",
		}
	} else {
		// Not stopped, proceed witch dataset check
		_, err := DatasetKeysIfExists(newState)
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
	}

	communication.NeuralState = newState
	log.Printf("New dataset state: %+v\n", communication.NeuralState)
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, newState)
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
	url := communication.ApiGatewayAddress
	url.Path = fmt.Sprintf("/service/%d", servReq.Service)
	resp, err := http.Get(url.String())
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
// changes the value of the current situation; this is used in the collecting mode of the dataset engine
func setSituation(w http.ResponseWriter, r *http.Request) {
	situationRequest := models.SituationRequest{}
	w.Header().Set("Content-Type", "application/json")
	if ok := models.MustValidate(r, &situationRequest); !ok {
		models.RespondWithError(w, http.StatusUnprocessableEntity, "bad syntax")
		return
	}

	url := communication.ApiGatewayAddress
	url.Path = fmt.Sprintf("/situation/%d", situationRequest.SituationId)
	resp, err := http.Get(url.String())
	if err != nil || resp.StatusCode == http.StatusInternalServerError {
		models.RespondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		models.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}
	newSituation := &models.Situation{}
	if err := models.ReadAndDecode(resp.Body, newSituation); err != nil {
		models.RespondWithError(w, http.StatusBadRequest, "server error")
		return
	}
	communication.Situation = newSituation

	// Respond with new situation
	situationResp := &models.SituationResponse{IsSet: true, Situation: communication.Situation}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, situationResp)
}

// [GET] /situation
// gets the value of the current situation
func getSituation(w http.ResponseWriter, r *http.Request) {
	var isSet bool
	if communication.Situation == nil {
		isSet = false
	} else {
		isSet = true
	}
	situationResp := &models.SituationResponse{
		IsSet:     isSet,
		Situation: communication.Situation,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, situationResp)
}

// [GET] /data_table
// returns an object containing the current readings relative to each active service
func getDataTable(w http.ResponseWriter, _ *http.Request) {
	table := make(map[string]string)
	for _, service := range communication.ActiveServices.AsSlice() {
		val, ok := communication.DataTable.Table()[service]
		if !ok {
			table[service] = "-"
		} else {
			// check if float is actually int
			if val == float64(int(val)) {
				table[service] = fmt.Sprintf("%d", int(val))
			} else {
				table[service] = fmt.Sprintf("%.1f", val)
			}
		}
	}
	resp := models.ServiceTableResponse{
		Table:        table,
		ServiceCount: len(table),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	models.MustEncode(w, resp)
}
