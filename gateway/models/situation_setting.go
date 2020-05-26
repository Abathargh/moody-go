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
table main.situation_setting {

    [0 ]  situation_id              integer              nullable: true len: 0
    [1 ]  action_id                 integer              nullable: true len: 0
}

*/

// SituationSetting struct is a row record of the situation_setting table in the main database
type SituationSetting struct {
	SituationID sql.NullInt64 `gorm:"column:situation_id;primary_key" json:"situation_id"`
	ActionID    sql.NullInt64 `gorm:"column:action_id" json:"action_id"`
}

// TableName sets the insert table name for this struct type
func (s *SituationSetting) TableName() string {
	return "situation_setting"
}
