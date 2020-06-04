package main

import (
	"gateway/communication"
	"gateway/model"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type NeuralState int64

const (
	Stopped NeuralState = iota
	Collecting
	Predicting
)

var (
	DataTable      *model.DataTable
	ActiveServices map[string]*model.Service
	neuralState    NeuralState
)

func main() {
	// Explicit logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set up a safe exit mechanism
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Load conf file
	conf := mustInitConf()

	if err := communication.StartCommInterface(conf); err != nil {
		log.Println("an error occurred while starting the communication interface")
		log.Fatal(err)
	}

	if err := communication.CommConnect(); err != nil {
		log.Println("an error occurred while connecting through the communication interface")
		log.Fatal(err)
	}

	<-quit
	communication.CommClose()
	log.Println("Gateway - Shutting Down")
}
