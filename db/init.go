package db

import (
	"fmt"
	"github.com/Abathargh/moody-go/db/dao"
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

	dao.DB = db
	return nil
}

/*
package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)
func main() {
	if err := DBInit(); err != nil {
		log.Println("Couldn't initalize the databse!")
		log.Fatal(err)
	}

	result, size, err := dao.GetAllNodes(nil, 0, 20, "")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found ", size, " records")

	for _, elem := range result {
		fmt.Println(elem)
	}

}
*/
