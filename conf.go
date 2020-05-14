package confinit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
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
