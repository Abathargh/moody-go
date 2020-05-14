package db

import (
	"github.com/Abathargh/moody-go/models"
	"github.com/smallnest/gen/dbmeta"
)

func GetService(id int64) (*models.Service, error) {
	var service models.Service

	if err := DB.First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func AddService(service models.Service) error {
	if err := DB.Save(service).Error; err != nil {
		return InsertFailedError
	}

	return nil
}

func UpdateService(id interface{}, updated *models.Node) error {
	service := &models.Service{}
	if err := DB.First(service, id).Error; err != nil {
		return NotFound
	}

	if err := dbmeta.Copy(service, updated); err != nil {
		return UpdateFailedError
	}

	if err := DB.Save(service).Error; err != nil {
		return UpdateFailedError
	}

	return nil
}

func DeleteService(id int64) error {
	service := &models.Service{}
	if err := DB.First(service, id).Error; err != nil {
		return NotFound
	}

	if err := DB.Delete(service).Error; err != nil {
		return DeleteFailedError
	}

	return nil
}
