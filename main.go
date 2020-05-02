package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"moody-go/communication"
	"moody-go/db"
	"moody-go/db/dao"
	"moody-go/models"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	conf, err := ConfInit()
	if err != nil {
		log.Println("an error occurred while reading the configuration file")
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Println("an error occurred while initializing the database")
		log.Fatal(err)
	}

	commIfc := &communication.CommInterface{
		ConnectedNodes: &models.ConnectedList{},
		DataTable:      &models.DataTable{},
	}
	if err := commIfc.Init(conf); err != nil {
		log.Println("an error occurred while initialing the communication interface")
		log.Fatal(err)
	}
	defer commIfc.Close()
	commIfc.Listen()

	<-quit
	if err := dao.DB.Close(); err != nil {
		log.Println("an error occurred while attempting to close the db connection")
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
