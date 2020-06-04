package main

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"os/signal"
	"service/model"
	"syscall"
)

const (
	dbDialect = "postgres"
)

var (
	conf *ServiceConf
)

func main() {
	// Explicit logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set up a safe exit mechanism
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Load conf file
	conf = mustInitConf()

	// Db connection
	postgresConn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		conf.DbHost, conf.DbPort, conf.DbUser, conf.DbName, conf.DbPass)
	db, err := gorm.Open(dbDialect, postgresConn)
	if err != nil {
		log.Fatal(err)
	}
	model.DB = db

	// Start the server
	server := HttpListenAndServer(conf.ServerPort)
	go func() { log.Fatal(server.ListenAndServe()) }()

	<-quit
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		log.Println("an error occurred while attempting to close the model connection")
		log.Fatal(err)
	}
	log.Println("Service service - Shutting Down")
}
