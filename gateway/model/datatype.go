package model

type DataType struct {
	Id   uint64 `gorm:"column:id;primary_key"`
	Type string `gorm:"column:name;unique" validate:"nonzero"`
}
