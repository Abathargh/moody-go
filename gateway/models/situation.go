package models

type SituationRequest struct {
	SituationId uint64 `json:"id" validate:"nonzero,min=0"`
}

type SituationResponse struct {
	IsSet     bool       `json:"isSet"` // IsSet == false => Situation = nil
	Situation *Situation `json:"situation"`
}

// Situation struct is a row record of the situation table in the main database
type Situation struct {
	SituationId   uint64 `gorm:"column:id;primary_key" json:"id"`
	SituationName string `gorm:"column:name;unique" json:"name" validate:"nonzero"`
}

// TableName sets the insert table name for this struct type
func (s *Situation) TableName() string {
	return "situation"
}
