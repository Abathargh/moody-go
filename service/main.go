package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"service/model"
	"syscall"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	dbDialect = "postgres"
	retries   = 5
)

var (
	conf *ServiceConf
)

func main() {
	// Explicit logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	// Set up a safe exit mechanism
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Load conf file
	conf = mustInitConf()

	// Db connection
	attempt := 0
	for attempt < retries {
		postgresConn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			conf.DbHost, conf.DbPort, conf.DbUser, conf.DbName, conf.DbPass)
		db, err := gorm.Open(dbDialect, postgresConn)
		if err == nil {
			model.DB = db
			break
		}
		attempt += 1
	}

	if attempt == retries {
		log.Fatal("Couldn't connect to the database")
	}

	// Start the server
	server := HttpListenAndServer(conf.ServerPort)
	go func() { log.Fatal(server.ListenAndServe()) }()

	<-quit
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if err := model.DB.Close(); err != nil {
		log.Println("an error occurred while attempting to close the models connection")
		log.Fatal(err)
	}
	log.Println("Service service - Shutting Down")
}
