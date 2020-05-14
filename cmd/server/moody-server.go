package main

import (
	"fmt"
	confinit "github.com/Abathargh/moody-go"
	"github.com/Abathargh/moody-go/communication"
	"github.com/Abathargh/moody-go/db"
	"github.com/Abathargh/moody-go/models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type MockObserver struct {
	DataChannel chan models.DataEvent
}

func (mo *MockObserver) ListenForUpdates() {
	for data := range mo.DataChannel {
		fmt.Println(data)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

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

	obs := MockObserver{
		DataChannel: make(chan models.DataEvent),
	}
	communication.DataTable.Attach(obs.DataChannel)
	go obs.ListenForUpdates()

	<-quit
	communication.CommClose()
	if err := db.DB.Close(); err != nil {
		log.Println("an error occurred while attempting to close the db connection")
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
