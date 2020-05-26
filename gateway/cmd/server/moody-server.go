package main

import (
	"context"
	"fmt"
	confinit "github.com/Abathargh/moody-go"
	"github.com/Abathargh/moody-go/communication"
	"github.com/Abathargh/moody-go/db"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	quit := make(chan os.Signal, 1)
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

	if err := confinit.LoadServices(); err != nil {
		log.Println("an error occurred while retrieving the activated services")
		log.Fatal(err)
	}

	server := communication.HttpListenAndServer()
	go func() { log.Fatal(server.ListenAndServe()) }()
	<-quit

	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}

	communication.CommClose()
	if err := db.DB.Close(); err != nil {
		log.Println("an error occurred while attempting to close the db connection")
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
