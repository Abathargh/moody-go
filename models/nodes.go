package models

import (
	"encoding/json"
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
		return Node{}, err
	}
	return node, nil
}

type NodeData struct {
	MacAddress string `json:"macaddress"`
	Data       string `json:"data"`
}

func DataFromJson(jsonData []byte) (NodeData, error) {
	// TODO Add validation
	//Parse json
	var data NodeData

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return NodeData{}, err
	}
	return data, nil
}
