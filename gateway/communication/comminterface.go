package communication

import (
	"errors"
	"gateway/models"
	"log"
	"net/url"
)

// This is the interface that defines the API for a Client
// To implement a new protocol, add a Client struct to the package
// and append it to the client mappings in the var section of this file
type Client interface {
	Init(conf interface{}) error
	Connect() error
	Forward(situation string) error
	SwitchToActuatorServer()
	StopTicker()
	Close()
}

// A ticker periodically sends a message to each actuator using a given protocol
// Its Tick method must be called in a separate goroutine since it's blocking
type Ticker interface {
	Tick()
	Done()
}

const (
	bufferSize = 5
)

var (
	clientMapping = map[string]Client{
		"mqtt": &MQTTClient{},
	}

	ApiGatewayAddress url.URL

	DataTable      *models.DataTable
	ActiveServices *models.SynchronizedStringSet
	NeuralState    models.NeuralState
	Situation      *models.Situation
	ActuatorIPs    *models.SynchronizedStringSet

	exposer   *WebSocketForwarder
	forwarder *DataForwarder

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

	// Check that the right type of data was passed through the conf file
	var ok bool
	gwAddr, ok := conf["apiGateway"].(string)
	if !ok {
		return errors.New("wrong syntax for the apiGateway field in conf.json")
	}

	var err error
	apiGWAddr, err := url.Parse(gwAddr)
	if err != nil {
		return err
	}
	ApiGatewayAddress = *apiGWAddr

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

func CommSwitchToActuatorServer() {
	for _, ifc := range clientMapping {
		ifc.SwitchToActuatorServer()
	}
}

func CommStopTickers() {
	for _, ifc := range clientMapping {
		ifc.StopTicker()
	}
}

func CommClose() {
	for _, ifc := range clientMapping {
		ifc.Close()
	}
}
