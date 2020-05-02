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
table main.node_meta {

    [0 ]  nodemeta_id               INTEGER              nullable: true len: 0
    [1 ]  nodemeta_group            text                 nullable: true len: 0
    [2 ]  nodemeta_datatype         text                 nullable: true len: 0
    [3 ]  nodemeta_valuenumber      int                  nullable: true len: 0
}

*/

// NodeMetum struct is a row record of the node_meta table in the main database
type NodeMetum struct {
	NodemetaID          sql.NullInt64  `gorm:"column:nodemeta_id;primary_key" json:"nodemeta_id"`
	NodemetaGroup       sql.NullString `gorm:"column:nodemeta_group" json:"nodemeta_group"`
	NodemetaDatatype    sql.NullString `gorm:"column:nodemeta_datatype" json:"nodemeta_datatype"`
	NodemetaValuenumber sql.NullInt64  `gorm:"column:nodemeta_valuenumber" json:"nodemeta_valuenumber"`
}

// TableName sets the insert table name for this struct type
func (n *NodeMetum) TableName() string {
	return "node_meta"
}
