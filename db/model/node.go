package model

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
table main.node {

    [0 ]  node_macaddress           text                  nullable: true len: 0
    [1 ]  node_type                 text                 nullable: true len: 0
    [2 ]  node_metadata             int                  nullable: true len: 0
}

*/

// Node struct is a row record of the node table in the main database
type Node struct {
	NodeMacaddress sql.NullString `gorm:"column:node_macaddress;primary_key" json:"node_macaddress"`
	NodeType       sql.NullString `gorm:"column:node_type" json:"node_type"`
	NodeMetadata   sql.NullInt64  `gorm:"column:node_metadata" json:"node_metadata"`
}

// TableName sets the insert table name for this struct type
func (n *Node) TableName() string {
	return "node"
}
