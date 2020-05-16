package db

import (
	"github.com/Abathargh/moody-go/models"
)

func GetAllServices() (services []*models.Service, totalRows int64, err error) {
	services = []*models.Service{}
	serviceOrm := DB.Model(&models.Service{})
	serviceOrm.Count(&totalRows)

	if err = serviceOrm.Find(&services).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return services, totalRows, nil
}

func GetService(id uint64) (*models.Service, error) {
	var service models.Service

	if err := DB.First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func AddService(service *models.Service) error {
	if err := DB.Save(service).Error; err != nil {
		return InsertFailedError
	}

	return nil
}

func DeleteService(s *models.Service) error {
	service := &models.Service{}
	if err := DB.First(service, s.Id).Error; err != nil {
		return NotFound
	}

	if err := DB.Delete(service).Error; err != nil {
		return DeleteFailedError
	}

	return nil
}
