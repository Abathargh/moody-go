package models

type Service struct {
	Id       uint64   `gorm:"column:id;primary_key" json:"id"`
	Name     string   `gorm:"column:name;unique" validate:"nonzero" json:"name"`
	DataType DataType `gorm:"foreignkey:Type" validate:"-" json:"-"`
	Type     uint64   `gorm:"column:datatype" validate:"min=1,max=3" json:"type"`
}

func (n *Service) TableName() string {
	return "service"
}
