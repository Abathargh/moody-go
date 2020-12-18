package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	dataFolder = "data"
	configFile = "conf.json"
)

func mustInitConf() map[string]interface{} {
	var conf map[string]interface{}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fullPath := filepath.Join(pwd, dataFolder, configFile)
	file, fileErr := ioutil.ReadFile(fullPath)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	if jsonErr := json.Unmarshal(file, &conf); jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return conf
}
