package modules

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
)

const (
	configFolder = ".moody"
	configFile   = "conf"
	dbName       = "moody.db"
)

var (
	connectionConfig map[string]interface{}
	connectedNodes   []Node
	nodesInfo        map[string]interface{}
)

func ConfInit() error {
	// Add check for connectionConfig != nil?
	homeDir, dirErr := homedir.Dir()
	if dirErr != nil {
		return dirErr
	}

	var jsonConfig string

	var path = fmt.Sprintf("%v%v%v%v", homeDir, string(os.PathSeparator), configFolder, string(os.PathSeparator))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dirErr := os.MkdirAll(path, os.ModePerm)
		if dirErr != nil {
			return dirErr
		}
	}

	file, fileErr := os.OpenFile(path+configFile, os.O_RDONLY|os.O_CREATE, 0666)
	if fileErr != nil {
		return fileErr
	}

	scanner := bufio.NewScanner(file)
	success := scanner.Scan()
	for success == true {
		jsonConfig += scanner.Text()
		success = scanner.Scan()
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return scanErr
	}

	if jsonErr := json.Unmarshal([]byte(jsonConfig), &connectionConfig); jsonErr != nil {
		return jsonErr
	}

	return nil
}

func GetConfig() map[string]interface{} {
	return connectionConfig
}

func GetConnected() []Node {
	return connectedNodes
}
