package models

import (
	"errors"
)

type Service struct {
	Id   uint64 `gorm:"column:id;primary_key" json:"id"`
	Name string `gorm:"column:name;unique" validate:"nonzero" json:"name"`
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
