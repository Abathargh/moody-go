package communication

import (
	"encoding/json"
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

func NewSocketioServer() (*SocketioServer, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	return &SocketioServer{Server: server}, nil
}

func (ss *SocketioServer) Serve() error {
	ss.Server.OnConnect("/", func(s socketio.Conn) error {
		ss.webappSocket = s
		log.Println(ss.webappSocket)
		return nil
	})
	if err := ss.Server.Serve(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ss *SocketioServer) Close() error {
	if err := ss.Server.Close(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ss *SocketioServer) ForwardServiceData(evt ServiceDataEvent) {
	jsonData, _ := json.Marshal(&evt)
	ss.webappSocket.Emit(eventData, jsonData)
}

func (ss *SocketioServer) ForwardActuatorData(evt ActuatorConnectedEvent) {
	jsonData, _ := json.Marshal(&evt)
	ss.webappSocket.Emit(eventActuator, jsonData)
}
