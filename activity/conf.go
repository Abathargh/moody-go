package main

import (
	"encoding/json"
	"gopkg.in/validator.v2"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	configFile = "conf.json"
)

type ServiceConf struct {
	ServerPort string `json:"server_port" validate:"nonzero"`
	DbHost     string `json:"db_host" validate:"nonzero"`
	DbPort     int    `json:"db_port" validate:"nonzero"`
	DbUser     string `json:"db_user" validate:"nonzero"`
	DbPass     string `json:"db_pass" validate:"nonzero"`
	DbName     string `json:"db_name" validate:"nonzero"`
}

func mustInitConf() *ServiceConf {
	var conf ServiceConf
	file, fileErr := ioutil.ReadFile(configFile)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	if jsonErr := json.Unmarshal(file, &conf); jsonErr != nil {
		log.Fatal(jsonErr)
	}
	if err := validator.Validate(conf); err != nil {
		log.Fatal(err)
	}
	return &conf
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Encode struct to json, if something goes wrong the application exits.
// Used internally to crate the body for the http responses.
func MustEncode(w http.ResponseWriter, response interface{}) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}
}

// Decode and validate json data to struct, returning a bool
// containing the outcome of the operation.
func MustValidate(r *http.Request, dest interface{}) (outcome bool) {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		log.Println(err)
		return false
	}
	if err := validator.Validate(dest); err != nil {
		log.Println(err)
		return false
	}
	return true
}
