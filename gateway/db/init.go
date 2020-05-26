package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/go-homedir"
	"os"
)

const (
	configFolder = ".moody"
	dbName       = "moody.db"
)

func Init() error {
	homeDir, dirErr := homedir.Dir()
	if dirErr != nil {
		return dirErr
	}

	var path = fmt.Sprintf("%v%v%v%v", homeDir, string(os.PathSeparator), configFolder, string(os.PathSeparator))

	if _, err := os.Stat(path + dbName); os.IsNotExist(err) {
		// moody.db does not exists: let's create it
		file, err := os.Create(path + dbName)
		if err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}

	db, err := gorm.Open("sqlite3", path+dbName)
	if err != nil {
		return err
	}
	DB = db
	return nil
}
