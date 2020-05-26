package models

import (
	"encoding/json"
)

// Node struct is a row record of the node table in the main database
type Node struct {
	MacAddress string `gorm:"column:node_macaddress;primary_key";json:"macaddress"`
	Type       string `gorm:"column:node_type";json:"type"`
	Group      string `gorm:"column:node_metadata";json:"group"`
}

func NodeFromJson(jsonData []byte) (Node, error) {
	//Parse json
	var node Node

	if err := json.Unmarshal(jsonData, &node); err != nil {
		return Node{}, err
	}
	return node, nil
}

// TableName sets the insert table name for this struct type
func (n *Node) TableName() string {
	return "node"
}
