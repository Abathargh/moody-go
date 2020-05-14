package models

type Service struct {
	Id   int64  `gorm:"column:id;primary_key"`
	Name string `gorm:"column:name"`
	Type string `gorm:"column:type"`
}

func (n *Service) TableName() string {
	return "service"
}
