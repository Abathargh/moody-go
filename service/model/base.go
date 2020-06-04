package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var (
	NotFound          = fmt.Errorf("record Not Found")
	UpdateFailedError = fmt.Errorf("model update error")
	InsertFailedError = fmt.Errorf("model insert error")
	DeleteFailedError = fmt.Errorf("model delete error")

	DB *gorm.DB
)
