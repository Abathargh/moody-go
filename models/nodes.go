package models

import (
	"encoding/json"
	"errors"
)

// Virtualization of the generic Node
type Node struct {
	MacAddress string   `json:"macaddress"`
	Type       string   `json:"type"`
	Group      string   `json:"group"`
	Datatypes  []string `json:"datatypes"`
}

func NodeFromJson(jsonData []byte) (Node, error) {
	//Parse json
	var node Node

	if err := json.Unmarshal(jsonData, &node); err != nil {
		return Node{}, errors.New("an error occurred while unmarshalling a greet packet")
	}
	return node, nil
}

type NodeData struct {
	MacAddress string `json:"macaddress"`
	Data       []byte `json:"data"`
}

func DataFromJson(jsonData []byte) (NodeData, error) {
	//Parse json
	var data NodeData

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return NodeData{}, errors.New("an error occurred while unmarshalling a data packet")
	}
	return data, nil
}
