package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
)

const (
	brokAddr = "localhost"
	port     = 1883
	hostId   = "Main"
	topic    = "test"
)

func main() {
	/*
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		modules.ConfInit()
		config := modules.GetConfig()

		fmt.Println(config)

		<-quit
		client.Disconnect(200)
		fmt.Println("Bye!")
	*/

	home, _ := homedir.Dir()
	db, err := gorm.Open("sqlite3", home+string(os.PathSeparator)+
		".moody"+string(os.PathSeparator)+"moody.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	var situations []Situation

	db.Find(&situations)
	fmt.Println("%v", situations)

}
