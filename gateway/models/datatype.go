package models

type DataType struct {
	Id   uint64 `gorm:"column:id;primary_key"`
	Type string `gorm:"column:name;unique" validate:"nonzero"`
}

func (n *DataType) TableName() string {
	return "datatype"
}
