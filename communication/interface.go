package communication

import (
	"github.com/Abathargh/moody-go/models"
	"log"
)

type Client interface {
	Init(conf interface{}) error
	Connect() error
	Forward(situation string) error
	Update(group, rule string) error
	Close()
}

var (
	clientMapping = map[string]Client{
		"mqtt": &MQTTClient{},
	}
	ConnectedNodes *models.ConnectedList
	DataTable      *models.DataTable
	Services       map[string]*models.Service
)

func StartCommInterface(conf map[string]interface{}) error {
	for name, commIfc := range clientMapping {
		// protocolConf => interface{} json object containing info about the specific protocol
		protocolConf, ok := conf[name]
		if !ok {
			log.Printf("protocol %s not supported\n", name)
			continue
		}

		err := commIfc.Init(protocolConf)
		if err != nil {
			log.Printf("an error occurred while trying to initalize client %s\n", name)
			return err
		}
		if err := commIfc.Connect(); err != nil {
			return err
		}
		ConnectedNodes = &models.ConnectedList{}
		DataTable = &models.DataTable{Data: make(map[string]string)}
		Services = make(map[string]*models.Service)
	}
	return nil
}

func CommConnect() error {
	for name, ifc := range clientMapping {
		if err := ifc.Connect(); err != nil {
			log.Printf("couldn't connect using the client for protocol %s\n", name)
			return err
		}
	}
	return nil
}

func CommForward(situation string) {
	for _, ifc := range clientMapping {
		ifc.Forward(situation)
	}
}

func CommUpdate(group, rule string) {
	for _, ifc := range clientMapping {
		ifc.Update(group, rule)
	}
}

func CommClose() {
	for _, ifc := range clientMapping {
		ifc.Close()
	}
}
