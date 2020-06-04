package model

type StateValue uint64

type Service struct {
	Id           uint64       `gorm:"column:id;primary_key" json:"id"`
	Name         string       `gorm:"column:name;unique" validate:"nonzero" json:"name"`
	DataType     DataType     `gorm:"foreignkey:Type" validate:"-" json:"-"`
	Type         uint64       `gorm:"column:datatype" validate:"min=1,max=3" json:"type"`
	ServiceState ServiceState `gorm:"foreignkey:State" validate:"-" json:"-"`
	State        StateValue   `gorm:"column:state" validate:"min=0,max=1" json:"state"`
}

type ServicesResponse struct {
	Services     []*Service `json:"services"`
	ServiceCount int64      `json:"count"`
}
