package models

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

/*
DB Table Details
table main.node_meta {

    [0 ]  nodemeta_id               INTEGER              nullable: true len: 0
    [1 ]  nodemeta_group            text                 nullable: true len: 0
    [2 ]  nodemeta_datatype         text                 nullable: true len: 0
    [3 ]  nodemeta_valuenumber      int                  nullable: true len: 0
}

*/

// NodeMeta struct is a row record of the node_meta table in the main database
type NodeMeta struct {
	Group       string `gorm:"column:nodemeta_group;primary_key" json:"nodemeta_group"`
	Datatype    string `gorm:"column:nodemeta_datatype" json:"nodemeta_datatype"`
	ValueNumber int64  `gorm:"column:nodemeta_valuenumber" json:"nodemeta_valuenumber"`
}

// TableName sets the insert table name for this struct type
func (n *NodeMeta) TableName() string {
	return "node_meta"
}
