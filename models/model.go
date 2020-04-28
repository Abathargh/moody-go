package models

import (
	"database/sql"
)

type Node struct {
	NodeMacaddress sql.NullInt64  `gorm:"column:node_macaddress;primary_key" json:"node_macaddress"`
	NodeType       sql.NullString `gorm:"column:node_type" json:"node_type"`
	NodeMetadata   sql.NullInt64  `gorm:"column:node_metadata" json:"node_metadata"`
}

func (n *Node) TableName() string {
	return "node"
}

type NodeMetum struct {
	NodemetaID          sql.NullInt64  `gorm:"column:nodemeta_id;primary_key" json:"nodemeta_id"`
	NodemetaGroup       sql.NullString `gorm:"column:nodemeta_group" json:"nodemeta_group"`
	NodemetaDatatype    sql.NullString `gorm:"column:nodemeta_datatype" json:"nodemeta_datatype"`
	NodemetaValuenumber sql.NullInt64  `gorm:"column:nodemeta_valuenumber" json:"nodemeta_valuenumber"`
}

func (n *NodeMetum) TableName() string {
	return "node_meta"
}

type NeuralDatum struct {
	NeuralID     sql.NullInt64  `gorm:"column:neural_id;primary_key" json:"neural_id"`
	NeuralData   sql.NullString `gorm:"column:neural_data" json:"neural_data"`
	NeuralMetaid sql.NullInt64  `gorm:"column:neural_metaid" json:"neural_metaid"`
}

func (n *NeuralDatum) TableName() string {
	return "neural_data"
}

type NeuralMetum struct {
	NeuralmetaID           sql.NullInt64  `gorm:"column:neuralmeta_id;primary_key" json:"neuralmeta_id"`
	NeuralmetaDatatypelist sql.NullString `gorm:"column:neuralmeta_datatypelist" json:"neuralmeta_datatypelist"`
	NeuralmetaSessionname  sql.NullString `gorm:"column:neuralmeta_sessionname" json:"neuralmeta_sessionname"`
}

func (n *NeuralMetum) TableName() string {
	return "neural_meta"
}

type Action struct {
	ActionID       sql.NullInt64  `gorm:"column:action_id;primary_key" json:"action_id"`
	ActionName     sql.NullString `gorm:"column:action_name" json:"action_name"`
	ActionValue    sql.NullString `gorm:"column:action_value" json:"action_value"`
	ActionMetadata sql.NullInt64  `gorm:"column:action_metadata" json:"action_metadata"`
}

func (a *Action) TableName() string {
	return "action"
}

type SituationSetting struct {
	SituationID sql.NullInt64 `gorm:"column:situation_id;primary_key" json:"situation_id"`
	ActionID    sql.NullInt64 `gorm:"column:action_id" json:"action_id"`
}

func (s *SituationSetting) TableName() string {
	return "situation_setting"
}

type Situation struct {
	SituationID   sql.NullInt64  `gorm:"column:situation_id;primary_key" json:"situation_id"`
	SituationName sql.NullString `gorm:"column:situation_name" json:"situation_name"`
}

func (s *Situation) TableName() string {
	return "situation"
}
