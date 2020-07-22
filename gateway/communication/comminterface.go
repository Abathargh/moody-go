package communication

import (
	"errors"
	"gateway/models"
	"log"
)

type Client interface {
	Init(conf interface{}) error
	Connect() error
	Forward(situation string) error
	Update(group, rule string) error
	Close()
}

const (
	bufferSize = 5
)

var (
	clientMapping = map[string]Client{
		"mqtt": &MQTTClient{},
	}

	ApiGatewayAddress string
	WebAppAddress string

	DataTable      *models.DataTable
	ActiveServices *models.SynchronizedStringSet
	NeuralState    models.NeuralState
	Situation      *models.Situation
	ActuatorIPs    *models.SynchronizedStringSet

	exposer   *SocketioExposer
	forwarder *DataForwarder

	SioServer *SocketioServer

	ActuatorMode = true
)

func StartCommInterface(conf map[string]interface{}) error {
	DataTable = models.NewDataTable()
	ActiveServices = models.NewSynchronizedStringSet()
	ActuatorIPs = models.NewSynchronizedStringSet()

	NeuralState = models.NeuralState{
		Mode:    models.Stopped,
		Dataset: "",
	}

	// Start the socketio Server
	var socketioErr error
	SioServer, socketioErr = NewSocketioServer()
	if socketioErr != nil {
		return socketioErr
	}

	// Check that the right type of data was passed through the conf file
	var ok bool
	ApiGatewayAddress, ok = conf["apiGateway"].(string)
	if !ok {
		return errors.New("wrong syntax for the apiGateway field in conf.json")
	}

	WebAppAddress, ok = conf["webAppAddr"].(string)
	if !ok {
		return errors.New("wrong syntax for the webAppAddr field in conf.json")
	}

	// Start each communication interface
	for name, commIfc := range clientMapping {
		// protocolConf => interface{} json object containing info about the specific protocol
		protocolConf, ok := conf[name]
		if !ok {
			log.Printf("protocol %s not supported\n", name)
			continue
		}
		if err := commIfc.Init(protocolConf); err != nil {
			log.Printf("an error occurred while trying to initalize client %s\n", name)
			return err
		}
		if err := commIfc.Connect(); err != nil {
			return err
		}
	}

	// Start the observers' goroutines
	exposer = NewSocketioExposer(bufferSize)
	forwarder = NewDataForwarder(bufferSize)
	go forwarder.ListenForUpdates()
	go exposer.ListenForUpdates()

	return nil
}

func CommConnect() error {
	for name, ifc := range clientMapping {
		if err := ifc.Connect(); err != nil {
			log.Printf("couldn't connect using the client for protocol %s\n", name)
			return err
		}
	}
	go SioServer.Serve()
	return nil
}

func CommForward(situation string) {
	for _, ifc := range clientMapping {
		err := ifc.Forward(situation)
		if err != nil {
			// Just for now
			continue
		}
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
	SioServer.Close()
}
