package model

import (
	"errors"
)

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

func GetAllServices() (services []*Service, totalRows int64, err error) {
	services = []*Service{}
	serviceOrm := DB.Model(&Service{})
	serviceOrm.Count(&totalRows)

	if err = serviceOrm.Find(&services).Error; err != nil {
		err := errors.New("not found")
		return nil, -1, err
	}

	return services, totalRows, nil
}

func GetActivatedServices() (services []*Service, err error) {
	services = []*Service{}
	serviceOrm := DB.Model(&Service{})

	if err = serviceOrm.Where("state = ?", Started).Find(&services).Error; err != nil {
		err = NotFound
		return nil, err
	}

	return services, nil
}

func GetService(id uint64) (*Service, error) {
	var service Service

	if err := DB.First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func AddService(service *Service) error {
	if err := DB.Save(service).Error; err != nil {
		return InsertFailedError
	}
	return nil
}

func DeleteService(s *Service) error {
	service := &Service{}
	if err := DB.First(service, s.Id).Error; err != nil {
		return NotFound
	}

	if err := DB.Delete(service).Error; err != nil {
		return DeleteFailedError
	}
	return nil
}

func PatchStateService(service *Service) error {
	if err := DB.Save(service).Error; err != nil {
		return UpdateFailedError
	}
	return nil
}
