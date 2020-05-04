package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"moody-go/communication"
	"moody-go/db"
	"moody-go/db/dao"
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
	fmt.Println(conf)

	if err := db.Init(); err != nil {
		log.Println("an error occurred while initializing the database")
		log.Fatal(err)
	}

	if err := communication.StartCommInterface(conf); err != nil {
		log.Println("an error occurred while starting the communication interface")
		log.Fatal(err)
	}
	defer communication.CommClose()

	if err := communication.CommConnect(); err != nil {
		log.Println("an error occurred while connecting thtough communication interface")
		log.Fatal(err)
	}

	<-quit
	if err := dao.DB.Close(); err != nil {
		log.Println("an error occurred while attempting to close the db connection")
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
