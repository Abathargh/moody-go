package model

type ServiceState struct {
	Id    uint64 `gorm:"column:id;primary_key"`
	State string `gorm:"column:state;unique" validate:"nonzero"`
}
