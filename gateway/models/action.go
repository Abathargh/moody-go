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
table main.action {

    [0 ]  action_id                 INTEGER              nullable: true len: 0
    [1 ]  action_name               text                 nullable: true len: 0
    [2 ]  action_value              text                 nullable: true len: 0
    [3 ]  action_metadata           int                  nullable: true len: 0
}

*/

// ActionPerformed struct is a row record of the action table in the main database
type Action struct {
	ActionID       sql.NullInt64  `gorm:"column:action_id;primary_key" json:"action_id"`
	ActionName     sql.NullString `gorm:"column:action_name" json:"action_name"`
	ActionValue    sql.NullString `gorm:"column:action_value" json:"action_value"`
	ActionMetadata sql.NullInt64  `gorm:"column:action_metadata" json:"action_metadata"`
}

// TableName sets the insert table name for this struct type
func (a *Action) TableName() string {
	return "action"
}
