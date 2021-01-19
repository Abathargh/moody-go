package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var (
	NotFound          = fmt.Errorf("record Not Found")
	UpdateFailedError = fmt.Errorf("models update error")
	InsertFailedError = fmt.Errorf("models insert error")
	DeleteFailedError = fmt.Errorf("models delete error")

	DB *gorm.DB
)
