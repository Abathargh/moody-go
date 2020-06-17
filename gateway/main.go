package main

import (
	"context"
	"gateway/communication"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// Start the server
	serverPort, ok := conf["serverPort"]
	if !ok {
		log.Fatal("serverPort parameter not present in conf.json")
	}

	port, ok := serverPort.(string)
	if !ok {
		log.Fatal("Wrong parameter type for serverPort")
	}
	server := HttpListenAndServer(port)
	go func() { log.Fatal(server.ListenAndServe()) }()

	<-quit
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}

	communication.CommClose()
	log.Println("Gateway - Shutting Down")
}
