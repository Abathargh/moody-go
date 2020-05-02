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
table main.situation {

    [0 ]  situation_id              INTEGER              nullable: true len: 0
    [1 ]  situation_name            text                 nullable: true len: 0
}

*/

// Situation struct is a row record of the situation table in the main database
type Situation struct {
	SituationID   sql.NullInt64  `gorm:"column:situation_id;primary_key" json:"situation_id"`
	SituationName sql.NullString `gorm:"column:situation_name" json:"situation_name"`
}

// TableName sets the insert table name for this struct type
func (s *Situation) TableName() string {
	return "situation"
}
