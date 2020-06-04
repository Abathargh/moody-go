package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	configFile = "conf.json"
)

func mustInitConf() map[string]interface{} {
	var conf map[string]interface{}
	file, fileErr := ioutil.ReadFile(configFile)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	if jsonErr := json.Unmarshal(file, &conf); jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return conf
}
