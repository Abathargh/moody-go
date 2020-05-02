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
table main.neural_meta {

    [0 ]  neuralmeta_id             INTEGER              nullable: true len: 0
    [1 ]  neuralmeta_datatypelist   varchar(100)         nullable: true len: 0
    [2 ]  neuralmeta_sessionname    varchar(50)          nullable: true len: 0
}

*/

// NeuralMetum struct is a row record of the neural_meta table in the main database
type NeuralMetum struct {
	NeuralmetaID           sql.NullInt64  `gorm:"column:neuralmeta_id;primary_key" json:"neuralmeta_id"`
	NeuralmetaDatatypelist sql.NullString `gorm:"column:neuralmeta_datatypelist" json:"neuralmeta_datatypelist"`
	NeuralmetaSessionname  sql.NullString `gorm:"column:neuralmeta_sessionname" json:"neuralmeta_sessionname"`
}

// TableName sets the insert table name for this struct type
func (n *NeuralMetum) TableName() string {
	return "neural_meta"
}
