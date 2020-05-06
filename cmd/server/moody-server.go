package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mitchellh/go-homedir"
	"log"
	"moody-go/communication"
	"moody-go/db"
	"moody-go/db/dao"
	"os"
	"os/signal"
	"syscall"
)

const (
	configFolder = ".moody"
	configFile   = "conf.json"
)

func ConfInit() (map[string]interface{}, error) {
	// Add check for connectionConfig != nil?
	homeDir, dirErr := homedir.Dir()
	if dirErr != nil {
		return nil, dirErr
	}

	var jsonConfig string

	var path = fmt.Sprintf("%v%v%v%v", homeDir, string(os.PathSeparator), configFolder, string(os.PathSeparator))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dirErr := os.MkdirAll(path, os.ModePerm)
		if dirErr != nil {
			return nil, dirErr
		}
	}

	file, fileErr := os.OpenFile(path+configFile, os.O_RDONLY|os.O_CREATE, 0666)

	if fileErr != nil {
		return nil, fileErr
	}

	scanner := bufio.NewScanner(file)
	success := scanner.Scan()
	for success {
		jsonConfig += scanner.Text()
		success = scanner.Scan()
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, scanErr
	}

	connectionConfig := make(map[string]interface{})
	if jsonErr := json.Unmarshal([]byte(jsonConfig), &connectionConfig); jsonErr != nil {
		return nil, jsonErr
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return connectionConfig, nil
}

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

	if err := communication.CommConnect(); err != nil {
		log.Println("an error occurred while connecting thtough communication interface")
		log.Fatal(err)
	}

	<-quit
	communication.CommClose()
	if err := dao.DB.Close(); err != nil {
		log.Println("an error occurred while attempting to close the db connection")
		log.Fatal(err)
	}
	fmt.Println("Bye!")
}
