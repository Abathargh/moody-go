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
table main.neural_data {

    [0 ]  neural_id                 INTEGER              nullable: true len: 0
    [1 ]  neural_data               varchar(100)         nullable: true len: 0
    [2 ]  neural_metaid             integer              nullable: true len: 0
}

*/

// NeuralDatum struct is a row record of the neural_data table in the main database
type NeuralDatum struct {
	NeuralID     sql.NullInt64  `gorm:"column:neural_id;primary_key" json:"neural_id"`
	NeuralData   sql.NullString `gorm:"column:neural_data" json:"neural_data"`
	NeuralMetaid sql.NullInt64  `gorm:"column:neural_metaid" json:"neural_metaid"`
}

// TableName sets the insert table name for this struct type
func (n *NeuralDatum) TableName() string {
	return "neural_data"
}
