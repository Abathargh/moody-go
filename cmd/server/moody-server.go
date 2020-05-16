package main

import (
	"fmt"
	confinit "github.com/Abathargh/moody-go"
	"github.com/Abathargh/moody-go/communication"
	"github.com/Abathargh/moody-go/db"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf, err := confinit.ConfInit()
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

	if err := communication.CommConnect(); err != nil {
		log.Println("an error occurred while connecting through the communication interface")
		log.Fatal(err)
	}

	communication.HttpListenAndServe()

	communication.CommClose()
	if err := db.DB.Close(); err != nil {
		log.Println("an error occurred while attempting to close the db connection")
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
