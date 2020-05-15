package models

// Situation struct is a row record of the situation table in the main database
type Situation struct {
	SituationName string `gorm:"column:situation_name;primary key" json:"name"`
}

// TableName sets the insert table name for this struct type
func (s *Situation) TableName() string {
	return "situation"
}
