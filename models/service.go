package models

type StateValue uint64

const (
	Stopped = iota
	Started
)

type Service struct {
	Id           uint64       `gorm:"column:id;primary_key" json:"id"`
	Name         string       `gorm:"column:name;unique" validate:"nonzero" json:"name"`
	DataType     DataType     `gorm:"foreignkey:Type" validate:"-" json:"-"`
	Type         uint64       `gorm:"column:datatype" validate:"min=1,max=3" json:"type"`
	ServiceState ServiceState `gorm:"foreignkey:State" validate:"-" json:"-"`
	State        StateValue   `gorm:"column:state" validate:"min=0,max=1" json:"state"`
}

func (n *Service) TableName() string {
	return "service"
}
