package communication

import (
	"encoding/json"
	"gateway/models"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

const (
	eventData     = "data"
	eventActuator = "actuator"
)

type SocketioServer struct {
	Server       *socketio.Server
	webappSocket socketio.Conn
}

type ServiceDataEvent struct {
	Service string  `json:"service"`
	Data    float64 `json:"data"`
}

type ActuatorConnectedEvent struct {
	IP string `json:"ip"`
}

// Creates a new SocketioServer instance
// Its API is used by the data and actuators listeners
// to forward messages to the webapp
func NewSocketioServer() (*SocketioServer, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	return &SocketioServer{Server: server}, nil
}

// Catches the connect event coming from the webapp and
// saves a reference to the connection within the object.
func (ss *SocketioServer) Serve() error {
	ss.Server.OnConnect("/", func(s socketio.Conn) error {
		ss.webappSocket = s
		return nil
	})
	if err := ss.Server.Serve(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Closes the Socketio Server
func (ss *SocketioServer) Close() error {
	if err := ss.Server.Close(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Marshals and forwards data referring a service event
func (ss *SocketioServer) ForwardServiceData(evt ServiceDataEvent) {
	if ss.webappSocket == nil {
		return
	}
	jsonData, _ := json.Marshal(&evt)
	ss.webappSocket.Emit(eventData, jsonData)
}

// Gets a mapping of an actuator server that just connected to the gateway and forwards
// it to the webapp
func (ss *SocketioServer) ForwardActuatorData(evt ActuatorConnectedEvent) {
	if ss.webappSocket == nil {
		return
	}
	mappings, err := models.GetActuatorMapping(evt.IP)
	if err != nil {
		log.Println(err)
		return
	}

	actuator := models.Actuator{
		IP:       evt.IP,
		Mappings: mappings.DataTable,
	}

	jsonData, _ := json.Marshal(&actuator)
	ss.webappSocket.Emit(eventActuator, jsonData)
}
