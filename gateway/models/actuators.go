package models

import (
	"fmt"
	"log"
	"net/http"
)

type IpToDelete struct {
	IP string `json:"ip"`
}

type Mapping struct {
	Situation int `json:"situation"`
	Action    int `json:"action"`
}

type IPMapping struct {
	IP        string `json:"ip"`
	Situation int    `json:"situation"`
	Action    int    `json:"action"`
}

type Mappings struct {
	DataTable []map[string]int `json:"mappings"`
}

type Actuator struct {
	IP       string           `json:"ip" validate:"nonzero"`
	Mappings []map[string]int `json:"mappings"`
}

type ActuatorState struct {
	Mode bool `json:"mode"`
}

type ConnectedServerActuatorList struct {
	Actuators       []Actuator `json:"actuatorList" validate:"nonzero"`
	ActuatorsNumber int        `json:"size" validate:"nonzero"`
}

// Helper function that returns the mappings of a specific actuator server
func GetActuatorMapping(ip string) (*Mappings, error) {
	mappingEndpoint := fmt.Sprintf("http://%s/mapping", ip)
	log.Println("Act IP: ", mappingEndpoint)
	resp, err := http.Get(mappingEndpoint)
	if err != nil {
		return nil, err
	}

	mappings := &Mappings{}
	if err := ReadAndDecode(resp.Body, mappings); err != nil {
		return nil, err
	}
	return mappings, nil
}
