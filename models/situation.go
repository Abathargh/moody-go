package models

// Situation struct is a row record of the situation table in the main database
type Situation struct {
	SituationId   uint64 `gorm:"column:id;primary_key" json:"id"`
	SituationName string `gorm:"column:name;unique" json:"name" validate:"nonzero"`
}

// TableName sets the insert table name for this struct type
func (s *Situation) TableName() string {
	return "situation"
}
